package viewmodel

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// MarketViewModel provides display information for market UI
type MarketViewModel struct {
	gameState          *core.GameState
	market             *core.Market
	nation             core.Nation
	itemViewModelCache *MarketItemViewModel
}

// NewMarketViewModel creates a new MarketViewModel
func NewMarketViewModel(gameState *core.GameState) *MarketViewModel {
	return &MarketViewModel{
		gameState: gameState,
	}
}

func (vm *MarketViewModel) SelectMarket(x, y int) {
	point, ok := vm.gameState.MapGrid.GetPoint(x, y)
	if !ok {
		return
	}

	marketPoint, ok := point.AsMarketPoint()
	if !ok {
		return
	}

	market, ok := vm.gameState.Markets[marketPoint.Nation().ID()]
	if !ok {
		return
	}

	vm.nation = marketPoint.Nation()
	vm.market = market
}

// Title returns the localized market title
func (vm *MarketViewModel) Title() string {
	if vm.nation == nil {
		return ""
	}

	return lang.Text(string(vm.nation.ID()))
}

// Level returns the current market level
func (vm *MarketViewModel) Level() float64 {
	if vm.market == nil {
		return 0.0
	}
	return float64(vm.market.Level)
}

// NumItems returns the number of market items
func (vm *MarketViewModel) NumItems() int {
	if vm.market == nil {
		return 0
	}
	return len(vm.market.Items)
}

// Item returns market item view model at the specified index
func (vm *MarketViewModel) Item(idx int) (*MarketItemViewModel, bool) {
	if vm.market == nil || idx < 0 || idx >= len(vm.market.Items) {
		return nil, false
	}

	if vm.itemViewModelCache == nil {
		vm.itemViewModelCache = &MarketItemViewModel{}
	}

	item := vm.market.Items[idx]
	vm.itemViewModelCache.Init(item, vm.market.Level, vm.gameState.Treasury)

	return vm.itemViewModelCache, true
}

// MarketItemViewModel provides display information for market items
type MarketItemViewModel struct {
	item        *core.MarketItem
	marketLevel core.MarketLevel
	treasury    *core.Treasury
}

func (m *MarketItemViewModel) Init(item *core.MarketItem, marketLevel core.MarketLevel, treasury *core.Treasury) {
	m.item = item
	m.marketLevel = marketLevel
	m.treasury = treasury
}

// ItemName returns the localized item name
func (vm *MarketItemViewModel) ItemName() string {
	if vm.item.CardPack() != nil {
		return lang.Text(string(vm.item.CardPack().CardPackID))
	}

	return ""
}

// RequiredLevel returns the required market level
func (vm *MarketItemViewModel) RequiredLevel() int {
	if vm.item == nil {
		return 0
	}
	return int(vm.item.RequiredLevel())
}

// Unlocked returns whether the item is unlocked (market level meets requirement)
func (vm *MarketItemViewModel) Unlocked() bool {
	if vm.item == nil {
		return false
	}
	return vm.marketLevel >= vm.item.RequiredLevel()
}

// CanPurchase returns whether the item can be purchased
func (vm *MarketItemViewModel) CanPurchase() bool {
	if vm.treasury == nil || vm.item == nil {
		return false
	}

	return vm.Unlocked() && vm.item.CanPurchase(vm.treasury)
}

// Price returns the item price
func (vm *MarketItemViewModel) Price() core.ResourceQuantity {
	if vm.item == nil {
		return core.ResourceQuantity{}
	}
	return vm.item.Price()
}

// ResourceSufficiency returns which resources are sufficient for purchase
func (vm *MarketItemViewModel) ResourceSufficiency() ResourceSufficiency {
	if vm.treasury == nil || vm.item == nil {
		return ResourceSufficiency{}
	}

	treasury := vm.treasury.Resources
	price := vm.item.Price()

	return ResourceSufficiency{
		Money: treasury.Money >= price.Money,
		Food:  treasury.Food >= price.Food,
		Wood:  treasury.Wood >= price.Wood,
		Iron:  treasury.Iron >= price.Iron,
		Mana:  treasury.Mana >= price.Mana,
	}
}

// ResourceSufficiency indicates which resources are sufficient for purchase
type ResourceSufficiency struct {
	Money bool
	Food  bool
	Wood  bool
	Iron  bool
	Mana  bool
}
