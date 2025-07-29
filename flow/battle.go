package flow

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// BattleFlow handles battle-related operations
type BattleFlow struct {
	gameState   *core.GameState
	battlefield *core.Battlefield
}

// NewBattleFlow creates a new BattleFlow
func NewBattleFlow(gameState *core.GameState, battlefield *core.Battlefield) *BattleFlow {
	return &BattleFlow{
		gameState:   gameState,
		battlefield: battlefield,
	}
}

// PlaceCard adds a battle card to the battlefield
func (bf *BattleFlow) PlaceCard(card *core.BattleCard) bool {
	if bf.battlefield == nil {
		return false
	}
	
	// Check if there's space for the card
	if len(bf.battlefield.BattleCards) >= bf.battlefield.CardSlot {
		return false
	}
	
	// Add card to battlefield
	bf.battlefield.BattleCards = append(bf.battlefield.BattleCards, card)
	
	// Remove from deck
	bf.gameState.CardDeck.Remove(card.CardID)
	
	return true
}

// RemoveFromBattle removes a card from battle at the specified index
func (bf *BattleFlow) RemoveFromBattle(cardIndex int) bool {
	if bf.battlefield == nil || cardIndex < 0 || cardIndex >= len(bf.battlefield.BattleCards) {
		return false
	}
	
	// Get the card to remove
	card := bf.battlefield.BattleCards[cardIndex]
	
	// Remove from battlefield
	bf.battlefield.BattleCards = append(
		bf.battlefield.BattleCards[:cardIndex],
		bf.battlefield.BattleCards[cardIndex+1:]...,
	)
	
	// Return to deck
	bf.gameState.CardDeck.Add(card.CardID)
	
	return true
}

// Conquer attempts to conquer the current battle point
func (bf *BattleFlow) Conquer() bool {
	if bf.battlefield == nil {
		return false
	}
	
	if !bf.battlefield.CanBeat() {
		return false
	}
	
	// Conquest logic - mark point as conquered
	// The actual conquest logic should be implemented based on the point type
	// For now, we return true to indicate successful conquest
	return true
}

// Rollback returns all cards from battlefield to deck and resets battlefield
func (bf *BattleFlow) Rollback() {
	if bf.battlefield == nil {
		return
	}
	
	// Return all cards to deck
	for _, card := range bf.battlefield.BattleCards {
		bf.gameState.CardDeck.Add(card.CardID)
	}
	
	// Clear battlefield
	bf.battlefield.BattleCards = make([]*core.BattleCard, 0)
}

// CanPlaceCard checks if a card can be placed
func (bf *BattleFlow) CanPlaceCard() bool {
	if bf.battlefield == nil {
		return false
	}
	return len(bf.battlefield.BattleCards) < bf.battlefield.CardSlot
} 