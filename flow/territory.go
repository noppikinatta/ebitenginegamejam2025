package flow

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// TerritoryFlow handles territory construction operations
type TerritoryFlow struct {
	gameState   *core.GameState
	territory   *core.Territory
	currentPlan *core.ConstructionPlan
}

// NewTerritoryFlow creates a new TerritoryFlow
func NewTerritoryFlow(gameState *core.GameState) *TerritoryFlow {
	return &TerritoryFlow{
		gameState: gameState,
	}
}

func (tf *TerritoryFlow) SelectTerritory(x, y int) (*viewmodel.TerritoryViewModel, bool) {
	point, ok := tf.gameState.MapGrid.GetPoint(x, y)
	if !ok {
		return nil, false
	}

	territoryPoint, ok := point.AsTerritoryPoint()
	if !ok {
		return nil, false
	}

	tf.territory = territoryPoint.Territory()
	tf.currentPlan = core.NewConstructionPlan(tf.territory)
	vm := viewmodel.NewTerritoryViewModel(tf.territory, tf.currentPlan)
	return vm, true
}

// PlaceCard adds a structure card to the construction plan
func (tf *TerritoryFlow) PlaceCard(card *core.StructureCard) bool {
	// Check if there's space for the card
	if !tf.CanPlaceCard() {
		return false
	}

	// Add card to territory (temporary placement)
	success := tf.currentPlan.AddCard(card)
	if !success {
		return false
	}

	// Remove from deck
	tf.gameState.CardDeck.Remove(card.ID())

	return true
}

// RemoveFromPlan removes a card from construction plan at the specified index
func (tf *TerritoryFlow) RemoveFromPlan(cardIndex int) bool {
	if cardIndex < 0 || cardIndex >= len(tf.currentPlan.Cards()) {
		return false
	}

	// Remove card from territory
	removedCard, ok := tf.currentPlan.RemoveCard(cardIndex)
	if !ok {
		return false
	}

	// Return to deck
	tf.gameState.CardDeck.Add(removedCard.ID())

	return true
}

// Commit applies the construction plan to the territory
func (tf *TerritoryFlow) Commit() {
	if tf.currentPlan == nil {
		return
	}

	tf.territory.ApplyConstructionPlan(tf.currentPlan)
	tf.currentPlan = nil
}

// Rollback reverts all changes to the original state
func (tf *TerritoryFlow) Rollback() {
	if tf.currentPlan == nil {
		return
	}

	delta := tf.currentPlan.GetRollbackCards()
	tf.gameState.CardDeck.ApplyDelta(delta)
	tf.currentPlan = nil
}

// CanPlaceCard checks if a card can be placed in the territory
func (tf *TerritoryFlow) CanPlaceCard() bool {
	if tf.currentPlan == nil {
		return false
	}
	return tf.currentPlan.CanPlaceCard()
}
