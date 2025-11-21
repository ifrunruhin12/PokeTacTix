package commands

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
)

// StatsCommand handles statistics-related commands
type StatsCommand struct {
	gameState *storage.GameState
	renderer  *ui.Renderer
	scanner   *bufio.Scanner
}

// NewStatsCommand creates a new stats command handler
func NewStatsCommand(gameState *storage.GameState, renderer *ui.Renderer, scanner *bufio.Scanner) *StatsCommand {
	return &StatsCommand{
		gameState: gameState,
		renderer:  renderer,
		scanner:   scanner,
	}
}

// ViewStats displays player statistics in a formatted layout
func (sc *StatsCommand) ViewStats() error {
	// Clear screen
	sc.renderer.Clear()

	// Display header
	fmt.Println(ui.RenderLogo())
	fmt.Println()
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println(ui.Colorize("PLAYER STATISTICS", ui.Bold+ui.ColorBrightCyan))
	fmt.Println(strings.Repeat("═", 80))
	fmt.Println()

	// Display player info
	fmt.Printf("Player: %s\n", ui.Colorize(sc.gameState.PlayerName, ui.Bold+ui.ColorBrightYellow))
	fmt.Printf("Coins: %s\n", ui.Colorize(fmt.Sprintf("%d", sc.gameState.Coins), ui.Bold+ui.ColorBrightGreen))
	fmt.Printf("Total Pokemon: %s\n", ui.Colorize(fmt.Sprintf("%d", len(sc.gameState.Collection)), ui.Bold+ui.ColorBrightCyan))
	fmt.Printf("Highest Level: %s\n", ui.Colorize(fmt.Sprintf("%d", sc.gameState.Stats.HighestLevel), ui.Bold+ui.ColorBrightMagenta))
	fmt.Println()

	// Display overall statistics
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println(ui.Colorize("OVERALL STATISTICS", ui.Bold))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()

	totalBattles := sc.gameState.Stats.TotalBattles1v1 + sc.gameState.Stats.TotalBattles5v5
	totalWins := sc.gameState.Stats.Wins1v1 + sc.gameState.Stats.Wins5v5
	totalLosses := sc.gameState.Stats.Losses1v1 + sc.gameState.Stats.Losses5v5
	totalDraws := sc.gameState.Stats.Draws1v1 + sc.gameState.Stats.Draws5v5

	var winRate float64
	if totalBattles > 0 {
		winRate = float64(totalWins) / float64(totalBattles) * 100
	}

	fmt.Printf("Total Battles:    %d\n", totalBattles)
	fmt.Printf("Wins:             %s\n", ui.Colorize(fmt.Sprintf("%d", totalWins), ui.ColorGreen))
	fmt.Printf("Losses:           %s\n", ui.Colorize(fmt.Sprintf("%d", totalLosses), ui.ColorRed))
	fmt.Printf("Draws:            %s\n", ui.Colorize(fmt.Sprintf("%d", totalDraws), ui.ColorYellow))
	fmt.Printf("Win Rate:         %s\n", sc.formatWinRate(winRate))
	fmt.Printf("Total Coins Earned: %s\n", ui.Colorize(fmt.Sprintf("%d", sc.gameState.Stats.TotalCoinsEarned), ui.ColorBrightGreen))
	fmt.Println()

	// Display 1v1 statistics
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println(ui.Colorize("1v1 BATTLE STATISTICS", ui.Bold))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()

	var winRate1v1 float64
	if sc.gameState.Stats.TotalBattles1v1 > 0 {
		winRate1v1 = float64(sc.gameState.Stats.Wins1v1) / float64(sc.gameState.Stats.TotalBattles1v1) * 100
	}

	fmt.Printf("Total Battles:    %d\n", sc.gameState.Stats.TotalBattles1v1)
	fmt.Printf("Wins:             %s\n", ui.Colorize(fmt.Sprintf("%d", sc.gameState.Stats.Wins1v1), ui.ColorGreen))
	fmt.Printf("Losses:           %s\n", ui.Colorize(fmt.Sprintf("%d", sc.gameState.Stats.Losses1v1), ui.ColorRed))
	fmt.Printf("Draws:            %s\n", ui.Colorize(fmt.Sprintf("%d", sc.gameState.Stats.Draws1v1), ui.ColorYellow))
	fmt.Printf("Win Rate:         %s\n", sc.formatWinRate(winRate1v1))
	fmt.Println()

	// Display 5v5 statistics
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println(ui.Colorize("5v5 BATTLE STATISTICS", ui.Bold))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()

	var winRate5v5 float64
	if sc.gameState.Stats.TotalBattles5v5 > 0 {
		winRate5v5 = float64(sc.gameState.Stats.Wins5v5) / float64(sc.gameState.Stats.TotalBattles5v5) * 100
	}

	fmt.Printf("Total Battles:    %d\n", sc.gameState.Stats.TotalBattles5v5)
	fmt.Printf("Wins:             %s\n", ui.Colorize(fmt.Sprintf("%d", sc.gameState.Stats.Wins5v5), ui.ColorGreen))
	fmt.Printf("Losses:           %s\n", ui.Colorize(fmt.Sprintf("%d", sc.gameState.Stats.Losses5v5), ui.ColorRed))
	fmt.Printf("Draws:            %s\n", ui.Colorize(fmt.Sprintf("%d", sc.gameState.Stats.Draws5v5), ui.ColorYellow))
	fmt.Printf("Win Rate:         %s\n", sc.formatWinRate(winRate5v5))
	fmt.Println()

	// Display battle history
	sc.displayBattleHistory()

	fmt.Println(strings.Repeat("═", 80))
	fmt.Println()
	fmt.Println("Press Enter to continue...")
	sc.scanner.Scan()

	return nil
}

