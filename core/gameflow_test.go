package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestGameState_AddYield(t *testing.T) {
	// テスト用のTerritories
	territory1 := &core.Territory{
		TerritoryID: "territory_1",
		Cards:       []*core.StructureCard{},
		CardSlot:    2,
		BaseYield:   core.ResourceQuantity{Money: 10, Food: 5},
	}

	territory2 := &core.Territory{
		TerritoryID: "territory_2",
		Cards:       []*core.StructureCard{},
		CardSlot:    3,
		BaseYield:   core.ResourceQuantity{Money: 8, Food: 4, Wood: 3},
	}

	// テスト用のWildernessPoint（制圧済み）
	wilderness1 := &core.WildernessPoint{
		Controlled: true,
		Territory:  territory1,
	}

	wilderness2 := &core.WildernessPoint{
		Controlled: true,
		Territory:  territory2,
	}

	// テスト用のWildernessPoint（未制圧）
	uncontrolledWilderness := &core.WildernessPoint{
		Controlled: false,
		Territory: &core.Territory{
			TerritoryID: "uncontrolled",
			BaseYield:   core.ResourceQuantity{Money: 100}, // 制圧されていないので収入なし
		},
	}

	// テスト用のMyNation
	myNation := &core.MyNation{
		BaseNation: core.BaseNation{
			NationID: "player",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
		BasicYield: core.ResourceQuantity{Money: 5, Food: 2, Mana: 1},
	}

	// テスト用のMapGrid
	points := []core.Point{
		&core.MyNationPoint{MyNation: myNation},
		wilderness1,
		wilderness2,
		uncontrolledWilderness,
	}

	mapGrid := &core.MapGrid{
		SizeX:  2,
		SizeY:  2,
		Points: points,
	}

	gameState := &core.GameState{
		MyNation: myNation,
		MapGrid:  mapGrid,
		Treasury: &core.Treasury{
			Resources: core.ResourceQuantity{Money: 100, Food: 50, Wood: 30, Iron: 20, Mana: 10},
		},
		CurrentTurn: 1,
	}

	// 初期状態の確認
	initialTreasury := gameState.Treasury.Resources

	// Yield加算実行
	gameState.AddYield()

	// 期待される結果：BasicYield + 制圧済みTerritoryのYield
	expectedYield := myNation.BasicYield.
		Add(territory1.Yield()).
		Add(territory2.Yield())

	expectedTreasury := initialTreasury.Add(expectedYield)

	if gameState.Treasury.Resources != expectedTreasury {
		t.Errorf("AddYield() treasury = %v, want %v", gameState.Treasury.Resources, expectedTreasury)
	}
}

func TestGameState_NextTurn(t *testing.T) {
	myNation := &core.MyNation{
		BaseNation: core.BaseNation{
			NationID: "player",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
		BasicYield: core.ResourceQuantity{Money: 5},
	}

	gameState := &core.GameState{
		MyNation: myNation,
		MapGrid:  &core.MapGrid{SizeX: 1, SizeY: 1, Points: []core.Point{&core.MyNationPoint{MyNation: myNation}}},
		Treasury: &core.Treasury{
			Resources: core.ResourceQuantity{Money: 100},
		},
		CurrentTurn: 5,
	}

	initialTurn := gameState.CurrentTurn
	initialTreasury := gameState.Treasury.Resources

	// ターン進行
	gameState.NextTurn()

	// ターンが進んでいることを確認
	if gameState.CurrentTurn != initialTurn+1 {
		t.Errorf("NextTurn() CurrentTurn = %v, want %v", gameState.CurrentTurn, initialTurn+1)
	}

	// Yieldが加算されていることを確認
	expectedTreasury := initialTreasury.Add(myNation.BasicYield)
	if gameState.Treasury.Resources != expectedTreasury {
		t.Errorf("NextTurn() treasury = %v, want %v", gameState.Treasury.Resources, expectedTreasury)
	}
}

func TestGameState_IsVictory(t *testing.T) {
	// テスト用のBoss
	boss := &core.Enemy{
		EnemyID:        "final_boss",
		EnemyType:      "dragon",
		Power:          100.0,
		BattleCardSlot: 4,
		Skills:         []core.EnemySkill{},
	}

	tests := []struct {
		name     string
		defeated bool
		expected bool
	}{
		{
			name:     "ボスが撃破されている場合は勝利",
			defeated: true,
			expected: true,
		},
		{
			name:     "ボスが撃破されていない場合は勝利ではない",
			defeated: false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bossPoint := &core.BossPoint{
				Boss:     boss,
				Defeated: tt.defeated,
			}

			myNation := &core.MyNation{
				BaseNation: core.BaseNation{
					NationID: "player",
					Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
				},
				BasicYield: core.ResourceQuantity{Money: 5},
			}

			mapGrid := &core.MapGrid{
				SizeX:  1,
				SizeY:  1,
				Points: []core.Point{bossPoint},
			}

			gameState := &core.GameState{
				MyNation:    myNation,
				MapGrid:     mapGrid,
				Treasury:    &core.Treasury{Resources: core.ResourceQuantity{}},
				CurrentTurn: 1,
			}

			result := gameState.IsVictory()
			if result != tt.expected {
				t.Errorf("IsVictory() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGameState_CanInteract(t *testing.T) {
	myNation := &core.MyNation{
		BaseNation: core.BaseNation{
			NationID: "player",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
		BasicYield: core.ResourceQuantity{Money: 5},
	}

	// 制圧済みWilderness
	controlledWilderness := &core.WildernessPoint{
		Controlled: true,
		Territory: &core.Territory{
			TerritoryID: "controlled",
			BaseYield:   core.ResourceQuantity{Money: 5},
		},
	}

	// 2x2のマップグリッド
	/*
		配置:
		(0,0) MyNation          (1,0) ControlledWilderness
		(0,1) nil               (1,1) nil
	*/
	points := []core.Point{
		&core.MyNationPoint{MyNation: myNation},
		controlledWilderness,
		&core.OtherNationPoint{OtherNation: &core.OtherNation{}},
		&core.BossPoint{Boss: &core.Enemy{}},
	}

	mapGrid := &core.MapGrid{
		SizeX:  2,
		SizeY:  2,
		Points: points,
	}

	gameState := &core.GameState{
		MyNation:    myNation,
		MapGrid:     mapGrid,
		Treasury:    &core.Treasury{Resources: core.ResourceQuantity{}},
		CurrentTurn: 1,
	}

	tests := []struct {
		name     string
		x, y     int
		expected bool
	}{
		{
			name:     "MyNationPoint(0,0)は操作可能",
			x:        0,
			y:        0,
			expected: true,
		},
		{
			name:     "隣接する制圧済みWilderness(1,0)は操作可能",
			x:        1,
			y:        0,
			expected: true,
		},
		{
			name:     "何もないポイント(0,1)は操作不可",
			x:        0,
			y:        1,
			expected: false,
		},
		{
			name:     "範囲外の座標は操作不可",
			x:        5,
			y:        5,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gameState.CanInteract(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("CanInteract(%d, %d) = %v, want %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

func TestGameState_GetPoint(t *testing.T) {
	myNation := &core.MyNation{
		BaseNation: core.BaseNation{
			NationID: "player",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
		BasicYield: core.ResourceQuantity{Money: 5},
	}

	myNationPoint := &core.MyNationPoint{MyNation: myNation}

	// 1x1のマップグリッド
	points := []core.Point{myNationPoint}

	mapGrid := &core.MapGrid{
		SizeX:  1,
		SizeY:  1,
		Points: points,
	}

	gameState := &core.GameState{
		MyNation:    myNation,
		MapGrid:     mapGrid,
		Treasury:    &core.Treasury{Resources: core.ResourceQuantity{}},
		CurrentTurn: 1,
	}

	tests := []struct {
		name        string
		x, y        int
		expectedNil bool
	}{
		{
			name:        "有効な座標(0,0)",
			x:           0,
			y:           0,
			expectedNil: false,
		},
		{
			name:        "無効な座標(1,0)",
			x:           1,
			y:           0,
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := gameState.GetPoint(tt.x, tt.y)

			if tt.expectedNil {
				if point != nil {
					t.Errorf("GetPoint(%d, %d) = %v, want nil", tt.x, tt.y, point)
				}
			} else {
				if point == nil {
					t.Errorf("GetPoint(%d, %d) = nil, want non-nil", tt.x, tt.y)
				} else if point != myNationPoint {
					t.Errorf("GetPoint(%d, %d) = %v, want %v", tt.x, tt.y, point, myNationPoint)
				}
			}
		})
	}
}
