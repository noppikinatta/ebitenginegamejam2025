package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestMapGrid_CanInteract(t *testing.T) {
	// テスト用のMyNation
	myNation := &core.MyNation{
		Nation: core.Nation{
			NationID: "player",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
		BasicYield: core.ResourceQuantity{Money: 10},
	}

	// テスト用のOtherNation
	otherNation := &core.OtherNation{
		Nation: core.Nation{
			NationID: "ally",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
	}

	// テスト用のEnemy
	enemy := &core.Enemy{
		EnemyID:        "test_orc",
		EnemyType:      "orc",
		Power:          15.0,
		BattleCardSlot: 2,
		Skills:         []core.EnemySkill{},
	}

	// テスト用のBoss
	boss := &core.Enemy{
		EnemyID:        "dragon_boss",
		EnemyType:      "dragon",
		Power:          100.0,
		BattleCardSlot: 4,
		Skills:         []core.EnemySkill{},
	}

	// テスト用のWildernessPoint（制圧済み）
	controlledWilderness := &core.WildernessPoint{
		Controlled: true,
		Enemy:      enemy,
		Territory: &core.Territory{
			TerritoryID: "controlled_territory",
			Cards:       []*core.StructureCard{},
			CardSlot:    2,
			BaseYield:   core.ResourceQuantity{Money: 5},
		},
	}

	// テスト用のWildernessPoint（未制圧）
	uncontrolledWilderness := &core.WildernessPoint{
		Controlled: false,
		Enemy:      enemy,
		Territory: &core.Territory{
			TerritoryID: "uncontrolled_territory",
			Cards:       []*core.StructureCard{},
			CardSlot:    2,
			BaseYield:   core.ResourceQuantity{Money: 5},
		},
	}

	// 3x3のマップグリッドを作成
	/*
		配置:
		(0,0) MyNation    (1,0) Controlled  (2,0) Uncontrolled
		(0,1) OtherNation (1,1) Controlled  (2,1) Boss
		(0,2) Uncontrolled(1,2) Controlled  (2,2) Uncontrolled
	*/
	points := []core.Point{
		// Row 0
		&core.MyNationPoint{MyNation: myNation},
		controlledWilderness,
		uncontrolledWilderness,
		// Row 1
		&core.OtherNationPoint{OtherNation: otherNation},
		controlledWilderness,
		&core.BossPoint{Boss: boss},
		// Row 2
		uncontrolledWilderness,
		controlledWilderness,
		uncontrolledWilderness,
	}

	mapGrid := &core.MapGrid{
		SizeX:  3,
		SizeY:  3,
		Points: points,
	}

	tests := []struct {
		name     string
		x, y     int
		expected bool
		reason   string
	}{
		{
			name:     "MyNationPoint(0,0)は常に操作可能",
			x:        0,
			y:        0,
			expected: true,
			reason:   "MyNationPointは常に操作可能",
		},
		{
			name:     "MyNationPointに隣接する制圧済みWilderness(1,0)は操作可能",
			x:        1,
			y:        0,
			expected: true,
			reason:   "MyNationPointに隣接し、制圧済み",
		},
		{
			name:     "MyNationPointに隣接するOtherNation(0,1)は操作可能",
			x:        0,
			y:        1,
			expected: true,
			reason:   "MyNationPointに隣接",
		},
		{
			name:     "MyNationPointに隣接する未制圧Wilderness(2,0)は操作不可",
			x:        2,
			y:        0,
			expected: false,
			reason:   "未制圧のWildernessを通る必要がある",
		},
		{
			name:     "制圧済みWilderness(1,0)に隣接する制圧済みWilderness(1,1)は操作可能",
			x:        1,
			y:        1,
			expected: true,
			reason:   "制圧済みのルートでアクセス可能",
		},
		{
			name:     "制圧済みWilderness(1,1)に隣接するBoss(2,1)は操作可能",
			x:        2,
			y:        1,
			expected: true,
			reason:   "制圧済みのルートでアクセス可能",
		},
		{
			name:     "遠く離れた未制圧Wilderness(2,2)は操作不可",
			x:        2,
			y:        2,
			expected: false,
			reason:   "制圧済みのルートでアクセス不可",
		},
		{
			name:     "制圧済みWilderness(1,2)は操作可能",
			x:        1,
			y:        2,
			expected: true,
			reason:   "制圧済みのルートでアクセス可能",
		},
		{
			name:     "範囲外の座標(-1,0)は操作不可",
			x:        -1,
			y:        0,
			expected: false,
			reason:   "マップ範囲外",
		},
		{
			name:     "範囲外の座標(3,0)は操作不可",
			x:        3,
			y:        0,
			expected: false,
			reason:   "マップ範囲外",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapGrid.CanInteract(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("CanInteract(%d, %d) = %v, want %v (%s)", tt.x, tt.y, result, tt.expected, tt.reason)
			}
		})
	}
}

func TestMyNationPoint_Location(t *testing.T) {
	myNation := &core.MyNation{
		Nation: core.Nation{
			NationID: "player",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
		BasicYield: core.ResourceQuantity{Money: 10},
	}

	point := &core.MyNationPoint{MyNation: myNation}

	// MyNationPointはLocationメソッドを持たないため、基本的な構造確認のみ
	if point.MyNation != myNation {
		t.Errorf("MyNationPoint.MyNation = %v, want %v", point.MyNation, myNation)
	}
}

func TestOtherNationPoint_Location(t *testing.T) {
	otherNation := &core.OtherNation{
		Nation: core.Nation{
			NationID: "ally",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
	}

	point := &core.OtherNationPoint{OtherNation: otherNation}

	// OtherNationPointはLocationメソッドを持たないため、基本的な構造確認のみ
	if point.OtherNation != otherNation {
		t.Errorf("OtherNationPoint.OtherNation = %v, want %v", point.OtherNation, otherNation)
	}
}

func TestWildernessPoint_Basic(t *testing.T) {
	enemy := &core.Enemy{
		EnemyID:        "test_goblin",
		EnemyType:      "goblin",
		Power:          10.0,
		BattleCardSlot: 1,
		Skills:         []core.EnemySkill{},
	}

	territory := &core.Territory{
		TerritoryID: "wilderness_territory",
		Cards:       []*core.StructureCard{},
		CardSlot:    3,
		BaseYield:   core.ResourceQuantity{Money: 8, Food: 4},
	}

	tests := []struct {
		name       string
		controlled bool
	}{
		{
			name:       "未制圧のWilderness",
			controlled: false,
		},
		{
			name:       "制圧済みのWilderness",
			controlled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := &core.WildernessPoint{
				Controlled: tt.controlled,
				Enemy:      enemy,
				Territory:  territory,
			}

			if point.Controlled != tt.controlled {
				t.Errorf("WildernessPoint.Controlled = %v, want %v", point.Controlled, tt.controlled)
			}
			if point.Enemy != enemy {
				t.Errorf("WildernessPoint.Enemy = %v, want %v", point.Enemy, enemy)
			}
			if point.Territory != territory {
				t.Errorf("WildernessPoint.Territory = %v, want %v", point.Territory, territory)
			}
		})
	}
}

func TestBossPoint_Basic(t *testing.T) {
	boss := &core.Enemy{
		EnemyID:        "final_boss",
		EnemyType:      "ancient_dragon",
		Power:          200.0,
		BattleCardSlot: 5,
		Skills:         []core.EnemySkill{},
	}

	point := &core.BossPoint{Boss: boss}

	if point.Boss != boss {
		t.Errorf("BossPoint.Boss = %v, want %v", point.Boss, boss)
	}
}

func TestMapGrid_GetPoint(t *testing.T) {
	// テスト用のポイント
	myNation := &core.MyNation{
		Nation: core.Nation{
			NationID: "player",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
		BasicYield: core.ResourceQuantity{Money: 10},
	}

	myNationPoint := &core.MyNationPoint{MyNation: myNation}

	// 2x2のマップグリッドを作成
	points := []core.Point{
		myNationPoint, nil,
		nil, nil,
	}

	mapGrid := &core.MapGrid{
		SizeX:  2,
		SizeY:  2,
		Points: points,
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
			name:        "有効な座標(1,0) - nil",
			x:           1,
			y:           0,
			expectedNil: true,
		},
		{
			name:        "無効な座標(-1,0)",
			x:           -1,
			y:           0,
			expectedNil: true,
		},
		{
			name:        "無効な座標(2,0)",
			x:           2,
			y:           0,
			expectedNil: true,
		},
		{
			name:        "無効な座標(0,-1)",
			x:           0,
			y:           -1,
			expectedNil: true,
		},
		{
			name:        "無効な座標(0,2)",
			x:           0,
			y:           2,
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point := mapGrid.GetPoint(tt.x, tt.y)

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
