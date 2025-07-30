package viewmodel

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// MarketViewModel provides display information for market UI
type MarketViewModel struct {
	gameState *core.GameState
	market    *core.Market
	nation    core.Nation
}

// NewMarketViewModel creates a new MarketViewModel
func NewMarketViewModel(gameState *core.GameState, market *core.Market, nation core.Nation) *MarketViewModel {
	return &MarketViewModel{
		gameState: gameState,
		market:    market,
		nation:    nation,
	}
}

// Title returns the localized market title
func (vm *MarketViewModel) Title() string {
	if vm.nation == nil {
		return lang.Text("market_title")
	}

	// Get localized market title based on nation
	return lang.Text("market_title_" + string(vm.nation.ID()))
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
func (vm *MarketViewModel) Item(idx int) *MarketItemViewModel {
	if vm.market == nil || idx < 0 || idx >= len(vm.market.Items) {
		return nil
	}

	item := vm.market.Items[idx]
	return NewMarketItemViewModel(vm.gameState, vm.market, item)
}

// MarketItemViewModel provides display information for market items
type MarketItemViewModel struct {
	gameState *core.GameState
	market    *core.Market
	item      *core.MarketItem
}

// NewMarketItemViewModel creates a new MarketItemViewModel
func NewMarketItemViewModel(gameState *core.GameState, market *core.Market, item *core.MarketItem) *MarketItemViewModel {
	return &MarketItemViewModel{
		gameState: gameState,
		market:    market,
		item:      item,
	}
}

// ItemName returns the localized item name
func (vm *MarketItemViewModel) ItemName() string {
	// Get localized market item name
	// This might depend on the card pack or investment type
	if vm.item.CardPack() != nil {
		return lang.Text("cardpack_name_" + string(vm.item.CardPack().CardPackID))
	}

	// For investment items without card packs
	return lang.Text("investment_item")
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
	if vm.market == nil || vm.item == nil {
		return false
	}
	return vm.market.Level >= vm.item.RequiredLevel()
}

// CanPurchase returns whether the item can be purchased
func (vm *MarketItemViewModel) CanPurchase() bool {
	if vm.gameState == nil || vm.item == nil {
		return false
	}

	return vm.Unlocked() && vm.item.CanPurchase(vm.gameState.Treasury)
}

// Price returns the item price
func (vm *MarketItemViewModel) Price() core.ResourceQuantity {
	if vm.item == nil {
		return core.ResourceQuantity{}
	}
	return vm.item.Price()
}

// ResourceSufficiency returns which resources are sufficient for purchase
func (vm *MarketItemViewModel) ResourceSufficiency() *ResourceSufficiency {
	if vm.gameState == nil || vm.item == nil {
		return &ResourceSufficiency{}
	}

	treasury := vm.gameState.Treasury.Resources
	price := vm.item.Price()

	return &ResourceSufficiency{
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
