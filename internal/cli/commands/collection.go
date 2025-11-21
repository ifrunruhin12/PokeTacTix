package commands

import (
	"bufio"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"pokemon-cli/internal/cli/storage"
	"pokemon-cli/internal/cli/ui"
)

// CollectionCommand handles collection-related commands
type CollectionCommand struct {
	gameState *storage.GameState
	renderer  *ui.Renderer
	scanner   *bufio.Scanner
}

// NewCollectionCommand creates a new collection command handler
func NewCollectionCommand(gameState *storage.GameState, renderer *ui.Renderer, scanner *bufio.Scanner) *CollectionCommand {
	return &CollectionCommand{
		gameState: gameState,
		renderer:  renderer,
		scanner:   scanner,
	}
}

// CollectionFilters represents filters for collection viewing
type CollectionFilters struct {
	TypeFilter   string // Filter by type (e.g., "fire", "water")
	RarityFilter string // Filter by rarity (common, uncommon, rare, legendary, mythical)
	MinLevel     int    // Minimum level
	MaxLevel     int    // Maximum level
	SearchName   string // Search by Pokemon name
	SortBy       string // Sort field (level, name, hp, attack, defense, speed)
	SortDesc     bool   // Sort descending
}

// ViewCollection displays all owned Pokemon with pagination
func (cc *CollectionCommand) ViewCollection() error {
	return cc.ViewCollectionWithFilters(CollectionFilters{})
}

// ViewCollectionWithFilters displays Pokemon collection with filters and sorting
func (cc *CollectionCommand) ViewCollectionWithFilters(filters CollectionFilters) error {
	// Check if collection is empty
	if len(cc.gameState.Collection) == 0 {
		fmt.Println(ui.Colorize("Your collection is empty!", ui.ColorYellow))
		fmt.Println("Win 5v5 battles to add Pokemon to your collection.")
		return nil
	}

	// Apply filters
	filteredCollection := cc.applyFilters(filters)

	if len(filteredCollection) == 0 {
		fmt.Println(ui.Colorize("No Pokemon match your filters.", ui.ColorYellow))
		fmt.Println("Press Enter to continue...")
		cc.scanner.Scan()
		return nil
	}

	// Apply sorting
	cc.applySorting(filteredCollection, filters)

	// Display with pagination
	return cc.displayPaginated(filteredCollection, filters)
}

// applyFilters applies filters to the collection
func (cc *CollectionCommand) applyFilters(filters CollectionFilters) []storage.PlayerCard {
	var filtered []storage.PlayerCard

	for _, card := range cc.gameState.Collection {
		// Type filter
		if filters.TypeFilter != "" {
			hasType := false
			for _, t := range card.Types {
				if strings.EqualFold(t, filters.TypeFilter) {
					hasType = true
					break
				}
			}
			if !hasType {
				continue
			}
		}

		// Rarity filter
		if filters.RarityFilter != "" {
			rarity := cc.getPokemonRarity(card)
			if !strings.EqualFold(rarity, filters.RarityFilter) {
				continue
			}
		}

		// Level range filter
		if filters.MinLevel > 0 && card.Level < filters.MinLevel {
			continue
		}
		if filters.MaxLevel > 0 && card.Level > filters.MaxLevel {
			continue
		}

		// Name search
		if filters.SearchName != "" {
			if !strings.Contains(strings.ToLower(card.Name), strings.ToLower(filters.SearchName)) {
				continue
			}
		}

		filtered = append(filtered, card)
	}

	return filtered
}

// applySorting sorts the collection based on filters
func (cc *CollectionCommand) applySorting(collection []storage.PlayerCard, filters CollectionFilters) {
	sort.Slice(collection, func(i, j int) bool {
		var less bool

		switch filters.SortBy {
		case "level":
			less = collection[i].Level < collection[j].Level
		case "name":
			less = strings.ToLower(collection[i].Name) < strings.ToLower(collection[j].Name)
		case "hp":
			statsI := collection[i].GetCurrentStats()
			statsJ := collection[j].GetCurrentStats()
			less = statsI.HP < statsJ.HP
		case "attack":
			statsI := collection[i].GetCurrentStats()
			statsJ := collection[j].GetCurrentStats()
			less = statsI.Attack < statsJ.Attack
		case "defense":
			statsI := collection[i].GetCurrentStats()
			statsJ := collection[j].GetCurrentStats()
			less = statsI.Defense < statsJ.Defense
		case "speed":
			statsI := collection[i].GetCurrentStats()
			statsJ := collection[j].GetCurrentStats()
			less = statsI.Speed < statsJ.Speed
		default:
			// Default sort by ID (acquisition order)
			less = collection[i].ID < collection[j].ID
		}

		if filters.SortDesc {
			return !less
		}
		return less
	})
}

