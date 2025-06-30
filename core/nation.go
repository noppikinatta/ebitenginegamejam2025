package core

import "fmt"

// Nation is an interface representing a nation.
// It has the necessary methods for display in MarketView.
type Nation interface {
	Name() string
	ID() NationID
	GetMarket() *Market
	VisibleMarketItems() []*MarketItem
	CanPurchase(index int, treasury *Treasury) bool
	Purchase(index int, treasury *Treasury) (*CardPack, bool)
}

// Old specification comments have been deleted and replaced with implementation.

// NationID is a unique identifier for a nation.
type NationID string

// BaseNation represents a nation.
type BaseNation struct {
	NationID NationID
	Market   *Market // A Market for purchasing card packs.
}

// ID returns the NationID.
func (n *BaseNation) ID() NationID {
	return n.NationID
}

// Name returns the nation's name.
func (n *BaseNation) Name() string {
	return fmt.Sprintf("Nation %s", n.NationID)
}

// GetMarket returns the Market.
func (n *BaseNation) GetMarket() *Market {
	return n.Market
}

// VisibleMarketItems returns a list of card packs visible in the Market.
func (n *BaseNation) VisibleMarketItems() []*MarketItem {
	return n.Market.VisibleMarketItems()
}

// CanPurchase returns whether the card pack at the given index can be purchased.
func (n *BaseNation) CanPurchase(index int, treasury *Treasury) bool {
	return n.Market.CanPurchase(index, treasury)
}

// Purchase buys the card pack at the given index. Returns false if the treasury is insufficient.
func (n *BaseNation) Purchase(index int, treasury *Treasury) (*CardPack, bool) {
	return n.Market.Purchase(index, treasury)
}

// MyNation represents the player's nation.
type MyNation struct {
	BaseNation                  // Embedded struct
	BasicYield ResourceQuantity // Basic Yield.
}

// Name returns the name of the player's nation.
func (mn *MyNation) Name() string {
	return "My Nation"
}

// AppendMarketItem adds the given item to Market.Items. This is a feature to increase the card packs that can be bought in one's own country as a reward for defeating an Enemy.
func (mn *MyNation) AppendMarketItem(item *MarketItem) {
	mn.Market.Items = append(mn.Market.Items, item)
}

// AppendLevel adds the given marketLevel to Market.Level. This is a feature to increase the card packs that can be bought in one's own country as a reward for defeating an Enemy.
func (mn *MyNation) AppendLevel(marketLevel MarketLevel) {
	mn.Market.Level += marketLevel
}

// OtherNation represents an NPC (card trading partner) nation.
type OtherNation struct {
	BaseNation // Embedded struct
}

// Purchase calls Nation.Purchase and adds 0.5 to Market.Level.
func (on *OtherNation) Purchase(index int, treasury *Treasury) (*CardPack, bool) {
	cardPack, ok := on.BaseNation.Purchase(index, treasury)
	if ok {
		on.Market.Level += 0.5
	}
	return cardPack, ok
}
