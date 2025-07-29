package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// CalendarView is a widget for displaying the calendar.
// Position: (1040,0,240,40).
// Displays the current Turn in year/month format.
type CalendarView struct {
	ViewModel *viewmodel.CalendarViewModel
}

// NewCalendarView creates a CalendarView.
func NewCalendarView(viewModel *viewmodel.CalendarViewModel) *CalendarView {
	return &CalendarView{
		ViewModel: viewModel,
	}
}

// HandleInput handles input (CalendarView does not accept input).
func (cv *CalendarView) HandleInput(input *Input) error {
	return nil
}

// Draw handles drawing.
func (cv *CalendarView) Draw(screen *ebiten.Image) {
	// Get formatted year/month from viewmodel
	text := cv.ViewModel.YearMonth()

	// Position: (1040,0,240,40).
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(1040, 0)
	drawing.DrawText(screen, text, 24, opt)
}
