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
	rand.Seed(time.Now().UnixNano())

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

	// Boss やその他のポイントは、パスが存在する場合にアクセス可能
	return pf.hasPathFromHome(x, y)
}

func (pf *Pathfinder) hasPathFromHome(destX, destY int) bool {
	homeX, homeY := 0, 6

	// BFS 探索
	type node struct{ x, y int }
	visited := make(map[[2]int]bool)
	queue := []node{{homeX, homeY}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.x == destX && cur.y == destY {
			return true
		}

		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for _, d := range dirs {
			nx, ny := cur.x+d[0], cur.y+d[1]
			key := [2]int{nx, ny}
			if visited[key] {
				continue
			}

			pt := pf.grid.GetPoint(nx, ny)
			if pt == nil {
				continue
			}

			// 通行可能条件
			passable := true
			if pt.Type == "Wild" && !pt.Defeated {
				passable = false
			}
			// NPC や Home は通行可能、Boss は目的地としてのみ許可
			if passable || (nx == destX && ny == destY) {
				visited[key] = true
				queue = append(queue, node{nx, ny})
			}
		}
	}

	return false
}

func (pf *Pathfinder) FindPath(fromX, fromY, toX, toY int) [][2]int {
	if !pf.IsPointAccessible(toX, toY) {
		return nil
	}

	type node struct {
		x, y int
		prev *node
	}
	visited := make(map[[2]int]bool)
	queue := []node{{fromX, fromY, nil}}
	var dest *node

	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.x == toX && cur.y == toY {
			dest = &cur
			break
		}

		for _, d := range dirs {
			nx, ny := cur.x+d[0], cur.y+d[1]
			key := [2]int{nx, ny}
			if visited[key] {
				continue
			}
			pt := pf.grid.GetPoint(nx, ny)
			if pt == nil {
				continue
			}
			passable := true
			if pt.Type == "Wild" && !pt.Defeated {
				passable = false
			}
			if passable || (nx == toX && ny == toY) {
				visited[key] = true
				queue = append(queue, node{nx, ny, &cur})
			}
		}
	}

	if dest == nil {
		return nil
	}
	// Reconstruct path
	var path [][2]int
	for n := dest; n != nil; n = n.prev {
		path = append([][2]int{{n.x, n.y}}, path...)
	}
	return path
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
