package core

// TerritoryID is a unique identifier for a territory.
type TerritoryID string

// TerrainID is a unique identifier for a terrain type.
type TerrainID string

// Terrain represents the immutable properties of a terrain type.
type Terrain struct {
	id        TerrainID
	baseYield ResourceQuantity
	cardSlot  int
}

// NewTerrain creates a new Terrain instance.
func NewTerrain(id TerrainID, baseYield ResourceQuantity, cardSlot int) *Terrain {
	return &Terrain{
		id:        id,
		baseYield: baseYield,
		cardSlot:  cardSlot,
	}
}

// ID returns the terrain ID.
func (t *Terrain) ID() TerrainID {
	return t.id
}

// BaseYield returns the base yield of the terrain.
func (t *Terrain) BaseYield() ResourceQuantity {
	return t.baseYield
}

// CardSlot returns the number of card slots available in this terrain.
func (t *Terrain) CardSlot() int {
	return t.cardSlot
}

// Territory is a conquered WildernessPoint.
// A Territory acquires Resources equal to its Yield each turn.
// StructureCards can be placed in a Territory.
type Territory struct {
	id      TerritoryID
	terrain *Terrain
	cards   []*StructureCard
}

// NewTerritory creates a new Territory instance.
func NewTerritory(id TerritoryID, terrain *Terrain) *Territory {
	return &Territory{
		id:      id,
		terrain: terrain,
		cards:   make([]*StructureCard, 0, terrain.CardSlot()),
	}
}

// ID returns the territory ID.
func (t *Territory) ID() TerritoryID {
	return t.id
}

// Terrain returns the terrain of this territory.
func (t *Territory) Terrain() *Terrain {
	return t.terrain
}

// Cards returns a defensive copy of the structure cards in this territory.
func (t *Territory) Cards() []*StructureCard {
	result := make([]*StructureCard, len(t.cards))
	copy(result, t.cards)
	return result
}

// AppendCard places a StructureCard in the territory.
func (t *Territory) AppendCard(card *StructureCard) bool {
	if len(t.cards) >= t.terrain.CardSlot() {
		return false // Slot limit reached
	}
	t.cards = append(t.cards, card)
	return true
}

// RemoveCard removes the StructureCard at the specified index from the territory.
func (t *Territory) RemoveCard(index int) (*StructureCard, bool) {
	if index < 0 || index >= len(t.cards) {
		return nil, false
	}

	card := t.cards[index]
	t.cards = append(t.cards[:index], t.cards[index+1:]...)
	return card, true
}

// Yield returns the result of passing the BaseYield through the yield effects of the placed StructureCards.
func (t *Territory) Yield() ResourceQuantity {
	yield := t.terrain.BaseYield()

	// Apply additive effects first
	for _, card := range t.cards {
		yield = yield.Add(card.YieldAdditiveValue())
	}

	// Apply multiplicative effects
	for _, card := range t.cards {
		yield = card.YieldModifier().Modify(yield)
	}

	return yield
}

// ConstructionPlan manages the planned construction of StructureCards in a Territory.
type ConstructionPlan struct {
	cards []*StructureCard
}

// NewConstructionPlan creates a new ConstructionPlan based on the current state of a Territory.
func NewConstructionPlan(territory *Territory) *ConstructionPlan {
	// Defensive copy of existing cards
	cards := make([]*StructureCard, len(territory.cards))
	copy(cards, territory.cards)
	return &ConstructionPlan{cards: cards}
}

// Cards returns a defensive copy of the cards in the construction plan.
func (cp *ConstructionPlan) Cards() []*StructureCard {
	result := make([]*StructureCard, len(cp.cards))
	copy(result, cp.cards)
	return result
}

// AddCard adds a StructureCard to the construction plan.
func (cp *ConstructionPlan) AddCard(card *StructureCard) bool {
	// Card slot limit check should be done by the caller
	cp.cards = append(cp.cards, card)
	return true
}

// RemoveCard removes a StructureCard at the specified index from the construction plan.
func (cp *ConstructionPlan) RemoveCard(index int) (*StructureCard, bool) {
	if index < 0 || index >= len(cp.cards) {
		return nil, false
	}
	card := cp.cards[index]
	cp.cards = append(cp.cards[:index], cp.cards[index+1:]...)
	return card, true
}

// ApplyConstructionPlan applies the construction plan to the territory.
func (t *Territory) ApplyConstructionPlan(plan *ConstructionPlan) {
	// Defensive copy to avoid memory sharing
	t.cards = make([]*StructureCard, len(plan.cards))
	copy(t.cards, plan.cards)
}
