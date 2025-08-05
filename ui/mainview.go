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

	m.Market = NewMarketView(flow.NewMarketFlow(gameState, intner), viewmodel.NewMarketViewModel(gameState))
	m.Battle = NewBattleView(onBack)
	m.Territory = NewTerritoryView(onBack)

	// Set GameState for each View.
	m.Battle.SetGameState(gameState)
	m.Territory.SetGameState(gameState)

	m.MapGrid = NewMapGridView(mapGridViewModel, mapGridFlow, func(point core.Point) {
		m.SetSelectedPoint(point)

		// Notify parent about point selection
		if m.OnPointSelected != nil {
			m.OnPointSelected(point)
		}

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
	switch m.CurrentView {
	case ViewTypeMarket:
		m.Market.SetNation(nation)
	}
}

// SetSelectedPoint sets the point to be displayed
func (m *MainView) SetSelectedPoint(point core.Point) {
	switch p := point.(type) {
	case *core.WildernessPoint:
		if p.Controlled() {
			// Set territory for TerritoryView
			territory := p.Territory()
			if territory != nil {
				m.Territory.SetTerritory(territory, "wilderness") // terrain type placeholder
			}
		} else {
			// Set battle point for BattleView
			m.Battle.SetBattlePoint(p)
		}
	case *core.BossPoint:
		// Set battle point for BattleView
		m.Battle.SetBattlePoint(p)
	}
}
