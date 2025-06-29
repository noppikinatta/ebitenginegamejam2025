package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// MockIntner はテスト用の乱数生成器のモック
type MockIntner struct {
	values []int
	count  int
}

func (m *MockIntner) Intn(n int) int {
	if m.count >= len(m.values) {
		return 0 // デフォルト値
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
			name: "基本的なカードパック開封",
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
			name: "同じカードが複数出る場合",
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
			containsAll:  []core.CardID{"card_a", "card_b"}, // この中から出現することを確認
		},
		{
			name: "単一カード",
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

			// 単一カードの場合は厳密にチェック
			if tt.name == "単一カード" {
				if result[0] != tt.containsAll[0] {
					t.Errorf("Open() = %v, want %v", result[0], tt.containsAll[0])
				}
				return
			}

			// 複数カードの場合は、期待されるカードのみが含まれることを確認
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
	// テスト用のカードデータベース
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
			"structure_1": {
				CardID: "structure_1",
			},
		},
	}

	tests := []struct {
		name        string
		cardIDs     []core.CardID
		expectOk    bool
		expectCards func(*core.Cards) bool
	}{
		{
			name:     "存在するカードのみ",
			cardIDs:  []core.CardID{"battle_1", "structure_1"},
			expectOk: true,
			expectCards: func(cards *core.Cards) bool {
				return len(cards.BattleCards) == 1 &&
					len(cards.StructureCards) == 1 &&
					cards.BattleCards[0].CardID == "battle_1" &&
					cards.StructureCards[0].CardID == "structure_1"
			},
		},
		{
			name:     "存在しないカードが含まれる",
			cardIDs:  []core.CardID{"battle_1", "nonexistent"},
			expectOk: false,
			expectCards: func(cards *core.Cards) bool {
				return cards == nil
			},
		},
		{
			name:     "空のリスト",
			cardIDs:  []core.CardID{},
			expectOk: true,
			expectCards: func(cards *core.Cards) bool {
				return len(cards.BattleCards) == 0 &&
					len(cards.StructureCards) == 0
			},
		},
		{
			name:     "同じタイプのカード複数",
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

// テスト用のモック実装
type mockBattleCardPowerModifier struct {
	canAffectFunc func(*core.BattleCard) bool
	modifyFunc    func(*core.BattleCard) float64
}

func (m *mockBattleCardPowerModifier) CanAffect(battleCard *core.BattleCard) bool {
	return m.canAffectFunc(battleCard)
}

func (m *mockBattleCardPowerModifier) Modify(battleCard *core.BattleCard) float64 {
	return m.modifyFunc(battleCard)
}
