package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// ResourceView is a widget for displaying resources.
// Position: (0,0,300,20).
// Displays 5 types of resources in 60x20 each.
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

	// Display 5 types of resources at 80x20 each.
	// Money (0, 0, 80, 20).
	DrawResource(screen, 0, 0, "resource-money", resources.Money, yield.Money)

	// Food (80, 0, 80, 20).
	DrawResource(screen, 80, 0, "resource-food", resources.Food, yield.Food)

	// Wood (160, 0, 80, 20).
	DrawResource(screen, 160, 0, "resource-wood", resources.Wood, yield.Wood)

	// Iron (240, 0, 80, 20).
	DrawResource(screen, 240, 0, "resource-iron", resources.Iron, yield.Iron)

	// Mana (320, 0, 80, 20).
	DrawResource(screen, 320, 0, "resource-mana", resources.Mana, yield.Mana)
}
