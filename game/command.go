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
	HaveCard        bool
	Round           int
	PlayerActiveIdx int
	AIActiveIdx     int
	CardMovePlayer  int
	CardMoveAI      int
	CurrentMovetype string
}

func CommandList(state *GameState) {
	if state.BattleStarted {
		fmt.Println("You are in a battle now. Use 'command --in-game' to see the available commands")
		return
	}

	fmt.Println("Available commands:")
	fmt.Println("1. search   - Search for a Pokémon by name and see its card")
	fmt.Println("2. battle   - Start a 5v5 match against AI (get 5 random cards)")
	fmt.Println("3. command  - Show this command list\n\t\t--in-battle (Type 'command --in-game' to see the commands available in a battle)")
	fmt.Println("4. exit     - Exit the game")
}

func CommandSearch(scanner *bufio.Scanner, state *GameState) {
	if state.BattleStarted {
		fmt.Println("No searching pokemons while a battle is in progress.")
		return
	}
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
	if state.BattleStarted {
		fmt.Println("A battle is already in progress.")
		return
	}

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
	fmt.Println("You are in a Battle with", state.AI.Name)
	fmt.Println("Use 'card all' to view your cards and use 'choose' to choose your card to play. (You cannot see AI's cards.)")
	state.InBattle = true
	state.Round = 1
}

func CommandExit() {
	fmt.Println("Thanks for playing! Goodbye.")
	os.Exit(0)
}

func CommandDefault() {
	fmt.Println("Unknown command. Type 'command' to see all available commands.")
}
