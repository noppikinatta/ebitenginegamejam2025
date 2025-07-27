package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// MockIntner is a mock random number generator for testing.
type MockIntner struct {
	values []int
	count  int
}

func (m *MockIntner) Intn(n int) int {
	if m.count >= len(m.values) {
		return 0 // Default value
	}
	result := m.values[m.count] % n
	m.count++
	return result
}

func TestCardPack_Open(t *testing.T) {
	tests := []struct {
		name         string
		cardPack     core.CardPack
		intner       core.Intner
		expectLength int
		containsAll  []core.CardID
	}{
		{
			name: "Basic card pack opening",
			cardPack: core.CardPack{
				CardPackID: "basic_pack",
				Ratios: map[core.CardID]int{
					"card_a": 10,
					"card_b": 20,
					"card_c": 30,
				},
				NumPerOpen: 3,
			},
			intner:       &MockIntner{values: []int{5, 35, 55}},
			expectLength: 3,
			containsAll:  []core.CardID{"card_a", "card_b", "card_c"},
		},
		{
			name: "When multiple of the same card are drawn",
			cardPack: core.CardPack{
				CardPackID: "duplicate_pack",
				Ratios: map[core.CardID]int{
					"card_a": 50,
					"card_b": 50,
				},
				NumPerOpen: 3,
			},
			intner:       &MockIntner{values: []int{25, 75, 25}},
			expectLength: 3,
			containsAll:  []core.CardID{"card_a", "card_b"}, // Check that the drawn cards are from this set
		},
		{
			name: "Single card",
			cardPack: core.CardPack{
				CardPackID: "single_pack",
				Ratios: map[core.CardID]int{
					"only_card": 100,
				},
				NumPerOpen: 1,
			},
			intner:       &MockIntner{values: []int{50}},
			expectLength: 1,
			containsAll:  []core.CardID{"only_card"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cardPack.Open(tt.intner)
			if len(result) != tt.expectLength {
				t.Errorf("Open() returned %d cards, want %d", len(result), tt.expectLength)
				return
			}

			// For a single card, check strictly
			if tt.name == "Single card" {
				if result[0] != tt.containsAll[0] {
					t.Errorf("Open() = %v, want %v", result[0], tt.containsAll[0])
				}
				return
			}

			// For multiple cards, ensure only expected cards are included
			for _, card := range result {
				found := false
				for _, expected := range tt.containsAll {
					if card == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Open() contains unexpected card %v", card)
				}
			}
		})
	}
}

func TestCardDatabase_GetCards(t *testing.T) {
	// Test card database
	gen := core.CardGenerator{
		BattleCards: map[core.CardID]*core.BattleCard{
			"battle_1": {
				CardID:    "battle_1",
				BasePower: 10.0,
				Type:      "warrior",
			},
			"battle_2": {
				CardID:    "battle_2",
				BasePower: 15.0,
				Type:      "mage",
			},
		},
		StructureCards: map[core.CardID]*core.StructureCard{
			"structure_1": core.NewStructureCard("structure_1", core.ResourceQuantity{}, core.NewResourceModifier(), 0.0, 0),
		},
	}

	tests := []struct {
		name        string
		cardIDs     []core.CardID
		expectOk    bool
		expectCards func(*core.Cards) bool
	}{
		{
			name:     "Only existing cards",
			cardIDs:  []core.CardID{"battle_1", "structure_1"},
			expectOk: true,
			expectCards: func(cards *core.Cards) bool {
				return len(cards.BattleCards) == 1 &&
					len(cards.StructureCards) == 1 &&
					cards.BattleCards[0].CardID == "battle_1" &&
					cards.StructureCards[0].ID() == "structure_1"
			},
		},
		{
			name:     "Contains non-existent cards",
			cardIDs:  []core.CardID{"battle_1", "nonexistent"},
			expectOk: false,
			expectCards: func(cards *core.Cards) bool {
				return cards == nil
			},
		},
		{
			name:     "Empty list",
			cardIDs:  []core.CardID{},
			expectOk: true,
			expectCards: func(cards *core.Cards) bool {
				return len(cards.BattleCards) == 0 &&
					len(cards.StructureCards) == 0
			},
		},
		{
			name:     "Multiple cards of the same type",
			cardIDs:  []core.CardID{"battle_1", "battle_2"},
			expectOk: true,
			expectCards: func(cards *core.Cards) bool {
				return len(cards.BattleCards) == 2 &&
					len(cards.StructureCards) == 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cards, ok := gen.Generate(tt.cardIDs)
			if ok != tt.expectOk {
				t.Errorf("GetCards() ok = %v, want %v", ok, tt.expectOk)
				return
			}
			if !tt.expectCards(cards) {
				t.Errorf("GetCards() returned unexpected cards")
			}
		})
	}
}

