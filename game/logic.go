package game

import (
	"bufio"
	"fmt"
	"math/rand"
	_"strings"
)

// ProcessTurn handles both player and AI moves, updates state, and returns a result string.
func ProcessTurn(playerMove string, aiMove string, player *Player, ai *Player) string {
	// TODO: Parse moves, call CalculateDamage, update HP/stamina, handle defend/surrender/sacrifice
	// For now, just print the moves chosen
	return fmt.Sprintf("Player chose: %s | AI chose: %s", playerMove, aiMove)
}

// CalculateDamage computes the damage dealt from attacker to defender using the chosen move and type multiplier.
func CalculateDamage(attacker *Player, defender *Player, moveIndex int) int {
	// TODO: Use attack stat, move power, RollDamagePercent, and TypeMultiplier
	return 0 // placeholder
}

// ChooseAIMove picks a move for the AI (random for now).
func ChooseAIMove(ai *Player, player *Player) string {
	// TODO: Add smarter logic later
	return "attack 1" // placeholder
}

func StartBattle(player *Player, ai *Player) {
	fmt.Println("[Battle logic will go here]")
}

// StartBattleLoop asks the player to choose a Pokémon and AI chooses randomly, then starts the turn loop.
func StartBattleLoop(scanner *bufio.Scanner, state *GameState) {
	chosenIdx := state.PlayerActiveIdx
	// AI chooses randomly
	state.AIActiveIdx = rand.Intn(len(state.AI.Deck))
	fmt.Printf("You chose %s. AI chose its Pokémon.\n", state.Player.Deck[chosenIdx].Name)
	StartTurnLoop(scanner, state)
}

// StartTurnLoop is a placeholder for the turn loop logic.
func StartTurnLoop(scanner *bufio.Scanner, state *GameState) {
	fmt.Println("[Turn loop will go here]")
	fmt.Println(state)
	fmt.Println(state.Player.Deck[state.PlayerActiveIdx].Name)
	fmt.Println(state.Player.Deck[state.PlayerActiveIdx].HP)
	fmt.Println(state.AI.Deck[state.AIActiveIdx].Name)
	fmt.Println(state.AI.Deck[state.AIActiveIdx].HP)
	fmt.Println(state.Player.Deck[state.PlayerActiveIdx].Moves[state.CardMovePlayer].Name) //print the attack move name for the pokemon
}
