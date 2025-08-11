package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/flow"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// ViewType is the type of View to be displayed in MainView.
type ViewType int

const (
	ViewTypeMapGrid ViewType = iota
	ViewTypeMarket
	ViewTypeBattle
	ViewTypeTerritory
)

type CenterViewModer interface {
	CurrentViewMode() ViewType
}

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
func NewMainView(gameState *core.GameState, intner core.Intner) *MainView {
	m := &MainView{
		CurrentView: ViewTypeMapGrid, // The initial display is MapGridView.
		GameState:   gameState,
	}

	// Construct child views
	m.Market = NewMarketView(flow.NewMarketFlow(gameState, intner), viewmodel.NewMarketViewModel(gameState))
	m.Battle = NewBattleView(flow.NewBattleFlow(gameState))
	m.Territory = NewTerritoryView(flow.NewTerritoryFlow(gameState))

	// No direct GameState injection to views; views use flow/viewmodel

	mapGridViewModel := viewmodel.NewMapGridViewModel(gameState)
	mapGridFlow := flow.NewMapGridFlow(gameState)
	m.MapGrid = NewMapGridView(mapGridViewModel, mapGridFlow, func(point core.Point) {
		m.SetSelectedPoint(point)

		switch p := point.(type) {
		case *core.MyNationPoint:
			x, y, ok := gameState.MapGrid.XYOfPoint(point)
			if ok {
				m.Market.Select(x, y)
			}
			m.SwitchView(ViewTypeMarket)
		case *core.OtherNationPoint:
			x, y, ok := gameState.MapGrid.XYOfPoint(point)
			if ok {
				m.Market.Select(x, y)
			}
			m.SwitchView(ViewTypeMarket)
		case *core.WildernessPoint:
			if p.Controlled() {
				x, y, ok := gameState.MapGrid.XYOfPoint(point)
				if ok {
					m.Territory.Select(x, y)
				}
				m.SwitchView(ViewTypeTerritory)
			} else {
				// Initialize battlefield and bind VM/Flow to Battle view
				x, y, ok := gameState.MapGrid.XYOfPoint(point)
				if ok {
					m.Battle.Select(x, y)
				}
				m.SwitchView(ViewTypeBattle)
			}
		case *core.BossPoint:
			x, y, ok := gameState.MapGrid.XYOfPoint(point)
			if ok {
				m.Battle.Select(x, y)
			}
			m.SwitchView(ViewTypeBattle)
		}
	})

	return m
}

// SwitchView switches the View to be displayed.
func (m *MainView) SwitchView(viewType ViewType) {
	m.CurrentView = viewType
}

func (m *MainView) CurrentViewMode() ViewType {
	return m.CurrentView
}

// HandleInput handles input.
func (m *MainView) HandleInput(input *Input) error {
	shouldBackToMapGrid, err := m.handleChildren(input)
	if err != nil {
		return err
	}

	if shouldBackToMapGrid {
		m.SwitchView(ViewTypeMapGrid)
	}

	return nil
}

func (m *MainView) handleChildren(input *Input) (bool, error) {
	switch m.CurrentView {
	case ViewTypeMapGrid:
		return false, m.MapGrid.HandleInput(input)
	case ViewTypeMarket:
		return m.Market.HandleInput(input)
	case ViewTypeBattle:
		return m.Battle.HandleInput(input)
	case ViewTypeTerritory:
		return m.Territory.HandleInput(input)
	}

	return false, nil
}

// Draw handles drawing.
func (m *MainView) Draw(screen *ebiten.Image) {
	switch m.CurrentView {
	case ViewTypeMapGrid:
		m.MapGrid.Draw(screen)
	case ViewTypeMarket:
		m.Market.Draw(screen)
	case ViewTypeBattle:
		m.Battle.Draw(screen)
	case ViewTypeTerritory:
		m.Territory.Draw(screen)
	}
}

// SetSelectedNation sets the nation to be displayed
func (m *MainView) SetSelectedNation(nation core.Nation) {
	// Deprecated: MarketView now determines nation from Select(x,y) via flow/viewmodel
}

// SetSelectedPoint sets the point to be displayed
func (m *MainView) SetSelectedPoint(point core.Point) {
	switch p := point.(type) {
	case *core.WildernessPoint:
		if p.Controlled() {
			// Territory selection is handled in MapGrid callback using coordinates
		} else {
			// Battle binding is handled in MapGrid callback using coordinates
		}
	case *core.BossPoint:
		// Battle binding is handled in MapGrid callback using coordinates
	}
}
