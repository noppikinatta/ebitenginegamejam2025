package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/geom"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// MapGridView is a widget for displaying the map grid.
// Position: (0,20,520,280) - Drawn within MainView.
// 5x5 Point arrangement (fixed), 520x280 divided into 5x5 (104x56 cells).
type MapGridView struct {
	GameState     *core.GameState
	TopLeft       geom.PointF
	CellSize      geom.PointF
	CellLocations []geom.PointF

	// View switching callback.
	OnPointClicked func(point core.Point)
}

// NewMapGridView creates a MapGridView.
func NewMapGridView(gameState *core.GameState, onPointClicked func(point core.Point)) *MapGridView {
	cellSize := geom.PointF{X: 520.0 / 5.0, Y: 280.0 / 5.0}
	cellLocations := make([]geom.PointF, 25)
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			// Calculate the drawing Y coordinate so that (0,0) is at the bottom left.
			// y is the logical coordinate (0 is the bottom), converted to drawing coordinate (0 is the top) with 4-y.
			cellLocations[y*5+x] = geom.PointF{X: float64(x) * cellSize.X, Y: float64(4-y) * cellSize.Y}
		}
	}
	return &MapGridView{
		GameState:      gameState,
		TopLeft:        geom.PointF{X: 0, Y: 20},
		CellSize:       cellSize,
		CellLocations:  cellLocations,
		OnPointClicked: onPointClicked,
	}
}

// HandleInput handles input.
func (m *MapGridView) HandleInput(input *Input) error {
	justReleased := input.Mouse.IsJustReleased(ebiten.MouseButtonLeft)
	if !justReleased {
		return nil
	}

	cursorX, cursorY := input.Mouse.CursorPosition()
	relativeX := float64(cursorX) - m.TopLeft.X
	relativeY := float64(cursorY) - m.TopLeft.Y

	viewWidth := m.CellSize.X * 5
	viewHeight := m.CellSize.Y * 5

	if relativeX < 0 || relativeX >= viewWidth || relativeY < 0 || relativeY >= viewHeight {
		return nil
	}

	drawGridX := int(relativeX / m.CellSize.X)
	drawGridY := int(relativeY / m.CellSize.Y)

	// Since (0,0) is at the bottom left, the Y coordinate is inverted.
	gridY := 4 - drawGridY

	point := m.GameState.MapGrid.GetPoint(drawGridX, gridY)
	if point == nil {
		return nil
	}

	// Calculate the drawing area of the Point image.
	cellTopLeft := m.CellLocations[gridY*5+drawGridX]
	imageX := cellTopLeft.X + (m.CellSize.X-24)/2
	imageY := cellTopLeft.Y + (m.CellSize.Y-24)/2 - 10 // Consider the offset in the Draw method.
	imageWidth := 24.0
	imageHeight := 24.0

	// Check if the click position is within the range of the Point image.
	if relativeX >= imageX && relativeX < imageX+imageWidth &&
		relativeY >= imageY && relativeY < imageY+imageHeight {

		// Check if the Point is reachable.
		if m.GameState.MapGrid.CanInteract(drawGridX, gridY) {
			if m.OnPointClicked != nil {
				m.OnPointClicked(point)
			}
		}
	}

	return nil
}

