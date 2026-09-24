package pokemon

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Minimal but structurally faithful excerpt of PokéAPI's evolution-chain
// payload for the Charmander family (chain id 2).
const charmanderChainJSON = `{
  "id": 2,
  "chain": {
    "species": { "name": "charmander", "url": "https://pokeapi.co/api/v2/pokemon-species/4/" },
    "evolution_details": [],
    "evolves_to": [
      {
        "species": { "name": "charmeleon", "url": "https://pokeapi.co/api/v2/pokemon-species/5/" },
        "evolution_details": [
          { "min_level": 16, "trigger": { "name": "level-up", "url": "https://pokeapi.co/api/v2/evolution-trigger/1/" } }
        ],
        "evolves_to": [
          {
            "species": { "name": "charizard", "url": "https://pokeapi.co/api/v2/pokemon-species/6/" },
            "evolution_details": [
              { "min_level": 36, "trigger": { "name": "level-up", "url": "https://pokeapi.co/api/v2/evolution-trigger/1/" } }
            ],
            "evolves_to": []
          }
        ]
      }
    ]
  }
}`

// Eevee-like branching chain with non-level triggers mixed in.
const branchingChainJSON = `{
  "id": 67,
  "chain": {
    "species": { "name": "eevee", "url": "https://pokeapi.co/api/v2/pokemon-species/133/" },
    "evolution_details": [],
    "evolves_to": [
      {
        "species": { "name": "jolteon", "url": "https://pokeapi.co/api/v2/pokemon-species/135/" },
        "evolution_details": [
          { "min_level": null, "trigger": { "name": "use-item", "url": "https://pokeapi.co/api/v2/evolution-trigger/3/" } }
        ],
        "evolves_to": []
      },
      {
        "species": { "name": "espeon", "url": "https://pokeapi.co/api/v2/pokemon-species/196/" },
        "evolution_details": [
          { "min_level": null, "trigger": { "name": "level-up", "url": "https://pokeapi.co/api/v2/evolution-trigger/1/" }, "time_of_day": "day" }
        ],
        "evolves_to": []
      }
    ]
  }
}`

func TestExtractEvolutionLinksCharmanderFamily(t *testing.T) {
	links, members, err := ExtractEvolutionLinks([]byte(charmanderChainJSON))
	require.NoError(t, err)

	assert.Equal(t, []int{4, 5, 6}, members)
	assert.Equal(t, []EvolutionLink{
		{FromSpeciesID: 4, ToSpeciesID: 5, MinLevel: 16, Trigger: "level-up"},
		{FromSpeciesID: 5, ToSpeciesID: 6, MinLevel: 36, Trigger: "level-up"},
	}, links)
}

func TestExtractEvolutionLinksBranchingChain(t *testing.T) {
	links, members, err := ExtractEvolutionLinks([]byte(branchingChainJSON))
	require.NoError(t, err)

	assert.Equal(t, []int{133, 135, 196}, members)
	require.Len(t, links, 2)
	assert.Equal(t, EvolutionLink{FromSpeciesID: 133, ToSpeciesID: 135, Trigger: "use-item"}, links[0])
	// level-up edge without a min_level is kept but unusable for our purposes
	assert.Equal(t, EvolutionLink{FromSpeciesID: 133, ToSpeciesID: 196, Trigger: "level-up", TimeOfDay: "day"}, links[1])
}

// Pichu family (chain id 10 extended): pichu evolves into pikachu via
// friendship (min_happiness, no min_level), then pikachu into raichu via a
// thunder stone.
const pichuChainJSON = `{
  "id": 10,
  "chain": {
    "species": { "name": "pichu", "url": "https://pokeapi.co/api/v2/pokemon-species/172/" },
    "evolution_details": [],
    "evolves_to": [
      {
        "species": { "name": "pikachu", "url": "https://pokeapi.co/api/v2/pokemon-species/25/" },
        "evolution_details": [
          { "min_level": null, "min_happiness": 220, "trigger": { "name": "level-up", "url": "https://pokeapi.co/api/v2/evolution-trigger/1/" } }
        ],
        "evolves_to": [
          {
            "species": { "name": "raichu", "url": "https://pokeapi.co/api/v2/pokemon-species/26/" },
            "evolution_details": [
              { "min_level": null, "trigger": { "name": "use-item", "url": "https://pokeapi.co/api/v2/evolution-trigger/3/" }, "item": { "name": "thunder-stone", "url": "https://pokeapi.co/api/v2/item/83/" } }
            ],
            "evolves_to": []
          }
        ]
      }
    ]
  }
}`

