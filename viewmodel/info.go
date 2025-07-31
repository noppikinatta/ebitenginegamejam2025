package viewmodel

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

type HistoryViewModel struct {
	gameState *core.GameState
	eventMap  map[string]any // Temporary map for template execution.
}

func NewHistoryViewModel(gameState *core.GameState) *HistoryViewModel {
	return &HistoryViewModel{gameState: gameState}
}

func (vm *HistoryViewModel) HistoryLen() int {
	return len(vm.gameState.Histories)
}

func (vm *HistoryViewModel) HistoryDateText(index int) string {
	history := vm.gameState.Histories[index]
	year, month := history.Turn.YearMonth()
	year += 1023
	return lang.ExecuteTemplate("ui-calendar", map[string]any{"year": year, "month": month})
}

func (vm *HistoryViewModel) HistoryEventText(index int) string {
	if vm.eventMap == nil {
		vm.eventMap = make(map[string]any)
	}

	history := vm.gameState.Histories[index]
	for k, v := range history.Data {
		if s, ok := v.(string); ok {
			vm.eventMap[k] = lang.Text(s)
		} else {
			vm.eventMap[k] = v
		}
	}
	return lang.ExecuteTemplate(history.Key, vm.eventMap)
}