// Draw handles drawing.
func (m *MapGridView) Draw(screen *ebiten.Image) {
	if m.GameState == nil || m.GameState.MapGrid == nil {
		return
	}

	mapGrid := m.GameState.MapGrid

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			point := mapGrid.GetPoint(x, y)
			if point == nil {
				continue
			}

			// Get the top-left coordinates of the cell.
			cellTopLeft := m.CellLocations[y*5+x]
			screenX := cellTopLeft.X + m.TopLeft.X
			screenY := cellTopLeft.Y + m.TopLeft.Y

			// Draw reachability lines.
			if m.GameState.CanInteract(x, y) {
				cellCenterX := screenX + m.CellSize.X/2
				cellCenterY := screenY + m.CellSize.Y/2
				m.drawConnectionLines(screen, x, y, cellCenterX, cellCenterY)
			}
		}
	}

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			point := mapGrid.GetPoint(x, y)
			if point == nil {
				continue
			}

			// Get the top-left coordinates of the cell.
			cellTopLeft := m.CellLocations[y*5+x]
			screenX := cellTopLeft.X + m.TopLeft.X
			screenY := cellTopLeft.Y + m.TopLeft.Y

			// Draw the Point image (24x24, center of the cell).
			imageX := screenX + (m.CellSize.X-24)/2
			imageY := screenY + (m.CellSize.Y-24)/2 - 10 // Consider the space for characters.
			interactive := m.GameState.CanInteract(x, y)
			m.drawPointImage(screen, imageX, imageY, point, interactive)

			// Draw the Point name (below the Point image).
			textX := screenX + 8
			textY := imageY + 24 + 2
			pointName := m.getPointName(x, y, point)

			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(textX, textY)
			drawing.DrawText(screen, pointName, 12, opt)

			// If not controlled, draw the enemy's power.
			if p, ok := point.(*core.WildernessPoint); ok && !p.Controlled && interactive {
				power := p.Enemy.Power
				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(imageX, imageY+8)
				powerIcon := drawing.Image("ui-power")
				screen.DrawImage(powerIcon, opt)
				opt.GeoM.Translate(16, 0)
				drawing.DrawText(screen, fmt.Sprintf("%.1f", power), 12, opt)
			}

			if p, ok := point.(*core.BossPoint); ok && interactive {
				power := p.Boss.Power
				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(imageX, imageY+8)
				powerIcon := drawing.Image("ui-power")
				screen.DrawImage(powerIcon, opt)
				opt.GeoM.Translate(16, 0)
				drawing.DrawText(screen, fmt.Sprintf("%.1f", power), 12, opt)
			}
		}
	}
}

// drawPointImage draws the image of the Point.
func (m *MapGridView) drawPointImage(screen *ebiten.Image, x, y float64, point core.Point, interactive bool) {
	switch typedPoint := point.(type) {
	case *core.MyNationPoint:
		pointImg := drawing.Image("point-mynation")
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(x, y)
		screen.DrawImage(pointImg, opt)
	case *core.OtherNationPoint:
		pointImg := drawing.Image("point-othernation")
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(x, y)
		if !interactive {
			opt.ColorScale.Scale(0.5, 0.5, 0.5, 1)
		}
		screen.DrawImage(pointImg, opt)
	case *core.WildernessPoint:
		pointImg := drawing.Image(typedPoint.TerrainType)
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(x, y)
		if !interactive {
			opt.ColorScale.Scale(0.5, 0.5, 0.5, 1)
		}
		screen.DrawImage(pointImg, opt)
	case *core.BossPoint:
		pointImg := drawing.Image("point-boss")
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(x, y)
		if !interactive {
			opt.ColorScale.Scale(0.5, 0.5, 0.5, 1)
		}
		screen.DrawImage(pointImg, opt)
	default:
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Scale(24, 24)
		opt.GeoM.Translate(x, y)
		screen.DrawImage(drawing.WhitePixel, opt)
	}
}

// getPointName gets the name of the Point.
func (m *MapGridView) getPointName(x, y int, point core.Point) string {
	switch p := point.(type) {
	case *core.MyNationPoint:
		return lang.Text("nation-mynation")
	case *core.OtherNationPoint:
		return lang.Text(string(p.OtherNation.NationID))
	case *core.WildernessPoint:
		return lang.Text(p.TerrainType)
	case *core.BossPoint:
		return lang.Text("point-boss")
	default:
		return lang.ExecuteTemplate("point-area", map[string]any{"x": x, "y": y})
	}
}

// drawConnectionLines draws lines to reachable Points.
func (m *MapGridView) drawConnectionLines(screen *ebiten.Image, x, y int, centerX, centerY float64) {
	// Check the 4 adjacent directions.
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for _, dir := range directions {
		nextX, nextY := x+dir[0], y+dir[1]

		// Check if within range.
		if nextX < 0 || nextX >= 5 || nextY < 0 || nextY >= 5 {
			continue
		}

		// 隣接Pointが到達可能かチェック
		if m.GameState.CanInteract(nextX, nextY) {
			// 線を描画
			nextCellTopLeft := m.CellLocations[nextY*5+nextX]
			nextScreenX := nextCellTopLeft.X + m.TopLeft.X
			nextScreenY := nextCellTopLeft.Y + m.TopLeft.Y
			nextCenterX := nextScreenX + m.CellSize.X/2
			nextCenterY := nextScreenY + m.CellSize.Y/2

			m.drawLine(screen, centerX, centerY, nextCenterX, nextCenterY)
		}
	}
}

// drawLine draws a line between two points.
func (m *MapGridView) drawLine(screen *ebiten.Image, x1, y1, x2, y2 float64) {
	vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 2, color.White, true)
}
