package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestTerritory_AppendCard(t *testing.T) {
	terrain := core.NewTerrain("test_terrain", core.ResourceQuantity{
		Money: 10, Food: 5, Wood: 3, Iron: 2, Mana: 1,
	}, 2)
	territory := core.NewTerritory("test_territory", terrain)

	card1 := core.NewStructureCard("structure_1", core.ResourceQuantity{}, core.NewResourceModifier(), 0.0, 0)
	card2 := core.NewStructureCard("structure_2", core.ResourceQuantity{}, core.NewResourceModifier(), 0.0, 0)
	card3 := core.NewStructureCard("structure_3", core.ResourceQuantity{}, core.NewResourceModifier(), 0.0, 0)

	tests := []struct {
		name     string
		card     *core.StructureCard
		expected bool
	}{
		{
			name:     "Add first card",
			card:     card1,
			expected: true,
		},
		{
			name:     "Add second card",
			card:     card2,
			expected: true,
		},
		{
			name:     "Exceed card slot limit",
			card:     card3,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := territory.AppendCard(tt.card)
			if result != tt.expected {
				t.Errorf("AppendCard() = %v, want %v", result, tt.expected)
			}
		})
	}

	// Check number of cards
	if len(territory.Cards()) != 2 {
		t.Errorf("Cards length = %v, want %v", len(territory.Cards()), 2)
	}
}

func TestTerritory_RemoveCard(t *testing.T) {
	card1 := core.NewStructureCard("structure_1", core.ResourceQuantity{}, core.NewResourceModifier(), 0.0, 0)
	card2 := core.NewStructureCard("structure_2", core.ResourceQuantity{}, core.NewResourceModifier(), 0.0, 0)

	terrain := core.NewTerrain("test_terrain", core.ResourceQuantity{
		Money: 10, Food: 5, Wood: 3, Iron: 2, Mana: 1,
	}, 3)
	territory := core.NewTerritory("test_territory", terrain)

	// Add cards to territory
	territory.AppendCard(card1)
	territory.AppendCard(card2)

	tests := []struct {
		name         string
		index        int
		expectedCard *core.StructureCard
		expectedOk   bool
	}{
		{
			name:         "Valid index (0)",
			index:        0,
			expectedCard: card1,
			expectedOk:   true,
		},
		{
			name:         "Valid index (1)",
			index:        1,
			expectedCard: card2,
			expectedOk:   true,
		},
		{
			name:         "Invalid index (-1)",
			index:        -1,
			expectedCard: nil,
			expectedOk:   false,
		},
		{
			name:         "Invalid index (out of range)",
			index:        5,
			expectedCard: nil,
			expectedOk:   false,
		},
	}

	// First test (remove index 0)
	card, ok := territory.RemoveCard(0)
	if !ok {
		t.Errorf("RemoveCard(0) ok = %v, want %v", ok, true)
	}
	if card != card1 {
		t.Errorf("RemoveCard(0) card = %v, want %v", card, card1)
	}
	if len(territory.Cards()) != 1 {
		t.Errorf("Cards length after removal = %v, want %v", len(territory.Cards()), 1)
	}

	// Remaining tests
	for i, tt := range tests[1:] {
		t.Run(tt.name, func(t *testing.T) {
			// Restore territory.Cards (only card2 remains)
			if i == 0 {
				// Test index 1 -> becomes index 0
				card, ok := territory.RemoveCard(0)
				if ok != tt.expectedOk {
					t.Errorf("RemoveCard() ok = %v, want %v", ok, tt.expectedOk)
				}
				if card != tt.expectedCard {
					t.Errorf("RemoveCard() card = %v, want %v", card, tt.expectedCard)
				}
			} else {
				// Test invalid index
				card, ok := territory.RemoveCard(tt.index)
				if ok != tt.expectedOk {
					t.Errorf("RemoveCard() ok = %v, want %v", ok, tt.expectedOk)
				}
				if card != tt.expectedCard {
					t.Errorf("RemoveCard() card = %v, want %v", card, tt.expectedCard)
				}
			}
		})
	}
}

