package pokemon

import (
	"fmt"
)

var stNames = map[string]string{
	"hp":      "HP",
	"attack":  "Attack",
	"defense": "Defense",
	"type":    "Type",
}

func DisplayPokemon(poke Pokemon) {
	fmt.Println("\nName:", poke.Name)

	for _, s := range poke.Stats {
		if label, ok := stNames[s.StName.Name]; ok {
			fmt.Printf("%s: %-3d %s\n", label, s.BaseSt, stBar(s.BaseSt))
		}
	}

	fmt.Print("Type(s): ")
	for i, t := range poke.Types {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Println(t.Type.Name)
	}

	fmt.Println("Sprite URL:", poke.Sprites.FrontDflt)
}

func stBar(value int) string {
	barLength := min(value/10, 20)
	return "[" + string(repeat('#', barLength)) + "]"
}

func repeat(char rune, count int) []rune {
	bar := make([]rune, count)
	for i := range bar {
		bar[i] = char
	}
	return bar
}
