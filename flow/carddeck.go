package flow

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// CardDeckFlow handles card deck operations
type CardDeckFlow struct {
	gameState *core.GameState
}

// NewCardDeckFlow creates a new CardDeckFlow
func NewCardDeckFlow(gameState *core.GameState) *CardDeckFlow {
	return &CardDeckFlow{
		gameState: gameState,
	}
}

func (f *CardDeckFlow) PlayBattleCardInBattle(id core.CardID) {
	battleCard, ok := f.gameState.CardDictionary.BattleCard(id)
	if !ok {
		return
	}
	countInHand := f.gameState.CardDeck.Count(id)
	if countInHand == 0 {
		return
	}

	battlefield, ok := f.gameState.Battlefield()
	if !ok {
		return
	}

	battlefield.AddBattleCard(battleCard)
	f.gameState.CardDeck.Remove(id)
}

func (f *CardDeckFlow) PlayStructureCardInTerritory(id core.CardID) {
	structureCard, ok := f.gameState.CardDictionary.StructureCard(id)
	if !ok {
		return
	}
	countInHand := f.gameState.CardDeck.Count(id)
	if countInHand == 0 {
		return
	}

	plan, ok := f.gameState.ConstructionPlan()
	if !ok {
		return
	}
	plan.AddCard(structureCard)
	f.gameState.CardDeck.Remove(id)
}
