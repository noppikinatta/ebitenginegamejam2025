package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/flow"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// GameUI is the controller that manages the entire game UI.
// It integrates all Widgets and is responsible for coordination between them.
type GameUI struct {
	// Widgets
	ResourceView *ResourceView
	CalendarView *CalendarView
	MainView     *MainView
	InfoView     *InfoView
	CardDeckView *CardDeckView

	// GameState reference (for creating flows and viewmodels)
	GameState *core.GameState

	// Flow instances for operations
	CardDeckFlow *flow.CardDeckFlow
	MapGridFlow  *flow.MapGridFlow

	// ViewModel instances for display
	ResourceViewModel *viewmodel.ResourceViewModel
	CalendarViewModel *viewmodel.CalendarViewModel
	CardDeckViewModel *viewmodel.CardDeckViewModel
	MapGridViewModel  *viewmodel.MapGridViewModel

	// Mouse position tracking
	MouseX, MouseY int
}

// NewGameUI creates a GameUI.
func NewGameUI(gameState *core.GameState) *GameUI {
	// Create flow instances
	cardDeckFlow := flow.NewCardDeckFlow(gameState)
	mapGridFlow := flow.NewMapGridFlow(gameState)

	// Create viewmodel instances
	resourceViewModel := viewmodel.NewResourceViewModel(gameState)
	calendarViewModel := viewmodel.NewCalendarViewModel(gameState)
	cardDeckViewModel := viewmodel.NewCardDeckViewModel(gameState)
	mapGridViewModel := viewmodel.NewMapGridViewModel(gameState)

	// Initialize each Widget with viewmodels
	resourceView := NewResourceView(resourceViewModel)
	calendarView := NewCalendarView(calendarViewModel)
	mainView := NewMainView(gameState, mapGridViewModel, mapGridFlow)
	infoView := NewInfoView(viewmodel.NewHistoryViewModel(gameState))

	cardDeckView := NewCardDeckView(mainView, cardDeckViewModel, cardDeckFlow)

	ui := &GameUI{
		ResourceView: resourceView,
		CalendarView: calendarView,
		MainView:     mainView,
		InfoView:     infoView,
		CardDeckView: cardDeckView,
		GameState:    gameState,

		// Store flows and viewmodels
		CardDeckFlow:      cardDeckFlow,
		MapGridFlow:       mapGridFlow,
		ResourceViewModel: resourceViewModel,
		CalendarViewModel: calendarViewModel,
		CardDeckViewModel: cardDeckViewModel,
		MapGridViewModel:  mapGridViewModel,
	}

	return ui
}

// HandleInput handles input for all Widgets.
func (gui *GameUI) HandleInput(input *Input) error {
	// Update mouse position
	gui.MouseX, gui.MouseY = input.Mouse.CursorPosition()

	// MainView processes input first (important processes such as View switching).
	if err := gui.MainView.HandleInput(input); err != nil {
		return err
	}

	// Handle input for CardDeckView
	if err := gui.CardDeckView.HandleInput(input); err != nil {
		return err
	}

	// Other Widgets are for display only and do not handle input.
	// ResourceView, CalendarView, and InfoView are basically for display only.

	return nil
}

// Update handles frame updates.
func (gui *GameUI) Update() error {
	return nil
}

// Draw draws all Widgets.
func (gui *GameUI) Draw(screen *ebiten.Image) {
	// Drawing order: from background to foreground

	// 1. ResourceView (top left)
	gui.ResourceView.Draw(screen)

	// 2. CalendarView (top right)
	gui.CalendarView.Draw(screen)

	// 3. MainView (center main)
	gui.MainView.Draw(screen)

	// 4. InfoView (right side info)
	gui.InfoView.Draw(screen)

	// 5. CardDeckView (bottom card deck)
	gui.CardDeckView.Draw(screen)
}

// GetCurrentMainViewType gets the current MainViewType.
func (gui *GameUI) GetCurrentMainViewType() ViewType {
	return gui.MainView.CurrentView
}

// SwitchMainView switches the MainView.
func (gui *GameUI) SwitchMainView(viewType ViewType) {
	gui.MainView.SwitchView(viewType)

	// Update InfoView when switching views
	switch viewType {
	case ViewTypeMapGrid:
		gui.InfoView.CurrentMode = InfoModeHistory
	case ViewTypeMarket:
		// Do nothing special for Market (updated when a card is selected)
	case ViewTypeBattle:
		// For Battle, may display Enemy information
	case ViewTypeTerritory:
		// For Territory, display Point information
	}
}
