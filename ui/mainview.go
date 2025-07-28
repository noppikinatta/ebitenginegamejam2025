package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// ViewType is the type of View to be displayed in MainView.
type ViewType int

const (
	ViewTypeMapGrid ViewType = iota
	ViewTypeMarket
	ViewTypeBattle
	ViewTypeTerritory
)

// MainView is the main view container Widget.
// Position: (0,40,1040,560).
// Switches between MapGridView, MarketView, BattleView, and TerritoryView.
type MainView struct {
	CurrentView ViewType
	MapGrid     *MapGridView
	Market      *MarketView
	Battle      *BattleView
	Territory   *TerritoryView

	// Game state
	GameState *core.GameState
}

// NewMainView creates a MainView.
func NewMainView(gameState *core.GameState) *MainView {
	m := &MainView{
		CurrentView: ViewTypeMapGrid, // The initial display is MapGridView.
		GameState:   gameState,
	}

	onBack := func() {
		m.SwitchView(ViewTypeMapGrid)
	}

	m.Market = NewMarketView(onBack)
	m.Battle = NewBattleView(onBack)
	m.Territory = NewTerritoryView(onBack)

	// Set GameState for each View.
	m.Market.SetGameState(gameState)
	m.Battle.SetGameState(gameState)
	m.Territory.SetGameState(gameState)

	m.MapGrid = NewMapGridView(gameState, func(point core.Point) {
		m.SetSelectedPoint(point)
		switch p := point.(type) {
		case *core.MyNationPoint:
			m.SetSelectedNation(p.MyNation)
			m.SwitchView(ViewTypeMarket)
		case *core.OtherNationPoint:
			m.SetSelectedNation(p.OtherNation)
			m.SwitchView(ViewTypeMarket)
		case *core.WildernessPoint:
			if p.Controlled() {
				m.SwitchView(ViewTypeTerritory)
			} else {
				m.SwitchView(ViewTypeBattle)
			}
		case *core.BossPoint:
			m.SwitchView(ViewTypeBattle)
		}
	})

	return m
}

// SwitchView switches the View to be displayed.
func (m *MainView) SwitchView(viewType ViewType) {
	m.CurrentView = viewType
}

// HandleInput handles input.
func (m *MainView) HandleInput(input *Input) error {
	// Forward input to the current View.
	switch m.CurrentView {
	case ViewTypeMapGrid:
		if m.MapGrid != nil {
			return m.MapGrid.HandleInput(input)
		}
	case ViewTypeMarket:
		if m.Market != nil {
			return m.Market.HandleInput(input)
		}
	case ViewTypeBattle:
		if m.Battle != nil {
			return m.Battle.HandleInput(input)
		}
	case ViewTypeTerritory:
		if m.Territory != nil {
			return m.Territory.HandleInput(input)
		}
	}
	return nil
}

// Draw handles drawing.
func (m *MainView) Draw(screen *ebiten.Image) {
	// Draw the current View.
	switch m.CurrentView {
	case ViewTypeMapGrid:
		if m.MapGrid != nil {
			m.MapGrid.Draw(screen)
		}
	case ViewTypeMarket:
		if m.Market != nil {
			m.Market.Draw(screen)
		}
	case ViewTypeBattle:
		if m.Battle != nil {
			m.Battle.Draw(screen)
		}
	case ViewTypeTerritory:
		if m.Territory != nil {
			m.Territory.Draw(screen)
		}
	}
}

// GetCurrentView gets the currently displayed View type.
func (m *MainView) GetCurrentView() ViewType {
	return m.CurrentView
}

// SetSelectedNation sets the nation to be displayed in MarketView.
func (m *MainView) SetSelectedNation(nation core.Nation) {
	if m.Market != nil {
		m.Market.SetNation(nation)
	}
}

// SetSelectedPoint sets the Point to be displayed in BattleView or TerritoryView.
func (m *MainView) SetSelectedPoint(point core.Point) {
	switch p := point.(type) {
	case *core.WildernessPoint:
		if p.Controlled() {
			if m.Territory != nil {
				territory := p.Territory()
				terrainType := ""
				if territory != nil && territory.Terrain() != nil {
					terrainType = string(territory.Terrain().ID())
				}
				m.Territory.SetTerritory(territory, terrainType)
			}
		} else {
			if m.Battle != nil {
				m.Battle.SetBattlePoint(p)
			}
		}
	case *core.BossPoint:
		if m.Battle != nil {
			m.Battle.SetBattlePoint(p)
		}
	}
}
