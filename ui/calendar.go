package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
)

// CalendarView is a widget for displaying the calendar.
// Position: (1040,0,240,40).
// Displays the current Turn in year/month format.
type CalendarView struct {
	GameState *core.GameState
}

// NewCalendarView creates a CalendarView.
func NewCalendarView(gameState *core.GameState) *CalendarView {
	return &CalendarView{
		GameState: gameState,
	}
}

// HandleInput handles input (CalendarView does not accept input).
func (cv *CalendarView) HandleInput(input *Input) error {
	return nil
}

// Draw handles drawing.
func (cv *CalendarView) Draw(screen *ebiten.Image) {
	// Get year and month from Turn.
	turn := core.Turn(cv.GameState.CurrentTurn)
	year, month := turn.YearMonth()

	// Display in YYYY/MM format.
	text := fmt.Sprintf("%04d/%02d", year, month)

	// Position: (1040,0,240,40).
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(1040, 0)
	drawing.DrawText(screen, text, 24, opt)
}
