package system

import (
	"fmt"
)

type Territory struct {
	X, Y       int
	Name       string
	Type       string         // "Home", "Wild", "NPC", "Boss"
	Resources  map[string]int // Resources generated per turn
	Controlled bool
}

type TerritoryManager struct {
	territories     map[string]*Territory // Key: "x,y"
	resourceManager *ResourceManager
}

func NewTerritoryManager(resourceManager *ResourceManager) *TerritoryManager {
	tm := &TerritoryManager{
		territories:     make(map[string]*Territory),
		resourceManager: resourceManager,
	}

	tm.initializeTerritories()
	return tm
}

func (tm *TerritoryManager) initializeTerritories() {
	// Initialize home territory as controlled
	home := &Territory{
		X: 0, Y: 6,
		Name: "Your Kingdom",
		Type: "Home",
		Resources: map[string]int{
			"Gold": 10, "Iron": 5, "Wood": 8, "Grain": 12, "Mana": 3,
		},
		Controlled: true,
	}
	tm.territories["0,6"] = home

	// Add some sample territories
	territories := []*Territory{
		{X: 1, Y: 6, Name: "Iron Mine", Type: "Wild",
			Resources: map[string]int{"Iron": 15, "Gold": 5}, Controlled: false},
		{X: 0, Y: 5, Name: "Forest", Type: "Wild",
			Resources: map[string]int{"Wood": 20, "Grain": 5}, Controlled: false},
		{X: 2, Y: 6, Name: "Farmland", Type: "Wild",
			Resources: map[string]int{"Grain": 25, "Gold": 3}, Controlled: false},
		{X: 3, Y: 4, Name: "Iron Republic", Type: "NPC",
			Resources: map[string]int{"Iron": 10, "Gold": 8}, Controlled: false},
	}

	for _, territory := range territories {
		key := fmt.Sprintf("%d,%d", territory.X, territory.Y)
		tm.territories[key] = territory
	}
}

func (tm *TerritoryManager) GetControlledTerritories() []*Territory {
	var controlled []*Territory
	for _, territory := range tm.territories {
		if territory.Controlled {
			controlled = append(controlled, territory)
		}
	}
	return controlled
}

func (tm *TerritoryManager) ConquerTerritory(x, y int) bool {
	key := fmt.Sprintf("%d,%d", x, y)
	if territory, exists := tm.territories[key]; exists {
		if !territory.Controlled && territory.Type != "NPC" {
			territory.Controlled = true
			return true
		}
	}
	return false
}

func (tm *TerritoryManager) GenerateResources() {
	for _, territory := range tm.territories {
		if territory.Controlled {
			for resourceType, amount := range territory.Resources {
				tm.resourceManager.AddResource(resourceType, amount)
			}
		}
	}
}

func (tm *TerritoryManager) GetTerritory(x, y int) *Territory {
	key := fmt.Sprintf("%d,%d", x, y)
	return tm.territories[key]
}

func (tm *TerritoryManager) GetAllTerritories() map[string]*Territory {
	return tm.territories
}
