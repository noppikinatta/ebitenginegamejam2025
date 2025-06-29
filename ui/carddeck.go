package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// CardDeckView カードデッキWidget
// 位置: (0,300,640,60)
// カードを40x60で最大16枚表示
type CardDeckView struct {
	CardDeck       *core.CardDeck    // 表示するカードデッキ
	SelectedIndex  int               // 選択中のカードインデックス (-1は未選択)
	OnCardSelected func(interface{}) // カード選択時のコールバック

	// 新しいコールバック
	OnBattleCardClicked    func(*core.BattleCard) bool    // BattleCardクリック時のコールバック
	OnStructureCardClicked func(*core.StructureCard) bool // StructureCardクリック時のコールバック

	// マウスカーソル位置（外部から設定）
	MouseX, MouseY int

	HoveredCard interface{}
}

// NewCardDeckView CardDeckViewを作成する
func NewCardDeckView(cardDeck *core.CardDeck) *CardDeckView {
	return &CardDeckView{
		CardDeck:      cardDeck,
		SelectedIndex: -1, // 初期は未選択
	}
}

// GetSelectedCard 選択中のカードを取得
func (c *CardDeckView) GetSelectedCard() interface{} {
	if c.CardDeck == nil || c.SelectedIndex < 0 {
		return nil
	}

	allCards := c.getAllCards()
	if c.SelectedIndex >= len(allCards) {
		return nil
	}

	return allCards[c.SelectedIndex]
}

// getAllCards 全てのカードを1つのスライスで取得
func (c *CardDeckView) getAllCards() []interface{} {
	if c.CardDeck == nil {
		return []interface{}{}
	}

	allCards := make([]interface{}, 0)

	// BattleCardsを追加
	for _, card := range c.CardDeck.BattleCards {
		allCards = append(allCards, card)
	}

	// StructureCardsを追加
	for _, card := range c.CardDeck.StructureCards {
		allCards = append(allCards, card)
	}

	return allCards
}

// SelectCard カードを選択
func (c *CardDeckView) SelectCard(index int) {
	if c.CardDeck == nil {
		return
	}

	allCards := c.getAllCards()
	if index < 0 || index >= len(allCards) {
		c.SelectedIndex = -1
		if c.OnCardSelected != nil {
			c.OnCardSelected(nil)
		}
		return
	}

	c.SelectedIndex = index
	if c.OnCardSelected != nil {
		c.OnCardSelected(allCards[index])
	}
}

// ClearSelection 選択をクリア
func (c *CardDeckView) ClearSelection() {
	c.SelectedIndex = -1
	if c.OnCardSelected != nil {
		c.OnCardSelected(nil)
	}
}

// RemoveSelectedCard 選択中のカードをデッキから除去
func (c *CardDeckView) RemoveSelectedCard() interface{} {
	if c.CardDeck == nil || c.SelectedIndex < 0 {
		return nil
	}

	allCards := c.getAllCards()
	if c.SelectedIndex >= len(allCards) {
		return nil
	}

	// 選択中のカードを取得
	selectedCard := allCards[c.SelectedIndex]

	// カードをデッキから除去
	switch card := selectedCard.(type) {
	case *core.BattleCard:
		c.removeBattleCard(card)
	case *core.StructureCard:
		c.removeStructureCard(card)
	}

	// 選択をクリア
	c.ClearSelection()

	return selectedCard
}

// removeBattleCard BattleCardをデッキから削除
func (c *CardDeckView) removeBattleCard(card *core.BattleCard) {
	for i, cardToRemove := range c.CardDeck.BattleCards {
		if cardToRemove == card {
			c.CardDeck.BattleCards = append(c.CardDeck.BattleCards[:i], c.CardDeck.BattleCards[i+1:]...)
			break
		}
	}
}

// removeStructureCard StructureCardをデッキから削除
func (c *CardDeckView) removeStructureCard(card *core.StructureCard) {
	for i, cardToRemove := range c.CardDeck.StructureCards {
		if cardToRemove == card {
			c.CardDeck.StructureCards = append(c.CardDeck.StructureCards[:i], c.CardDeck.StructureCards[i+1:]...)
			break
		}
	}
}

// AddCard カードをデッキに追加
func (c *CardDeckView) AddCard(card interface{}) {
	if c.CardDeck == nil {
		return
	}

	switch newCard := card.(type) {
	case *core.BattleCard:
		c.CardDeck.BattleCards = append(c.CardDeck.BattleCards, newCard)
	case *core.StructureCard:
		c.CardDeck.StructureCards = append(c.CardDeck.StructureCards, newCard)
	}
}

