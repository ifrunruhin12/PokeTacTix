package game

import (
	"bufio"
	"fmt"
	"os"
	"pokemon-cli/pokemon"
	"regexp"
	"strings"
)

type GameState struct {
	BattleStarted     bool
	Player            *Player
	AI                *Player
	PlayerName        string
	InBattle          bool
	HaveCard          bool
	Round             int
	PlayerActiveIdx   int
	AIActiveIdx       int
	CardMovePlayer    int
	CardMoveAI        int
	CurrentMovetype   string
	RoundStarted      bool
	SwitchedThisRound bool
	BattleOver        bool
	RoundOver         bool
	SacrificeCount    map[int]int // key: PlayerActiveIdx, value: number of sacrifices for that Pokémon
	LastHpLost        int
	LastStaminaLost   int
	LastDamageDealt   int
	PlayerSurrendered bool // Track if player surrendered the whole battle
	JustSwitched      bool // true if player just switched and hasn't played a round yet
	HasPlayedRound    bool
	TurnNumber        int
	BattleMode        string // "1v1" or "5v5"
}

func CommandList(state *GameState) {
	if state.BattleStarted {
		fmt.Println("You are in a battle now. Use 'command --in-battle' to see the available commands")
		return
	}

	fmt.Println("Available commands:")
	fmt.Println("1. search   - Search for a Pokémon by name and see its card")
	fmt.Println("2. version  - Show the version of the game")
	fmt.Println("3. battle   - Start a 1v1 or 5v5 match against AI (get random cards)")
	fmt.Println("4. command  - Show this command list\n\t\t--in-battle (Type 'command --in-battle' to see the commands available in a battle)")
	fmt.Println("5. exit     - Exit the game")
}

func CommandVersion() {
	fmt.Println("Version 0.0.1(alpha)")
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

func validatePlayerName(name string) bool {
	if len(name) == 0 || len(name) > 10{
		return false
	}
	validName := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return validName.MatchString(name)
}

func CommandBattle(scanner *bufio.Scanner, state *GameState) {
	if state.BattleStarted {
		fmt.Println("A battle is already in progress.")
		return
	}

	var playerName string
	for {
		fmt.Print("Enter your name (max 10 characters, letters/numbers only, no spaces): ")
		if !scanner.Scan() {
			return
		}
		playerName = strings.TrimSpace(scanner.Text())
		
		if validatePlayerName(playerName) {
			break
		}
		
		if len(playerName) == 0 {
			fmt.Println("Name cannot be empty. Please try again.")
		} else if len(playerName) > 10 {
			fmt.Println("Name is too long. Maximum 10 characters allowed. Please try again.")
		} else {
			fmt.Println("Invalid name. Only letters and numbers are allowed (no spaces or special characters). Please try again.")
		}
	}

	for {
		fmt.Print("Choose battle mode (1v1 or 5v5): ")
		if !scanner.Scan() {
			return
		}
		battleMode := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if battleMode == "1v1" || battleMode == "5v5" {
			state.BattleMode = battleMode
			break
		}
		fmt.Println("Invalid battle mode. Please enter '1v1' or '5v5'.")
	}

	if state.BattleMode == "1v1" {
		fmt.Printf("%s, do you want to start the 1v1 battle? (y/n): ", playerName)
	} else {
		fmt.Printf("%s, do you want to start the 5v5 battle? (y/n): ", playerName)
	}

	if !scanner.Scan() {
		return
	}

	confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("Battle cancelled.")
		return
	}

	state.BattleStarted = true
	state.PlayerName = playerName
	state.InBattle = true
	state.Round = 1

	if state.BattleMode == "1v1" {
		playerDeck := []pokemon.Card{FetchRandomCard()}
		aiDeck := []pokemon.Card{FetchRandomCard()}
		state.Player = &Player{Name: playerName, Deck: playerDeck}
		state.AI = &Player{Name: "AI", Deck: aiDeck}
		fmt.Println("1v1 Battle started! You and AI each have 1 random Pokémon card.")
		fmt.Println("Use 'card' to view your card and get ready to battle!")
		fmt.Println("You are in a Battle with", state.AI.Name)
		StartTurnLoop(scanner, state)
	} else {
		state.Player = NewPlayer(playerName, FetchRandomDeck())
		state.AI = NewPlayer("AI", FetchRandomDeck())
		fmt.Println("5v5 Battle started! You and AI each have 5 random Pokémon cards.")
		fmt.Println("Use 'card all' to view your cards and use 'choose' to choose your card to play. (You cannot see AI's cards.)")
		fmt.Println("You are in a Battle with", state.AI.Name)
		// Note: 5v5 battles don't immediately start the turn loop, they wait for card selection
	}
}


func CommandExit(state *GameState) {
	if state != nil && (state.BattleStarted || state.InBattle) {
		fmt.Println("You need to finish the battle first before exiting. Use 'surrender all' to immediately lose and exit.")
		return
	}
	fmt.Println("Thanks for playing! Goodbye.")
	os.Exit(0)
}

func CommandDefault() {
	fmt.Println("Unknown command. Type 'command' to see all available commands.")
}

// Helper function to fetch a single random card
func FetchRandomCard() pokemon.Card {
	deck := FetchRandomDeck()
	return deck[0] // Return the first card from a random deck
}
