package core

// MarketLevel is the level of the Market. It is used to determine if a MarketItem is visible to the player.
type MarketLevel float64

// Market is where card packs can be purchased.
type Market struct {
	Level MarketLevel   // The level of this Market.
	Items []*MarketItem // A list of card packs.
}

// VisibleMarketItems returns a list of MarketItems that are visible in the Market.
func (m *Market) VisibleMarketItems() []*MarketItem {
	var visibleItems []*MarketItem

	for _, item := range m.Items {
		if m.Level >= item.RequiredLevel {
			visibleItems = append(visibleItems, item)
		}
	}

	return visibleItems
}

// CanPurchase returns whether the card pack at the given index can be purchased.
func (m *Market) CanPurchase(index int, treasury *Treasury) bool {
	if index < 0 || index >= len(m.Items) {
		return false
	}

	item := m.Items[index]

	// Market level check
	if m.Level < item.RequiredLevel {
		return false
	}

	// Check if the item itself can be purchased
	return item.CanPurchase(treasury)
}

// Purchase buys the card pack at the given index. Returns false if the treasury is insufficient.
func (m *Market) Purchase(index int, treasury *Treasury) (*CardPack, bool) {
	if !m.CanPurchase(index, treasury) {
		return nil, false
	}

	item := m.Items[index]

	// Subtract the price from the treasury
	if !treasury.Sub(item.Price) {
		return nil, false
	}

	return item.CardPack, true
}

// MarketItem represents a card pack.
type MarketItem struct {
	CardPack      *CardPack
	Price         ResourceQuantity // The price of the card pack.
	RequiredLevel MarketLevel      // The Market level required to purchase the card pack.
}

// CanPurchase returns whether the card pack can be purchased with the given treasury.
func (mi *MarketItem) CanPurchase(treasury *Treasury) bool {
	return treasury.Resources.CanPurchase(mi.Price)
}

// Treasury represents the treasury.
type Treasury struct {
	Resources ResourceQuantity // The amount of Resources in the treasury.
}

// Add adds the given other to the treasury.
func (t *Treasury) Add(other ResourceQuantity) {
	t.Resources = t.Resources.Add(other)
}

// Sub subtracts the given other from the treasury. If the treasury is insufficient, it does not perform the subtraction and returns false.
func (t *Treasury) Sub(other ResourceQuantity) bool {
	if !t.Resources.CanPurchase(other) {
		return false
	}
	t.Resources = t.Resources.Sub(other)
	return true
}
