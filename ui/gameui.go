package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// GameUI ゲーム全体のUIを管理するコントローラ
// 全Widget を統合し、Widget間の連携を担当
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

// NewGameUI GameUIを作成する
func NewGameUI(gameState *core.GameState) *GameUI {
	// 各Widgetを初期化
	var resourceView *ResourceView
	if gameState != nil && gameState.Treasury != nil {
		resourceView = NewResourceView(gameState.Treasury)
	} else {
		// デフォルトのTreasuryで初期化
		resourceView = NewResourceView(&core.Treasury{})
	}

	calendarView := NewCalendarView()
	mainView := NewMainView(gameState)
	infoView := NewInfoView()

	// CardDeckViewは最初はnilでも可（後でSetCardDeckで設定）
	cardDeckView := NewCardDeckView(nil)

	ui := &GameUI{
		ResourceView: resourceView,
		CalendarView: calendarView,
		MainView:     mainView,
		InfoView:     infoView,
		CardDeckView: cardDeckView,
		GameState:    gameState,
	}

	// Widget間の連携を設定
	ui.setupWidgetConnections()

	return ui
}

// setupWidgetConnections Widget間の連携を設定
func (gui *GameUI) setupWidgetConnections() {
	// CardDeckViewからInfoViewへの連携
	gui.CardDeckView.OnCardSelected = func(card interface{}) {
		gui.InfoView.SetSelectedCard(card)
	}

	// MainViewからInfoViewへの連携は、Update()で動的に処理
}

// SetGameState GameStateを設定
func (gui *GameUI) SetGameState(gameState *core.GameState) {
	gui.GameState = gameState
	gui.InfoView.SetGameState(gameState)

	// ResourceViewの更新
	if gameState != nil && gameState.Treasury != nil {
		gui.ResourceView.Treasury = gameState.Treasury

		// 次ターンの収入予測を設定
		increment := gui.calculateYieldIncrement()
		gui.ResourceView.SetIncrement(increment)
	}

	// CalendarViewの更新
	if gameState != nil {
		gui.CalendarView.SetCurrentTurn(gameState.CurrentTurn)
	}

	// CardDeckViewの更新は将来的に実装
	// 現在はMyNationにCardDeckフィールドが存在しないため省略
}

// calculateYieldIncrement 次ターンの収入予測を計算
func (gui *GameUI) calculateYieldIncrement() core.ResourceQuantity {
	if gui.GameState == nil {
		return core.ResourceQuantity{}
	}

	// TODO: ここでGameState.CalculateYield()を呼び出す
	// 現在はダミー実装
	return core.ResourceQuantity{
		Money: 5,
		Food:  3,
		Wood:  2,
		Iron:  1,
		Mana:  1,
	}
}

// AddHistoryEvent 履歴イベントを追加
func (gui *GameUI) AddHistoryEvent(event string) {
	gui.InfoView.AddHistoryEvent(event)
}

// SetMousePosition マウス位置を更新
func (gui *GameUI) SetMousePosition(x, y int) {
	gui.MouseX = x
	gui.MouseY = y

	// 各Widgetのマウス位置更新は将来的に実装
	// 現在はWidgetにMouseX/MouseYフィールドが存在しないため省略
}

// HandleInput 統一Input処理
func (gui *GameUI) HandleInput(input *Input) error {
	// MainViewが最初に入力を処理（Viewの切り替えなど重要な処理）
	if err := gui.MainView.HandleInput(input); err != nil {
		return err
	}

	// CardDeckViewの入力処理
	if err := gui.CardDeckView.HandleInput(input); err != nil {
		return err
	}

	// 他のWidgetは表示専用のため入力処理なし
	// ResourceView, CalendarView, InfoViewは基本的に表示のみ

	return nil
}

// Update フレーム更新処理
func (gui *GameUI) Update() error {
	// Widget間の動的連携処理
	gui.updateWidgetConnections()

	// 必要に応じてGameStateから最新データを反映
	if gui.GameState != nil {
		// ResourceViewの増分更新
		increment := gui.calculateYieldIncrement()
		gui.ResourceView.SetIncrement(increment)
	}

	return nil
}

