// This file defines the Player model for the Pokémon CLI application.

package models

import (
	"pokemon-cli/internal/pokemon"
)

type Player struct {
	Name       string
	Deck       []pokemon.Card
	CurrentIdx int // index of the current card in use
}

func NewPlayer(name string, deck []pokemon.Card) *Player {
	return &Player{
		Name:       name,
		Deck:       deck,
		CurrentIdx: 0,
	}
}

func (p *Player) CurrentCard() *pokemon.Card {
	if len(p.Deck) == 0 || p.CurrentIdx < 0 || p.CurrentIdx >= len(p.Deck) {
		return nil
	}
	return &p.Deck[p.CurrentIdx]
}

func (p *Player) AllCards() []pokemon.Card {
	return p.Deck
}
