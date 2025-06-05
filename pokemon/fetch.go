package pokemon

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

func GetMoves(rawMoves []RawMove) []Move {
	const maxMoves = 4
	perm := rand.Perm(len(rawMoves))
	var gameMoves []Move

	for _, i := range perm {
		moveURL := rawMoves[i].Move.URL
		resp, err := http.Get(moveURL)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var data struct {
			Name  string `json:"name"`
			Power int    `json:"power"`
			Type  struct {
				Name string `json:"name"`
			} `json:"type"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			continue
		}

		if data.Power <= 0 {
			continue
		}

		gameMoves = append(gameMoves, Move{
			Name:        data.Name,
			Power:       data.Power,
			StaminaCost: data.Power / 3, // will go to game logic later
			Type:        data.Type.Name,
		})

		if len(gameMoves) == maxMoves {
			break
		}
	}

	return gameMoves
}

func FetchPokemon(name string) (Pokemon, []Move, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + strings.ToLower(name)
	resp, err := http.Get(url)
	if err != nil {
		return Pokemon{}, nil, fmt.Errorf("failed to fetch pokemon data: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return Pokemon{}, nil, fmt.Errorf("Pokemon \"%s\" not found. Please check the name and try again", name)
	}

	defer resp.Body.Close()

	var poke Pokemon

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&poke); err != nil {
		return Pokemon{}, nil, fmt.Errorf("failed to decode pokemon data: %w", err)
	}

	pokeMoves := GetMoves(poke.Moves)
	return poke, pokeMoves, nil
}
