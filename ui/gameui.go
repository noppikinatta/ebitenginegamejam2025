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
	cardDeckFlow := flow.NewCardDeckFlow(gameState, gameState.CardDictionary)
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

	// Set up coordination between Widgets
	ui.CardDeckView.OnBattleCardClicked = ui.onBattleCardClicked
	ui.CardDeckView.OnStructureCardClicked = ui.onStructureCardClicked

	return ui
}

// onBattleCardClicked handles the click on a BattleCard.
func (gui *GameUI) onBattleCardClicked(card *core.BattleCard) bool {
	if gui.MainView.CurrentView != ViewTypeBattle {
		return false
	}

	// Use BattleFlow through MainView.Battle
	if gui.MainView.Battle.CanPlaceCard() {
		return gui.MainView.Battle.PlaceCard(card)
	}

	return false
}

// onStructureCardClicked handles the click on a StructureCard.
func (gui *GameUI) onStructureCardClicked(card *core.StructureCard) bool {
	if gui.MainView.CurrentView != ViewTypeTerritory {
		return false
	}

	// Use TerritoryFlow through MainView.Territory
	if gui.MainView.Territory.CanPlaceCard() {
		return gui.MainView.Territory.PlaceCard(card)
	}

	return false
}

// HandleInput handles input for all Widgets.
func (gui *GameUI) HandleInput(input *Input) error {
	// Update mouse position
	gui.MouseX, gui.MouseY = input.Mouse.CursorPosition()

	// Pass mouse position to child Widgets
	gui.CardDeckView.MouseX = gui.MouseX
	gui.CardDeckView.MouseY = gui.MouseY

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

// SelectCardFromDeck selects a specific card from the CardDeck.
func (gui *GameUI) SelectCardFromDeck(index int) {
	// Use CardDeckFlow for selection
	gui.CardDeckFlow.Select(index)

	// Update CardDeckView to reflect the selection
	gui.CardDeckView.SetSelectedIndex(index)
}

// MoveCardToTerritory moves the selected card to the TerritoryView.
func (gui *GameUI) MoveCardToTerritory() bool {
	selectedCard := gui.CardDeckFlow.GetSelectedCard()
	if selectedCard == nil {
		return false
	}

	// Only StructureCards can be moved
	if structureCard, ok := selectedCard.(*core.StructureCard); ok {
		// This will be handled by TerritoryFlow in the TerritoryView
		return gui.MainView.Territory.PlaceCard(structureCard)
	}

	return false
}

// ReturnCardToDeck returns a card to the CardDeck.
func (gui *GameUI) ReturnCardToDeck(card interface{}) {
	// This operation will be handled by the appropriate flow
	// For now, we'll delegate to CardDeckView to handle the UI update
	gui.CardDeckView.AddCard(card)
}
