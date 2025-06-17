package component_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

// TestCalendarShowsCorrectYearMonthFormat validates starting date and display string
func TestCalendarShowsCorrectYearMonthFormat(t *testing.T) {
    g := scene.CreateSequence()
    g.SetCurrentScene("ingame")
    cal := g.GetInGameScene().GetCalendar()
    if cal == nil {
        t.Fatal("calendar nil")
    }
    if y := cal.GetCurrentYear(); y != 1000 {
        t.Errorf("expected year 1000 got %d", y)
    }
    if m := cal.GetCurrentMonth(); m != 4 {
        t.Errorf("expected month 4 got %d", m)
    }
    if cal.GetDisplayText() == "" {
        t.Error("display text empty")
    }
} 