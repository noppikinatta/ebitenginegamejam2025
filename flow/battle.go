package flow

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// BattleFlow handles battle-related operations
type BattleFlow struct {
	gameState   *core.GameState
	battlefield *core.Battlefield
}

// NewBattleFlow creates a new BattleFlow
func NewBattleFlow(gameState *core.GameState) *BattleFlow {
	return &BattleFlow{
		gameState: gameState,
	}
}

func (bf *BattleFlow) Select(x, y int) (*viewmodel.BattleViewModel, bool) {
	ok := bf.gameState.InitBattlefield(x, y)
	if !ok {
		return nil, false
	}

	bf.battlefield, ok = bf.gameState.Battlefield()
	if !ok {
		return nil, false
	}

	return viewmodel.NewBattleViewModel(bf.gameState, bf.battlefield, bf.battlefield.Point), true
}

// RemoveFromBattle removes a card from battle at the specified index
func (bf *BattleFlow) RemoveFromBattle(cardIndex int) bool {
	battlefield, ok := bf.gameState.Battlefield()
	if !ok || cardIndex < 0 || cardIndex >= len(battlefield.BattleCards) {
		return false
	}

	// Get the card to remove
	card := battlefield.BattleCards[cardIndex]

	// Remove from battlefield
	battlefield.BattleCards = append(
		battlefield.BattleCards[:cardIndex],
		battlefield.BattleCards[cardIndex+1:]...,
	)

	// Return to deck
	bf.gameState.CardDeck.Add(card.CardID)

	return true
}

// Conquer attempts to conquer the current battle point
func (bf *BattleFlow) Conquer() bool {
	battlefield, ok := bf.gameState.Battlefield()
	if !ok {
		return false
	}

	if !battlefield.CanBeat() {
		return false
	}

	// Mark conquest on game state
	bf.gameState.Conquer()
	return true
}

// Rollback returns all cards from battlefield to deck and resets battlefield
func (bf *BattleFlow) Rollback() {
	battlefield, ok := bf.gameState.Battlefield()
	if !ok {
		return
	}

	// Return all cards to deck
	for _, card := range battlefield.BattleCards {
		bf.gameState.CardDeck.Add(card.CardID)
	}

	// Clear battlefield
	battlefield.BattleCards = make([]*core.BattleCard, 0)
}
