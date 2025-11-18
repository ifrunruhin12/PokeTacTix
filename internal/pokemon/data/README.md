# Pokemon Offline Data

This directory contains the embedded Pokemon database used by the CLI version of PokeTacTix.

## Generating the Data

To generate the `pokemon_data.json` file, run:

```bash
go run scripts/generate_pokemon_data.go
```

This will:
- Fetch Pokemon data for Gen 1-5 (IDs 1-649) from PokeAPI
- Apply rate limiting (1 request per 100ms) to respect API limits
- Retry failed requests up to 3 times with exponential backoff
- Store Pokemon data including: ID, name, stats, types, moves, sprites, and rarity flags
- Save the data to `internal/pokemon/data/pokemon_data.json`

**Note:** The generation process takes approximately 1-2 hours due to rate limiting.

## Data Structure

The generated JSON file contains:

```json
{
  "pokemon": [
    {
      "id": 1,
      "name": "bulbasaur",
      "hp": 68,
      "attack": 49,
      "defense": 49,
      "speed": 45,
      "types": ["grass", "poison"],
      "moves": [
        {
          "name": "razor-leaf",
          "power": 55,
          "stamina_cost": 18,
          "attack_type": "grass"
        }
      ],
      "sprite": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/1.png",
      "is_legendary": false,
      "is_mythical": false
    }
  ],
  "generated": "2025-11-18T12:00:00Z",
  "version": "1.0.0"
}
```

## Usage in CLI

The CLI automatically uses offline data when `POKEMON_OFFLINE_MODE=true` is set:

```bash
export POKEMON_OFFLINE_MODE=true
./poketactix-cli
```

The offline data is embedded in the binary using Go's `//go:embed` directive, so no external files are needed at runtime.

## File Size

Expected file size: 2-5 MB uncompressed

The data is embedded directly in the binary, adding approximately 2-5 MB to the executable size.
