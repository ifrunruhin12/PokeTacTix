package pokemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type PokeAPIClient interface {
	FetchPokemonRaw(ctx context.Context, idOrName string) ([]byte, error)
	FetchSpeciesRaw(ctx context.Context, idOrName string) ([]byte, error)
	FetchEvolutionChainRaw(ctx context.Context, chainID int) ([]byte, error)
}

type pokeAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPokeAPIClient(baseURL string, timeout time.Duration) PokeAPIClient {
	if baseURL == "" {
		baseURL = "https://pokeapi.co/api/v2"
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &pokeAPIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *pokeAPIClient) FetchPokemonRaw(ctx context.Context, idOrName string) ([]byte, error) {
	target := fmt.Sprintf("%s/pokemon/%s", c.baseURL, url.PathEscape(strings.ToLower(idOrName)))
	return c.get(ctx, target)
}

func (c *pokeAPIClient) FetchSpeciesRaw(ctx context.Context, idOrName string) ([]byte, error) {
	target := fmt.Sprintf("%s/pokemon-species/%s", c.baseURL, url.PathEscape(strings.ToLower(idOrName)))
	return c.get(ctx, target)
}

func (c *pokeAPIClient) FetchEvolutionChainRaw(ctx context.Context, chainID int) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/evolution-chain/%d", c.baseURL, chainID)
	return c.get(ctx, endpoint)
}

func (c *pokeAPIClient) get(ctx context.Context, endpoint string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed for %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("resource not found at %s", endpoint)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, endpoint)
	}

	var data json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}
	return data, nil
}

// Helpers for parsing PokéAPI payloads

func parseEvolutionChainIDFromURL(rawURL string) (int, error) {
	trimmed := strings.TrimRight(rawURL, "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid evolution chain url: %s", rawURL)
	}
	return strconv.Atoi(parts[len(parts)-1])
}

func parseGenerationNumber(genName string) int {
	// generation-i -> 1, generation-ii -> 2, etc.
	romanMap := map[string]int{
		"generation-i":    1,
		"generation-ii":   2,
		"generation-iii":  3,
		"generation-iv":   4,
		"generation-v":    5,
		"generation-vi":   6,
		"generation-vii":  7,
		"generation-viii": 8,
		"generation-ix":   9,
	}
	if val, ok := romanMap[strings.ToLower(genName)]; ok {
		return val
	}
	// fallback if it's numeric or unknown
	parts := strings.Split(genName, "-")
	if len(parts) > 1 {
		if num, err := strconv.Atoi(parts[1]); err == nil {
			return num
		}
	}
	return 1
}

// ExtractMemberSpeciesIDs extracts all species IDs recursively from evolution chain JSON
func ExtractMemberSpeciesIDs(chainJSON []byte) ([]int, error) {
	var payload struct {
		Chain struct {
			Species struct {
				URL string `json:"url"`
			} `json:"species"`
			EvolvesTo []json.RawMessage `json:"evolves_to"`
		} `json:"chain"`
	}

	if err := json.Unmarshal(chainJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal evolution chain payload: %w", err)
	}

	var memberIDs []int
	var walk func(rawNode json.RawMessage)
	walk = func(rawNode json.RawMessage) {
		var node struct {
			Species struct {
				URL string `json:"url"`
			} `json:"species"`
			EvolvesTo []json.RawMessage `json:"evolves_to"`
		}
		if err := json.Unmarshal(rawNode, &node); err != nil {
			return
		}
		if id, err := parseEvolutionChainIDFromURL(node.Species.URL); err == nil {
			memberIDs = append(memberIDs, id)
		}
		for _, child := range node.EvolvesTo {
			walk(child)
		}
	}

	// Process root
	if id, err := parseEvolutionChainIDFromURL(payload.Chain.Species.URL); err == nil {
		memberIDs = append(memberIDs, id)
	}
	for _, child := range payload.Chain.EvolvesTo {
		walk(child)
	}

	return memberIDs, nil
}

// evolutionChainNode mirrors one node of the PokéAPI evolution chain tree.
type evolutionChainNode struct {
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	EvolutionDetails []evolutionDetail    `json:"evolution_details"`
	EvolvesTo        []evolutionChainNode `json:"evolves_to"`
}

type namedRef struct {
	Name string `json:"name"`
}

// evolutionDetail holds the trigger conditions for a single evolution edge.
type evolutionDetail struct {
	MinLevel *int `json:"min_level"`
	Trigger  struct {
		Name string `json:"name"`
	} `json:"trigger"`
	Item *struct {
		Name string `json:"name"`
	} `json:"item"`
	MinHappiness *int   `json:"min_happiness"`
	TimeOfDay    string `json:"time_of_day"`
	// Form/region fields are objects in PokéAPI; the game models base forms
	// only, so non-base names here make the edge unsupported.
	Region              *namedRef                  `json:"region"`
	RequiredPokemonForm *namedRef                  `json:"required_pokemon_form"`
	EvolvedPokemonForm  *namedRef                  `json:"evolved_pokemon_form"`
	Conditions          map[string]json.RawMessage `json:"-"`
}

func (d *evolutionDetail) UnmarshalJSON(data []byte) error {
	type detail evolutionDetail
	if err := json.Unmarshal(data, (*detail)(d)); err != nil {
		return err
	}
	return json.Unmarshal(data, &d.Conditions)
}

