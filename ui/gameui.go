package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
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

	// GameState reference
	GameState *core.GameState

	// Mouse position tracking
	MouseX, MouseY int
}

// NewGameUI creates a GameUI.
func NewGameUI(gameState *core.GameState) *GameUI {
	// Initialize each Widget
	resourceView := NewResourceView(gameState)
	calendarView := NewCalendarView(gameState)
	mainView := NewMainView(gameState)
	infoView := NewInfoView()

	cardDeckView := NewCardDeckView(gameState.CardDeck)

	ui := &GameUI{
		ResourceView: resourceView,
		CalendarView: calendarView,
		MainView:     mainView,
		InfoView:     infoView,
		CardDeckView: cardDeckView,
		GameState:    gameState,
	}

	// Set up coordination between Widgets
	ui.setupWidgetConnections()
	ui.CardDeckView.OnBattleCardClicked = ui.onBattleCardClicked
	ui.CardDeckView.OnStructureCardClicked = ui.onStructureCardClicked

	return ui
}

// onBattleCardClicked handles the click on a BattleCard.
func (gui *GameUI) onBattleCardClicked(card *core.BattleCard) bool {
	if gui.MainView.CurrentView != ViewTypeBattle {
		return false
	}

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

	if gui.MainView.Territory.CanPlaceCard() {
		return gui.MainView.Territory.PlaceCard(card)
	}

	return false
}

// setupWidgetConnections sets up coordination between Widgets.
func (gui *GameUI) setupWidgetConnections() {
	// Coordination from CardDeckView to InfoView
	gui.CardDeckView.OnCardSelected = func(card interface{}) {
		gui.InfoView.SetSelectedCard(card)
	}

	// Coordination from MainView to InfoView is handled dynamically in Update()
}

// AddHistoryEvent adds a history event.
func (gui *GameUI) AddHistoryEvent(event string) {
	gui.InfoView.AddHistoryEvent(event)
}

// SetMousePosition updates the mouse position.
func (gui *GameUI) SetMousePosition(x, y int) {
	gui.MouseX = x
	gui.MouseY = y

	// Updating the mouse position for each Widget will be implemented in the future.
	// Omitted for now as Widgets do not currently have MouseX/MouseY fields.
}

// HandleInput handles unified input.
func (gui *GameUI) HandleInput(input *Input) error {
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

// SelectPoint selects a Point and reflects it in the InfoView.
func (gui *GameUI) SelectPoint(point core.Point) {
	gui.InfoView.SetSelectedPoint(point)
}

// SelectCard selects a card and reflects it in the InfoView.
func (gui *GameUI) SelectCard(card interface{}) {
	gui.InfoView.SetSelectedCard(card)
}

// SelectCardFromDeck selects a specific card from the CardDeck.
func (gui *GameUI) SelectCardFromDeck(index int) {
	gui.CardDeckView.SelectCard(index)
}

// MoveCardToTerritory moves the selected card to the TerritoryView.
func (gui *GameUI) MoveCardToTerritory() bool {
	selectedCard := gui.CardDeckView.GetSelectedCard()
	if selectedCard == nil {
		return false
	}

	// Only StructureCards can be moved
	if structureCard, ok := selectedCard.(*core.StructureCard); ok {
		// Remove from CardDeck
		gui.CardDeckView.RemoveSelectedCard()

		// Add to TerritoryView
		gui.MainView.Territory.AddStructureCard(structureCard)

		gui.AddHistoryEvent("Card moved to territory")
		return true
	}

	return false
}

// ReturnCardToDeck returns a card to the CardDeck.
func (gui *GameUI) ReturnCardToDeck(card interface{}) {
	gui.CardDeckView.AddCard(card)
	gui.AddHistoryEvent("Card returned to deck")
}
