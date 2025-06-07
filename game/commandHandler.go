package game

import (
	"bufio"
	"strings"
)

func HandleCommand(input string, scanner *bufio.Scanner, state *GameState) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "command":
		CommandList(state)
	case "command --in-battle":
		CommandListInBattle(state)
	case "search":
		CommandSearch(scanner, state)
	case "battle":
		CommandBattle(scanner, state)
	case "card all":
		CommandCardAll(state)
	case "card":
		CommandCurrentCard(state)
	case "choose":
		CommandCardChooser(scanner, state)
	case "attack":
		CommandMovesAttack(scanner, state)
	case "defend":
		CommandDefendMove(state)
	case "exit":
		CommandExit(state)
	case "switch":
		CommandSwitch(scanner, state)
	case "surrender":
		CommandSurrender(scanner, state, false)
	case "surrender all":
		CommandSurrender(scanner, state, true)
	default:
		CommandDefault()
	}
}