func TestCardDeck_Add(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[core.CardID]int
		cardID   core.CardID
		expected map[core.CardID]int
	}{
		{
			name:     "新しいカードを追加",
			initial:  map[core.CardID]int{},
			cardID:   "card1",
			expected: map[core.CardID]int{"card1": 1},
		},
		{
			name:     "既存カードの枚数増加",
			initial:  map[core.CardID]int{"card1": 1},
			cardID:   "card1",
			expected: map[core.CardID]int{"card1": 2},
		},
		{
			name:     "複数の異なるカード",
			initial:  map[core.CardID]int{"card1": 2, "card2": 1},
			cardID:   "card3",
			expected: map[core.CardID]int{"card1": 2, "card2": 1, "card3": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deck := &core.CardDeck{}
			// Initialize with reflection to access private field for testing
			if len(tt.initial) > 0 {
				deck = core.NewCardDeck()
				for cardID, count := range tt.initial {
					for i := 0; i < count; i++ {
						deck.Add(cardID)
					}
				}
			} else {
				deck = core.NewCardDeck()
			}

			deck.Add(tt.cardID)

			// Verify the result by checking counts
			for expectedCardID, expectedCount := range tt.expected {
				if actualCount := deck.Count(expectedCardID); actualCount != expectedCount {
					t.Errorf("Count(%s) = %d, want %d", expectedCardID, actualCount, expectedCount)
				}
			}
		})
	}
}

func TestCardDeck_Remove(t *testing.T) {
	tests := []struct {
		name         string
		initial      map[core.CardID]int
		cardID       core.CardID
		expectedOk   bool
		expectedDeck map[core.CardID]int
	}{
		{
			name:         "存在しないカードを削除",
			initial:      map[core.CardID]int{},
			cardID:       "card1",
			expectedOk:   false,
			expectedDeck: map[core.CardID]int{},
		},
		{
			name:         "1枚のカードを削除",
			initial:      map[core.CardID]int{"card1": 1},
			cardID:       "card1",
			expectedOk:   true,
			expectedDeck: map[core.CardID]int{},
		},
		{
			name:         "複数枚のうち1枚を削除",
			initial:      map[core.CardID]int{"card1": 3},
			cardID:       "card1",
			expectedOk:   true,
			expectedDeck: map[core.CardID]int{"card1": 2},
		},
		{
			name:         "0枚のカードを削除",
			initial:      map[core.CardID]int{"card1": 0},
			cardID:       "card1",
			expectedOk:   false,
			expectedDeck: map[core.CardID]int{},
		},
		{
			name:         "複数カード中の1つを削除",
			initial:      map[core.CardID]int{"card1": 2, "card2": 1, "card3": 3},
			cardID:       "card2",
			expectedOk:   true,
			expectedDeck: map[core.CardID]int{"card1": 2, "card3": 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deck := core.NewCardDeck()
			// Initialize deck
			for cardID, count := range tt.initial {
				for i := 0; i < count; i++ {
					deck.Add(cardID)
				}
			}

			ok := deck.Remove(tt.cardID)

			if ok != tt.expectedOk {
				t.Errorf("Remove(%s) = %v, want %v", tt.cardID, ok, tt.expectedOk)
			}

			// Verify the deck state
			for expectedCardID, expectedCount := range tt.expectedDeck {
				if actualCount := deck.Count(expectedCardID); actualCount != expectedCount {
					t.Errorf("After Remove, Count(%s) = %d, want %d", expectedCardID, actualCount, expectedCount)
				}
			}

			// Verify no unexpected cards exist
			allCounts := deck.GetAllCardCounts()
			for cardID := range allCounts {
				if _, expected := tt.expectedDeck[cardID]; !expected {
					t.Errorf("Unexpected card %s found in deck after Remove", cardID)
				}
			}
		})
	}
}

func TestCardDeck_Count(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[core.CardID]int
		cardID   core.CardID
		expected int
	}{
		{
			name:     "存在しないカード",
			initial:  map[core.CardID]int{},
			cardID:   "card1",
			expected: 0,
		},
		{
			name:     "1枚のカード",
			initial:  map[core.CardID]int{"card1": 1},
			cardID:   "card1",
			expected: 1,
		},
		{
			name:     "複数枚のカード",
			initial:  map[core.CardID]int{"card1": 5},
			cardID:   "card1",
			expected: 5,
		},
		{
			name:     "複数カード中の1つ",
			initial:  map[core.CardID]int{"card1": 2, "card2": 3, "card3": 1},
			cardID:   "card2",
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deck := core.NewCardDeck()
			// Initialize deck
			for cardID, count := range tt.initial {
				for i := 0; i < count; i++ {
					deck.Add(cardID)
				}
			}

			actual := deck.Count(tt.cardID)
			if actual != tt.expected {
				t.Errorf("Count(%s) = %d, want %d", tt.cardID, actual, tt.expected)
			}
		})
	}
}