// updateWidgetConnections Widget間の動的連携を更新
func (gui *GameUI) updateWidgetConnections() {
	// MainViewの現在選択状態に応じてInfoViewを更新
	// これは将来的にマウス位置や選択状態から判定する

	// TODO: マウス位置判定によるInfoView表示切り替え
	// 例：マウスがMapGrid上のPointにあるときはそのPointの情報を表示
}

// Draw 全Widget描画
func (gui *GameUI) Draw(screen *ebiten.Image) {
	// 描画順序：背景から前景へ

	// 1. ResourceView (上部左)
	gui.ResourceView.Draw(screen)

	// 2. CalendarView (上部右)
	gui.CalendarView.Draw(screen)

	// 3. MainView (中央メイン)
	gui.MainView.Draw(screen)

	// 4. InfoView (右側情報)
	gui.InfoView.Draw(screen)

	// 5. CardDeckView (下部カードデッキ)
	gui.CardDeckView.Draw(screen)
}

// GetCurrentMainViewType 現在のMainViewTypeを取得
func (gui *GameUI) GetCurrentMainViewType() ViewType {
	return gui.MainView.CurrentView
}

// SwitchMainView MainViewを切り替え
func (gui *GameUI) SwitchMainView(viewType ViewType) {
	gui.MainView.SwitchView(viewType)

	// View切り替え時のInfoView更新
	switch viewType {
	case ViewTypeMapGrid:
		gui.InfoView.CurrentMode = InfoModeHistory
	case ViewTypeMarket:
		// Marketの場合は特に何もしない（カード選択時に更新）
	case ViewTypeBattle:
		// Battleの場合はEnemyの情報を表示する可能性
	case ViewTypeTerritory:
		// Territoryの場合はPointの情報を表示
	}
}

// SelectPoint Pointを選択してInfoViewに反映
func (gui *GameUI) SelectPoint(point core.Point) {
	gui.InfoView.SetSelectedPoint(point)
}

// SelectCard カードを選択してInfoViewに反映
func (gui *GameUI) SelectCard(card interface{}) {
	gui.InfoView.SetSelectedCard(card)
}

// SelectCardFromDeck CardDeckから特定のカードを選択
func (gui *GameUI) SelectCardFromDeck(index int) {
	gui.CardDeckView.SelectCard(index)
}

// MoveCardToBattle 選択中のカードをBattleViewに移動
func (gui *GameUI) MoveCardToBattle() bool {
	selectedCard := gui.CardDeckView.GetSelectedCard()
	if selectedCard == nil {
		return false
	}

	// BattleCardのみ移動可能
	if battleCard, ok := selectedCard.(*core.BattleCard); ok {
		// CardDeckから除去
		gui.CardDeckView.RemoveSelectedCard()

		// BattleViewに追加
		gui.MainView.Battle.AddBattleCard(battleCard)

		gui.AddHistoryEvent("Card moved to battle")
		return true
	}

	return false
}

// MoveCardToTerritory 選択中のカードをTerritoryViewに移動
func (gui *GameUI) MoveCardToTerritory() bool {
	selectedCard := gui.CardDeckView.GetSelectedCard()
	if selectedCard == nil {
		return false
	}

	// StructureCardのみ移動可能
	if structureCard, ok := selectedCard.(*core.StructureCard); ok {
		// CardDeckから除去
		gui.CardDeckView.RemoveSelectedCard()

		// TerritoryViewに追加
		gui.MainView.Territory.AddStructureCard(structureCard)

		gui.AddHistoryEvent("Card moved to territory")
		return true
	}

	return false
}

// ReturnCardToDeck カードをCardDeckに戻す
func (gui *GameUI) ReturnCardToDeck(card interface{}) {
	gui.CardDeckView.AddCard(card)
	gui.AddHistoryEvent("Card returned to deck")
}
