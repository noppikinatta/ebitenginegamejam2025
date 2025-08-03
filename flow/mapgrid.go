package flow

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// MapGridFlow handles map grid point selection operations
type MapGridFlow struct {
	gameState     *core.GameState
	selectedPoint core.Point
}

// NewMapGridFlow creates a new MapGridFlow
func NewMapGridFlow(gameState *core.GameState) *MapGridFlow {
	return &MapGridFlow{
		gameState: gameState,
	}
}

// SelectPoint selects a point at the specified coordinates
func (mf *MapGridFlow) SelectPoint(x, y int) bool {
	if mf.gameState == nil || mf.gameState.MapGrid == nil {
		return false
	}

	point, ok := mf.gameState.MapGrid.GetPoint(x, y)
	if !ok {
		return false
	}

	mf.selectedPoint = point
	return true
}

// GetSelectedPoint returns the currently selected point
func (mf *MapGridFlow) GetSelectedPoint() core.Point {
	return mf.selectedPoint
}

// ClearSelection clears the current point selection
func (mf *MapGridFlow) ClearSelection() {
	mf.selectedPoint = nil
}

// GetPointType returns the type of the selected point
func (mf *MapGridFlow) GetPointType() core.PointType {
	if mf.selectedPoint == nil {
		return core.PointTypeUnknown
	}

	switch mf.selectedPoint.(type) {
	case *core.MyNationPoint:
		return core.PointTypeMyNation
	case *core.OtherNationPoint:
		return core.PointTypeOtherNation
	case *core.WildernessPoint:
		return core.PointTypeWilderness
	case *core.BossPoint:
		return core.PointTypeBoss
	default:
		return core.PointTypeUnknown
	}
}

// GetNationFromSelectedPoint returns the nation from the selected point if applicable
func (mf *MapGridFlow) GetNationFromSelectedPoint() core.Nation {
	if mf.selectedPoint == nil {
		return nil
	}

	switch p := mf.selectedPoint.(type) {
	case *core.MyNationPoint:
		return p.MyNation
	case *core.OtherNationPoint:
		return p.OtherNation
	default:
		return nil
	}
}

// GetTerritoryFromSelectedPoint returns the territory from the selected point if applicable
func (mf *MapGridFlow) GetTerritoryFromSelectedPoint() *core.Territory {
	if mf.selectedPoint == nil {
		return nil
	}

	if p, ok := mf.selectedPoint.(*core.WildernessPoint); ok {
		if p.Controlled() {
			return p.Territory()
		}
	}

	return nil
}

// GetBattlePointFromSelectedPoint returns the battle point if applicable
func (mf *MapGridFlow) GetBattlePointFromSelectedPoint() core.BattlePoint {
	if mf.selectedPoint == nil {
		return nil
	}

	switch p := mf.selectedPoint.(type) {
	case *core.WildernessPoint:
		if !p.Controlled() {
			return p
		}
	case *core.BossPoint:
		return p
	}

	return nil
}

// IsSelectedPointControllable checks if the selected point can be controlled
func (mf *MapGridFlow) IsSelectedPointControllable() bool {
	if p, ok := mf.selectedPoint.(*core.WildernessPoint); ok {
		return !p.Controlled()
	}

	if _, ok := mf.selectedPoint.(*core.BossPoint); ok {
		return true // Boss points can be challenged
	}

	return false
}

// IsSelectedPointControlled checks if the selected point is already controlled
func (mf *MapGridFlow) IsSelectedPointControlled() bool {
	if p, ok := mf.selectedPoint.(*core.WildernessPoint); ok {
		return p.Controlled()
	}

	return false
}