// conditionMetadataFields are PokéAPI evolution-detail fields that carry
// pure entry metadata (version coverage and the is_default marker present on
// every entry) rather than a requirement this game would have to enforce.
// They appear on ordinary edges — e.g. Pikachu→Raichu lists version_group and
// is_default — and must not mark an edge as unsupported. Form/region fields
// are handled separately as typed fields: see regionalVariantEdge.
var conditionMetadataFields = map[string]bool{
	"is_default":    true,
	"version_group": true,
}

// regionalVariantEdge reports whether the detail describes a regional or
// form-variant evolution this game does not model (base forms only): an edge
// locked to a specific region, restricted to a non-base source form, or
// producing a non-base target form. Form names equal to the plain species
// name (e.g. required_pokemon_form "pikachu" on Pikachu→Raichu) are the base
// form and therefore fine.
func (d evolutionDetail) regionalVariantEdge(fromSpecies, toSpecies string) bool {
	if d.Region != nil && d.Region.Name != "" {
		return true
	}
	if d.RequiredPokemonForm != nil && d.RequiredPokemonForm.Name != "" && d.RequiredPokemonForm.Name != fromSpecies {
		return true
	}
	if d.EvolvedPokemonForm != nil && d.EvolvedPokemonForm.Name != "" && d.EvolvedPokemonForm.Name != toSpecies {
		return true
	}
	return false
}

// nonEmptyCondition reports whether an unknown evolution-detail field holds
// an actual (non-null, non-empty, non-false) value, ignoring metadata fields.
func nonEmptyCondition(name string, value json.RawMessage) bool {
	if conditionMetadataFields[name] {
		return false
	}
	return string(value) != "null" && string(value) != `""` && string(value) != "false"
}

func (d evolutionDetail) friendshipOnly(fromSpecies, toSpecies string) bool {
	if d.Trigger.Name != TriggerLevelUp || d.MinLevel != nil || d.MinHappiness == nil || *d.MinHappiness <= 0 || d.TimeOfDay != "" {
		return false
	}
	if d.regionalVariantEdge(fromSpecies, toSpecies) {
		return false
	}
	for name, value := range d.Conditions {
		switch name {
		case "trigger", "min_level", "min_happiness", "time_of_day",
			"region", "required_pokemon_form", "evolved_pokemon_form":
		default:
			if nonEmptyCondition(name, value) {
				return false
			}
		}
	}
	return true
}

func (d evolutionDetail) hasUnsupportedConditions(fromSpecies, toSpecies string) bool {
	if d.MinHappiness != nil && *d.MinHappiness > 0 && !d.friendshipOnly(fromSpecies, toSpecies) {
		return true
	}
	if d.Item != nil && d.Item.Name != "" && d.Trigger.Name != TriggerUseItem {
		return true
	}
	if d.MinLevel != nil && *d.MinLevel > 0 && d.Trigger.Name != TriggerLevelUp {
		return true
	}
	if d.regionalVariantEdge(fromSpecies, toSpecies) {
		return true
	}
	for name, value := range d.Conditions {
		switch name {
		case "trigger", "min_level", "min_happiness", "item", "time_of_day",
			"region", "required_pokemon_form", "evolved_pokemon_form":
		default:
			if nonEmptyCondition(name, value) {
				return true
			}
		}
	}
	return false
}

// ExtractEvolutionLinks walks the evolution chain tree and returns every edge
// with its trigger details. Each evolves_to entry may carry multiple distinct
// requirements; retain each branch without inventing an unconditional edge.
func ExtractEvolutionLinks(chainJSON []byte) ([]EvolutionLink, []int, error) {
	var payload struct {
		Chain evolutionChainNode `json:"chain"`
	}
	if err := json.Unmarshal(chainJSON, &payload); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal evolution chain payload: %w", err)
	}

	var links []EvolutionLink
	seen := make(map[EvolutionLink]bool)
	var memberIDs []int

	speciesID := func(n evolutionChainNode) (int, bool) {
		id, err := parseEvolutionChainIDFromURL(n.Species.URL)
		return id, err == nil
	}

	var walk func(n evolutionChainNode)
	walk = func(n evolutionChainNode) {
		if fromID, ok := speciesID(n); ok {
			memberIDs = append(memberIDs, fromID)
		}
		for _, child := range n.EvolvesTo {
			toID, ok := speciesID(child)
			if !ok {
				continue
			}

			fromID, hasFrom := speciesID(n)
			if hasFrom {
				details := child.EvolutionDetails
				if len(details) == 0 {
					details = []evolutionDetail{{}}
				}
				for _, detail := range details {
					link := EvolutionLink{
						FromSpeciesID:         fromID,
						ToSpeciesID:           toID,
						Trigger:               detail.Trigger.Name,
						TimeOfDay:             detail.TimeOfDay,
						UnsupportedConditions: detail.hasUnsupportedConditions(n.Species.Name, child.Species.Name),
					}
					if detail.MinLevel != nil {
						link.MinLevel = *detail.MinLevel
					}
					if detail.friendshipOnly(n.Species.Name, child.Species.Name) {
						link.MinLevel = FriendshipEvolutionLevel
					}
					if detail.Item != nil {
						link.Item = detail.Item.Name
					}
					if !seen[link] {
						links = append(links, link)
						seen[link] = true
					}
				}
			}
			walk(child)
		}
	}
	walk(payload.Chain)

	return links, memberIDs, nil
}
