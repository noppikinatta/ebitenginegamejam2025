package system

import (
	"fmt"
)

type TurnManager struct {
	currentYear  int
	currentMonth int
	turnCount    int
}

func NewTurnManager() *TurnManager {
	return &TurnManager{
		currentYear:  1000, // Kingdom Year 1000
		currentMonth: 4,    // Month 4
		turnCount:    0,
	}
}

func (tm *TurnManager) GetCurrentYear() int {
	return tm.currentYear
}

func (tm *TurnManager) GetCurrentMonth() int {
	return tm.currentMonth
}

func (tm *TurnManager) GetTurnCount() int {
	return tm.turnCount
}

func (tm *TurnManager) AdvanceTurn() {
	tm.turnCount++
	tm.currentMonth++

	if tm.currentMonth > 12 {
		tm.currentMonth = 1
		tm.currentYear++
	}
}

func (tm *TurnManager) GetDisplayText() string {
	return fmt.Sprintf("Kingdom Year %d, Month %d", tm.currentYear, tm.currentMonth)
}
