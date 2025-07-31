package viewmodel

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// ResourceViewModel provides display information for resource UI
type ResourceViewModel struct {
	gameState *core.GameState
}

// NewResourceViewModel creates a new ResourceViewModel
func NewResourceViewModel(gameState *core.GameState) *ResourceViewModel {
	return &ResourceViewModel{
		gameState: gameState,
	}
}

// Quantity returns the current resource quantity
func (vm *ResourceViewModel) Quantity() core.ResourceQuantity {
	if vm.gameState == nil || vm.gameState.Treasury == nil {
		return core.ResourceQuantity{}
	}
	return vm.gameState.Treasury.Resources
}

func (vm *ResourceViewModel) Yield() core.ResourceQuantity {
	if vm.gameState == nil {
		return core.ResourceQuantity{}
	}
	return vm.gameState.GetYield()
}
