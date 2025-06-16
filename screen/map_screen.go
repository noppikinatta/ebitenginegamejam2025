package screen

import (
	"github.com/noppikinatta/ebitenginegamejam2025/system"
)

type MapScreen struct {
	grid       *system.MapGrid
	pathfinder *system.Pathfinder
}

func NewMapScreen() *MapScreen {
	grid := system.NewMapGrid()
	pathfinder := system.NewPathfinder(grid)

	return &MapScreen{
		grid:       grid,
		pathfinder: pathfinder,
	}
}

func (ms *MapScreen) GetGrid() *system.MapGrid {
	return ms.grid
}

func (ms *MapScreen) GetPathfinder() *system.Pathfinder {
	return ms.pathfinder
}