func TestExtractEvolutionLinksFriendshipSynthesizesLevel(t *testing.T) {
	links, _, err := ExtractEvolutionLinks([]byte(pichuChainJSON))
	require.NoError(t, err)

	require.Len(t, links, 2)
	// Friendship edges become ordinary level-up evolutions at the fixed level.
	assert.Equal(t, EvolutionLink{
		FromSpeciesID: 172,
		ToSpeciesID:   25,
		MinLevel:      FriendshipEvolutionLevel,
		Trigger:       "level-up",
	}, links[0])
	assert.Equal(t, EvolutionLink{
		FromSpeciesID: 25,
		ToSpeciesID:   26,
		Trigger:       "use-item",
		Item:          "thunder-stone",
	}, links[1])
}

func TestExtractEvolutionLinksRetainsDistinctDetails(t *testing.T) {
	chain := `{"chain":{"species":{"url":"https://pokeapi.co/api/v2/pokemon-species/1/"},"evolves_to":[{"species":{"url":"https://pokeapi.co/api/v2/pokemon-species/2/"},"evolution_details":[{"trigger":{"name":"use-item"},"item":{"name":"fire-stone"}},{"trigger":{"name":"use-item"},"item":{"name":"water-stone"}},{"trigger":{"name":"level-up"},"min_happiness":220},{"trigger":{"name":"level-up"},"min_happiness":220,"time_of_day":"day"},{"trigger":{"name":"level-up"},"min_happiness":220,"known_move":{"name":"tackle"}},{"trigger":{"name":"level-up"},"min_happiness":220,"time_of_day":"night"}]}]}}`
	links, _, err := ExtractEvolutionLinks([]byte(chain))
	require.NoError(t, err)
	require.Len(t, links, 6)
	assert.Equal(t, "fire-stone", links[0].Item)
	assert.Equal(t, "water-stone", links[1].Item)
	assert.Equal(t, FriendshipEvolutionLevel, links[2].MinLevel)
	assert.Equal(t, "day", links[3].TimeOfDay)
	assert.Zero(t, links[3].MinLevel)
	assert.Zero(t, links[4].MinLevel)
	assert.True(t, links[4].UnsupportedConditions)
	assert.Equal(t, "night", links[5].TimeOfDay)
	assert.Nil(t, PickEvolutionLink(links[3:], 1, FriendshipEvolutionLevel))
	assert.Equal(t, "water-stone", PickEvolutionLinkForItem(links, 1, "water-stone").Item)
}

