package commands

import (
	"bufio"
	"fmt"
	"strings"

	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
)

// CommandHandler handles routing of commands to appropriate handlers
type CommandHandler struct {
	gameState *storage.GameState
	renderer  *ui.Renderer
	scanner   *bufio.Scanner

	// Command handlers
	battleCmd     *BattleCommand
	collectionCmd *CollectionCommand
	deckCmd       *DeckCommand
	shopCmd       *ShopCommand
	statsCmd      *StatsCommand
	settingsCmd   *SettingsCommand
}

// NewCommandHandler creates a new command handler with all sub-handlers
func NewCommandHandler(gameState *storage.GameState, renderer *ui.Renderer, scanner *bufio.Scanner) *CommandHandler {
	return &CommandHandler{
		gameState:     gameState,
		renderer:      renderer,
		scanner:       scanner,
		battleCmd:     NewBattleCommand(gameState, renderer, scanner),
		collectionCmd: NewCollectionCommand(gameState, renderer, scanner),
		deckCmd:       NewDeckCommand(gameState, renderer, scanner),
		shopCmd:       NewShopCommand(gameState, renderer, scanner),
		statsCmd:      NewStatsCommand(gameState, renderer, scanner),
		settingsCmd:   NewSettingsCommand(gameState, renderer, scanner),
	}
}

// HandleCommand routes commands to appropriate handlers
// Supports commands: battle, collection, deck, shop, stats, save, help, quit
// Also supports aliases: b, c, d, s, st, h, q
func (ch *CommandHandler) HandleCommand(cmd string, args []string) error {
	// Normalize command to lowercase
	cmd = strings.ToLower(strings.TrimSpace(cmd))

	// Route to appropriate handler
	switch cmd {
	case "battle", "b":
		return ch.battleCmd.StartBattle()

	case "collection", "c":
		if len(args) > 0 && args[0] == "filter" {
			// Future: handle filter arguments
			return ch.collectionCmd.ViewCollection()
		}
		return ch.collectionCmd.ViewCollection()

	case "deck", "d":
		if len(args) > 0 && args[0] == "edit" {
			return ch.deckCmd.EditDeck()
		}
		return ch.deckCmd.ViewDeck()

	case "shop", "s":
		return ch.shopCmd.ViewShop()

	case "stats", "st":
		return ch.statsCmd.ViewStats()

	case "settings", "config":
		return ch.settingsCmd.ViewSettings()

	case "save":
		return ch.handleSave()

	case "help", "h", "?":
		return ch.ShowHelp()

	case "tutorial", "t":
		return ch.ShowTutorial()

	case "quit", "q", "exit":
		return ch.handleQuit()

	default:
		return ch.handleUnknownCommand(cmd)
	}
}

// handleSave manually saves the game state
func (ch *CommandHandler) handleSave() error {
	fmt.Println()
	fmt.Println("Saving game...")

	err := storage.SaveGameState(ch.gameState)
	if err != nil {
		fmt.Println(ui.Colorize(fmt.Sprintf("✗ Failed to save game: %v", err), ui.ColorRed))
		return err
	}

	fmt.Println(ui.Colorize("✓ Game saved successfully!", ui.ColorGreen))
	fmt.Println()
	fmt.Println("Press Enter to continue...")
	ch.scanner.Scan()

	return nil
}

// handleQuit handles the quit command with confirmation
func (ch *CommandHandler) handleQuit() error {
	fmt.Println()
	
	// Ask if user wants to save before quitting
	if ui.ConfirmationPrompt(ch.scanner, "Save before quitting?", false) {
		err := storage.SaveGameState(ch.gameState)
		if err != nil {
			fmt.Println(ui.Colorize(fmt.Sprintf("Warning: Failed to save game: %v", err), ui.ColorRed))
			
			// Ask if they want to quit anyway
			if !ui.ConfirmationPrompt(ch.scanner, "Quit anyway?", true) {
				fmt.Println("Quit cancelled.")
				return nil
			}
		} else {
			fmt.Println(ui.Colorize("✓ Game saved!", ui.ColorGreen))
		}
	}

	fmt.Println()
	fmt.Println(ui.Colorize("Thanks for playing PokeTacTix!", ui.Bold+ui.ColorBrightCyan))
	fmt.Println()

	// Return a special error to signal quit
	return fmt.Errorf("QUIT")
}

