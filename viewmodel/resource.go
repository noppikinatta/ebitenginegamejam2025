package viewmodel

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// ResourceViewModel provides display information for resource UI
type ResourceViewModel struct {
	gameState *core.GameState
}

// NewResourceViewModel creates a new ResourceViewModel
func NewResourceViewModel(gameState *core.GameState) *ResourceViewModel {
	return &ResourceViewModel{
		gameState: gameState,
	}
}

// Quantity returns the current resource quantity
func (vm *ResourceViewModel) Quantity() core.ResourceQuantity {
	if vm.gameState == nil || vm.gameState.Treasury == nil {
		return core.ResourceQuantity{}
	}
	return vm.gameState.Treasury.Resources
}

// CalendarViewModel provides display information for calendar UI
type CalendarViewModel struct {
	gameState *core.GameState
}

// NewCalendarViewModel creates a new CalendarViewModel
func NewCalendarViewModel(gameState *core.GameState) *CalendarViewModel {
	return &CalendarViewModel{
		gameState: gameState,
	}
}

// YearMonth returns the formatted year and month string
func (vm *CalendarViewModel) YearMonth() string {
	if vm.gameState == nil {
		return ""
	}

	// Extract year and month from current turn
	// This depends on how Turn is implemented in core
	turn := vm.gameState.CurrentTurn

	// Convert turn to year/month based on game logic
	// For now, assume simple conversion
	year := int(turn)/12 + 1  // Start from year 1
	month := int(turn)%12 + 1 // 1-12

	// Format year/month using template
	return lang.ExecuteTemplate("year_month", map[string]any{
		"year":  year,
		"month": month,
	})
}
