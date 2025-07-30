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
		if m.Level >= item.RequiredLevel() {
			visibleItems = append(visibleItems, item)
		}
	}

	return visibleItems
}

// CanPurchase returns whether the item at the given index can be purchased.
func (m *Market) CanPurchase(index int, treasury *Treasury) bool {
	if index < 0 || index >= len(m.Items) {
		return false
	}

	item := m.Items[index]

	// Market level check
	if m.Level < item.RequiredLevel() {
		return false
	}

	// Check if the item itself can be purchased
	return item.CanPurchase(treasury)
}

// Purchase buys the item at the given index. Returns the card pack if applicable, false if purchase failed.
func (m *Market) Purchase(index int, treasury *Treasury) (*CardPack, bool) {
	if !m.CanPurchase(index, treasury) {
		return nil, false
	}

	item := m.Items[index]

	// Subtract the price from the treasury
	if !treasury.Sub(item.Price()) {
		return nil, false
	}

	// Apply level effect automatically
	m.Level += item.LevelEffect()

	// Handle resource trading
	if item.ResourceQuantity() != nil {
		treasury.Add(*item.ResourceQuantity())
	}

	return item.CardPack(), true
}

// MarketItem represents a card pack or investment item.
type MarketItem struct {
	cardPack         *CardPack         // The card pack (nil for investment items)
	price            ResourceQuantity  // The price of the item
	requiredLevel    MarketLevel       // The Market level required to purchase the item
	levelEffect      MarketLevel       // The level effect applied to market when purchased
	resourceQuantity *ResourceQuantity // Resource quantity given when purchased (nil if none)
}

// NewMarketItem creates a new MarketItem instance.
func NewMarketItem(cardPack *CardPack, price ResourceQuantity, requiredLevel MarketLevel, levelEffect MarketLevel) *MarketItem {
	return &MarketItem{
		cardPack:      cardPack,
		price:         price,
		requiredLevel: requiredLevel,
		levelEffect:   levelEffect,
	}
}

// NewMarketItemWithResources creates a new MarketItem with resource trading functionality.
func NewMarketItemWithResources(cardPack *CardPack, price ResourceQuantity, requiredLevel MarketLevel, levelEffect MarketLevel, resourceQuantity ResourceQuantity) *MarketItem {
	return &MarketItem{
		cardPack:         cardPack,
		price:            price,
		requiredLevel:    requiredLevel,
		levelEffect:      levelEffect,
		resourceQuantity: &resourceQuantity,
	}
}

// CardPack returns the card pack (may be nil for investment items).
func (mi *MarketItem) CardPack() *CardPack {
	return mi.cardPack
}

// Price returns the price of the item.
func (mi *MarketItem) Price() ResourceQuantity {
	return mi.price
}

// RequiredLevel returns the required market level.
func (mi *MarketItem) RequiredLevel() MarketLevel {
	return mi.requiredLevel
}

// LevelEffect returns the level effect applied when purchased.
func (mi *MarketItem) LevelEffect() MarketLevel {
	return mi.levelEffect
}

// ResourceQuantity returns the resource quantity given when purchased (may be nil).
func (mi *MarketItem) ResourceQuantity() *ResourceQuantity {
	return mi.resourceQuantity
}

// CanPurchase returns whether the item can be purchased with the given treasury.
func (mi *MarketItem) CanPurchase(treasury *Treasury) bool {
	return treasury.Resources.CanPurchase(mi.price)
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
