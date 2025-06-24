package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestMarketItem_CanPurchase(t *testing.T) {
	// テスト用のカードパック
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
			name: "十分な資源がある場合",
			item: &core.MarketItem{
				CardPack: cardPack,
				Price: core.ResourceQuantity{
					Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
				},
				RequiredLevel: 1.0,
			},
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{
					Money: 150, Food: 80, Wood: 50, Iron: 30, Mana: 20,
				},
			},
			expected: true,
		},
		{
			name: "ちょうど同じ資源量の場合",
			item: &core.MarketItem{
				CardPack: cardPack,
				Price: core.ResourceQuantity{
					Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
				},
				RequiredLevel: 1.0,
			},
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{
					Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
				},
			},
			expected: true,
		},
		{
			name: "資源が不足している場合",
			item: &core.MarketItem{
				CardPack: cardPack,
				Price: core.ResourceQuantity{
					Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
				},
				RequiredLevel: 1.0,
			},
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{
					Money: 50, Food: 30, Wood: 20, Iron: 10, Mana: 5,
				},
			},
			expected: false,
		},
		{
			name: "一部リソースが不足している場合",
			item: &core.MarketItem{
				CardPack: cardPack,
				Price: core.ResourceQuantity{
					Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
				},
				RequiredLevel: 1.0,
			},
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{
					Money: 150, Food: 30, Wood: 50, Iron: 30, Mana: 20, // Foodが不足
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
	// テスト用のカードパック
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

	// テスト用のマーケットアイテム
	items := []*core.MarketItem{
		{
			CardPack: pack1,
			Price: core.ResourceQuantity{
				Money: 50,
			},
			RequiredLevel: 1.0,
		},
		{
			CardPack: pack2,
			Price: core.ResourceQuantity{
				Money: 100,
			},
			RequiredLevel: 2.0,
		},
		{
			CardPack: pack3,
			Price: core.ResourceQuantity{
				Money: 200,
			},
			RequiredLevel: 3.0,
		},
	}

	tests := []struct {
		name          string
		marketLevel   core.MarketLevel
		expectedCount int
		expectedPacks []string // CardPackIDのリスト
	}{
		{
			name:          "レベル1.0 - 基本パックのみ見える",
			marketLevel:   1.0,
			expectedCount: 1,
			expectedPacks: []string{"basic_pack"},
		},
		{
			name:          "レベル2.0 - 基本パック＋上級パック",
			marketLevel:   2.0,
			expectedCount: 2,
			expectedPacks: []string{"basic_pack", "advanced_pack"},
		},
		{
			name:          "レベル3.0 - 全パック見える",
			marketLevel:   3.0,
			expectedCount: 3,
			expectedPacks: []string{"basic_pack", "advanced_pack", "premium_pack"},
		},
		{
			name:          "レベル0.5 - 何も見えない",
			marketLevel:   0.5,
			expectedCount: 0,
			expectedPacks: []string{},
		},
		{
			name:          "レベル1.5 - 基本パックのみ",
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

			// 期待されるCardPackIDが含まれているかチェック
			for i, item := range result {
				if i >= len(tt.expectedPacks) {
					t.Errorf("Unexpected pack at index %d: %v", i, item.CardPack.CardPackID)
					continue
				}
				if string(item.CardPack.CardPackID) != tt.expectedPacks[i] {
					t.Errorf("VisibleCardPacks()[%d] = %v, want %v", i, item.CardPack.CardPackID, tt.expectedPacks[i])
				}
			}
		})
	}
}

func TestMarket_CanPurchase(t *testing.T) {
	// テスト用のカードパック
	pack := &core.CardPack{
		CardPackID: "test_pack",
		Ratios: map[core.CardID]int{
			"card_a": 100,
		},
		NumPerOpen: 1,
	}

	// テスト用のマーケットアイテム
	items := []*core.MarketItem{
		{
			CardPack: pack,
			Price: core.ResourceQuantity{
				Money: 100,
			},
			RequiredLevel: 1.0,
		},
		{
			CardPack: pack,
			Price: core.ResourceQuantity{
				Money: 200,
			},
			RequiredLevel: 2.0,
		},
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
			name:  "有効なインデックス、十分な資源",
			index: 0,
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{
					Money: 150,
				},
			},
			expected: true,
		},
		{
			name:  "有効なインデックス、資源不足",
			index: 1,
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{
					Money: 100, // 200必要だが100しかない
				},
			},
			expected: false,
		},
		{
			name:  "無効なインデックス（負の値）",
			index: -1,
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{
					Money: 1000,
				},
			},
			expected: false,
		},
		{
			name:  "無効なインデックス（範囲外）",
			index: 5,
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{
					Money: 1000,
				},
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
	// テスト用のカードパック
	pack := &core.CardPack{
		CardPackID: "purchase_pack",
		Ratios: map[core.CardID]int{
			"card_a": 100,
		},
		NumPerOpen: 1,
	}

	// テスト用のマーケットアイテム
	items := []*core.MarketItem{
		{
			CardPack: pack,
			Price: core.ResourceQuantity{
				Money: 100, Food: 50,
			},
			RequiredLevel: 1.0,
		},
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
			name:  "正常な購入",
			index: 0,
			initialTreasury: core.ResourceQuantity{
				Money: 200, Food: 100, Wood: 50, Iron: 30, Mana: 20,
			},
			expectedOk: true,
			finalTreasury: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 50, Iron: 30, Mana: 20, // Money: 200-100, Food: 100-50
			},
		},
		{
			name:  "資源不足で購入失敗",
			index: 0,
			initialTreasury: core.ResourceQuantity{
				Money: 50, Food: 25, Wood: 10, Iron: 5, Mana: 2,
			},
			expectedOk: false,
			finalTreasury: core.ResourceQuantity{
				Money: 50, Food: 25, Wood: 10, Iron: 5, Mana: 2, // 変化なし
			},
		},
		{
			name:  "無効なインデックス",
			index: 5,
			initialTreasury: core.ResourceQuantity{
				Money: 1000, Food: 1000, Wood: 1000, Iron: 1000, Mana: 1000,
			},
			expectedOk: false,
			finalTreasury: core.ResourceQuantity{
				Money: 1000, Food: 1000, Wood: 1000, Iron: 1000, Mana: 1000, // 変化なし
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
