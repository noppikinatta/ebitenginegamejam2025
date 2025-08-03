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

	f.gameState.CardDeck.Remove(id)
	//f.gameState. // TODO: can reference battle field
}

func (f *CardDeckFlow) PlayStructureCardInTerritory(id core.CardID) {

}
