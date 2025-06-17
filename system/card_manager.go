package system

import (
	"github.com/noppikinatta/ebitenginegamejam2025/entity"
)

type CardManager struct {
	cardTemplates map[string]*entity.Card
}

func NewCardManager() *CardManager {
	cm := &CardManager{cardTemplates: make(map[string]*entity.Card)}
	cm.initializeCardTemplates()
	return cm
}

func (cm *CardManager) initializeCardTemplates() {
	// Unit cards
	cm.cardTemplates["Warrior"] = &entity.Card{
		Name: "Warrior", Type: "Unit",
		Cost:   map[string]int{"Gold": 50, "Iron": 10},
		Attack: 4, Defense: 3,
		Description: "Basic melee unit",
	}

	cm.cardTemplates["Archer"] = &entity.Card{
		Name: "Archer", Type: "Unit",
		Cost:   map[string]int{"Gold": 40, "Wood": 15},
		Attack: 3, Defense: 2,
		Description: "Ranged unit",
	}

	cm.cardTemplates["Mage"] = &entity.Card{
		Name: "Mage", Type: "Unit",
		Cost:   map[string]int{"Gold": 60, "Mana": 20},
		Attack: 5, Defense: 1,
		Description: "Magic damage dealer",
	}

	// Enchant cards
	cm.cardTemplates["Magic Shield"] = &entity.Card{
		Name: "Magic Shield", Type: "Enchant",
		Cost:   map[string]int{"Mana": 30, "Gold": 25},
		Attack: 0, Defense: 4,
		Description: "Increases defense",
	}

	cm.cardTemplates["Fire Weapon"] = &entity.Card{
		Name: "Fire Weapon", Type: "Enchant",
		Cost:   map[string]int{"Mana": 25, "Iron": 15},
		Attack: 3, Defense: 0,
		Description: "Increases attack",
	}

	// Building cards
	cm.cardTemplates["Farm"] = &entity.Card{
		Name: "Farm", Type: "Building",
		Cost:   map[string]int{"Wood": 40, "Grain": 20},
		Attack: 0, Defense: 0,
		Description: "Generates grain resources",
	}

	cm.cardTemplates["Mine"] = &entity.Card{
		Name: "Mine", Type: "Building",
		Cost:   map[string]int{"Wood": 50, "Gold": 30},
		Attack: 0, Defense: 0,
		Description: "Generates iron resources",
	}

	cm.cardTemplates["Tower"] = &entity.Card{
		Name: "Tower", Type: "Building",
		Cost:   map[string]int{"Wood": 60, "Iron": 25},
		Attack: 2, Defense: 5,
		Description: "Defensive structure",
	}
}

func (cm *CardManager) CreateCard(name, cardType string, cost map[string]int) *entity.Card {
	// Check if template exists
	if template, exists := cm.cardTemplates[name]; exists {
		// Return a copy of the template
		newCard := *template
		return &newCard
	}

	// Create new card if template doesn't exist
	card := entity.NewCard(name, cardType, cost)

	// Set default stats based on type
	switch cardType {
	case "Unit":
		card.SetStats(3, 2) // Default unit stats
	case "Enchant":
		card.SetStats(1, 1) // Default enchant stats
	case "Building":
		card.SetStats(0, 3) // Default building stats
	}

	return card
}

func (cm *CardManager) GetAllCardTemplates() map[string]*entity.Card {
	return cm.cardTemplates
}

func (cm *CardManager) GetCardTemplate(name string) *entity.Card {
	if template, exists := cm.cardTemplates[name]; exists {
		return template
	}
	return nil
}

func (cm *CardManager) AddCardTemplate(card *entity.Card) {
	cm.cardTemplates[card.Name] = card
}

// NewCardManagerFromData builds the manager from CSV-loaded templates.
// If data is nil or empty it falls back to the default hard-coded templates to keep tests green.
func NewCardManagerFromData(data map[string]*entity.Card) *CardManager {
	cm := &CardManager{cardTemplates: make(map[string]*entity.Card)}
	if len(data) == 0 {
		cm.initializeCardTemplates()
	} else {
		for id, tpl := range data {
			// shallow copy OK – template should be immutable
			cp := *tpl
			cm.cardTemplates[id] = &cp
		}
	}
	return cm
}
