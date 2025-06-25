package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// ViewType MainViewで表示するViewの種類
type ViewType int

const (
	ViewTypeMapGrid ViewType = iota
	ViewTypeMarket
	ViewTypeBattle
	ViewTypeTerritory
)

// MainView メインビューコンテナWidget
// 位置: (0,20,520,280)
// MapGridView, MarketView, BattleView, TerritoryViewを切り替える
type MainView struct {
	CurrentView ViewType
	MapGrid     *MapGridView
	Market      *MarketView
	Battle      *BattleView
	Territory   *TerritoryView

	// ゲーム状態
	GameState *core.GameState
}

// NewMainView MainViewを作成する
func NewMainView(gameState *core.GameState) *MainView {
	m := &MainView{
		CurrentView: ViewTypeMapGrid, // 初期表示はMapGridView
		GameState:   gameState,
	}

	onBack := func() {
		m.SwitchView(ViewTypeMapGrid)
	}

	m.Market = NewMarketView(onBack)
	m.Battle = NewBattleView(onBack)
	m.Territory = NewTerritoryView(onBack)

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
			if p.Controlled {
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

// SwitchView 表示するViewを切り替える
func (m *MainView) SwitchView(viewType ViewType) {
	m.CurrentView = viewType
}

// HandleInput 入力処理
func (m *MainView) HandleInput(input *Input) error {
	// 現在のViewに入力を転送
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

// Draw 描画処理
func (m *MainView) Draw(screen *ebiten.Image) {
	// 現在のViewを描画
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

// GetCurrentView 現在表示中のViewタイプを取得
func (m *MainView) GetCurrentView() ViewType {
	return m.CurrentView
}

// SetSelectedNation MarketViewで表示する国家を設定
func (m *MainView) SetSelectedNation(nation core.Nation) {
	if m.Market != nil {
		m.Market.SetNation(nation)
	}
}

// SetSelectedPoint BattleViewやTerritoryViewで表示するPointを設定
func (m *MainView) SetSelectedPoint(point core.Point) {
	switch p := point.(type) {
	case *core.WildernessPoint:
		if p.Controlled {
			if m.Territory != nil {
				m.Territory.SetTerritory(p.Territory)
			}
		} else {
			if m.Battle != nil {
				m.Battle.SetEnemy(p.Enemy)
			}
		}
	case *core.BossPoint:
		if m.Battle != nil {
			m.Battle.SetEnemy(p.Boss)
		}
	}
}
