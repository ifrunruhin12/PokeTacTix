package battle

import (
	"context"
	"encoding/json"
	"fmt"

	"pokemon-cli/internal/pokemon"

	"github.com/jackc/pgx/v5"
)

// evolutionEvent records a completed evolution so the caller can emit the
// metric and log entry only after the transaction commits — otherwise a later
// statement failure would roll the evolution back while the counter and log
// permanently record it.
type evolutionEvent struct {
	userID     int
	cardID     int
	from, into string
	level      int
}

// evolvedBaseStats returns the base stat values to store on player_cards after
// evolution. Delegates to the same helpers that pokemon.Pokemon.ToCard() uses so
// the two code paths can never diverge.
func evolvedBaseStats(target *pokemon.Pokemon) (hp, attack, defense, speed int) {
	return target.CardBaseStats()
}

// applyEvolution checks whether the player card should evolve at newLevel and,
// if so, rewrites the card row (name, pokemon_id, sprite, types, base stats) in
// the rewards transaction. The card keeps its current moves and level/XP.
// pokemonID must be non-zero; callers are responsible for resolving it before
// entering the transaction (to avoid holding a DB connection during a PokeAPI
// fetch). Returns the target Pokemon, or nil when no evolution applies.
// Metric incrementing and logging are the caller's job, post-commit.
func applyEvolution(
	ctx context.Context,
	tx pgx.Tx,
	pokemonSvc pokemon.PokemonService,
	cardID, userID, pokemonID, newLevel int,
	currentPokemonName string,
) (*pokemon.Pokemon, error) {
	target, err := pokemonSvc.GetEvolutionForLevel(ctx, pokemonID, newLevel)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, nil
	}

	// Derive base stats the same way new cards are built from pokemon data.
	baseHP, baseAttack, baseDefense, baseSpeed := evolvedBaseStats(target)

	typesJSON, err := json.Marshal(target.Types)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal evolved types: %w", err)
	}

	// Moves intentionally carry over: this game's evolution is a species change
	// with new base stats, not a move-set reset.
	_, err = tx.Exec(ctx, `
		UPDATE player_cards
		SET pokemon_name = $1, pokemon_id = $2, sprite = $3, types = $4,
		    base_hp = $5, base_attack = $6, base_defense = $7, base_speed = $8
		WHERE id = $9 AND user_id = $10
	`,
		target.Name, target.ID, target.SpriteURL, typesJSON,
		baseHP, baseAttack, baseDefense, baseSpeed,
		cardID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to evolve card %d into %s: %w", cardID, target.Name, err)
	}

	return target, nil
}
