# Architecture Decision: Web vs CLI Separation

## Question

Should the CLI and Web applications share the same codebase and Pokemon fetching logic?

## Context

**Web App:**
- Online multiplayer game
- Requires user accounts and authentication
- Uses PostgreSQL for persistence
- Fetches Pokemon data from PokeAPI (online)
- Deployed as a web service

**CLI App:**
- Single-player offline game
- No account required - instant play
- Standalone downloadable binary
- Uses embedded Pokemon data (offline)
- Runs locally on user's machine

## Decision

**Use build tags to conditionally compile different implementations while sharing core logic.**

## Rationale

### Why NOT Separate Repositories?

❌ **Full Separation** (two separate repos)
- Duplicates all game logic (battle system, card mechanics, etc.)
- Maintenance nightmare - bug fixes need to be applied twice
- Diverging implementations over time
- No code reuse

### Why NOT Single Binary with Runtime Checks?

❌ **Runtime Environment Variables** (original approach)
- Web binary includes unnecessary 5MB embedded data
- CLI binary includes unnecessary network code
- Runtime overhead for mode checking
- Confusing deployment (same binary, different modes)
- Larger binaries for both use cases

### Why Build Tags? ✅

✅ **Build Tags** (current approach)
- **Shared Core Logic**: Battle system, card mechanics, types all shared
- **Optimized Binaries**: Each binary only includes what it needs
- **Clean Separation**: Web code stays in web files, CLI code in CLI files
- **No Runtime Overhead**: Compile-time selection, zero runtime cost
- **Easy Maintenance**: Bug fixes in shared code benefit both
- **Clear Architecture**: Build tags make intent explicit

## Implementation

### File Structure

```
internal/pokemon/
├── types.go           # Shared (always compiled)
├── builder.go         # Shared (always compiled)
├── rarity.go          # Shared (always compiled)
├── fetcher_web.go     # Web-only (tag: !cli)
├── fetcher_cli.go     # CLI-only (tag: cli)
├── offline_data.go    # CLI-only (tag: cli)
└── data/
    └── pokemon_data.json  # CLI-only (embedded)
```

### Build Commands

```bash
# Web API (default - no tags needed)
go build -o bin/api cmd/api/main.go
# Result: ~15-20 MB, uses PokeAPI

# CLI Binary (with cli tag)
go build -tags cli -o bin/poketactix-cli cmd/cli/main.go
# Result: ~20-25 MB, includes embedded data
```

### Code Interface

Both implementations expose the same interface:

```go
// Works in both web and CLI builds
card := pokemon.FetchRandomPokemonCard(false)
```

- **Web build**: Calls PokeAPI
- **CLI build**: Uses embedded data

## Benefits

### 1. Optimal Binary Sizes
- Web: No embedded data bloat
- CLI: No unnecessary network code

### 2. Shared Game Logic
- Battle system
- Card mechanics
- Type effectiveness
- XP/leveling
- All shared - single source of truth

### 3. Independent Deployment
- Web can be updated without affecting CLI
- CLI can be distributed as standalone binary
- Different release cycles possible

### 4. Clear Separation of Concerns
- Web-specific: Auth, database, API routes
- CLI-specific: Offline data, local saves, terminal UI
- Shared: Game mechanics, battle logic

### 5. Future Flexibility
- Easy to add more build variants (mobile, embedded, etc.)
- Can optimize each target independently
- No runtime configuration complexity

## Trade-offs

### Pros
- ✅ Smaller binaries
- ✅ Faster compilation (only compiles needed code)
- ✅ Zero runtime overhead
- ✅ Clear architecture
- ✅ Shared core logic

### Cons
- ⚠️ Need to remember build tags when building CLI
- ⚠️ Slightly more complex build process
- ⚠️ Need to test both build variants

## Alternatives Considered

### 1. Separate Repositories
```
pokemon-web/     # Web app
pokemon-cli/     # CLI app
pokemon-core/    # Shared library
```

**Rejected because:**
- Overhead of managing 3 repos
- Versioning complexity
- Still need to publish/consume shared library
- Overkill for current project size

### 2. Monorepo with Separate Modules
```
/web
  /cmd/api
  /internal/...
/cli
  /cmd/cli
  /internal/...
/shared
  /battle
  /pokemon
```

**Rejected because:**
- More complex than build tags
- Still need to manage module dependencies
- Harder to share code
- More directory structure overhead

### 3. Single Binary with Modes
```go
if os.Getenv("MODE") == "cli" {
    // CLI logic
} else {
    // Web logic
}
```

**Rejected because:**
- Both binaries include all code
- Runtime overhead
- Confusing deployment
- Larger binaries

## Future Considerations

### Adding Mobile Support

Easy to add with build tags:

```go
// fetcher_mobile.go
//go:build mobile
// +build mobile
```

### Adding Embedded/IoT Support

```go
// fetcher_embedded.go
//go:build embedded
// +build embedded
```

### Multiplayer CLI

If CLI ever needs multiplayer:
- Can add network code to CLI build
- Use build tags to conditionally include
- Keep offline mode as default

## Conclusion

**Build tags provide the best balance of:**
- Code reuse (shared game logic)
- Binary optimization (only include what's needed)
- Clear architecture (explicit separation)
- Maintainability (single codebase)
- Flexibility (easy to extend)

This approach allows the web and CLI apps to coexist in the same repository while maintaining optimal binaries and clear separation of concerns.

## References

- [Go Build Tags Documentation](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- `internal/pokemon/BUILD_TAGS.md` - Implementation details
- `docs/cli-offline-mode.md` - CLI usage guide
