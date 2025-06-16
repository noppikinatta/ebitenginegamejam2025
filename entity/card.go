package entity

type Card struct {
	Name        string
	Type        string         // "Unit", "Enchant", "Building"
	Cost        map[string]int // Resource costs
	Description string
	Attack      int
	Defense     int
}

type NPCInfo struct {
	Name         string
	Location     [2]int // x, y coordinates
	Specialty    string // "Military", "Magic", "Trade", "Farming"
	Relationship int    // 0-100, higher is better
}

type Cost struct {
	Gold  int
	Iron  int
	Wood  int
	Grain int
	Mana  int
}

func NewCard(name, cardType string, cost map[string]int) *Card {
	return &Card{
		Name:        name,
		Type:        cardType,
		Cost:        cost,
		Description: "",
		Attack:      0,
		Defense:     0,
	}
}

func (c *Card) GetCost() map[string]int {
	return c.Cost
}

func (c *Card) SetStats(attack, defense int) {
	c.Attack = attack
	c.Defense = defense
}

func NewNPCInfo(name string, x, y int, specialty string) *NPCInfo {
	return &NPCInfo{
		Name:         name,
		Location:     [2]int{x, y},
		Specialty:    specialty,
		Relationship: 50, // Neutral start
	}
}

func NewCost(gold, iron, wood, grain, mana int) *Cost {
	return &Cost{
		Gold:  gold,
		Iron:  iron,
		Wood:  wood,
		Grain: grain,
		Mana:  mana,
	}
}
