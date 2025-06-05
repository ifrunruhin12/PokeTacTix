package game

import (
	"bufio"
	"fmt"
	"os"
	"pokemon-cli/pokemon"
	"strings"
)

type GameState struct {
	BattleStarted   bool
	Player          *Player
	AI              *Player
	PlayerName      string
	InBattle        bool
	Round           int
	PlayerActiveIdx int
	AIActiveIdx     int
}

func CommandList() {
	fmt.Println("Available commands:")
	fmt.Println("1. search   - Search for a Pokémon by name and see its card")
	fmt.Println("2. battle   - Start a 5v5 match against AI (get 5 random cards)")
	fmt.Println("3. card all - Show all your cards (only after starting battle)")
	fmt.Println("4. command  - Show this command list")
	fmt.Println("5. exit     - Exit the game")
}

func CommandSearch(scanner *bufio.Scanner) {
	fmt.Print("Enter the name of the Pokémon: ")
	if !scanner.Scan() {
		return
	}
	name := strings.TrimSpace(scanner.Text())
	poke, moves, err := pokemon.FetchPokemon(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	card := pokemon.BuildCardFromPokemon(poke, moves)
	PrintCard(card)
}

func CommandBattle(scanner *bufio.Scanner, state *GameState) {
	fmt.Print("Enter your name: ")
	if !scanner.Scan() {
		return
	}
	playerName := strings.TrimSpace(scanner.Text())
	fmt.Printf("%s, do you want to start the battle? (y/n): ", playerName)
	if !scanner.Scan() {
		return
	}
	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("Battle cancelled.")
		return
	}
	state.Player = NewPlayer(playerName, FetchRandomDeck())
	state.AI = NewPlayer("AI", FetchRandomDeck())
	state.BattleStarted = true
	state.PlayerName = playerName
	fmt.Println("Battle started! You and AI each have 5 random Pokémon cards.")
	fmt.Println("Use 'card all' to view your cards. (You cannot see AI's cards.)")
	state.InBattle = true
	state.Round = 1
	StartBattleLoop(scanner, state)
}

func CommandCardAll(state *GameState) {
	if !state.BattleStarted {
		fmt.Println("You need to start a battle first with the 'battle' command.")
		return
	}
	fmt.Printf("%s, here are your cards:\n", state.PlayerName)
	for _, card := range state.Player.AllCards() {
		PrintCard(card)
	}
}

func CommandExit() {
	fmt.Println("Thanks for playing! Goodbye.")
	os.Exit(0)
}

func CommandDefault() {
	fmt.Println("Unknown command. Type 'command' to see all available commands.")
}

func HandleCommand(input string, scanner *bufio.Scanner, state *GameState) {
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "command":
		CommandList()
	case "search":
		CommandSearch(scanner)
	case "battle":
		CommandBattle(scanner, state)
	case "card all":
		CommandCardAll(state)
	case "exit":
		CommandExit()
	default:
		CommandDefault()
	}
}