// PokéAPI decorates every evolution_details entry with metadata fields
// (is_default, version_group, and often form/region data). Regression test:
// they must not flag ordinary edges as unsupported.
func TestExtractEvolutionLinksIgnoresPokeAPIMetadata(t *testing.T) {
	chain := `{
	  "chain": {
	    "species": { "name": "pichu", "url": "https://pokeapi.co/api/v2/pokemon-species/172/" },
	    "evolution_details": [],
	    "evolves_to": [
	      {
	        "species": { "name": "pikachu", "url": "https://pokeapi.co/api/v2/pokemon-species/25/" },
	        "evolution_details": [
	          {
	            "gender": null, "held_item": null, "item": null, "known_move": null,
	            "known_move_type": null, "location": null, "max_happiness": null,
	            "max_beauty": null, "max_affection": null, "min_affection": null,
	            "min_beauty": null, "min_happiness": 220, "min_level": null,
	            "needs_overworld_rain": false, "party_species": null, "party_type": null,
	            "relative_physical_stats": null, "time_of_day": "", "trade_species": null,
	            "trigger": { "name": "level-up", "url": "https://pokeapi.co/api/v2/evolution-trigger/1/" },
	            "turn_upside_down": false,
	            "is_default": true,
	            "version_group": { "name": "gold-silver", "url": "https://pokeapi.co/api/v2/version-group/3/" }
	          }
	        ],
	        "evolves_to": [
	          {
	            "species": { "name": "raichu", "url": "https://pokeapi.co/api/v2/pokemon-species/26/" },
	            "evolution_details": [
	              {
	                "gender": null, "held_item": null, "item": { "name": "thunder-stone", "url": "https://pokeapi.co/api/v2/item/83/" },
	                "known_move": null, "known_move_type": null, "location": null,
	                "max_happiness": null, "max_beauty": null, "max_affection": null,
	                "min_affection": null, "min_beauty": null, "min_happiness": null,
	                "min_level": null, "needs_overworld_rain": false,
	                "party_species": null, "party_type": null,
	                "relative_physical_stats": null, "time_of_day": "", "trade_species": null,
	                "trigger": { "name": "use-item", "url": "https://pokeapi.co/api/v2/evolution-trigger/3/" },
	                "turn_upside_down": false,
	                "is_default": true,
	                "version_group": { "name": "red-blue", "url": "https://pokeapi.co/api/v2/version-group/1/" },
	                "required_pokemon_form": { "name": "pikachu", "url": "https://pokeapi.co/api/v2/pokemon-form/25/" }
	              },
	              {
	                "gender": null, "held_item": null, "item": { "name": "thunder-stone", "url": "https://pokeapi.co/api/v2/item/83/" },
	                "known_move": null, "known_move_type": null, "location": null,
	                "max_happiness": null, "max_beauty": null, "max_affection": null,
	                "min_affection": null, "min_beauty": null, "min_happiness": null,
	                "min_level": null, "needs_overworld_rain": false,
	                "party_species": null, "party_type": null,
	                "relative_physical_stats": null, "time_of_day": "", "trade_species": null,
	                "trigger": { "name": "use-item", "url": "https://pokeapi.co/api/v2/evolution-trigger/3/" },
	                "turn_upside_down": false,
	                "is_default": true,
	                "version_group": { "name": "sun-moon", "url": "https://pokeapi.co/api/v2/version-group/17/" },
	                "region": { "name": "alola", "url": "https://pokeapi.co/api/v2/region/7/" },
	                "required_pokemon_form": { "name": "pikachu", "url": "https://pokeapi.co/api/v2/pokemon-form/25/" },
	                "evolved_pokemon_form": { "name": "raichu-alola", "url": "https://pokeapi.co/api/v2/pokemon-form/10202/" }
	              }
	            ],
	            "evolves_to": []
	          }
	        ]
	      }
	    ]
	  }
	}`
	links, _, err := ExtractEvolutionLinks([]byte(chain))
	require.NoError(t, err)

	// The base edges stay usable; the Alolan-form entry is retained but
	// flagged unsupported (this game models base forms only).
	require.Len(t, links, 3)
	assert.Equal(t, EvolutionLink{
		FromSpeciesID: 172,
		ToSpeciesID:   25,
		MinLevel:      FriendshipEvolutionLevel,
		Trigger:       "level-up",
	}, links[0])
	baseRaichu := EvolutionLink{
		FromSpeciesID: 25,
		ToSpeciesID:   26,
		Trigger:       "use-item",
		Item:          "thunder-stone",
	}
	assert.Equal(t, baseRaichu, links[1])
	assert.Equal(t, EvolutionLink{
		FromSpeciesID:         25,
		ToSpeciesID:           26,
		Trigger:               "use-item",
		Item:                  "thunder-stone",
		UnsupportedConditions: true, // region: alola + evolved_pokemon_form: raichu-alola
	}, links[2])

	// Base edges stay usable for level/item evolution.
	assert.NotNil(t, PickEvolutionLink(links, 172, FriendshipEvolutionLevel))
	assert.Equal(t, "thunder-stone", PickEvolutionLinkForItem(links, 25, "thunder-stone").Item)
	options := EvolutionOptionsFromLinks(links, 25)
	require.Len(t, options, 1)
	assert.Equal(t, MethodItem, options[0].Method)
}

