package core_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

func TestGameState_AddYield(t *testing.T) {
	// Territories for testing
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

	// WildernessPoint for testing (controlled)
	wilderness1 := &core.WildernessPoint{
		Controlled: true,
		Territory:  territory1,
	}

	wilderness2 := &core.WildernessPoint{
		Controlled: true,
		Territory:  territory2,
	}

	// WildernessPoint for testing (uncontrolled)
	uncontrolledWilderness := &core.WildernessPoint{
		Controlled: false,
		Territory: &core.Territory{
			TerritoryID: "uncontrolled",
			BaseYield:   core.ResourceQuantity{Money: 100}, // No income as it is not controlled
		},
	}

	// MyNation for testing
	myNation := &core.MyNation{
		BaseNation: core.BaseNation{
			NationID: "player",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
		BasicYield: core.ResourceQuantity{Money: 5, Food: 2, Mana: 1},
	}

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

	// Expected result: BasicYield + Yield of controlled Territories
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
	expectedTreasury := initialTreasury.Add(myNation.BasicYield)
	if gameState.Treasury.Resources != expectedTreasury {
		t.Errorf("NextTurn() treasury = %v, want %v", gameState.Treasury.Resources, expectedTreasury)
	}
}

func TestGameState_IsVictory(t *testing.T) {
	// Boss for testing
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
	myNation := &core.MyNation{
		BaseNation: core.BaseNation{
			NationID: "player",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
		BasicYield: core.ResourceQuantity{Money: 5},
	}

	points := []core.Point{
		&core.MyNationPoint{MyNation: myNation},
		&core.WildernessPoint{
			Controlled: true,
			Territory: &core.Territory{
				TerritoryID: "controlled",
				BaseYield:   core.ResourceQuantity{Money: 5},
			},
		},
		&core.OtherNationPoint{OtherNation: &core.OtherNation{}},
		&core.WildernessPoint{
			Controlled: false,
			Territory: &core.Territory{
				TerritoryID: "uncontrolled",
				BaseYield:   core.ResourceQuantity{Money: 5},
			},
		},
		&core.WildernessPoint{
			Controlled: false,
			Territory: &core.Territory{
				TerritoryID: "uncontrolled",
				BaseYield:   core.ResourceQuantity{Money: 5},
			},
		},
		&core.WildernessPoint{
			Controlled: true,
			Territory: &core.Territory{
				TerritoryID: "controlled",
				BaseYield:   core.ResourceQuantity{Money: 5},
			},
		},
		&core.BossPoint{Boss: &core.Enemy{}},
		&core.WildernessPoint{
			Controlled: true,
			Territory: &core.Territory{
				TerritoryID: "controlled",
				BaseYield:   core.ResourceQuantity{Money: 5},
			},
		},
		&core.OtherNationPoint{OtherNation: &core.OtherNation{}},
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
	myNation := &core.MyNation{
		BaseNation: core.BaseNation{
			NationID: "player",
			Market:   &core.Market{Level: 1.0, Items: []*core.MarketItem{}},
		},
		BasicYield: core.ResourceQuantity{Money: 5},
	}

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
