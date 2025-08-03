package viewmodel

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// MapGridViewModel provides display information for map grid UI
type MapGridViewModel struct {
	gameState *core.GameState
}

// NewMapGridViewModel creates a new MapGridViewModel
func NewMapGridViewModel(gameState *core.GameState) *MapGridViewModel {
	return &MapGridViewModel{
		gameState: gameState,
	}
}

// Size returns the map grid size
func (vm *MapGridViewModel) Size() core.MapGridSize {
	if vm.gameState == nil || vm.gameState.MapGrid == nil {
		return core.MapGridSize{}
	}
	return vm.gameState.MapGrid.Size
}

// Point returns point view model at the specified coordinates
func (vm *MapGridViewModel) Point(x, y int) *PointViewModel {
	if vm.gameState == nil || vm.gameState.MapGrid == nil {
		return nil
	}

	point, ok := vm.gameState.MapGrid.GetPoint(x, y)
	if !ok {
		return nil
	}

	return NewPointViewModel(vm.gameState, point)
}

// ShouldDrawLineToRight determines if a line should be drawn to the right
func (vm *MapGridViewModel) ShouldDrawLineToRight(x, y int) bool {
	if vm.gameState == nil || vm.gameState.MapGrid == nil {
		return false
	}

	// Check if there's a connection to the right
	// For now, return a simple implementation based on grid bounds
	size := vm.gameState.MapGrid.Size
	return x < size.X-1
}

// ShouldDrawLineToUpper determines if a line should be drawn upward
func (vm *MapGridViewModel) ShouldDrawLineToUpper(x, y int) bool {
	if vm.gameState == nil || vm.gameState.MapGrid == nil {
		return false
	}

	// Check if there's a connection upward
	// For now, return a simple implementation based on grid bounds
	return y > 0
}

// PointViewModel provides display information for individual points
type PointViewModel struct {
	gameState *core.GameState
	point     core.Point
}

// NewPointViewModel creates a new PointViewModel
func NewPointViewModel(gameState *core.GameState, point core.Point) *PointViewModel {
	return &PointViewModel{
		gameState: gameState,
		point:     point,
	}
}

// Image returns the point image
func (vm *PointViewModel) Image() *ebiten.Image {
	switch vm.point.PointType() {
	case core.PointTypeMyNation:
		return drawing.Image("point-mynation")
	case core.PointTypeOtherNation:
		return drawing.Image("point-othernation")
	case core.PointTypeWilderness:
		return drawing.Image("point-wilderness")
	case core.PointTypeBoss:
		return drawing.Image("point-boss")
	default:
		return drawing.Image("dummy")
	}
}

// Name returns the localized point name
func (vm *PointViewModel) Name() string {
	// Get localized point name based on point type and context
	switch vm.point.PointType() {
	case core.PointTypeMyNation:
		return lang.Text("point_mynation")
	case core.PointTypeOtherNation:
		return lang.Text("point_othernation")
	case core.PointTypeWilderness:
		return lang.Text("point_wilderness")
	case core.PointTypeBoss:
		return lang.Text("point_boss")
	default:
		return lang.Text("point_unknown")
	}
}

// HasEnemy returns whether the point has an enemy
func (vm *PointViewModel) HasEnemy() bool {
	if battlePoint, ok := vm.point.AsBattlePoint(); ok {
		return battlePoint.Enemy() != nil
	}
	return false
}

// EnemyPower returns the enemy power if present
func (vm *PointViewModel) EnemyPower() float64 {
	if battlePoint, ok := vm.point.AsBattlePoint(); ok {
		if enemy := battlePoint.Enemy(); enemy != nil {
			return enemy.Power()
		}
	}
	return 0.0
}
