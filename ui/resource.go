package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// ResourceView is a widget for displaying resources.
// Position: (0,0,600,40).
// Displays 5 types of resources in 120x40 each.
type ResourceView struct {
	GameState *core.GameState
}

// NewResourceView creates a ResourceView.
func NewResourceView(gameState *core.GameState) *ResourceView {
	return &ResourceView{
		GameState: gameState,
	}
}

// HandleInput handles input (ResourceView does not accept input).
func (rv *ResourceView) HandleInput(input *Input) error {
	return nil
}

// Draw handles drawing.
func (rv *ResourceView) Draw(screen *ebiten.Image) {
	resources := rv.GameState.Treasury.Resources
	yield := rv.GameState.GetYield()

	// Display 5 types of resources at 120x40 each.
	// Money (0, 0, 120, 40).
	DrawResource(screen, 0, 0, "resource-money", resources.Money, yield.Money)

	// Food (120, 0, 120, 40).
	DrawResource(screen, 120, 0, "resource-food", resources.Food, yield.Food)

	// Wood (240, 0, 120, 40).
	DrawResource(screen, 240, 0, "resource-wood", resources.Wood, yield.Wood)

	// Iron (360, 0, 120, 40).
	DrawResource(screen, 360, 0, "resource-iron", resources.Iron, yield.Iron)

	// Mana (480, 0, 120, 40).
	DrawResource(screen, 480, 0, "resource-mana", resources.Mana, yield.Mana)
}
