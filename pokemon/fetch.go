package pokemon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func FetchPokemon(name string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(name)
	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, fmt.Errorf("failed to fetch pokemon data: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Pokemon{}, fmt.Errorf("Pokemon \"%s\" not found. Please check the name and try again", name)
	}

	defer resp.Body.Close()

	var poke Pokemon

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&poke); err != nil {
		return Pokemon{}, fmt.Errorf("failed to decode pokemon data: %w", err)
	}

	return poke, nil
}

