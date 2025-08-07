package viewmodel

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// TerritoryViewModel provides display information for territory UI
type TerritoryViewModel struct {
	territory          *core.Territory
	constructionPlan   *core.ConstructionPlan
	cardViewModelCache *CardViewModel
}

// NewTerritoryViewModel creates a new TerritoryViewModel
func NewTerritoryViewModel(territory *core.Territory, constructionPlan *core.ConstructionPlan) *TerritoryViewModel {
	return &TerritoryViewModel{
		territory:        territory,
		constructionPlan: constructionPlan,
	}
}

// Title returns the localized territory title
func (vm *TerritoryViewModel) Title() string {
	// Get localized territory title based on territory ID
	return lang.Text(string(vm.territory.Terrain().ID()))
}

// CardSlot returns the maximum number of cards that can be placed
func (vm *TerritoryViewModel) CardSlot() int {
	return vm.territory.Terrain().CardSlot()
}

// NumCards returns the current number of cards in the territory
func (vm *TerritoryViewModel) NumCards() int {
	return len(vm.territory.Cards())
}

// Card returns structure card view model at the specified index
func (vm *TerritoryViewModel) Card(idx int) (*CardViewModel, bool) {
	cards := vm.territory.Cards()
	if idx < 0 || idx >= len(cards) {
		return nil, false
	}

	card := cards[idx]
	if vm.cardViewModelCache == nil {
		vm.cardViewModelCache = &CardViewModel{}
	}
	vm.cardViewModelCache.FromStructureCard(card)
	return vm.cardViewModelCache, true
}

// Yield returns the total yield of the territory including card effects
func (vm *TerritoryViewModel) CurrentYield() core.ResourceQuantity {
	return vm.territory.Yield()
}

// PredictedYield returns the predicted yield of the territory including card effects
func (vm *TerritoryViewModel) PredictedYield() core.ResourceQuantity {
	return vm.constructionPlan.Yield()
}

// SupportPower returns the total support power provided by structure cards
func (vm *TerritoryViewModel) CurrentSupportPower() float64 {
	return vm.territory.SupportPower()
}

// PredictedSupportPower returns the predicted support power provided by structure cards
func (vm *TerritoryViewModel) PredictedSupportPower() float64 {
	return vm.constructionPlan.SupportPower()
}

// SupportCardSlot returns the additional card slots provided by structure cards
func (vm *TerritoryViewModel) CurrentSupportCardSlot() int {
	return vm.territory.SupportCardSlot()
}

// PredictedSupportCardSlot returns the additional card slots provided by structure cards
func (vm *TerritoryViewModel) PredictedSupportCardSlot() int {
	return vm.constructionPlan.SupportCardSlot()
}