// displayPaginated displays collection with pagination (10 per page)
func (cc *CollectionCommand) displayPaginated(collection []storage.PlayerCard, filters CollectionFilters) error {
	const itemsPerPage = 10
	totalPages := (len(collection) + itemsPerPage - 1) / itemsPerPage
	currentPage := 0

	for {
		// Clear screen
		cc.renderer.Clear()

		// Display header
		fmt.Println(ui.RenderLogo())
		fmt.Println()
		fmt.Println(strings.Repeat("═", 80))
		fmt.Println(ui.Colorize("POKEMON COLLECTION", ui.Bold+ui.ColorBrightCyan))
		fmt.Println(strings.Repeat("═", 80))
		fmt.Println()

		// Display active filters
		if cc.hasActiveFilters(filters) {
			fmt.Print(ui.Colorize("Active Filters: ", ui.Bold))
			filterParts := []string{}
			if filters.TypeFilter != "" {
				filterParts = append(filterParts, fmt.Sprintf("Type=%s", filters.TypeFilter))
			}
			if filters.RarityFilter != "" {
				filterParts = append(filterParts, fmt.Sprintf("Rarity=%s", filters.RarityFilter))
			}
			if filters.MinLevel > 0 || filters.MaxLevel > 0 {
				if filters.MaxLevel > 0 {
					filterParts = append(filterParts, fmt.Sprintf("Level=%d-%d", filters.MinLevel, filters.MaxLevel))
				} else {
					filterParts = append(filterParts, fmt.Sprintf("Level>=%d", filters.MinLevel))
				}
			}
			if filters.SearchName != "" {
				filterParts = append(filterParts, fmt.Sprintf("Name=%s", filters.SearchName))
			}
			if filters.SortBy != "" {
				sortDir := "asc"
				if filters.SortDesc {
					sortDir = "desc"
				}
				filterParts = append(filterParts, fmt.Sprintf("Sort=%s(%s)", filters.SortBy, sortDir))
			}
			fmt.Println(strings.Join(filterParts, ", "))
			fmt.Println()
		}

		// Display total count
		fmt.Printf("Total Pokemon: %d", len(collection))
		if len(collection) != len(cc.gameState.Collection) {
			fmt.Printf(" (filtered from %d)", len(cc.gameState.Collection))
		}
		fmt.Printf(" | Page %d/%d\n", currentPage+1, totalPages)
		fmt.Println()

		// Calculate page bounds
		startIdx := currentPage * itemsPerPage
		endIdx := startIdx + itemsPerPage
		if endIdx > len(collection) {
			endIdx = len(collection)
		}

		// Display Pokemon table
		cc.displayPokemonTable(collection[startIdx:endIdx], startIdx)

		fmt.Println()
		fmt.Println(strings.Repeat("═", 80))
		fmt.Println()

		// Display navigation options
		fmt.Println("Options:")
		if currentPage > 0 {
			fmt.Println("  [P] Previous page")
		}
		if currentPage < totalPages-1 {
			fmt.Println("  [N] Next page")
		}
		fmt.Println("  [F] Filter/Sort")
		fmt.Println("  [C] Clear filters")
		fmt.Println("  [Q] Back to menu")
		fmt.Println()
		fmt.Print("Enter your choice: ")

		// Get user input
		if !cc.scanner.Scan() {
			return fmt.Errorf("failed to read input")
		}

		input := strings.ToUpper(strings.TrimSpace(cc.scanner.Text()))

		switch input {
		case "P":
			if currentPage > 0 {
				currentPage--
			}
		case "N":
			if currentPage < totalPages-1 {
				currentPage++
			}
		case "F":
			// Show filter/sort menu
			newFilters, err := cc.showFilterMenu(filters)
			if err != nil {
				return err
			}
			return cc.ViewCollectionWithFilters(newFilters)
		case "C":
			// Clear all filters
			return cc.ViewCollection()
		case "Q":
			return nil
		default:
			// Invalid input, just redisplay
			continue
		}
	}
}

