package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestMapGrid_CanInteract(t *testing.T) {
	// MyNation for testing
	myNation := core.NewMyNation("player", "Player Nation")

	// OtherNation for testing
	otherNation := core.NewOtherNation("ally", "Ally Nation")

	// Enemy for testing
	enemy := core.NewEnemy("test_orc", "orc", 15.0, []*core.EnemySkill{}, 2)

	// Boss for testing
	boss := core.NewEnemy("dragon_boss", "dragon", 100.0, []*core.EnemySkill{}, 4)

	// Terrain for territories
	controlledTerrain := core.NewTerrain("controlled_terrain", core.ResourceQuantity{Money: 5}, 2)
	uncontrolledTerrain := core.NewTerrain("uncontrolled_terrain", core.ResourceQuantity{Money: 5}, 2)

	// Territory for testing
	controlledTerritory := core.NewTerritory("controlled_territory", controlledTerrain)
	uncontrolledTerritory := core.NewTerritory("uncontrolled_territory", uncontrolledTerrain)

	// WildernessPoint for testing (controlled)
	controlledWilderness := &core.WildernessPoint{}
	controlledWilderness.SetControlledForTest(true)
	controlledWilderness.SetEnemyForTest(enemy)
	controlledWilderness.SetTerritoryForTest(controlledTerritory)

	// WildernessPoint for testing (uncontrolled)
	uncontrolledWilderness := &core.WildernessPoint{}
	uncontrolledWilderness.SetControlledForTest(false)
	uncontrolledWilderness.SetEnemyForTest(enemy)
	uncontrolledWilderness.SetTerritoryForTest(uncontrolledTerritory)

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
		func() *core.BossPoint {
			bp := &core.BossPoint{}
			bp.SetBossForTest(boss)
			return bp
		}(),
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
	myNation := core.NewMyNation("player", "Player Nation")

	point := &core.MyNationPoint{MyNation: myNation}

	// MyNationPoint does not have a Location method, so only basic structure is checked
	if point.MyNation != myNation {
		t.Errorf("MyNationPoint.MyNation = %v, want %v", point.MyNation, myNation)
	}
}

func TestOtherNationPoint_Location(t *testing.T) {
	otherNation := core.NewOtherNation("ally", "Ally Nation")

	point := &core.OtherNationPoint{OtherNation: otherNation}

	// OtherNationPoint does not have a Location method, so only basic structure is checked
	if point.OtherNation != otherNation {
		t.Errorf("OtherNationPoint.OtherNation = %v, want %v", point.OtherNation, otherNation)
	}
}

func TestWildernessPoint_Basic(t *testing.T) {
	enemy := core.NewEnemy("test_goblin", "goblin", 10.0, []*core.EnemySkill{}, 1)

	terrain := core.NewTerrain("wilderness_terrain", core.ResourceQuantity{Money: 8, Food: 4}, 3)
	territory := core.NewTerritory("wilderness_territory", terrain)

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
			point := &core.WildernessPoint{}
			point.SetControlledForTest(tt.controlled)
			point.SetEnemyForTest(enemy)
			point.SetTerritoryForTest(territory)

			if point.Controlled() != tt.controlled {
				t.Errorf("WildernessPoint.Controlled() = %v, want %v", point.Controlled(), tt.controlled)
			}
			if point.Enemy() != enemy {
				t.Errorf("WildernessPoint.Enemy() = %v, want %v", point.Enemy(), enemy)
			}
			if point.Territory() != territory {
				t.Errorf("WildernessPoint.Territory() = %v, want %v", point.Territory(), territory)
			}
		})
	}
}

func TestBossPoint_Basic(t *testing.T) {
	boss := core.NewEnemy("final_boss", "ancient_dragon", 200.0, []*core.EnemySkill{}, 5)

	point := &core.BossPoint{}
	point.SetBossForTest(boss)

	if point.Boss() != boss {
		t.Errorf("BossPoint.Boss() = %v, want %v", point.Boss(), boss)
	}
}

func TestMapGrid_GetPoint(t *testing.T) {
	// Points for testing
	myNation := core.NewMyNation("player", "Player Nation")

	myNationPoint := &core.MyNationPoint{MyNation: myNation}

	// Create 2x2 map grid
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
			name:        "Valid coordinate (0,0)",
			x:           0,
			y:           0,
			expectedNil: false,
		},
		{
			name:        "Valid coordinate (1,0) - nil",
			x:           1,
			y:           0,
			expectedNil: true,
		},
		{
			name:        "Invalid coordinate (-1,0)",
			x:           -1,
			y:           0,
			expectedNil: true,
		},
		{
			name:        "Invalid coordinate (2,0)",
			x:           2,
			y:           0,
			expectedNil: true,
		},
		{
			name:        "Invalid coordinate (0,-1)",
			x:           0,
			y:           -1,
			expectedNil: true,
		},
		{
			name:        "Invalid coordinate (0,2)",
			x:           0,
			y:           2,
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			point, ok := mapGrid.GetPoint(tt.x, tt.y)

			if tt.expectedNil {
				if ok {
					t.Errorf("GetPoint(%d, %d) = %v, want nil", tt.x, tt.y, point)
				}
			} else {
				if !ok {
					t.Errorf("GetPoint(%d, %d) = nil, want non-nil", tt.x, tt.y)
				} else if point != myNationPoint {
					t.Errorf("GetPoint(%d, %d) = %v, want %v", tt.x, tt.y, point, myNationPoint)
				}
			}
		})
	}
}
