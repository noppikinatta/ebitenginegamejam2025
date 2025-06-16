package screen

import (
	"github.com/noppikinatta/ebitenginegamejam2025/entity"
)

type DiplomacyScreen struct {
	availableCards []*entity.Card
	currentNPC     *entity.NPCInfo
}

func NewDiplomacyScreen() *DiplomacyScreen {
	// Create sample available cards
	cards := []*entity.Card{
		entity.NewCard("Iron Warrior", "Unit", map[string]int{"Gold": 50, "Iron": 20}),
		entity.NewCard("Wooden Spear", "Unit", map[string]int{"Wood": 25, "Gold": 10}),
		entity.NewCard("Magic Shield", "Enchant", map[string]int{"Mana": 30, "Gold": 40}),
		entity.NewCard("Farm Building", "Building", map[string]int{"Wood": 50, "Grain": 30}),
	}

	// Set card stats
	cards[0].SetStats(5, 3)
	cards[1].SetStats(3, 1)
	cards[2].SetStats(0, 4)
	cards[3].SetStats(0, 0)

	// Create sample NPC
	npc := entity.NewNPCInfo("Iron Republic", 3, 4, "Military")

	return &DiplomacyScreen{
		availableCards: cards,
		currentNPC:     npc,
	}
}

func (ds *DiplomacyScreen) GetAvailableCards() []*entity.Card {
	return ds.availableCards
}

func (ds *DiplomacyScreen) GetCardCost(cardName string) *entity.Cost {
	for _, card := range ds.availableCards {
		if card.Name == cardName {
			cost := card.GetCost()
			return entity.NewCost(
				cost["Gold"],
				cost["Iron"],
				cost["Wood"],
				cost["Grain"],
				cost["Mana"],
			)
		}
	}
	return nil
}

func (ds *DiplomacyScreen) GetCurrentNPCInfo() *entity.NPCInfo {
	return ds.currentNPC
}

func (ds *DiplomacyScreen) SetCurrentNPC(npc *entity.NPCInfo) {
	ds.currentNPC = npc
	// Could update available cards based on NPC specialty
}