func TestCardDeck_GetAllCardIDs(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[core.CardID]int
		expected []core.CardID
	}{
		{
			name:     "空のデッキ",
			initial:  map[core.CardID]int{},
			expected: []core.CardID{},
		},
		{
			name:     "1枚のカード",
			initial:  map[core.CardID]int{"card1": 1},
			expected: []core.CardID{"card1"},
		},
		{
			name:     "複数枚の同じカード",
			initial:  map[core.CardID]int{"card1": 3},
			expected: []core.CardID{"card1", "card1", "card1"},
		},
		{
			name:     "複数の異なるカード",
			initial:  map[core.CardID]int{"card1": 2, "card2": 1},
			expected: []core.CardID{"card1", "card1", "card2"}, // Order may vary
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deck := core.NewCardDeck()
			// Initialize deck
			for cardID, count := range tt.initial {
				for i := 0; i < count; i++ {
					deck.Add(cardID)
				}
			}

			actual := deck.GetAllCardIDs()

			if len(actual) != len(tt.expected) {
				t.Errorf("GetAllCardIDs() length = %d, want %d", len(actual), len(tt.expected))
				return
			}

			// Count occurrences of each card ID
			actualCounts := make(map[core.CardID]int)
			for _, cardID := range actual {
				actualCounts[cardID]++
			}

			expectedCounts := make(map[core.CardID]int)
			for _, cardID := range tt.expected {
				expectedCounts[cardID]++
			}

			for cardID, expectedCount := range expectedCounts {
				if actualCount := actualCounts[cardID]; actualCount != expectedCount {
					t.Errorf("GetAllCardIDs() contains %d of %s, want %d", actualCount, cardID, expectedCount)
				}
			}
		})
	}
}

func TestCardDeck_GetAllCardCounts(t *testing.T) {
	tests := []struct {
		name     string
		initial  map[core.CardID]int
		expected map[core.CardID]int
	}{
		{
			name:     "空のデッキ",
			initial:  map[core.CardID]int{},
			expected: map[core.CardID]int{},
		},
		{
			name:     "1枚のカード",
			initial:  map[core.CardID]int{"card1": 1},
			expected: map[core.CardID]int{"card1": 1},
		},
		{
			name:     "複数の異なるカード",
			initial:  map[core.CardID]int{"card1": 2, "card2": 3, "card3": 1},
			expected: map[core.CardID]int{"card1": 2, "card2": 3, "card3": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deck := core.NewCardDeck()
			// Initialize deck
			for cardID, count := range tt.initial {
				for i := 0; i < count; i++ {
					deck.Add(cardID)
				}
			}

			actual := deck.GetAllCardCounts()

			if len(actual) != len(tt.expected) {
				t.Errorf("GetAllCardCounts() length = %d, want %d", len(actual), len(tt.expected))
				return
			}

			for cardID, expectedCount := range tt.expected {
				if actualCount := actual[cardID]; actualCount != expectedCount {
					t.Errorf("GetAllCardCounts()[%s] = %d, want %d", cardID, actualCount, expectedCount)
				}
			}

			// Verify no unexpected cards
			for cardID := range actual {
				if _, expected := tt.expected[cardID]; !expected {
					t.Errorf("GetAllCardCounts() contains unexpected card %s", cardID)
				}
			}
		})
	}
}

func TestCardDeck_Integration(t *testing.T) {
	t.Run("複数操作の統合テスト", func(t *testing.T) {
		deck := core.NewCardDeck()

		// Add multiple cards
		deck.Add("card1")
		deck.Add("card1")
		deck.Add("card2")
		deck.Add("card3")
		deck.Add("card3")
		deck.Add("card3")

		// Verify counts
		if count := deck.Count("card1"); count != 2 {
			t.Errorf("Count(card1) = %d, want 2", count)
		}
		if count := deck.Count("card2"); count != 1 {
			t.Errorf("Count(card2) = %d, want 1", count)
		}
		if count := deck.Count("card3"); count != 3 {
			t.Errorf("Count(card3) = %d, want 3", count)
		}

		// Remove some cards
		if ok := deck.Remove("card1"); !ok {
			t.Error("Remove(card1) should succeed")
		}
		if ok := deck.Remove("card3"); !ok {
			t.Error("Remove(card3) should succeed")
		}
		if ok := deck.Remove("card3"); !ok {
			t.Error("Remove(card3) should succeed")
		}

		// Verify counts after removal
		if count := deck.Count("card1"); count != 1 {
			t.Errorf("After removal, Count(card1) = %d, want 1", count)
		}
		if count := deck.Count("card2"); count != 1 {
			t.Errorf("After removal, Count(card2) = %d, want 1", count)
		}
		if count := deck.Count("card3"); count != 1 {
			t.Errorf("After removal, Count(card3) = %d, want 1", count)
		}

		// Try to remove non-existent card
		if ok := deck.Remove("card4"); ok {
			t.Error("Remove(card4) should fail for non-existent card")
		}

		// Remove all remaining cards
		deck.Remove("card1")
		deck.Remove("card2")
		deck.Remove("card3")

		// Verify empty deck
		if count := deck.Count("card1"); count != 0 {
			t.Errorf("After removing all, Count(card1) = %d, want 0", count)
		}
		if cardIDs := deck.GetAllCardIDs(); len(cardIDs) != 0 {
			t.Errorf("After removing all, GetAllCardIDs() length = %d, want 0", len(cardIDs))
		}
		if counts := deck.GetAllCardCounts(); len(counts) != 0 {
			t.Errorf("After removing all, GetAllCardCounts() length = %d, want 0", len(counts))
		}
	})
}
