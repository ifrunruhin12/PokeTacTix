//go:build cli
// +build cli

package pokemon

// FetchRandomPokemonCard returns a random Card from offline data (CLI version)
// This version is used when building with -tags cli
func FetchRandomPokemonCard(_ bool) Card {
	return FetchRandomPokemonCardOffline()
}
