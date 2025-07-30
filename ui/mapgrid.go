package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/flow"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// MapGridView is a Widget for displaying the MapGrid.
// Position: Drawn within MainView
type MapGridView struct {
	ViewModel *viewmodel.MapGridViewModel
	Flow      *flow.MapGridFlow

	// Callbacks
	OnPointClicked func(point core.Point)
}

// NewMapGridView creates a MapGridView
func NewMapGridView(viewModel *viewmodel.MapGridViewModel, flow *flow.MapGridFlow, onPointClicked func(point core.Point)) *MapGridView {
	return &MapGridView{
		ViewModel:      viewModel,
		Flow:           flow,
		OnPointClicked: onPointClicked,
	}
}

// HandleInput processes input
func (m *MapGridView) HandleInput(input *Input) error {
	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		cursorX, cursorY := input.Mouse.CursorPosition()

		// Calculate grid coordinates from cursor position
		x, y := m.getGridCoordinates(cursorX, cursorY)
		if x >= 0 && y >= 0 {
			// Select point using flow
			if m.Flow.SelectPoint(x, y) {
				// Get selected point and notify callback
				selectedPoint := m.Flow.GetSelectedPoint()
				if selectedPoint != nil && m.OnPointClicked != nil {
					m.OnPointClicked(selectedPoint)
				}
			}
		}
	}
	return nil
}

// getGridCoordinates converts screen coordinates to grid coordinates
func (m *MapGridView) getGridCoordinates(screenX, screenY int) (int, int) {
	// This is a simplified implementation
	// The actual conversion depends on the specific layout of the map grid

	// Main view area: (0,40,1040,560)
	if screenX < 0 || screenX >= 1040 || screenY < 40 || screenY >= 600 {
		return -1, -1
	}

	// Convert to grid coordinates (this is a placeholder implementation)
	gridX := screenX / 100        // assuming each grid cell is 100px wide
	gridY := (screenY - 40) / 100 // assuming each grid cell is 100px tall

	size := m.ViewModel.Size()
	if gridX >= size.X || gridY >= size.Y {
		return -1, -1
	}

	return gridX, gridY
}

// Draw handles the drawing process
func (m *MapGridView) Draw(screen *ebiten.Image) {
	size := m.ViewModel.Size()

	// Draw grid points
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			pointVM := m.ViewModel.Point(x, y)
			if pointVM != nil {
				m.drawPoint(screen, x, y, pointVM)
			}
		}
	}

	// Draw connections between points
	m.drawConnections(screen, size)
}

// drawPoint draws a single point on the grid
func (m *MapGridView) drawPoint(screen *ebiten.Image, x, y int, pointVM *viewmodel.PointViewModel) {
	// Calculate screen position
	screenX := float64(x*100 + 50) // Center of grid cell
	screenY := float64(y*100 + 90) // Center of grid cell (offset by 40 for main view)

	// Draw point image
	image := pointVM.Image()
	if image != nil {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(screenX-16, screenY-16) // Center the 32x32 image
		screen.DrawImage(image, opt)
	}

	// Draw point name
	name := pointVM.Name()
	if name != "" {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(screenX-20, screenY+20)
		drawing.DrawText(screen, name, 12, opt)
	}

	// Draw enemy power if applicable
	if pointVM.HasEnemy() {
		power := pointVM.EnemyPower()
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(screenX+15, screenY-15)
		drawing.DrawText(screen, fmt.Sprintf("%.0f", power), 10, opt)
	}
}

// drawConnections draws lines connecting adjacent points
func (m *MapGridView) drawConnections(screen *ebiten.Image, size core.MapGridSize) {
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			// Draw line to the right
			if m.ViewModel.ShouldDrawLineToRight(x, y) {
				startX := float64(x*100 + 66)
				startY := float64(y*100 + 90)
				endX := float64((x+1)*100 + 34)
				endY := float64(y*100 + 90)
				m.drawLine(screen, startX, startY, endX, endY)
			}

			// Draw line upward
			if m.ViewModel.ShouldDrawLineToUpper(x, y) {
				startX := float64(x*100 + 50)
				startY := float64(y*100 + 74)
				endX := float64(x*100 + 50)
				endY := float64((y-1)*100 + 106)
				m.drawLine(screen, startX, startY, endX, endY)
			}
		}
	}
}

// drawLine draws a line between two points
func (m *MapGridView) drawLine(screen *ebiten.Image, x1, y1, x2, y2 float64) {
	// Simple line drawing using rectangles
	length := ((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
	if length > 0 {
		drawing.DrawRect(screen, x1, y1, x2-x1, 2, 0.5, 0.5, 0.5, 1.0)
	}
}
