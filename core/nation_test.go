package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestTreasury_Add(t *testing.T) {
	tests := []struct {
		name     string
		initial  core.ResourceQuantity
		toAdd    core.ResourceQuantity
		expected core.ResourceQuantity
	}{
		{
			name: "正常な加算",
			initial: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			toAdd: core.ResourceQuantity{
				Money: 50, Food: 25, Wood: 15, Iron: 10, Mana: 5,
			},
			expected: core.ResourceQuantity{
				Money: 150, Food: 75, Wood: 45, Iron: 30, Mana: 15,
			},
		},
		{
			name: "ゼロリソースの加算",
			initial: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			toAdd: core.ResourceQuantity{},
			expected: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			treasury := &core.Treasury{
				Resources: tt.initial,
			}

			treasury.Add(tt.toAdd)

			if treasury.Resources != tt.expected {
				t.Errorf("Add() result = %v, want %v", treasury.Resources, tt.expected)
			}
		})
	}
}

func TestTreasury_Sub(t *testing.T) {
	tests := []struct {
		name       string
		initial    core.ResourceQuantity
		toSub      core.ResourceQuantity
		expected   core.ResourceQuantity
		expectedOk bool
	}{
		{
			name: "正常な減算",
			initial: core.ResourceQuantity{
				Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10,
			},
			toSub: core.ResourceQuantity{
				Money: 30, Food: 20, Wood: 10, Iron: 5, Mana: 3,
			},
			expected: core.ResourceQuantity{
				Money: 70, Food: 30, Wood: 20, Iron: 15, Mana: 7,
			},
			expectedOk: true,
		},
		{
			name: "資源不足で減算失敗",
			initial: core.ResourceQuantity{
				Money: 50, Food: 30, Wood: 20, Iron: 10, Mana: 5,
			},
			toSub: core.ResourceQuantity{
				Money: 100, Food: 20, Wood: 10, Iron: 5, Mana: 2, // Moneyが不足
			},
			expected: core.ResourceQuantity{
				Money: 50, Food: 30, Wood: 20, Iron: 10, Mana: 5, // 変化なし
			},
			expectedOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			treasury := &core.Treasury{
				Resources: tt.initial,
			}

			ok := treasury.Sub(tt.toSub)

			if ok != tt.expectedOk {
				t.Errorf("Sub() ok = %v, want %v", ok, tt.expectedOk)
			}
			if treasury.Resources != tt.expected {
				t.Errorf("Sub() result = %v, want %v", treasury.Resources, tt.expected)
			}
		})
	}
}

func TestNation_VisibleCardPacks(t *testing.T) {
	// テスト用のカードパック
	pack1 := &core.CardPack{
		CardPackID: "basic_pack",
		Ratios: map[core.CardID]int{
			"card_a": 100,
		},
		NumPerOpen: 1,
	}

	pack2 := &core.CardPack{
		CardPackID: "advanced_pack",
		Ratios: map[core.CardID]int{
			"card_b": 100,
		},
		NumPerOpen: 1,
	}

	// テスト用のマーケットアイテム
	items := []*core.MarketItem{
		{
			CardPack:      pack1,
			Price:         core.ResourceQuantity{Money: 50},
			RequiredLevel: 1.0,
		},
		{
			CardPack:      pack2,
			Price:         core.ResourceQuantity{Money: 100},
			RequiredLevel: 2.0,
		},
	}

	market := &core.Market{
		Level: 1.5,
		Items: items,
	}

	nation := &core.Nation{
		NationID: "test_nation",
		Market:   market,
	}

	visiblePacks := nation.VisibleCardPacks()

	// レベル1.5なので、RequiredLevel 1.0の基本パックのみ見える
	if len(visiblePacks) != 1 {
		t.Errorf("VisibleCardPacks() returned %d packs, want 1", len(visiblePacks))
		return
	}

	if visiblePacks[0].CardPackID != "basic_pack" {
		t.Errorf("VisibleCardPacks()[0] = %v, want basic_pack", visiblePacks[0].CardPackID)
	}
}

