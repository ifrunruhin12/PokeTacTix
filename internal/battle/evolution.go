package battle

import (
	"context"
	"encoding/json"
	"fmt"

	"pokemon-cli/internal/pokemon"

	"github.com/jackc/pgx/v5/pgconn"
)

type evolutionExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

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

// resolveEvolutionTargets resolves every level-eligible link before the rewards
// transaction starts. Each lookup uses the newly evolved target ID so a large
// level jump can traverse multiple stages in one payout.
func resolveEvolutionTargets(
	ctx context.Context,
	pokemonSvc pokemon.PokemonService,
	pokemonID, newLevel int,
) ([]*pokemon.Pokemon, error) {
	seen := map[int]struct{}{pokemonID: {}}
	var targets []*pokemon.Pokemon

	for {
		target, err := pokemonSvc.GetEvolutionForLevel(ctx, pokemonID, newLevel)
		if err != nil {
			return nil, err
		}
		if target == nil {
			return targets, nil
		}
		if target.ID <= 0 {
			return nil, fmt.Errorf("evolution target has invalid pokemon ID %d", target.ID)
		}
		if _, ok := seen[target.ID]; ok {
			return nil, fmt.Errorf("evolution chain contains a cycle at pokemon ID %d", target.ID)
		}

		seen[target.ID] = struct{}{}
		targets = append(targets, target)
		pokemonID = target.ID
	}
}

// applyEvolution applies targets that were successfully resolved before the
// rewards transaction began. The card keeps its current moves and level/XP.
// Returns the final target Pokemon, or nil when no evolution applies. Metric
// incrementing and logging are the caller's job, post-commit.
func applyEvolution(
	ctx context.Context,
	tx evolutionExecutor,
	cardID, userID int,
	targets []*pokemon.Pokemon,
) (*pokemon.Pokemon, error) {
	if len(targets) == 0 {
		return nil, nil
	}

	for _, target := range targets {
		if target == nil {
			return nil, fmt.Errorf("resolved evolution target is nil")
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
	}

	return targets[len(targets)-1], nil
}