// Form- and region-conditional edges are unsupported unless they name the
// base form. required_pokemon_form matching the source species is fine;
// anything else (regional forms, alternate forms) flags the edge.
func TestExtractEvolutionLinksFormConditionalEdges(t *testing.T) {
	chain := `{
	  "chain": {
	    "species": { "name": "exeggcute", "url": "https://pokeapi.co/api/v2/pokemon-species/102/" },
	    "evolution_details": [],
	    "evolves_to": [
	      {
	        "species": { "name": "exeggutor", "url": "https://pokeapi.co/api/v2/pokemon-species/103/" },
	        "evolution_details": [
	          {
	            "item": { "name": "leaf-stone" },
	            "trigger": { "name": "use-item" },
	            "is_default": true,
	            "version_group": { "name": "red-blue" },
	            "required_pokemon_form": { "name": "exeggcute" }
	          },
	          {
	            "item": { "name": "leaf-stone" },
	            "trigger": { "name": "use-item" },
	            "is_default": true,
	            "version_group": { "name": "sun-moon" },
	            "region": { "name": "alola" },
	            "required_pokemon_form": { "name": "exeggcute" },
	            "evolved_pokemon_form": { "name": "exeggutor-alola" }
	          }
	        ],
	        "evolves_to": []
	      },
	      {
	        "species": { "name": "runerigus", "url": "https://pokeapi.co/api/v2/pokemon-species/867/" },
	        "evolution_details": [
	          {
	            "trigger": { "name": "use-item" },
	            "item": { "name": "everstone" },
	            "required_pokemon_form": { "name": "yamask-galar" }
	          }
	        ],
	        "evolves_to": []
	      }
	    ]
	  }
	}`
	links, _, err := ExtractEvolutionLinks([]byte(chain))
	require.NoError(t, err)
	require.Len(t, links, 3)

	// Base Kanto edge: usable.
	assert.False(t, links[0].UnsupportedConditions)
	assert.Equal(t, "leaf-stone", links[0].Item)
	// Alolan variant edge: flagged.
	assert.True(t, links[1].UnsupportedConditions)
	// Non-base required form (yamask-galar): flagged.
	assert.True(t, links[2].UnsupportedConditions)

	assert.Equal(t, "leaf-stone", PickEvolutionLinkForItem(links, 102, "leaf-stone").Item)
	assert.Nil(t, PickEvolutionLinkForItem(links, 102, "everstone"))
}

func TestFriendshipEdgeFeedsLevelHelpers(t *testing.T) {
	links, _, err := ExtractEvolutionLinks([]byte(pichuChainJSON))
	require.NoError(t, err)

	// The synthesized edge is a normal level option, so the deck UI shows it.
	options := EvolutionOptionsFromLinks(links, 172)
	require.Len(t, options, 1)
	assert.Equal(t, MethodLevel, options[0].Method)
	assert.Equal(t, FriendshipEvolutionLevel, options[0].MinLevel)
	assert.Equal(t, 25, options[0].ToSpeciesID)

	// Battle-reward auto-evolution picks it up once the level is reached.
	assert.Nil(t, PickEvolutionLink(links, 172, FriendshipEvolutionLevel-1), "not yet eligible")
	picked := PickEvolutionLink(links, 172, FriendshipEvolutionLevel)
	require.NotNil(t, picked)
	assert.Equal(t, 25, picked.ToSpeciesID)
}

// Pikachu family (chain id 10): pikachu evolves into raichu via use-item with
// a thunder stone. Structurally faithful to PokéAPI's payload shape.
const pikachuChainJSON = `{
  "id": 10,
  "chain": {
    "species": { "name": "pikachu", "url": "https://pokeapi.co/api/v2/pokemon-species/25/" },
    "evolution_details": [],
    "evolves_to": [
      {
        "species": { "name": "raichu", "url": "https://pokeapi.co/api/v2/pokemon-species/26/" },
        "evolution_details": [
          { "min_level": null, "item": { "name": "thunder-stone", "url": "https://pokeapi.co/api/v2/item/83/" }, "trigger": { "name": "use-item", "url": "https://pokeapi.co/api/v2/evolution-trigger/3/" } }
        ],
        "evolves_to": []
      }
    ]
  }
}`

func TestExtractEvolutionLinksCapturesItem(t *testing.T) {
	links, members, err := ExtractEvolutionLinks([]byte(pikachuChainJSON))
	require.NoError(t, err)

	assert.Equal(t, []int{25, 26}, members)
	assert.Equal(t, []EvolutionLink{
		{FromSpeciesID: 25, ToSpeciesID: 26, Trigger: "use-item", Item: "thunder-stone"},
	}, links)
}

