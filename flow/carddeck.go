package flow

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// CardDeckFlow handles card deck operations
type CardDeckFlow struct {
	gameState     *core.GameState
	cardGenerator *core.CardGenerator
	selectedIndex int
}

// NewCardDeckFlow creates a new CardDeckFlow
func NewCardDeckFlow(gameState *core.GameState, cardGenerator *core.CardGenerator) *CardDeckFlow {
	return &CardDeckFlow{
		gameState:     gameState,
		cardGenerator: cardGenerator,
		selectedIndex: -1, // No selection initially
	}
}

// Select selects a card at the specified index
func (cf *CardDeckFlow) Select(cardIndex int) bool {
	allCards := cf.getAllCards()
	
	if cardIndex < 0 || cardIndex >= len(allCards) {
		cf.selectedIndex = -1
		return false
	}
	
	cf.selectedIndex = cardIndex
	return true
}

// GetSelectedCard returns the currently selected card
func (cf *CardDeckFlow) GetSelectedCard() interface{} {
	if cf.selectedIndex < 0 {
		return nil
	}
	
	allCards := cf.getAllCards()
	if cf.selectedIndex >= len(allCards) {
		return nil
	}
	
	return allCards[cf.selectedIndex]
}

// GetSelectedIndex returns the currently selected index
func (cf *CardDeckFlow) GetSelectedIndex() int {
	return cf.selectedIndex
}

// ClearSelection clears the current selection
func (cf *CardDeckFlow) ClearSelection() {
	cf.selectedIndex = -1
}

// GetAllCards returns all cards in the deck
func (cf *CardDeckFlow) GetAllCards() []interface{} {
	return cf.getAllCards()
}

// getAllCards gets all cards in a single slice
func (cf *CardDeckFlow) getAllCards() []interface{} {
	if cf.gameState == nil || cf.gameState.CardDeck == nil || cf.cardGenerator == nil {
		return []interface{}{}
	}
	
	cardIDs := cf.gameState.CardDeck.GetAllCardIDs()
	if len(cardIDs) == 0 {
		return []interface{}{}
	}
	
	// Generate cards from CardIDs
	cards, ok := cf.cardGenerator.Generate(cardIDs)
	if !ok {
		return []interface{}{}
	}
	
	var allCards []interface{}
	
	// Add battle cards
	for _, card := range cards.BattleCards {
		allCards = append(allCards, card)
	}
	
	// Add structure cards
	for _, card := range cards.StructureCards {
		allCards = append(allCards, card)
	}
	
	return allCards
}

// GetBattleCards returns only battle cards from the deck
func (cf *CardDeckFlow) GetBattleCards() []*core.BattleCard {
	if cf.gameState == nil || cf.gameState.CardDeck == nil || cf.cardGenerator == nil {
		return []*core.BattleCard{}
	}
	
	cardIDs := cf.gameState.CardDeck.GetAllCardIDs()
	if len(cardIDs) == 0 {
		return []*core.BattleCard{}
	}
	
	cards, ok := cf.cardGenerator.Generate(cardIDs)
	if !ok {
		return []*core.BattleCard{}
	}
	
	return cards.BattleCards
}

// GetStructureCards returns only structure cards from the deck
func (cf *CardDeckFlow) GetStructureCards() []*core.StructureCard {
	if cf.gameState == nil || cf.gameState.CardDeck == nil || cf.cardGenerator == nil {
		return []*core.StructureCard{}
	}
	
	cardIDs := cf.gameState.CardDeck.GetAllCardIDs()
	if len(cardIDs) == 0 {
		return []*core.StructureCard{}
	}
	
	cards, ok := cf.cardGenerator.Generate(cardIDs)
	if !ok {
		return []*core.StructureCard{}
	}
	
	return cards.StructureCards
} 