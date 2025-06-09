package pokemon

func BuildCardFromPokemon(poke Pokemon, moves []Move) Card {
	var hp, defense, attack, speed int
	for _, stat := range poke.Stats {
		switch stat.StName.Name {
		case "hp":
			hp = stat.BaseSt
		case "defense":
			defense = stat.BaseSt
		case "attack":
			attack = stat.BaseSt
		case "speed":
			speed = stat.BaseSt
		}
	}

	stamina := speed * 2

	hp = hp + int(float64(hp)*0.5) // 50% of HP is added to the card

	types := make([]string, len(poke.Types))
	for i, t := range poke.Types {
		types[i] = t.Type.Name
	}

	// Sprite
	sprite := poke.Sprites.FrontDflt

	return Card{
		Name:    poke.Name,
		HP:      hp,
		HPMax:   hp,
		Stamina: stamina,
		Defense: defense,
		Attack:  attack,
		Moves:   moves,
		Types:   types,
		Sprite:  sprite,
	}
}

