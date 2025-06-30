package core

// TerritoryID is a unique identifier for a territory.
type TerritoryID string

// Territory is a conquered WildernessPoint.
// A Territory acquires Resources equal to its Yield each turn.
// StructureCards can be placed in a Territory.
type Territory struct {
	TerritoryID TerritoryID
	Cards       []*StructureCard
	CardSlot    int
	BaseYield   ResourceQuantity
}

// AppendCard places a StructureCard in the territory.
func (t *Territory) AppendCard(card *StructureCard) bool {
	if len(t.Cards) >= t.CardSlot {
		return false // Slot limit reached
	}
	t.Cards = append(t.Cards, card)
	return true
}

// RemoveCard removes the StructureCard at the specified index from the territory.
func (t *Territory) RemoveCard(index int) (*StructureCard, bool) {
	if index < 0 || index >= len(t.Cards) {
		return nil, false
	}

	card := t.Cards[index]
	t.Cards = append(t.Cards[:index], t.Cards[index+1:]...)
	return card, true
}

// Yield returns the result of passing the BaseYield through the YieldModifiers of the placed StructureCards.
func (t *Territory) Yield() ResourceQuantity {
	yield := t.BaseYield

	// Apply the YieldModifiers of the placed StructureCards in order.
	for _, card := range t.Cards {
		if card.YieldModifier != nil {
			yield = card.YieldModifier.Modify(yield)
		}
	}

	return yield
}