// HandleInput 入力処理
func (c *CardDeckView) HandleInput(input *Input) error {
	cursorX, cursorY := input.Mouse.CursorPosition()
	c.MouseX = cursorX
	c.MouseY = cursorY

	cardIndex := c.CardIndex(cursorX, cursorY)

	if cardIndex != -1 {
		c.HoveredCard = c.getAllCards()[cardIndex]
	} else {
		c.HoveredCard = nil
	}

	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		// CardDeckView領域内かチェック (0,300,640,60)
		if cursorY >= 300 && cursorY < 360 && cursorX >= 0 && cursorX < 640 {
			c.handleCardClick(cursorX, cursorY)
		}
	}
	return nil
}

func (c *CardDeckView) CardIndex(cursorX, cursorY int) int {
	if cursorX < 0 || cursorX >= 640 || cursorY < 300 || cursorY >= 360 {
		return -1
	}

	allCards := c.getAllCards()
	cardIndex := cursorX / 40

	if cardIndex < 0 || cardIndex >= len(allCards) || cardIndex >= 16 {
		return -1
	}

	return cardIndex
}

// handleCardClick カードクリック処理
func (c *CardDeckView) handleCardClick(cursorX, cursorY int) {
	if c.CardDeck == nil {
		return
	}

	allCards := c.getAllCards()
	cardIndex := c.CardIndex(cursorX, cursorY)
	if cardIndex == -1 {
		return
	}

	card := allCards[cardIndex]

	// カードタイプに応じてコールバックを呼び出し
	switch cardData := card.(type) {
	case *core.BattleCard:
		if c.OnBattleCardClicked != nil {
			if c.OnBattleCardClicked(cardData) {
				// trueが返された場合、カードをデッキから削除
				c.removeBattleCard(cardData)
			}
		}
	case *core.StructureCard:
		if c.OnStructureCardClicked != nil {
			if c.OnStructureCardClicked(cardData) {
				// trueが返された場合、カードをデッキから削除
				c.removeStructureCard(cardData)
			}
		}
	}
}

// Draw 描画処理
func (c *CardDeckView) Draw(screen *ebiten.Image) {
	// 背景描画
	c.drawBackground(screen)

	// カード描画
	c.drawCards(screen)

	c.drawHoveredCardTooltip(screen)
}

// drawBackground 背景を描画
func (c *CardDeckView) drawBackground(screen *ebiten.Image) {
	// CardDeckView背景 (0,300,640,60)
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 300, SrcX: 0, SrcY: 0, ColorR: 0.1, ColorG: 0.1, ColorB: 0.15, ColorA: 1},
		{DstX: 640, DstY: 300, SrcX: 0, SrcY: 0, ColorR: 0.1, ColorG: 0.1, ColorB: 0.15, ColorA: 1},
		{DstX: 640, DstY: 360, SrcX: 0, SrcY: 0, ColorR: 0.1, ColorG: 0.1, ColorB: 0.15, ColorA: 1},
		{DstX: 0, DstY: 360, SrcX: 0, SrcY: 0, ColorR: 0.1, ColorG: 0.1, ColorB: 0.15, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
}

// drawCards カードを描画
func (c *CardDeckView) drawCards(screen *ebiten.Image) {
	if c.CardDeck == nil {
		// カードデッキがない場合のメッセージ
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(320, 320)
		drawing.DrawText(screen, lang.Text("card-no-deck"), 12, opt)
		return
	}

	allCards := c.getAllCards()

	if len(allCards) == 0 {
		// カードがない場合のメッセージ
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(300, 320)
		drawing.DrawText(screen, lang.Text("card-no-cards"), 12, opt)
		return
	}

	// カードを40x60サイズで描画（最大16枚）
	for i, card := range allCards {
		if i >= 16 { // 最大16枚まで
			break
		}

		x := float64(i * 40)
		y := float64(300)

		c.drawCard(screen, card, x, y, i == c.SelectedIndex)
	}

	// 16枚を超える場合は省略表示
	if len(allCards) > 16 {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(580, 345)
		drawing.DrawText(screen, fmt.Sprintf("+%d", len(allCards)-16), 9, opt)
	}
}

// drawCard 個別のカードを描画
func (c *CardDeckView) drawCard(screen *ebiten.Image, card interface{}, x, y float64, selected bool) {
	switch typedCard := card.(type) {
	case *core.BattleCard:
		DrawBattleCard(screen, x, y, typedCard)
	case *core.StructureCard:
		DrawCard(screen, x, y, string(typedCard.CardID))
	}
}

func (c *CardDeckView) drawHoveredCardTooltip(screen *ebiten.Image) {
	if c.HoveredCard == nil {
		return
	}

	DrawCardDescriptionTooltip(screen, c.HoveredCard, c.MouseX, c.MouseY)
}