func TestPickEvolutionLinkForItem(t *testing.T) {
	links := []EvolutionLink{
		{FromSpeciesID: 4, ToSpeciesID: 5, MinLevel: 16, Trigger: "level-up"},
		{FromSpeciesID: 25, ToSpeciesID: 26, Trigger: "use-item", Item: "thunder-stone"},
		{FromSpeciesID: 133, ToSpeciesID: 134, Trigger: "use-item", Item: "water-stone"},
	}

	tests := []struct {
		name      string
		speciesID int
		itemID    string
		wantTo    int
		wantFound bool
	}{
		{name: "pikachu with thunder stone evolves", speciesID: 25, itemID: "thunder-stone", wantTo: 26, wantFound: true},
		{name: "pikachu with wrong item rejected", speciesID: 25, itemID: "fire-stone", wantFound: false},
		{name: "pikachu with empty item rejected", speciesID: 25, itemID: "", wantFound: false},
		{name: "level-up species never matches an item", speciesID: 4, itemID: "thunder-stone", wantFound: false},
		{name: "final form has no item evolution", speciesID: 26, itemID: "thunder-stone", wantFound: false},
		{name: "branching chain matches its own item", speciesID: 133, itemID: "water-stone", wantTo: 134, wantFound: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PickEvolutionLinkForItem(links, tt.speciesID, tt.itemID)
			if !tt.wantFound {
				assert.Nil(t, got)
				return
			}
			require.NotNil(t, got)
			assert.Equal(t, tt.wantTo, got.ToSpeciesID)
		})
	}
}

func TestEvolutionOptionsFromLinks(t *testing.T) {
	links := []EvolutionLink{
		{FromSpeciesID: 4, ToSpeciesID: 5, MinLevel: 16, Trigger: "level-up"},
		{FromSpeciesID: 25, ToSpeciesID: 26, Trigger: "use-item", Item: "thunder-stone"},
		// unusable edges are left out
		{FromSpeciesID: 133, ToSpeciesID: 196, Trigger: "level-up"},           // no min level
		{FromSpeciesID: 133, ToSpeciesID: 700, Trigger: "use-item"},           // no item name
		{FromSpeciesID: 1, ToSpeciesID: 2, Trigger: "trade"},                  // unsupported mechanism
		{FromSpeciesID: 5, ToSpeciesID: 6, MinLevel: 36, Trigger: "level-up"}, // different species
	}

	options := EvolutionOptionsFromLinks(links, 25)
	assert.Equal(t, []EvolutionOption{
		{Method: MethodItem, ToSpeciesID: 26, Item: "thunder-stone"},
	}, options)

	options = EvolutionOptionsFromLinks(links, 4)
	assert.Equal(t, []EvolutionOption{
		{Method: MethodLevel, ToSpeciesID: 5, MinLevel: 16},
	}, options)

	assert.Empty(t, EvolutionOptionsFromLinks(links, 26))
}

func TestPickEvolutionLink(t *testing.T) {
	links := []EvolutionLink{
		{FromSpeciesID: 4, ToSpeciesID: 5, MinLevel: 16, Trigger: "level-up"},
		{FromSpeciesID: 5, ToSpeciesID: 6, MinLevel: 36, Trigger: "level-up"},
		{FromSpeciesID: 133, ToSpeciesID: 135, Trigger: "use-item"},
	}

	t.Run("below threshold", func(t *testing.T) {
		assert.Nil(t, PickEvolutionLink(links, 4, 15))
	})

	t.Run("at threshold", func(t *testing.T) {
		got := PickEvolutionLink(links, 4, 16)
		require.NotNil(t, got)
		assert.Equal(t, 5, got.ToSpeciesID)
	})

	t.Run("second stage", func(t *testing.T) {
		got := PickEvolutionLink(links, 5, 40)
		require.NotNil(t, got)
		assert.Equal(t, 6, got.ToSpeciesID)
	})

	t.Run("final form has no further evolution", func(t *testing.T) {
		assert.Nil(t, PickEvolutionLink(links, 6, 50))
	})

	t.Run("non-level-up triggers are ignored", func(t *testing.T) {
		assert.Nil(t, PickEvolutionLink(links, 133, 50))
	})
}
