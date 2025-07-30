package viewmodel

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// CardDeckViewModel provides display information for card deck UI
type CardDeckViewModel struct {
	gameState     *core.GameState
	cardGenerator *core.CardGenerator
}

// NewCardDeckViewModel creates a new CardDeckViewModel
func NewCardDeckViewModel(gameState *core.GameState, cardGenerator *core.CardGenerator) *CardDeckViewModel {
	return &CardDeckViewModel{
		gameState:     gameState,
		cardGenerator: cardGenerator,
	}
}

// NumBattleCards returns the number of unique battle card types in the deck
func (vm *CardDeckViewModel) NumBattleCards() int {
	if vm.gameState == nil || vm.gameState.CardDeck == nil {
		return 0
	}

	cardIDs := vm.gameState.CardDeck.GetAllCardIDs()

	// Use a set to count unique battle card types
	uniqueBattleCards := make(map[core.CardID]bool)

	if vm.cardGenerator != nil {
		cards, ok := vm.cardGenerator.Generate(cardIDs)
		if ok {
			for _, card := range cards.BattleCards {
				uniqueBattleCards[card.CardID] = true
			}
		}
	}

	return len(uniqueBattleCards)
}

// NumStructureCards returns the number of unique structure card types in the deck
func (vm *CardDeckViewModel) NumStructureCards() int {
	if vm.gameState == nil || vm.gameState.CardDeck == nil {
		return 0
	}

	cardIDs := vm.gameState.CardDeck.GetAllCardIDs()

	// Use a set to count unique structure card types
	uniqueStructureCards := make(map[core.CardID]bool)

	if vm.cardGenerator != nil {
		cards, ok := vm.cardGenerator.Generate(cardIDs)
		if ok {
			for _, card := range cards.StructureCards {
				uniqueStructureCards[card.ID()] = true
			}
		}
	}

	return len(uniqueStructureCards)
}

// BattleCard returns battle card view model at the specified index
func (vm *CardDeckViewModel) BattleCard(idx int) *BattleCardViewModel {
	battleCards := vm.getBattleCards()
	if idx < 0 || idx >= len(battleCards) {
		return nil
	}

	card := battleCards[idx]
	return NewBattleCardViewModel(vm.gameState, card, float64(card.Power()))
}

// StructureCard returns structure card view model at the specified index
func (vm *CardDeckViewModel) StructureCard(idx int) *StructureCardViewModel {
	structureCards := vm.getStructureCards()
	if idx < 0 || idx >= len(structureCards) {
		return nil
	}

	card := structureCards[idx]
	return NewStructureCardViewModel(vm.gameState, card)
}

// getBattleCards returns all battle cards in the deck
func (vm *CardDeckViewModel) getBattleCards() []*core.BattleCard {
	if vm.gameState == nil || vm.gameState.CardDeck == nil || vm.cardGenerator == nil {
		return []*core.BattleCard{}
	}

	cardIDs := vm.gameState.CardDeck.GetAllCardIDs()
	if len(cardIDs) == 0 {
		return []*core.BattleCard{}
	}

	cards, ok := vm.cardGenerator.Generate(cardIDs)
	if !ok {
		return []*core.BattleCard{}
	}

	return cards.BattleCards
}

// getStructureCards returns all structure cards in the deck
func (vm *CardDeckViewModel) getStructureCards() []*core.StructureCard {
	if vm.gameState == nil || vm.gameState.CardDeck == nil || vm.cardGenerator == nil {
		return []*core.StructureCard{}
	}

	cardIDs := vm.gameState.CardDeck.GetAllCardIDs()
	if len(cardIDs) == 0 {
		return []*core.StructureCard{}
	}

	cards, ok := vm.cardGenerator.Generate(cardIDs)
	if !ok {
		return []*core.StructureCard{}
	}

	return cards.StructureCards
}

// GetDuplicateCount returns the number of duplicates for a specific card
func (vm *CardDeckViewModel) GetDuplicateCount(card interface{}) int {
	if vm.gameState == nil || vm.gameState.CardDeck == nil {
		return 1
	}

	switch c := card.(type) {
	case *core.BattleCard:
		return vm.gameState.CardDeck.Count(c.CardID)
	case *core.StructureCard:
		return vm.gameState.CardDeck.Count(c.ID())
	default:
		return 1
	}
}
