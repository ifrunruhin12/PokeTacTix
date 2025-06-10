package game

import (
	"bufio"
	"fmt"
	"strings"
)

func HandleCommand(input string, scanner *bufio.Scanner, state *GameState) {
	command := strings.ToLower(strings.TrimSpace(input))
	
	switch command {
	case "command":
		CommandList(state)
	case "command --in-battle":
		CommandListInBattle(state)
	case "search":
		CommandSearch(scanner, state)
	case "version":
		CommandVersion()
	case "battle":
		CommandBattle(scanner, state)
	case "card all":
		// Only allow in 5v5 battles
		if state.BattleStarted && state.BattleMode == "1v1" {
			fmt.Println("'card all' is not available in 1v1 battles. Use 'card' to see your single card.")
			return
		}
		CommandCardAll(state)
	case "card":
		CommandCurrentCard(state)
	case "choose":
		// Only allow in 5v5 battles
		if state.BattleStarted && state.BattleMode == "1v1" {
			fmt.Println("'choose' is not available in 1v1 battles. You only have one card to battle with.")
			return
		}
		CommandCardChooser(scanner, state)
	case "attack":
		CommandMovesAttack(scanner, state)
	case "defend":
		CommandDefendMove(state)
	case "exit":
		CommandExit(state)
	case "switch":
		// Only allow in 5v5 battles
		if state.BattleStarted && state.BattleMode == "1v1" {
			fmt.Println("'switch' is not available in 1v1 battles. You only have one Pokémon.")
			return
		}
		CommandSwitch(scanner, state)
	case "surrender":
		CommandSurrender(scanner, state, false)
	case "surrender all":
		CommandSurrender(scanner, state, true)
	default:
		CommandDefault()
	}
}