func TestTerritory_Yield(t *testing.T) {
	// ResourceModifiers for testing
	doubleMoneyModifier := core.ResourceModifier{
		Money: 1.0, // 100% increase = double
		Food:  0.0,
		Wood:  0.0,
		Iron:  0.0,
		Mana:  0.0,
	}

	// StructureCard without modifier
	simpleCard := core.NewStructureCard("simple_card", core.ResourceQuantity{}, core.NewResourceModifier(), 0.0, 0)

	// StructureCard with multiplicative modifier
	bonusCard := core.NewStructureCard("bonus_card", core.ResourceQuantity{}, doubleMoneyModifier, 0.0, 0)

	// StructureCard with additive value
	addCard := core.NewStructureCard("add_card", core.ResourceQuantity{
		Money: 5, Food: 3, Wood: 2, Iron: 1, Mana: 1,
	}, core.NewResourceModifier(), 0.0, 0)

	tests := []struct {
		name     string
		cards    []*core.StructureCard
		expected core.ResourceQuantity
	}{
		{
			name:  "No cards",
			cards: []*core.StructureCard{},
			expected: core.ResourceQuantity{
				Money: 10, Food: 5, Wood: 3, Iron: 2, Mana: 1,
			},
		},
		{
			name:  "Card without modifier",
			cards: []*core.StructureCard{simpleCard},
			expected: core.ResourceQuantity{
				Money: 10, Food: 5, Wood: 3, Iron: 2, Mana: 1,
			},
		},
		{
			name:  "Modifier that doubles money",
			cards: []*core.StructureCard{bonusCard},
			expected: core.ResourceQuantity{
				Money: 20, Food: 5, Wood: 3, Iron: 2, Mana: 1, // Money: 10 * 2 = 20
			},
		},
		{
			name:  "Additive card that adds resources",
			cards: []*core.StructureCard{addCard},
			expected: core.ResourceQuantity{
				Money: 15, Food: 8, Wood: 5, Iron: 3, Mana: 2, // Base + (5,3,2,1,1)
			},
		},
		{
			name:  "Multiple modifiers (additive first, then multiplicative)",
			cards: []*core.StructureCard{bonusCard, addCard},
			expected: core.ResourceQuantity{
				Money: 30, Food: 8, Wood: 5, Iron: 3, Mana: 2, // (10 + 5) * 2 = 30 for Money
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terrain := core.NewTerrain("test_terrain", core.ResourceQuantity{
				Money: 10, Food: 5, Wood: 3, Iron: 2, Mana: 1,
			}, 3)
			territory := core.NewTerritory("test_territory", terrain)

			// Add cards to territory
			for _, card := range tt.cards {
				territory.AppendCard(card)
			}

			result := territory.Yield()
			if result != tt.expected {
				t.Errorf("Yield() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTerritory_Basic(t *testing.T) {
	terrain := core.NewTerrain("basic_terrain", core.ResourceQuantity{
		Money: 15, Food: 10, Wood: 5, Iron: 3, Mana: 2,
	}, 2)
	territory := core.NewTerritory("basic_territory", terrain)

	if territory.ID() != "basic_territory" {
		t.Errorf("ID() = %v, want %v", territory.ID(), "basic_territory")
	}
	if territory.Terrain().CardSlot() != 2 {
		t.Errorf("Terrain().CardSlot() = %v, want %v", territory.Terrain().CardSlot(), 2)
	}
	if len(territory.Cards()) != 0 {
		t.Errorf("Cards() length = %v, want %v", len(territory.Cards()), 0)
	}

	expectedYield := core.ResourceQuantity{
		Money: 15, Food: 10, Wood: 5, Iron: 3, Mana: 2,
	}
	yield := territory.Yield()
	if yield != expectedYield {
		t.Errorf("Yield() = %v, want %v", yield, expectedYield)
	}
}

// Note: mockYieldModifier is no longer needed as we use ResourceModifier directly
