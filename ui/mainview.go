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
	mv := &MainView{
		CurrentView: ViewTypeMapGrid, // 初期表示はMapGridView
		GameState:   gameState,
	}

	onBack := func() {
		mv.SwitchView(ViewTypeMapGrid)
	}

	mv.Market = NewMarketView(onBack)
	mv.Battle = NewBattleView(onBack)
	mv.Territory = NewTerritoryView(onBack)

	mv.MapGrid = NewMapGridView(gameState, func(point core.Point) {
		mv.SetSelectedPoint(point)
		switch p := point.(type) {
		case *core.MyNationPoint:
			mv.SetSelectedNation(p.MyNation)
			mv.SwitchView(ViewTypeMarket)
		case *core.OtherNationPoint:
			mv.SetSelectedNation(p.OtherNation)
			mv.SwitchView(ViewTypeMarket)
		case *core.WildernessPoint:
			if p.Controlled {
				mv.SwitchView(ViewTypeTerritory)
			} else {
				mv.SwitchView(ViewTypeBattle)
			}
		case *core.BossPoint:
			mv.SwitchView(ViewTypeBattle)
		}
	})

	return mv
}

// SwitchView 表示するViewを切り替える
func (mv *MainView) SwitchView(viewType ViewType) {
	mv.CurrentView = viewType
}

// HandleInput 入力処理
func (mv *MainView) HandleInput(input *Input) error {
	// 現在のViewに入力を転送
	switch mv.CurrentView {
	case ViewTypeMapGrid:
		if mv.MapGrid != nil {
			return mv.MapGrid.HandleInput(input)
		}
	case ViewTypeMarket:
		if mv.Market != nil {
			return mv.Market.HandleInput(input)
		}
	case ViewTypeBattle:
		if mv.Battle != nil {
			return mv.Battle.HandleInput(input)
		}
	case ViewTypeTerritory:
		if mv.Territory != nil {
			return mv.Territory.HandleInput(input)
		}
	}
	return nil
}

// Draw 描画処理
func (mv *MainView) Draw(screen *ebiten.Image) {
	// 現在のViewを描画
	switch mv.CurrentView {
	case ViewTypeMapGrid:
		if mv.MapGrid != nil {
			mv.MapGrid.Draw(screen)
		}
	case ViewTypeMarket:
		if mv.Market != nil {
			mv.Market.Draw(screen)
		}
	case ViewTypeBattle:
		if mv.Battle != nil {
			mv.Battle.Draw(screen)
		}
	case ViewTypeTerritory:
		if mv.Territory != nil {
			mv.Territory.Draw(screen)
		}
	}
}

// GetCurrentView 現在表示中のViewタイプを取得
func (mv *MainView) GetCurrentView() ViewType {
	return mv.CurrentView
}

// SetSelectedNation MarketViewで表示する国家を設定
func (mv *MainView) SetSelectedNation(nation interface{}) {
	if mv.Market != nil {
		mv.Market.SetNation(nation)
	}
}

// SetSelectedPoint BattleViewやTerritoryViewで表示するPointを設定
func (mv *MainView) SetSelectedPoint(point core.Point) {
	switch p := point.(type) {
	case *core.WildernessPoint:
		if p.Controlled {
			if mv.Territory != nil {
				mv.Territory.SetTerritory(p.Territory)
			}
		} else {
			if mv.Battle != nil {
				mv.Battle.SetEnemy(p.Enemy)
			}
		}
	case *core.BossPoint:
		if mv.Battle != nil {
			mv.Battle.SetEnemy(p.Boss)
		}
	}
}
