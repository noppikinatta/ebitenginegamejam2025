package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestMarketItem_CanPurchase(t *testing.T) {
	// Card pack for testing
	cardPack := &core.CardPack{
		CardPackID: "test_pack",
		Ratios: map[core.CardID]int{
			"card_a": 50,
			"card_b": 50,
		},
		NumPerOpen: 1,
	}

	tests := []struct {
		name     string
		item     *core.MarketItem
		treasury *core.Treasury
		expected bool
	}{
		{
			name: "Sufficient resources",
			item: core.NewMarketItem(cardPack, core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			}, 1.0, 0.0),
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{
					Money: 150, Food: 80, Wood: 50, Iron: 30, Mana: 20,
				},
			},
			expected: true,
		},
		{
			name: "Exactly the same amount of resources",
			item: core.NewMarketItem(cardPack, core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			}, 1.0, 0.0),
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{
					Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
				},
			},
			expected: true,
		},
		{
			name: "Insufficient resources",
			item: core.NewMarketItem(cardPack, core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			}, 1.0, 0.0),
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{
					Money: 50, Food: 30, Wood: 20, Iron: 10, Mana: 5,
				},
			},
			expected: false,
		},
		{
			name: "Insufficient partial resources",
			item: core.NewMarketItem(cardPack, core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			}, 1.0, 0.0),
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{
					Money: 150, Food: 30, Wood: 50, Iron: 30, Mana: 20, // Insufficient Food
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.CanPurchase(tt.treasury)
			if result != tt.expected {
				t.Errorf("CanPurchase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMarket_VisibleCardPacks(t *testing.T) {
	// Card pack for testing
	pack1 := &core.CardPack{
		CardPackID: "basic_pack",
		Ratios: map[core.CardID]int{
			"card_a": 70,
			"card_b": 30,
		},
		NumPerOpen: 2,
	}

	pack2 := &core.CardPack{
		CardPackID: "advanced_pack",
		Ratios: map[core.CardID]int{
			"card_c": 50,
			"card_d": 50,
		},
		NumPerOpen: 3,
	}

	pack3 := &core.CardPack{
		CardPackID: "premium_pack",
		Ratios: map[core.CardID]int{
			"card_e": 40,
			"card_f": 60,
		},
		NumPerOpen: 1,
	}

	// Market items for testing
	items := []*core.MarketItem{
		core.NewMarketItem(pack1, core.ResourceQuantity{Money: 50}, 1.0, 0.0),
		core.NewMarketItem(pack2, core.ResourceQuantity{Money: 100}, 2.0, 0.0),
		core.NewMarketItem(pack3, core.ResourceQuantity{Money: 200}, 3.0, 0.0),
	}

	tests := []struct {
		name          string
		marketLevel   core.MarketLevel
		expectedCount int
		expectedPacks []string // List of CardPackIDs
	}{
		{
			name:          "Level 1.0 - Only basic pack is visible",
			marketLevel:   1.0,
			expectedCount: 1,
			expectedPacks: []string{"basic_pack"},
		},
		{
			name:          "Level 2.0 - Basic and advanced packs are visible",
			marketLevel:   2.0,
			expectedCount: 2,
			expectedPacks: []string{"basic_pack", "advanced_pack"},
		},
		{
			name:          "Level 3.0 - All packs are visible",
			marketLevel:   3.0,
			expectedCount: 3,
			expectedPacks: []string{"basic_pack", "advanced_pack", "premium_pack"},
		},
		{
			name:          "Level 0.5 - Nothing is visible",
			marketLevel:   0.5,
			expectedCount: 0,
			expectedPacks: []string{},
		},
		{
			name:          "Level 1.5 - Only basic pack is visible",
			marketLevel:   1.5,
			expectedCount: 1,
			expectedPacks: []string{"basic_pack"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			market := &core.Market{
				Level: tt.marketLevel,
				Items: items,
			}

			result := market.VisibleMarketItems()
			if len(result) != tt.expectedCount {
				t.Errorf("VisibleMarketItems() returned %d packs, want %d", len(result), tt.expectedCount)
				return
			}

			// Check if the expected CardPackIDs are included
			for i, item := range result {
				if i >= len(tt.expectedPacks) {
					t.Errorf("Unexpected pack at index %d: %v", i, item.CardPack().CardPackID)
					continue
				}
				if string(item.CardPack().CardPackID) != tt.expectedPacks[i] {
					t.Errorf("VisibleCardPacks()[%d] = %v, want %v", i, item.CardPack().CardPackID, tt.expectedPacks[i])
				}
			}
		})
	}
}

func TestMarket_CanPurchase(t *testing.T) {
	// Card pack for testing
	pack := &core.CardPack{
		CardPackID: "test_pack",
		Ratios: map[core.CardID]int{
			"card_a": 100,
		},
		NumPerOpen: 1,
	}

	// Market items for testing
	items := []*core.MarketItem{
		core.NewMarketItem(pack, core.ResourceQuantity{Money: 100}, 1.0, 0.0),
		core.NewMarketItem(pack, core.ResourceQuantity{Money: 200}, 2.0, 0.0),
	}

	market := &core.Market{
		Level: 2.0,
		Items: items,
	}

	tests := []struct {
		name     string
		index    int
		treasury *core.Treasury
		expected bool
	}{
		{
			name:  "Valid index, sufficient resources",
			index: 0,
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{Money: 150},
			},
			expected: true,
		},
		{
			name:  "Valid index, insufficient resources",
			index: 1,
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{Money: 100}, // Requires 200, but only have 100
			},
			expected: false,
		},
		{
			name:  "Invalid index (negative)",
			index: -1,
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{Money: 1000},
			},
			expected: false,
		},
		{
			name:  "Invalid index (out of range)",
			index: 5,
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{Money: 1000},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := market.CanPurchase(tt.index, tt.treasury)
			if result != tt.expected {
				t.Errorf("CanPurchase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMarket_Purchase(t *testing.T) {
	// Card pack for testing
	pack := &core.CardPack{
		CardPackID: "purchase_test_pack",
		Ratios: map[core.CardID]int{
			"card_a": 100,
		},
		NumPerOpen: 1,
	}

	// Market items for testing
	items := []*core.MarketItem{
		core.NewMarketItem(pack, core.ResourceQuantity{Money: 100, Food: 50}, 1.0, 0.0),
	}

	market := &core.Market{
		Level: 1.0,
		Items: items,
	}

	tests := []struct {
		name            string
		index           int
		initialTreasury core.ResourceQuantity
		expectedOk      bool
		finalTreasury   core.ResourceQuantity
	}{
		{
			name:  "Normal purchase",
			index: 0,
			initialTreasury: core.ResourceQuantity{
				Money: 200, Food: 100, Wood: 50, Iron: 30, Mana: 20,
			},
			expectedOk: true,
			finalTreasury: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 50, Iron: 30, Mana: 20,
			},
		},
		{
			name:  "Purchase failed due to insufficient resources",
			index: 0,
			initialTreasury: core.ResourceQuantity{
				Money: 50, Food: 25, Wood: 10, Iron: 5, Mana: 2,
			},
			expectedOk: false,
			finalTreasury: core.ResourceQuantity{
				Money: 50, Food: 25, Wood: 10, Iron: 5, Mana: 2,
			},
		},
		{
			name:  "Invalid index",
			index: 1,
			initialTreasury: core.ResourceQuantity{
				Money: 1000, Food: 1000, Wood: 1000, Iron: 1000, Mana: 1000,
			},
			expectedOk: false,
			finalTreasury: core.ResourceQuantity{
				Money: 1000, Food: 1000, Wood: 1000, Iron: 1000, Mana: 1000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			treasury := &core.Treasury{
				Resources: tt.initialTreasury,
			}

			cardPack, ok := market.Purchase(tt.index, treasury)

			if ok != tt.expectedOk {
				t.Errorf("Purchase() ok = %v, want %v", ok, tt.expectedOk)
			}

			if tt.expectedOk && cardPack == nil {
				t.Errorf("Purchase() returned nil cardPack when expecting success")
			}

			if !tt.expectedOk && cardPack != nil {
				t.Errorf("Purchase() returned cardPack when expecting failure")
			}

			if treasury.Resources != tt.finalTreasury {
				t.Errorf("Treasury after purchase = %v, want %v", treasury.Resources, tt.finalTreasury)
			}
		})
	}
}