// displayPokemonTable displays Pokemon in a formatted table
func (cc *CollectionCommand) displayPokemonTable(pokemon []storage.PlayerCard, startIdx int) {
	// Table header
	fmt.Printf("%-4s %-15s %-6s %-8s %-6s %-6s %-6s %-6s %-20s\n",
		"#", "NAME", "LEVEL", "XP", "HP", "ATK", "DEF", "SPD", "TYPES")
	fmt.Println(strings.Repeat("-", 80))

	// Table rows
	for i, card := range pokemon {
		stats := card.GetCurrentStats()
		
		// Format number
		num := fmt.Sprintf("%d.", startIdx+i+1)
		
		// Format name with color if legendary/mythical
		name := card.Name
		if cc.renderer.ColorSupport {
			if card.IsMythical {
				name = ui.Colorize(name, ui.ColorMagenta)
			} else if card.IsLegendary {
				name = ui.Colorize(name, ui.ColorYellow)
			}
		}
		
		// Format XP progress
		xpProgress := fmt.Sprintf("%d/100", card.XP)
		
		// Format types
		typeStr := ""
		for j, t := range card.Types {
			if j > 0 {
				typeStr += "/"
			}
			if cc.renderer.ColorSupport {
				typeStr += ui.ColorizeType(strings.ToUpper(t), t)
			} else {
				typeStr += strings.ToUpper(t)
			}
		}

		fmt.Printf("%-4s %-15s %-6d %-8s %-6d %-6d %-6d %-6d %s\n",
			num, name, card.Level, xpProgress,
			stats.HP, stats.Attack, stats.Defense, stats.Speed,
			typeStr)
	}
}

