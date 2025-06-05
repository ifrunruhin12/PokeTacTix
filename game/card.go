package game

import (
	"fmt"
	"pokemon-cli/pokemon"
)

// BuildCardFromPokemon converts a pokemon.Pokemon and its moves into a pokemon.Card for gameplay
func BuildCardFromPokemon(poke pokemon.Pokemon, moves []pokemon.Move) pokemon.Card {
	var hp, defense, attack int
	for _, stat := range poke.Stats {
		switch stat.StName.Name {
		case "hp":
			hp = stat.BaseSt
		case "defense":
			defense = stat.BaseSt
		case "attack":
			attack = stat.BaseSt
		}
	}

	// Default stamina: could be based on HP or a fixed value (will be tweaked as needed)
	stamina := int(float64(hp) * 2.5)
	hp = hp + int(float64(hp)*0.5) // 50% of HP is added to the card

	// Extract types
	types := make([]string, len(poke.Types))
	for i, t := range poke.Types {
		types[i] = t.Type.Name
	}

	// Sprite
	sprite := poke.Sprites.FrontDflt

	return pokemon.Card{
		Name:    poke.Name,
		HP:      hp,
		Stamina: stamina,
		Defense: defense,
		Attack:  attack,
		Moves:   moves,
		Types:   types,
		Sprite:  sprite,
	}
}

// PrintCard pretty-prints a pokemon.Card for debugging and user display
func PrintCard(card pokemon.Card) {
	fmt.Println("====================")
	fmt.Println("Name:", card.Name)
	fmt.Println("HP:", card.HP)
	fmt.Println("Stamina:", card.Stamina)
	fmt.Println("Attack:", card.Attack)
	fmt.Println("Defense:", card.Defense)
	fmt.Println("Types:", card.Types)
	fmt.Println("Sprite URL:", card.Sprite)
	fmt.Println("Moves:")
	for i, move := range card.Moves {
		fmt.Printf("  %d. %s (Power: %d, Stamina Cost: %d, Type: %s)\n", i+1, move.Name, move.Power, move.StaminaCost, move.Type)
	}
	fmt.Println("====================")
}
