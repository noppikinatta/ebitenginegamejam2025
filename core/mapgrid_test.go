package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestMapGrid_CanInteract(t *testing.T) {
	// MyNation for testing
	myNation := &core.MyNation{
		BaseNation: core.BaseNation{
			NationID: "player",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
		BasicYield: core.ResourceQuantity{Money: 10},
	}

	// OtherNation for testing
	otherNation := &core.OtherNation{
		BaseNation: core.BaseNation{
			NationID: "ally",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
	}

	// Enemy for testing
	enemy := &core.Enemy{
		EnemyID:        "test_orc",
		EnemyType:      "orc",
		Power:          15.0,
		BattleCardSlot: 2,
		Skills:         []core.EnemySkill{},
	}

	// Boss for testing
	boss := &core.Enemy{
		EnemyID:        "dragon_boss",
		EnemyType:      "dragon",
		Power:          100.0,
		BattleCardSlot: 4,
		Skills:         []core.EnemySkill{},
	}

	// WildernessPoint for testing (controlled)
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

	// WildernessPoint for testing (uncontrolled)
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

	points := []core.Point{
		// Row 0
		&core.MyNationPoint{MyNation: myNation},
		controlledWilderness,
		&core.OtherNationPoint{OtherNation: otherNation},
		controlledWilderness,
		&core.OtherNationPoint{OtherNation: otherNation},
		// Row 1
		uncontrolledWilderness,
		uncontrolledWilderness,
		uncontrolledWilderness,
		uncontrolledWilderness,
		controlledWilderness,
		// Row 2
		&core.OtherNationPoint{OtherNation: otherNation},
		uncontrolledWilderness,
		&core.OtherNationPoint{OtherNation: otherNation},
		controlledWilderness,
		&core.OtherNationPoint{OtherNation: otherNation},
		// Row 3
		uncontrolledWilderness,
		uncontrolledWilderness,
		uncontrolledWilderness,
		uncontrolledWilderness,
		uncontrolledWilderness,
		// Row 4
		&core.OtherNationPoint{OtherNation: otherNation},
		uncontrolledWilderness,
		&core.OtherNationPoint{OtherNation: otherNation},
		uncontrolledWilderness,
		&core.BossPoint{Boss: boss},
	}

	mapGrid := &core.MapGrid{
		Size:   core.MapGridSize{X: 5, Y: 5},
		Points: points,
	}

	tests := []struct {
		name     string
		x, y     int
		expected bool
		reason   string
	}{
		{
			name:     "MyNationPoint(0,0) is always interactable",
			x:        0,
			y:        0,
			expected: true,
			reason:   "MyNationPoint is always interactable",
		},
		{
			name:     "Controlled Wilderness(1,0) adjacent to MyNationPoint is interactable",
			x:        1,
			y:        0,
			expected: true,
			reason:   "Adjacent to MyNationPoint and controlled",
		},
		{
			name:     "Uncontrolled Wilderness(0,1) adjacent to MyNationPoint is interactable",
			x:        0,
			y:        1,
			expected: true,
			reason:   "Adjacent to MyNationPoint and uncontrolled",
		},
		{
			name:     "OtherNation(0,2) unreachable from MyNationPoint is not interactable",
			x:        0,
			y:        2,
			expected: false,
			reason:   "Unreachable from MyNationPoint",
		},
		{
			name:     "Uncontrolled Wilderness(1,2) adjacent to reachable OtherNation(2,2) is interactable",
			x:        1,
			y:        2,
			expected: true,
			reason:   "Accessible via a controlled route",
		},
		{
			name:     "Uncontrolled Wilderness(3,3) adjacent to reachable controlled Wilderness(3,2) is interactable",
			x:        3,
			y:        3,
			expected: true,
			reason:   "Accessible via a controlled route",
		},
		{
			name:     "Distant uncontrolled Wilderness(3,4) is not interactable",
			x:        3,
			y:        4,
			expected: false,
			reason:   "Inaccessible via a controlled route",
		},
		{
			name:     "Out-of-bounds coordinate (-1,0) is not interactable",
			x:        -1,
			y:        0,
			expected: false,
			reason:   "Out of map bounds",
		},
		{
			name:     "Out-of-bounds coordinate (5,0) is not interactable",
			x:        5,
			y:        0,
			expected: false,
			reason:   "Out of map bounds",
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
		BaseNation: core.BaseNation{
			NationID: "player",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
		BasicYield: core.ResourceQuantity{Money: 10},
	}

	point := &core.MyNationPoint{MyNation: myNation}

	// MyNationPoint does not have a Location method, so only basic structure is checked
	if point.MyNation != myNation {
		t.Errorf("MyNationPoint.MyNation = %v, want %v", point.MyNation, myNation)
	}
}

func TestOtherNationPoint_Location(t *testing.T) {
	otherNation := &core.OtherNation{
		BaseNation: core.BaseNation{
			NationID: "ally",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
	}

	point := &core.OtherNationPoint{OtherNation: otherNation}

	// OtherNationPoint does not have a Location method, so only basic structure is checked
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
			name:       "Uncontrolled Wilderness",
			controlled: false,
		},
		{
			name:       "Controlled Wilderness",
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
		BaseNation: core.BaseNation{
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
		Size:   core.MapGridSize{X: 2, Y: 2},
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
