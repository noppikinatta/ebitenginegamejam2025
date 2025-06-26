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
	CardDeck    *CardDeckView

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

	// CardDeckViewを作成し、コールバックを設定
	m.CardDeck = NewCardDeckView(gameState.CardDeck)
	m.CardDeck.OnBattleCardClicked = m.onBattleCardClicked
	m.CardDeck.OnStructureCardClicked = m.onStructureCardClicked

	// 各ViewにGameStateを設定
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

// onBattleCardClicked BattleCardクリック時の処理
func (m *MainView) onBattleCardClicked(card *core.BattleCard) bool {
	if m.CurrentView != ViewTypeBattle {
		return false
	}

	if m.Battle.CanPlaceCard() {
		return m.Battle.PlaceCard(card)
	}

	return false
}

// onStructureCardClicked StructureCardクリック時の処理
func (m *MainView) onStructureCardClicked(card *core.StructureCard) bool {
	if m.CurrentView != ViewTypeTerritory {
		return false
	}

	if m.Territory.CanPlaceCard() {
		return m.Territory.PlaceCard(card)
	}

	return false
}

// SwitchView 表示するViewを切り替える
func (m *MainView) SwitchView(viewType ViewType) {
	m.CurrentView = viewType
}

// HandleInput 入力処理
func (m *MainView) HandleInput(input *Input) error {
	// CardDeckViewの入力処理（常に有効）
	if m.CardDeck != nil {
		if err := m.CardDeck.HandleInput(input); err != nil {
			return err
		}
	}

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

	// CardDeckViewを常に描画
	if m.CardDeck != nil {
		m.CardDeck.Draw(screen)
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
				// TODO: 実際の座標を渡すようにMapGridViewを改修する必要がある
				m.Battle.SetEnemy(p.Enemy, 0, 0) // 暫定的に0,0を設定
			}
		}
	case *core.BossPoint:
		if m.Battle != nil {
			// TODO: 実際の座標を渡すようにMapGridViewを改修する必要がある
			m.Battle.SetEnemy(p.Boss, 0, 0) // 暫定的に0,0を設定
		}
	}
}
