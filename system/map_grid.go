package system

import (
	"math/rand"
	"time"

	"github.com/noppikinatta/ebitenginegamejam2025/entity"
)

type MapGrid struct {
	width  int
	height int
	points [][]*entity.Point
}

type Pathfinder struct {
	grid *MapGrid
}

func NewMapGrid() *MapGrid {
	grid := &MapGrid{
		width:  13,
		height: 7,
		points: make([][]*entity.Point, 13),
	}

	// Initialize grid
	for x := 0; x < 13; x++ {
		grid.points[x] = make([]*entity.Point, 7)
		for y := 0; y < 7; y++ {
			grid.points[x][y] = grid.generatePoint(x, y)
		}
	}

	return grid
}

func (mg *MapGrid) generatePoint(x, y int) *entity.Point {
	// Home point at bottom-left (0,6)
	if x == 0 && y == 6 {
		point := entity.NewPoint(x, y, "Home")
		point.SetName("Your Kingdom")
		return point
	}

	// Boss point at top-right (12,0)
	if x == 12 && y == 0 {
		point := entity.NewPoint(x, y, "Boss")
		point.SetName("Dark Fortress")
		return point
	}

	// Generate random points for others
	rand.Seed(time.Now().UnixNano() + int64(x*y))
	pointType := "Wild"
	if rand.Float32() < 0.3 { // 30% chance of NPC
		pointType = "NPC"
	}

	point := entity.NewPoint(x, y, pointType)
	if pointType == "NPC" {
		npcNames := []string{"Iron Republic", "Forest Alliance", "Desert Emirate", "Mountain Clans"}
		point.SetName(npcNames[rand.Intn(len(npcNames))])
	} else {
		wildNames := []string{"Goblin Camp", "Bandit Hideout", "Ancient Ruins", "Cursed Woods"}
		point.SetName(wildNames[rand.Intn(len(wildNames))])
	}

	return point
}

func (mg *MapGrid) GetWidth() int {
	return mg.width
}

func (mg *MapGrid) GetHeight() int {
	return mg.height
}

func (mg *MapGrid) GetTotalPoints() int {
	return mg.width * mg.height
}

func (mg *MapGrid) GetPoint(x, y int) *entity.Point {
	if x < 0 || x >= mg.width || y < 0 || y >= mg.height {
		return nil
	}
	return mg.points[x][y]
}

func NewPathfinder(grid *MapGrid) *Pathfinder {
	return &Pathfinder{
		grid: grid,
	}
}

func (pf *Pathfinder) IsPointAccessible(x, y int) bool {
	// Home is always accessible
	if x == 0 && y == 6 {
		return true
	}

	// Boss is not accessible until path is clear
	if x == 12 && y == 0 {
		return false // For now, boss is not accessible
	}

	// Other points are accessible if there's a clear path from Home
	return pf.hasPathFromHome(x, y)
}

func (pf *Pathfinder) hasPathFromHome(x, y int) bool {
	// Simple path check - adjacent to home or defeated wild points create accessible paths
	// For now, make points within 2 steps of home accessible
	homeX, homeY := 0, 6
	distance := abs(x-homeX) + abs(y-homeY)
	return distance <= 3
}

func (pf *Pathfinder) FindPath(fromX, fromY, toX, toY int) [][2]int {
	// Simple pathfinding - return direct path if accessible
	if !pf.IsPointAccessible(toX, toY) {
		return nil
	}

	// Return simple path (for now just direct connection)
	return [][2]int{{fromX, fromY}, {toX, toY}}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
