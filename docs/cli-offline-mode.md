# CLI Offline Mode Guide

## Quick Start

### 1. Generate Pokemon Data (One-time setup)

```bash
# This takes 1-2 hours due to API rate limiting
go run scripts/generate_pokemon_data.go
```

This will create `internal/pokemon/data/pokemon_data.json` with 649 Pokemon from Gen 1-5.

### 2. Build CLI with Embedded Data

```bash
# The data is automatically embedded during build
go build -o bin/poketactix-cli cmd/cli/main.go
```

### 3. Run in Offline Mode

```bash
# Set environment variable to enable offline mode
export POKEMON_OFFLINE_MODE=true
./bin/poketactix-cli
```

## How It Works

### Automatic Mode Detection

The system automatically detects offline mode via the `POKEMON_OFFLINE_MODE` environment variable:

```go
// In your code
card := pokemon.FetchRandomPokemonCard(false)
// Automatically uses offline data if POKEMON_OFFLINE_MODE=true
```

### Manual Offline Mode

You can also explicitly use offline functions:

```go
// Always use offline data
card := pokemon.FetchRandomPokemonCardOffline()

// Get specific Pokemon by ID
pikachu, err := pokemon.GetPokemonByID(25)

// Get random Pokemon with filters
normalPokemon, err := pokemon.GetRandomPokemon(true, true) // exclude legendary & mythical
```

## Testing Offline Mode

### Test the Data System

```bash
go run scripts/test_offline_data.go
```

### Test with Web API

You can test offline mode with the web API:

```bash
POKEMON_OFFLINE_MODE=true go run cmd/api/main.go
```

Then start a battle - it will use offline data instead of calling PokeAPI.

## Data Generation Details

### What Gets Fetched

- **Pokemon IDs**: 1-649 (Gen 1-5)
- **Data per Pokemon**:
  - ID, name, base stats (HP, Attack, Defense, Speed)
  - Types (e.g., fire, water, grass)
  - 4 random moves with power and stamina cost
  - Sprite URL
  - Legendary/Mythical flags

### Rate Limiting

- **1 request per 100ms** to respect PokeAPI limits
- **3 retry attempts** with exponential backoff
- **Fallback data** for failed requests

### Expected Output

```
Fetching Pokemon 1/649...
Fetching Pokemon 2/649...
...
Progress: 50/649 (Success: 48, Failed: 2)
...
Progress: 649/649 (Success: 645, Failed: 4)
Successfully generated internal/pokemon/data/pokemon_data.json (3.2 MB)
Total Pokemon: 649
Generation complete!
```

## File Structure

```
internal/pokemon/
├── data/
│   ├── pokemon_data.json      # Generated Pokemon database
│   └── README.md              # Data documentation
├── offline_data.go            # Offline data loading
├── fetcher.go                 # Online/offline routing
├── builder.go                 # Card building
└── types.go                   # Type definitions

scripts/
├── generate_pokemon_data.go   # Data generation script
└── test_offline_data.go       # Testing script
```

## Performance

- **Startup**: < 1 second (lazy loading)
- **Lookup**: < 10ms (in-memory hash map)
- **Memory**: ~5-10 MB (parsed data)
- **Binary Size**: +2-5 MB (embedded JSON)

## Troubleshooting

### "Failed to parse embedded Pokemon data"

The `pokemon_data.json` file is missing or corrupted. Regenerate it:

```bash
go run scripts/generate_pokemon_data.go
```

### "Pokemon with ID X not found"

The Pokemon ID is not in the database (only Gen 1-5 are included). Use IDs 1-649.

### Offline mode not working

Make sure the environment variable is set:

```bash
# Check if set
echo $POKEMON_OFFLINE_MODE

# Set it
export POKEMON_OFFLINE_MODE=true
```

### Generation script fails

- Check internet connection
- PokeAPI might be down - try again later
- Some Pokemon might fail - the script continues with fallback data

## Development Tips

### Adding More Generations

To include Gen 6-9, modify the generation script:

```go
// In scripts/generate_pokemon_data.go
for id := 1; id <= 1025; id++ { // Up to Gen 9
    // ...
}
```

### Customizing Rarity

Modify the rarity logic in `offline_data.go`:

```go
mythicalOdds := 0.001  // 0.1% instead of 0.01%
legendaryOdds := 0.001 // 0.1% instead of 0.01%
```

### Using Different Data Sources

You can modify the generation script to fetch from other sources or use local data files.

## CLI vs Web Mode

| Feature | Web Mode | CLI Mode |
|---------|----------|----------|
| Data Source | PokeAPI (online) | Embedded JSON |
| Network Required | Yes | No |
| Pokemon Available | All generations | Gen 1-5 (649) |
| Startup Time | Varies (network) | < 1 second |
| Binary Size | Smaller | +2-5 MB |
| Mode Switch | `POKEMON_OFFLINE_MODE=false` | `POKEMON_OFFLINE_MODE=true` |