// showFilterMenu displays the filter/sort menu and returns new filters
func (cc *CollectionCommand) showFilterMenu(currentFilters CollectionFilters) (CollectionFilters, error) {
	cc.renderer.Clear()
	
	fmt.Println(ui.RenderLogo())
	fmt.Println()
	fmt.Println(strings.Repeat("═", 60))
	fmt.Println(ui.Colorize("FILTER & SORT OPTIONS", ui.Bold+ui.ColorBrightCyan))
	fmt.Println(strings.Repeat("═", 60))
	fmt.Println()

	options := []ui.MenuOption{
		{Label: "Filter by Type", Description: "Show only Pokemon of a specific type", Value: "type"},
		{Label: "Filter by Rarity", Description: "Show only Pokemon of a specific rarity", Value: "rarity"},
		{Label: "Filter by Level Range", Description: "Show Pokemon within a level range", Value: "level"},
		{Label: "Search by Name", Description: "Find Pokemon by name", Value: "name"},
		{Label: "Sort Collection", Description: "Change sort order", Value: "sort"},
		{Label: "Back", Description: "Return to collection", Value: "back"},
	}

	fmt.Println(cc.renderer.RenderBorderedMenu(options, 0, "SELECT FILTER OPTION"))
	fmt.Print("Enter your choice (1-6): ")

	if !cc.scanner.Scan() {
		return currentFilters, fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(cc.scanner.Text())
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > 6 {
		return currentFilters, nil
	}

	switch choice {
	case 1:
		return cc.filterByType(currentFilters)
	case 2:
		return cc.filterByRarity(currentFilters)
	case 3:
		return cc.filterByLevel(currentFilters)
	case 4:
		return cc.searchByName(currentFilters)
	case 5:
		return cc.sortCollection(currentFilters)
	case 6:
		return currentFilters, nil
	}

	return currentFilters, nil
}

// filterByType prompts for type filter
func (cc *CollectionCommand) filterByType(filters CollectionFilters) (CollectionFilters, error) {
	fmt.Println()
	fmt.Println("Available types: fire, water, grass, electric, ice, fighting, poison,")
	fmt.Println("                 ground, flying, psychic, bug, rock, ghost, dragon,")
	fmt.Println("                 dark, steel, fairy, normal")
	fmt.Println()
	fmt.Print("Enter type (or press Enter to clear filter): ")

	if !cc.scanner.Scan() {
		return filters, fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(strings.ToLower(cc.scanner.Text()))
	filters.TypeFilter = input
	return filters, nil
}

// filterByRarity prompts for rarity filter
func (cc *CollectionCommand) filterByRarity(filters CollectionFilters) (CollectionFilters, error) {
	fmt.Println()
	fmt.Println("Available rarities:")
	fmt.Println("  1. Common")
	fmt.Println("  2. Uncommon")
	fmt.Println("  3. Rare")
	fmt.Println("  4. Legendary")
	fmt.Println("  5. Mythical")
	fmt.Println("  6. Clear filter")
	fmt.Println()
	fmt.Print("Enter choice (1-6): ")

	if !cc.scanner.Scan() {
		return filters, fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(cc.scanner.Text())
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > 6 {
		return filters, nil
	}

	rarities := []string{"common", "uncommon", "rare", "legendary", "mythical", ""}
	filters.RarityFilter = rarities[choice-1]
	return filters, nil
}

// filterByLevel prompts for level range filter
func (cc *CollectionCommand) filterByLevel(filters CollectionFilters) (CollectionFilters, error) {
	fmt.Println()
	fmt.Print("Enter minimum level (or 0 for no minimum): ")

	if !cc.scanner.Scan() {
		return filters, fmt.Errorf("failed to read input")
	}

	minInput := strings.TrimSpace(cc.scanner.Text())
	minLevel, err := strconv.Atoi(minInput)
	if err != nil {
		minLevel = 0
	}

	fmt.Print("Enter maximum level (or 0 for no maximum): ")

	if !cc.scanner.Scan() {
		return filters, fmt.Errorf("failed to read input")
	}

	maxInput := strings.TrimSpace(cc.scanner.Text())
	maxLevel, err := strconv.Atoi(maxInput)
	if err != nil {
		maxLevel = 0
	}

	filters.MinLevel = minLevel
	filters.MaxLevel = maxLevel
	return filters, nil
}

// searchByName prompts for name search
func (cc *CollectionCommand) searchByName(filters CollectionFilters) (CollectionFilters, error) {
	fmt.Println()
	fmt.Print("Enter Pokemon name to search (or press Enter to clear): ")

	if !cc.scanner.Scan() {
		return filters, fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(cc.scanner.Text())
	filters.SearchName = input
	return filters, nil
}

// sortCollection prompts for sort options
func (cc *CollectionCommand) sortCollection(filters CollectionFilters) (CollectionFilters, error) {
	fmt.Println()
	fmt.Println("Sort by:")
	fmt.Println("  1. Level")
	fmt.Println("  2. Name")
	fmt.Println("  3. HP")
	fmt.Println("  4. Attack")
	fmt.Println("  5. Defense")
	fmt.Println("  6. Speed")
	fmt.Println("  7. Acquisition order (default)")
	fmt.Println()
	fmt.Print("Enter choice (1-7): ")

	if !cc.scanner.Scan() {
		return filters, fmt.Errorf("failed to read input")
	}

	input := strings.TrimSpace(cc.scanner.Text())
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > 7 {
		return filters, nil
	}

	sortFields := []string{"level", "name", "hp", "attack", "defense", "speed", ""}
	filters.SortBy = sortFields[choice-1]

	if filters.SortBy != "" {
		fmt.Println()
		fmt.Print("Sort order (1=Ascending, 2=Descending): ")

		if !cc.scanner.Scan() {
			return filters, fmt.Errorf("failed to read input")
		}

		orderInput := strings.TrimSpace(cc.scanner.Text())
		orderChoice, err := strconv.Atoi(orderInput)
		if err == nil && orderChoice == 2 {
			filters.SortDesc = true
		} else {
			filters.SortDesc = false
		}
	}

	return filters, nil
}

// hasActiveFilters checks if any filters are active
func (cc *CollectionCommand) hasActiveFilters(filters CollectionFilters) bool {
	return filters.TypeFilter != "" ||
		filters.RarityFilter != "" ||
		filters.MinLevel > 0 ||
		filters.MaxLevel > 0 ||
		filters.SearchName != "" ||
		filters.SortBy != ""
}

// getPokemonRarity determines the rarity of a Pokemon
func (cc *CollectionCommand) getPokemonRarity(card storage.PlayerCard) string {
	if card.IsMythical {
		return "mythical"
	}
	if card.IsLegendary {
		return "legendary"
	}

	// Calculate rarity based on base stat total
	statTotal := card.BaseHP + card.BaseAttack + card.BaseDefense + card.BaseSpeed

	if statTotal >= 500 {
		return "rare"
	} else if statTotal >= 400 {
		return "uncommon"
	}
	return "common"
}
