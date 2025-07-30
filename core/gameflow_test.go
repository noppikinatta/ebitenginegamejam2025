package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestGameState_AddYield(t *testing.T) {
	// Terrains for testing
	terrain1 := core.NewTerrain("terrain_1", core.ResourceQuantity{Money: 10, Food: 5}, 2)
	terrain2 := core.NewTerrain("terrain_2", core.ResourceQuantity{Money: 8, Food: 4, Wood: 3}, 3)

	// Territories for testing
	territory1 := core.NewTerritory("territory_1", terrain1)
	territory2 := core.NewTerritory("territory_2", terrain2)

	// WildernessPoint for testing (controlled)
	wilderness1 := &core.WildernessPoint{}
	wilderness1.SetControlledForTest(true)
	wilderness1.SetTerritoryForTest(territory1)

	wilderness2 := &core.WildernessPoint{}
	wilderness2.SetControlledForTest(true)
	wilderness2.SetTerritoryForTest(territory2)

	// WildernessPoint for testing (uncontrolled)
	uncontrolledTerrain := core.NewTerrain("uncontrolled", core.ResourceQuantity{Money: 100}, 1)
	uncontrolledTerritory := core.NewTerritory("uncontrolled", uncontrolledTerrain)
	uncontrolledWilderness := &core.WildernessPoint{}
	uncontrolledWilderness.SetControlledForTest(false)
	uncontrolledWilderness.SetTerritoryForTest(uncontrolledTerritory)

	// MyNation for testing
	myNation := core.NewMyNation("player", "My Nation")

	// MapGrid for testing
	points := []core.Point{
		&core.MyNationPoint{MyNation: myNation},
		wilderness1,
		wilderness2,
		uncontrolledWilderness,
	}

	mapGrid := &core.MapGrid{
		Size:   core.MapGridSize{X: 2, Y: 2},
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

	// Check initial state
	initialTreasury := gameState.Treasury.Resources

	// Execute AddYield
	gameState.AddYield()

	// Expected result: Yield of controlled Territories only (no BasicYield)
	expectedYield := territory1.Yield().Add(territory2.Yield())

	expectedTreasury := initialTreasury.Add(expectedYield)

	if gameState.Treasury.Resources != expectedTreasury {
		t.Errorf("AddYield() treasury = %v, want %v", gameState.Treasury.Resources, expectedTreasury)
	}
}

func TestGameState_NextTurn(t *testing.T) {
	myNation := core.NewMyNation("player", "My Nation")

	gameState := &core.GameState{
		MyNation: myNation,
		MapGrid:  &core.MapGrid{Size: core.MapGridSize{X: 1, Y: 1}, Points: []core.Point{&core.MyNationPoint{MyNation: myNation}}},
		Treasury: &core.Treasury{
			Resources: core.ResourceQuantity{Money: 100},
		},
		CurrentTurn: 5,
	}

	initialTurn := gameState.CurrentTurn
	initialTreasury := gameState.Treasury.Resources

	// Advance turn
	gameState.NextTurn()

	// Check if the turn has advanced
	if gameState.CurrentTurn != initialTurn+1 {
		t.Errorf("NextTurn() CurrentTurn = %v, want %v", gameState.CurrentTurn, initialTurn+1)
	}

	// Check if Yield has been added
	expectedTreasury := initialTreasury.Add(core.ResourceQuantity{})
	if gameState.Treasury.Resources != expectedTreasury {
		t.Errorf("NextTurn() treasury = %v, want %v", gameState.Treasury.Resources, expectedTreasury)
	}
}

func TestGameState_IsVictory(t *testing.T) {
	// Boss for testing
	boss := core.NewEnemy("final_boss", "dragon", 100.0, []*core.EnemySkill{}, 4)

	tests := []struct {
		name     string
		defeated bool
		expected bool
	}{
		{
			name:     "Victory if the boss is defeated",
			defeated: true,
			expected: true,
		},
		{
			name:     "Not a victory if the boss is not defeated",
			defeated: false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bossPoint := &core.BossPoint{}
			bossPoint.SetBossForTest(boss)
			bossPoint.SetDefeatedForTest(tt.defeated)

			myNation := core.NewMyNation("player", "Player Nation")

			mapGrid := &core.MapGrid{
				Size:   core.MapGridSize{X: 1, Y: 1},
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
	myNation := core.NewMyNation("player", "Player Nation")

	// Create terrains and territories for testing
	controlledTerrain := core.NewTerrain("controlled_terrain", core.ResourceQuantity{Money: 5}, 1)
	controlledTerritory := core.NewTerritory("controlled", controlledTerrain)

	uncontrolledTerrain := core.NewTerrain("uncontrolled_terrain", core.ResourceQuantity{Money: 5}, 1)
	uncontrolledTerritory := core.NewTerritory("uncontrolled", uncontrolledTerrain)

	// Create other nation
	otherNation := core.NewOtherNation("other", "Other Nation")

	// Create enemy for boss point
	bossEnemy := core.NewEnemy("boss", "dragon", 100.0, []*core.EnemySkill{}, 4)

	// Create wilderness points
	controlledWilderness1 := &core.WildernessPoint{}
	controlledWilderness1.SetControlledForTest(true)
	controlledWilderness1.SetTerritoryForTest(controlledTerritory)

	uncontrolledWilderness1 := &core.WildernessPoint{}
	uncontrolledWilderness1.SetControlledForTest(false)
	uncontrolledWilderness1.SetTerritoryForTest(uncontrolledTerritory)

	uncontrolledWilderness2 := &core.WildernessPoint{}
	uncontrolledWilderness2.SetControlledForTest(false)
	uncontrolledWilderness2.SetTerritoryForTest(uncontrolledTerritory)

	controlledWilderness2 := &core.WildernessPoint{}
	controlledWilderness2.SetControlledForTest(true)
	controlledWilderness2.SetTerritoryForTest(controlledTerritory)

	controlledWilderness3 := &core.WildernessPoint{}
	controlledWilderness3.SetControlledForTest(true)
	controlledWilderness3.SetTerritoryForTest(controlledTerritory)

	// Create boss point
	bossPoint := &core.BossPoint{}
	bossPoint.SetBossForTest(bossEnemy)

	points := []core.Point{
		&core.MyNationPoint{MyNation: myNation},
		controlledWilderness1,
		&core.OtherNationPoint{OtherNation: otherNation},
		uncontrolledWilderness1,
		uncontrolledWilderness2,
		controlledWilderness2,
		bossPoint,
		controlledWilderness3,
		&core.OtherNationPoint{OtherNation: otherNation},
	}

	mapGrid := &core.MapGrid{
		Size:   core.MapGridSize{X: 3, Y: 3},
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
			name:     "Own country's square is interactable",
			x:        0,
			y:        0,
			expected: true,
		},
		{
			name:     "Controlled Wilderness(1,0) adjacent to MyNationPoint is interactable",
			x:        1,
			y:        0,
			expected: true,
		},
		{
			name:     "Boss point is interactable",
			x:        0,
			y:        2,
			expected: true,
		},
		{
			name:     "Out-of-bounds coordinate is not interactable",
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
	myNation := core.NewMyNation("player", "Player Nation")

	myNationPoint := &core.MyNationPoint{MyNation: myNation}

	// 1x1 map grid
	points := []core.Point{myNationPoint}

	mapGrid := &core.MapGrid{
		Size:   core.MapGridSize{X: 1, Y: 1},
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
			name:        "Valid coordinates",
			x:           0,
			y:           0,
			expectedNil: false,
		},
		{
			name:        "Invalid coordinates",
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
