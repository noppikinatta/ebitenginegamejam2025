package viewmodel

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// TerritoryViewModel provides display information for territory UI
type TerritoryViewModel struct {
	gameState *core.GameState
	territory *core.Territory
}

// NewTerritoryViewModel creates a new TerritoryViewModel
func NewTerritoryViewModel(gameState *core.GameState, territory *core.Territory) *TerritoryViewModel {
	return &TerritoryViewModel{
		gameState: gameState,
		territory: territory,
	}
}

// Title returns the localized territory title
func (vm *TerritoryViewModel) Title() string {
	if vm.territory == nil {
		return lang.Text("territory_title")
	}
	
	// Get localized territory title based on territory ID
	return lang.Text("territory_title_" + string(vm.territory.ID()))
}

// CardSlot returns the maximum number of cards that can be placed
func (vm *TerritoryViewModel) CardSlot() int {
	if vm.territory == nil || vm.territory.Terrain() == nil {
		return 0
	}
	return vm.territory.Terrain().CardSlot()
}

// NumCards returns the current number of cards in the territory
func (vm *TerritoryViewModel) NumCards() int {
	if vm.territory == nil {
		return 0
	}
	return len(vm.territory.Cards())
}

// Card returns structure card view model at the specified index
func (vm *TerritoryViewModel) Card(idx int) *StructureCardViewModel {
	if vm.territory == nil {
		return nil
	}
	
	cards := vm.territory.Cards()
	if idx < 0 || idx >= len(cards) {
		return nil
	}
	
	card := cards[idx]
	return NewStructureCardViewModel(vm.gameState, card)
}

// Yield returns the total yield of the territory including card effects
func (vm *TerritoryViewModel) Yield() core.ResourceQuantity {
	if vm.territory == nil {
		return core.ResourceQuantity{}
	}
	
	// Start with base yield from terrain
	totalYield := vm.territory.Terrain().BaseYield()
	
	// Apply effects from placed structure cards
	for _, card := range vm.territory.Cards() {
		// Add additive yield value
		totalYield = totalYield.Add(card.YieldAdditiveValue())
		
		// Apply yield modifier
		totalYield = card.YieldModifier().Modify(totalYield)
	}
	
	return totalYield
}

// SupportPower returns the total support power provided by structure cards
func (vm *TerritoryViewModel) SupportPower() float64 {
	if vm.territory == nil {
		return 0.0
	}
	
	totalSupportPower := 0.0
	for _, card := range vm.territory.Cards() {
		totalSupportPower += card.SupportPower()
	}
	
	return totalSupportPower
}

// SupportCardSlot returns the additional card slots provided by structure cards
func (vm *TerritoryViewModel) SupportCardSlot() int {
	if vm.territory == nil {
		return 0
	}
	
	totalSupportCardSlot := 0
	for _, card := range vm.territory.Cards() {
		totalSupportCardSlot += card.SupportCardSlot()
	}
	
	return totalSupportCardSlot
} 