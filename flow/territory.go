package flow

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// TerritoryFlow handles territory construction operations
type TerritoryFlow struct {
	gameState     *core.GameState
	territory     *core.Territory
	originalCards []*core.StructureCard // Backup for rollback
}

// NewTerritoryFlow creates a new TerritoryFlow
func NewTerritoryFlow(gameState *core.GameState, territory *core.Territory) *TerritoryFlow {
	// Create backup of original cards for rollback functionality
	originalCards := make([]*core.StructureCard, len(territory.Cards()))
	copy(originalCards, territory.Cards())

	return &TerritoryFlow{
		gameState:     gameState,
		territory:     territory,
		originalCards: originalCards,
	}
}

// PlaceCard adds a structure card to the construction plan
func (tf *TerritoryFlow) PlaceCard(card *core.StructureCard) bool {
	if tf.territory == nil {
		return false
	}

	// Check if there's space for the card
	if !tf.CanPlaceCard() {
		return false
	}

	// Add card to territory (temporary placement)
	success := tf.territory.AppendCard(card)
	if !success {
		return false
	}

	// Remove from deck
	tf.gameState.CardDeck.Remove(card.ID())

	return true
}

// RemoveFromPlan removes a card from construction plan at the specified index
func (tf *TerritoryFlow) RemoveFromPlan(cardIndex int) bool {
	if tf.territory == nil {
		return false
	}

	cards := tf.territory.Cards()
	if cardIndex < 0 || cardIndex >= len(cards) {
		return false
	}

	// Remove card from territory
	removedCard, ok := tf.territory.RemoveCard(cardIndex)
	if !ok {
		return false
	}

	// Return to deck
	tf.gameState.CardDeck.Add(removedCard.ID())

	return true
}

// Commit applies the construction plan to the territory
func (tf *TerritoryFlow) Commit() {
	// In the current implementation, changes are applied immediately
	// This method could be used for batch operations in the future
	// For now, it's a no-op since cards are placed directly
}

// Rollback reverts all changes to the original state
func (tf *TerritoryFlow) Rollback() {
	if tf.territory == nil {
		return
	}

	// Return any cards that were added back to deck
	currentCards := tf.territory.Cards()
	for _, currentCard := range currentCards {
		found := false
		for _, originalCard := range tf.originalCards {
			if currentCard == originalCard {
				found = true
				break
			}
		}
		if !found {
			// This card was added during the session, return it to deck
			tf.gameState.CardDeck.Add(currentCard.ID())
		}
	}

	// Re-add any cards that were removed
	for _, originalCard := range tf.originalCards {
		found := false
		for _, currentCard := range currentCards {
			if currentCard == originalCard {
				found = true
				break
			}
		}
		if !found {
			// This card was removed during the session, remove it from deck and add back to territory
			tf.gameState.CardDeck.Remove(originalCard.ID())
			tf.territory.AppendCard(originalCard)
		}
	}
}

// CanPlaceCard checks if a card can be placed in the territory
func (tf *TerritoryFlow) CanPlaceCard() bool {
	if tf.territory == nil {
		return false
	}
	return len(tf.territory.Cards()) < tf.territory.Terrain().CardSlot()
}

// GetCurrentCardCount returns the current number of cards in the territory
func (tf *TerritoryFlow) GetCurrentCardCount() int {
	if tf.territory == nil {
		return 0
	}
	return len(tf.territory.Cards())
}

// GetMaxCardSlot returns the maximum number of cards that can be placed
func (tf *TerritoryFlow) GetMaxCardSlot() int {
	if tf.territory == nil {
		return 0
	}
	return tf.territory.Terrain().CardSlot()
}
