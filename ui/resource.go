package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// ResourceView is a widget for displaying resources.
// Position: (0,0,600,40).
// Displays 5 types of resources in 120x40 each.
type ResourceView struct {
	ViewModel *viewmodel.ResourceViewModel
}

// NewResourceView creates a ResourceView.
func NewResourceView(viewModel *viewmodel.ResourceViewModel) *ResourceView {
	return &ResourceView{
		ViewModel: viewModel,
	}
}

// HandleInput handles input (ResourceView does not accept input).
func (rv *ResourceView) HandleInput(input *Input) error {
	return nil
}

// Draw handles drawing.
func (rv *ResourceView) Draw(screen *ebiten.Image) {
	resources := rv.ViewModel.Quantity()

	// TODO: Get yield information from viewmodel (currently not available)
	// For now, display without yield increment

	// Display 5 types of resources at 120x40 each.
	// Money (0, 0, 120, 40).
	DrawResource(screen, 0, 0, "resource-money", resources.Money, 0)

	// Food (120, 0, 120, 40).
	DrawResource(screen, 120, 0, "resource-food", resources.Food, 0)

	// Wood (240, 0, 120, 40).
	DrawResource(screen, 240, 0, "resource-wood", resources.Wood, 0)

	// Iron (360, 0, 120, 40).
	DrawResource(screen, 360, 0, "resource-iron", resources.Iron, 0)

	// Mana (480, 0, 120, 40).
	DrawResource(screen, 480, 0, "resource-mana", resources.Mana, 0)
}