// handleUnknownCommand displays an error for unknown commands
func (ch *CommandHandler) handleUnknownCommand(cmd string) error {
	ch.renderer.Clear()
	fmt.Println()
	fmt.Println(ui.Colorize(fmt.Sprintf("Unknown command: '%s'", cmd), ui.ColorRed))
	fmt.Println()
	fmt.Println("Type 'help' or 'h' to see available commands.")
	fmt.Println()
	fmt.Println("Press Enter to continue...")
	ch.scanner.Scan()

	return nil
}

// ShowHelp displays all available commands with descriptions and usage
func (ch *CommandHandler) ShowHelp() error {
	ch.renderer.Clear()

	fmt.Println(ui.RenderLogo())
	fmt.Println()
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println(ui.Colorize("HELP - AVAILABLE COMMANDS", ui.Bold+ui.ColorBrightCyan))
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println()

	// Define command sections
	sections := []struct {
		Title    string
		Commands []struct {
			Name        string
			Aliases     string
			Description string
			Usage       string
		}
	}{
		{
			Title: "GAMEPLAY COMMANDS",
			Commands: []struct {
				Name        string
				Aliases     string
				Description string
				Usage       string
			}{
				{
					Name:        "battle",
					Aliases:     "b",
					Description: "Start a battle (1v1 or 5v5 mode)",
					Usage:       "battle",
				},
				{
					Name:        "collection",
					Aliases:     "c",
					Description: "View your Pokemon collection with filtering and sorting",
					Usage:       "collection",
				},
				{
					Name:        "deck",
					Aliases:     "d",
					Description: "View your current battle deck",
					Usage:       "deck",
				},
				{
					Name:        "deck edit",
					Aliases:     "d edit",
					Description: "Edit your battle deck (add, remove, reorder Pokemon)",
					Usage:       "deck edit",
				},
				{
					Name:        "shop",
					Aliases:     "s",
					Description: "Visit the shop to buy Pokemon with coins",
					Usage:       "shop",
				},
				{
					Name:        "stats",
					Aliases:     "st",
					Description: "View your battle statistics and history",
					Usage:       "stats",
				},
			},
		},
		{
			Title: "SYSTEM COMMANDS",
			Commands: []struct {
				Name        string
				Aliases     string
				Description string
				Usage       string
			}{
				{
					Name:        "save",
					Aliases:     "",
					Description: "Manually save your game progress",
					Usage:       "save",
				},
				{
					Name:        "help",
					Aliases:     "h, ?",
					Description: "Display this help message",
					Usage:       "help",
				},
				{
					Name:        "tutorial",
					Aliases:     "t",
					Description: "Show the game tutorial for beginners",
					Usage:       "tutorial",
				},
				{
					Name:        "quit",
					Aliases:     "q, exit",
					Description: "Exit the game (with save prompt)",
					Usage:       "quit",
				},
			},
		},
	}

	// Display each section
	for _, section := range sections {
		fmt.Println(ui.Colorize(section.Title, ui.Bold+ui.ColorBrightYellow))
		fmt.Println(strings.Repeat("─", 80))
		fmt.Println()

		for _, cmd := range section.Commands {
			// Command name
			cmdName := cmd.Name
			if ch.renderer.ColorSupport {
				cmdName = ui.Colorize(cmdName, ui.Bold+ui.ColorBrightGreen)
			}
			fmt.Printf("  %s", cmdName)

			// Aliases
			if cmd.Aliases != "" {
				aliasText := fmt.Sprintf(" (aliases: %s)", cmd.Aliases)
				if ch.renderer.ColorSupport {
					aliasText = ui.Colorize(aliasText, ui.ColorGray)
				}
				fmt.Print(aliasText)
			}
			fmt.Println()

			// Description
			fmt.Printf("    %s\n", cmd.Description)

			// Usage
			usageText := fmt.Sprintf("    Usage: %s", cmd.Usage)
			if ch.renderer.ColorSupport {
				usageText = ui.Colorize(usageText, ui.ColorCyan)
			}
			fmt.Println(usageText)
			fmt.Println()
		}

		fmt.Println()
	}

	// Display examples section
	fmt.Println(ui.Colorize("EXAMPLES", ui.Bold+ui.ColorBrightYellow))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()

	examples := []struct {
		Command     string
		Description string
	}{
		{"battle", "Start a battle and choose between 1v1 or 5v5 mode"},
		{"b", "Quick shortcut to start a battle"},
		{"deck edit", "Open the deck editor to customize your team"},
		{"collection", "Browse your Pokemon with filtering options"},
		{"shop", "Buy new Pokemon with your earned coins"},
		{"stats", "Check your win rate and battle history"},
		{"save", "Manually save your progress (auto-saves after battles)"},
	}

	for _, example := range examples {
		cmdText := example.Command
		if ch.renderer.ColorSupport {
			cmdText = ui.Colorize(cmdText, ui.Bold+ui.ColorBrightGreen)
		}
		fmt.Printf("  %s\n", cmdText)
		fmt.Printf("    → %s\n\n", example.Description)
	}

	// Display tips section
	fmt.Println(ui.Colorize("TIPS", ui.Bold+ui.ColorBrightYellow))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()

	tips := []string{
		"Win 5v5 battles to select an AI Pokemon to add to your collection",
		"The shop refreshes every 10 battles with new Pokemon",
		"Pokemon gain XP and level up after battles, increasing their stats",
		"Your deck must have exactly 5 Pokemon to battle",
		"Use type advantages in battle for bonus damage",
		"The game auto-saves after every battle and deck change",
	}

	for i, tip := range tips {
		bullet := "•"
		if ch.renderer.ColorSupport {
			bullet = ui.Colorize("•", ui.ColorBrightCyan)
		}
		fmt.Printf("  %s %s\n", bullet, tip)
		if i < len(tips)-1 {
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println()
	fmt.Println("Press Enter to continue...")
	ch.scanner.Scan()

	return nil
}

// ShowHint displays a helpful hint based on game state
// Provides contextual tips for new players
func (ch *CommandHandler) ShowHint() string {
	// Determine which hint to show based on game state
	hints := ch.getContextualHints()

	if len(hints) == 0 {
		return ""
	}

	// Pick a random hint from available hints
	// For simplicity, just return the first one
	hint := hints[0]

	if ch.renderer.ColorSupport {
		return ui.Colorize(fmt.Sprintf("💡 Hint: %s", hint), ui.ColorBrightYellow)
	}
	return fmt.Sprintf("Hint: %s", hint)
}

// getContextualHints returns hints based on current game state
func (ch *CommandHandler) getContextualHints() []string {
	var hints []string

	// Check various game state conditions and provide relevant hints

	// New player hints (low total battles)
	totalBattles := ch.gameState.Stats.TotalBattles1v1 + ch.gameState.Stats.TotalBattles5v5
	if totalBattles == 0 {
		hints = append(hints, "Type 'battle' or 'b' to start your first battle!")
		hints = append(hints, "Use 'help' to see all available commands")
	} else if totalBattles < 5 {
		hints = append(hints, "Try both 1v1 and 5v5 battle modes to see which you prefer")
		hints = append(hints, "Win 5v5 battles to add AI Pokemon to your collection")
	}

	// Deck hints
	if len(ch.gameState.Deck) < 5 {
		hints = append(hints, "Your deck needs exactly 5 Pokemon. Use 'deck edit' to complete it")
	}

	// Collection hints
	if len(ch.gameState.Collection) <= 5 {
		hints = append(hints, "Win 5v5 battles to expand your Pokemon collection")
	} else if len(ch.gameState.Collection) > 10 {
		hints = append(hints, "Use 'collection' to view and filter your growing Pokemon roster")
	}

	// Coins hints
	if ch.gameState.Coins >= 500 {
		hints = append(hints, "You have enough coins to buy rare Pokemon from the shop!")
	} else if ch.gameState.Coins >= 250 {
		hints = append(hints, "Visit the shop to buy uncommon Pokemon with your coins")
	} else if ch.gameState.Coins >= 100 {
		hints = append(hints, "You can afford common Pokemon in the shop now")
	}

	// Shop hints
	if ch.gameState.ShopState.BattlesSinceRefresh >= 8 {
		hints = append(hints, "The shop will refresh in 2 more battles with new Pokemon")
	}

	// Level hints
	if ch.gameState.Stats.HighestLevel >= 10 {
		hints = append(hints, "Your Pokemon are getting strong! Higher levels mean better stats")
	}

	// Win rate hints
	totalWins := ch.gameState.Stats.Wins1v1 + ch.gameState.Stats.Wins5v5
	if totalBattles > 0 {
		winRate := float64(totalWins) / float64(totalBattles)
		if winRate < 0.3 && totalBattles >= 5 {
			hints = append(hints, "Try using type advantages in battle for bonus damage")
			hints = append(hints, "Use the Defend action when your Pokemon is low on HP")
		} else if winRate >= 0.7 && totalBattles >= 10 {
			hints = append(hints, "You're doing great! Keep up the winning streak")
		}
	}

	// General gameplay hints
	generalHints := []string{
		"Type advantages deal 1.5x damage - learn the type chart!",
		"Use Pass to regenerate stamina when you're running low",
		"The game auto-saves after battles, but you can manually save anytime",
		"Check 'stats' to see your battle history and win rate",
		"Build a balanced deck with different types for better coverage",
	}

	// Add a general hint if we don't have many contextual ones
	if len(hints) < 2 && totalBattles > 0 {
		hints = append(hints, generalHints[totalBattles%len(generalHints)])
	}

	return hints
}

// DisplayCommandHints displays available commands in context
// Shows hints based on current game state
func (ch *CommandHandler) DisplayCommandHints() {
	fmt.Println()
	fmt.Println(strings.Repeat("─", 60))

	// Show contextual hint
	hint := ch.ShowHint()
	if hint != "" {
		fmt.Println(hint)
		fmt.Println()
	}

	// Show available commands
	fmt.Print("Available commands: ")
	commands := []string{"battle", "collection", "deck", "shop", "stats", "help", "quit"}

	for i, cmd := range commands {
		if ch.renderer.ColorSupport {
			fmt.Print(ui.Colorize(cmd, ui.ColorBrightCyan))
		} else {
			fmt.Print(cmd)
		}

		if i < len(commands)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println()
	fmt.Println(strings.Repeat("─", 60))
}

// ShowTutorial displays the game tutorial for new players
// Explains game basics: battles, collection, deck, shop
// Shows command examples and battle mechanics
func (ch *CommandHandler) ShowTutorial() error {
	ch.renderer.Clear()

	fmt.Println(ui.RenderLogo())
	fmt.Println()
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println(ui.Colorize("GAME TUTORIAL", ui.Bold+ui.ColorBrightCyan))
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println()

	// Game Overview
	fmt.Println(ui.Colorize("WELCOME TO POKETACTIX!", ui.Bold+ui.ColorBrightYellow))
	fmt.Println()
	fmt.Println("PokeTacTix is an offline Pokemon battle game where you collect, train,")
	fmt.Println("and battle with Pokemon from Generations 1-5. Build your ultimate team")
	fmt.Println("and become a Pokemon master!")
	fmt.Println()

	// Battle Modes
	fmt.Println(ui.Colorize("BATTLE MODES", ui.Bold+ui.ColorBrightGreen))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()
	fmt.Println("  • 1v1 Mode: Quick battles with one random Pokemon from your deck")
	fmt.Println("    - Rewards: 50 coins on win, 20 XP to your Pokemon")
	fmt.Println("    - Perfect for quick play sessions")
	fmt.Println()
	fmt.Println("  • 5v5 Mode: Strategic battles using all 5 Pokemon in your deck")
	fmt.Println("    - Rewards: 150 coins on win, 15 XP to each Pokemon")
	fmt.Println("    - Win bonus: Select one AI Pokemon to add to your collection!")
	fmt.Println("    - Requires more strategy and Pokemon management")
	fmt.Println()

	// Battle Actions
	fmt.Println(ui.Colorize("BATTLE ACTIONS", ui.Bold+ui.ColorBrightGreen))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()
	fmt.Println("  • Attack: Use one of your Pokemon's 4 moves to deal damage")
	fmt.Println("    - Each move costs stamina and has different power/type")
	fmt.Println("    - Type advantages give bonus damage!")
	fmt.Println()
	fmt.Println("  • Defend: Reduce incoming damage by 50% for one turn")
	fmt.Println("    - Use when your Pokemon is low on HP")
	fmt.Println()
	fmt.Println("  • Pass: Skip your turn to regenerate 20% stamina")
	fmt.Println("    - Useful when you're out of stamina for powerful moves")
	fmt.Println()
	fmt.Println("  • Sacrifice: (5v5 only) Sacrifice current Pokemon to fully heal next one")
	fmt.Println("    - Strategic option when a Pokemon is about to faint")
	fmt.Println()
	fmt.Println("  • Surrender: Give up the battle")
	fmt.Println("    - You'll still earn some coins (10 for 1v1, 25 for 5v5)")
	fmt.Println()

	// Progression System
	fmt.Println(ui.Colorize("PROGRESSION & REWARDS", ui.Bold+ui.ColorBrightGreen))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()
	fmt.Println("  • Win battles to earn coins and XP for your Pokemon")
	fmt.Println("  • Pokemon level up at 100 XP, increasing their stats:")
	fmt.Println("    - HP increases by 3% per level")
	fmt.Println("    - Attack and Defense increase by 2% per level")
	fmt.Println("    - Speed increases by 1% per level")
	fmt.Println()
	fmt.Println("  • Use coins to purchase new Pokemon from the shop")
	fmt.Println("    - Common Pokemon: 100 coins")
	fmt.Println("    - Uncommon Pokemon: 250 coins")
	fmt.Println("    - Rare Pokemon: 500 coins")
	fmt.Println()
	fmt.Println("  • Shop refreshes every 10 battles with new Pokemon")
	fmt.Println()

	// Collection & Deck
	fmt.Println(ui.Colorize("COLLECTION & DECK MANAGEMENT", ui.Bold+ui.ColorBrightGreen))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()
	fmt.Println("  • Collection: All Pokemon you've acquired")
	fmt.Println("    - View with 'collection' command")
	fmt.Println("    - Filter by type, rarity, level")
	fmt.Println("    - Sort by various stats")
	fmt.Println()
	fmt.Println("  • Deck: Your active battle team (exactly 5 Pokemon)")
	fmt.Println("    - View with 'deck' command")
	fmt.Println("    - Edit with 'deck edit' command")
	fmt.Println("    - Must have 5 Pokemon to battle")
	fmt.Println()

	// Type Effectiveness
	fmt.Println(ui.Colorize("TYPE EFFECTIVENESS", ui.Bold+ui.ColorBrightGreen))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()
	fmt.Println("  • Fire beats Grass, Ice, Bug, Steel")
	fmt.Println("  • Water beats Fire, Ground, Rock")
	fmt.Println("  • Grass beats Water, Ground, Rock")
	fmt.Println("  • Electric beats Water, Flying")
	fmt.Println("  • Fighting beats Normal, Rock, Steel, Ice, Dark")
	fmt.Println("  • And many more type matchups to discover!")
	fmt.Println()
	fmt.Println("  Super effective attacks deal 1.5x damage!")
	fmt.Println()

	// Commands Quick Reference
	fmt.Println(ui.Colorize("ESSENTIAL COMMANDS", ui.Bold+ui.ColorBrightGreen))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()

	commands := []struct {
		Command string
		Alias   string
		Purpose string
	}{
		{"battle", "b", "Start a battle (1v1 or 5v5)"},
		{"collection", "c", "View your Pokemon collection"},
		{"deck", "d", "View your battle deck"},
		{"deck edit", "d edit", "Edit your battle deck"},
		{"shop", "s", "Buy Pokemon with coins"},
		{"stats", "st", "View battle statistics"},
		{"help", "h", "Show all commands"},
		{"quit", "q", "Exit the game"},
	}

	for _, cmd := range commands {
		cmdText := cmd.Command
		if ch.renderer.ColorSupport {
			cmdText = ui.Colorize(cmdText, ui.Bold+ui.ColorBrightCyan)
		}

		aliasText := ""
		if cmd.Alias != "" {
			aliasText = fmt.Sprintf(" (%s)", cmd.Alias)
			if ch.renderer.ColorSupport {
				aliasText = ui.Colorize(aliasText, ui.ColorGray)
			}
		}

		fmt.Printf("  %s%s - %s\n", cmdText, aliasText, cmd.Purpose)
	}

	fmt.Println()

	// Tips for Success
	fmt.Println(ui.Colorize("TIPS FOR SUCCESS", ui.Bold+ui.ColorBrightGreen))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()

	tips := []string{
		"Start with 1v1 battles to learn the mechanics and earn coins",
		"Build a balanced deck with different types for type coverage",
		"Save your strongest Pokemon for later rounds in 5v5 battles",
		"Use the Pass action strategically to regenerate stamina",
		"Check the shop regularly - it refreshes every 10 battles",
		"Level up your Pokemon by using them in battles",
		"Type advantages can turn the tide of battle - learn them!",
		"The game auto-saves after battles, but you can manually save anytime",
	}

	for _, tip := range tips {
		bullet := "•"
		if ch.renderer.ColorSupport {
			bullet = ui.Colorize("•", ui.ColorBrightYellow)
		}
		fmt.Printf("  %s %s\n", bullet, tip)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println()
	fmt.Println(ui.Colorize("Ready to begin your journey? Type 'battle' to start!", ui.Bold+ui.ColorBrightGreen))
	fmt.Println()
	fmt.Println("Press Enter to continue...")
	ch.scanner.Scan()

	return nil
}