// formatWinRate formats the win rate with color coding
func (sc *StatsCommand) formatWinRate(winRate float64) string {
	winRateStr := fmt.Sprintf("%.1f%%", winRate)
	
	if winRate >= 70 {
		return ui.Colorize(winRateStr, ui.Bold+ui.ColorBrightGreen)
	} else if winRate >= 50 {
		return ui.Colorize(winRateStr, ui.ColorGreen)
	} else if winRate >= 30 {
		return ui.Colorize(winRateStr, ui.ColorYellow)
	} else {
		return ui.Colorize(winRateStr, ui.ColorRed)
	}
}

// displayBattleHistory displays the last 10 battles
func (sc *StatsCommand) displayBattleHistory() {
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println(ui.Colorize("RECENT BATTLE HISTORY", ui.Bold))
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()

	if len(sc.gameState.BattleHistory) == 0 {
		fmt.Println(ui.Colorize("No battles yet. Start your first battle!", ui.ColorYellow))
		fmt.Println()
		return
	}

	// Get last 10 battles (or fewer if less than 10)
	startIdx := 0
	if len(sc.gameState.BattleHistory) > 10 {
		startIdx = len(sc.gameState.BattleHistory) - 10
	}
	recentBattles := sc.gameState.BattleHistory[startIdx:]

	// Display table header
	fmt.Printf("%-4s %-20s %-6s %-10s %-12s", "#", "DATE", "MODE", "RESULT", "COINS")
	if sc.hasDurationData() {
		fmt.Printf(" %-10s", "DURATION")
	}
	fmt.Println()
	fmt.Println(strings.Repeat("-", 80))

	// Display battles in reverse order (most recent first)
	for i := len(recentBattles) - 1; i >= 0; i-- {
		battle := recentBattles[i]
		num := len(recentBattles) - i

		// Format date
		dateStr := sc.formatDate(battle.Timestamp)

		// Format mode
		modeStr := strings.ToUpper(battle.Mode)

		// Format result with color
		resultStr := sc.formatResult(battle.Result)

		// Format coins
		coinsStr := fmt.Sprintf("+%d", battle.CoinsEarned)
		if sc.renderer.ColorSupport {
			coinsStr = ui.Colorize(coinsStr, ui.ColorBrightGreen)
		}

		// Format duration if available
		durationStr := ""
		if battle.Duration > 0 {
			durationStr = sc.formatDuration(battle.Duration)
		}

		// Print row
		fmt.Printf("%-4d %-20s %-6s %-10s %-12s", num, dateStr, modeStr, resultStr, coinsStr)
		if sc.hasDurationData() {
			fmt.Printf(" %-10s", durationStr)
		}
		fmt.Println()
	}

	fmt.Println()
}

// formatDate formats a timestamp into a readable date string
func (sc *StatsCommand) formatDate(t time.Time) string {
	now := time.Now()
	
	// If today, show time
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return fmt.Sprintf("Today %s", t.Format("15:04"))
	}
	
	// If yesterday
	yesterday := now.AddDate(0, 0, -1)
	if t.Year() == yesterday.Year() && t.YearDay() == yesterday.YearDay() {
		return fmt.Sprintf("Yesterday %s", t.Format("15:04"))
	}
	
	// If within last week, show day name
	weekAgo := now.AddDate(0, 0, -7)
	if t.After(weekAgo) {
		return t.Format("Mon 15:04")
	}
	
	// Otherwise show date
	return t.Format("Jan 02, 2006")
}

// formatResult formats the battle result with color coding
func (sc *StatsCommand) formatResult(result string) string {
	resultUpper := strings.ToUpper(result)
	
	if !sc.renderer.ColorSupport {
		return resultUpper
	}
	
	switch result {
	case "win":
		return ui.Colorize(resultUpper, ui.Bold+ui.ColorBrightGreen)
	case "loss":
		return ui.Colorize(resultUpper, ui.Bold+ui.ColorRed)
	case "draw":
		return ui.Colorize(resultUpper, ui.Bold+ui.ColorYellow)
	default:
		return resultUpper
	}
}

// formatDuration formats battle duration in seconds to a readable string
func (sc *StatsCommand) formatDuration(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	
	minutes := seconds / 60
	remainingSeconds := seconds % 60
	
	if remainingSeconds == 0 {
		return fmt.Sprintf("%dm", minutes)
	}
	
	return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
}

// hasDurationData checks if any battle has duration data
func (sc *StatsCommand) hasDurationData() bool {
	for _, battle := range sc.gameState.BattleHistory {
		if battle.Duration > 0 {
			return true
		}
	}
	return false
}