func TestNation_CanPurchase(t *testing.T) {
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
			CardPack:      pack,
			Price:         core.ResourceQuantity{Money: 100},
			RequiredLevel: 1.0,
		},
	}

	market := &core.Market{
		Level: 1.0,
		Items: items,
	}

	nation := &core.Nation{
		NationID: "test_nation",
		Market:   market,
	}

	tests := []struct {
		name     string
		index    int
		treasury *core.Treasury
		expected bool
	}{
		{
			name:  "購入可能",
			index: 0,
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{Money: 150},
			},
			expected: true,
		},
		{
			name:  "資源不足",
			index: 0,
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{Money: 50},
			},
			expected: false,
		},
		{
			name:  "無効なインデックス",
			index: 5,
			treasury: &core.Treasury{
				Resources: core.ResourceQuantity{Money: 1000},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nation.CanPurchase(tt.index, tt.treasury)
			if result != tt.expected {
				t.Errorf("CanPurchase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNation_Purchase(t *testing.T) {
	// テスト用のカードパック
	pack := &core.CardPack{
		CardPackID: "purchase_test_pack",
		Ratios: map[core.CardID]int{
			"card_a": 100,
		},
		NumPerOpen: 1,
	}

	// テスト用のマーケットアイテム
	items := []*core.MarketItem{
		{
			CardPack:      pack,
			Price:         core.ResourceQuantity{Money: 100, Food: 50},
			RequiredLevel: 1.0,
		},
	}

	market := &core.Market{
		Level: 1.0,
		Items: items,
	}

	nation := &core.Nation{
		NationID: "test_nation",
		Market:   market,
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
				Money: 100, Food: 50, Wood: 50, Iron: 30, Mana: 20,
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
				Money: 50, Food: 25, Wood: 10, Iron: 5, Mana: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			treasury := &core.Treasury{
				Resources: tt.initialTreasury,
			}

			cardPack, ok := nation.Purchase(tt.index, treasury)

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

func TestMyNation_AppendMarketItem(t *testing.T) {
	// テスト用のカードパック
	pack := &core.CardPack{
		CardPackID: "new_pack",
		Ratios: map[core.CardID]int{
			"card_new": 100,
		},
		NumPerOpen: 1,
	}

	newItem := &core.MarketItem{
		CardPack:      pack,
		Price:         core.ResourceQuantity{Money: 150},
		RequiredLevel: 2.0,
	}

	market := &core.Market{
		Level: 2.0,
		Items: []*core.MarketItem{},
	}

	myNation := &core.MyNation{
		Nation: core.Nation{
			NationID: "my_nation",
			Market:   market,
		},
		BasicYield: core.ResourceQuantity{Money: 10, Food: 5},
	}

	// 初期状態の確認
	if len(myNation.Market.Items) != 0 {
		t.Errorf("Initial market items = %d, want 0", len(myNation.Market.Items))
	}

	// アイテム追加
	myNation.AppendMarketItem(newItem)

	// 追加後の確認
	if len(myNation.Market.Items) != 1 {
		t.Errorf("Market items after append = %d, want 1", len(myNation.Market.Items))
		return
	}

	if myNation.Market.Items[0] != newItem {
		t.Errorf("Appended item = %v, want %v", myNation.Market.Items[0], newItem)
	}
}

func TestMyNation_AppendLevel(t *testing.T) {
	market := &core.Market{
		Level: 1.0,
		Items: []*core.MarketItem{},
	}

	myNation := &core.MyNation{
		Nation: core.Nation{
			NationID: "my_nation",
			Market:   market,
		},
		BasicYield: core.ResourceQuantity{Money: 10, Food: 5},
	}

	// 初期レベルの確認
	if myNation.Market.Level != 1.0 {
		t.Errorf("Initial market level = %v, want 1.0", myNation.Market.Level)
	}

	// レベル追加
	myNation.AppendLevel(0.5)

	// 追加後の確認
	if myNation.Market.Level != 1.5 {
		t.Errorf("Market level after append = %v, want 1.5", myNation.Market.Level)
	}

	// さらに追加
	myNation.AppendLevel(1.0)

	if myNation.Market.Level != 2.5 {
		t.Errorf("Market level after second append = %v, want 2.5", myNation.Market.Level)
	}
}

func TestOtherNation_Purchase(t *testing.T) {
	// テスト用のカードパック
	pack := &core.CardPack{
		CardPackID: "other_pack",
		Ratios: map[core.CardID]int{
			"card_other": 100,
		},
		NumPerOpen: 1,
	}

	// テスト用のマーケットアイテム
	items := []*core.MarketItem{
		{
			CardPack:      pack,
			Price:         core.ResourceQuantity{Money: 100},
			RequiredLevel: 1.0,
		},
	}

	market := &core.Market{
		Level: 1.0,
		Items: items,
	}

	otherNation := &core.OtherNation{
		Nation: core.Nation{
			NationID: "other_nation",
			Market:   market,
		},
	}

	treasury := &core.Treasury{
		Resources: core.ResourceQuantity{Money: 150},
	}

	// 購入前のマーケットレベル
	initialLevel := otherNation.Market.Level

	// 購入実行
	cardPack, ok := otherNation.Purchase(0, treasury)

	// 購入が成功することを確認
	if !ok {
		t.Errorf("Purchase() ok = %v, want true", ok)
	}

	if cardPack == nil {
		t.Errorf("Purchase() returned nil cardPack")
	}

	// 国庫が正しく減っていることを確認
	expectedTreasury := core.ResourceQuantity{Money: 50}
	if treasury.Resources != expectedTreasury {
		t.Errorf("Treasury after purchase = %v, want %v", treasury.Resources, expectedTreasury)
	}

	// OtherNationの特殊機能：Market.Levelが0.2加算されることを確認
	expectedLevel := initialLevel + 0.2
	if otherNation.Market.Level != expectedLevel {
		t.Errorf("Market level after purchase = %v, want %v", otherNation.Market.Level, expectedLevel)
	}
}
