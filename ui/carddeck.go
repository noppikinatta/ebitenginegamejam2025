package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
)

// CardDeckView カードデッキWidget
// 位置: (0,300,640,60)
// カードを40x60で最大16枚表示
type CardDeckView struct {
	CardDeck       *core.CardDeck    // 表示するカードデッキ
	SelectedIndex  int               // 選択中のカードインデックス (-1は未選択)
	OnCardSelected func(interface{}) // カード選択時のコールバック

	// マウスカーソル位置（外部から設定）
	MouseX, MouseY int
}

// NewCardDeckView CardDeckViewを作成する
func NewCardDeckView(cardDeck *core.CardDeck) *CardDeckView {
	return &CardDeckView{
		CardDeck:      cardDeck,
		SelectedIndex: -1, // 初期は未選択
	}
}

// SetCardDeck カードデッキを設定
func (c *CardDeckView) SetCardDeck(cardDeck *core.CardDeck) {
	c.CardDeck = cardDeck
	c.SelectedIndex = -1 // デッキが変わったら選択も解除
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

	// ResourceCardsを追加
	for _, card := range c.CardDeck.ResourceCards {
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
	case *core.ResourceCard:
		c.removeResourceCard(card)
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

// removeResourceCard ResourceCardをデッキから削除
func (c *CardDeckView) removeResourceCard(card *core.ResourceCard) {
	for i, cardToRemove := range c.CardDeck.ResourceCards {
		if cardToRemove == card {
			c.CardDeck.ResourceCards = append(c.CardDeck.ResourceCards[:i], c.CardDeck.ResourceCards[i+1:]...)
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
	case *core.ResourceCard:
		c.CardDeck.ResourceCards = append(c.CardDeck.ResourceCards, newCard)
	}
}

// HandleInput 入力処理
func (c *CardDeckView) HandleInput(input *Input) error {
	// TODO: マウスクリックでカード選択
	// 現在はマウス操作が困難なため、後回し
	return nil
}

// Draw 描画処理
func (c *CardDeckView) Draw(screen *ebiten.Image) {
	// 背景描画
	c.drawBackground(screen)

	// カード描画
	c.drawCards(screen)
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
		drawing.DrawText(screen, "No card deck", 12, opt)
		return
	}

	allCards := c.getAllCards()

	if len(allCards) == 0 {
		// カードがない場合のメッセージ
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(300, 320)
		drawing.DrawText(screen, "No cards in deck", 12, opt)
		return
	}

	// カードを40x60サイズで描画（最大16枚）
	for i, card := range allCards {
		if i >= 16 { // 最大16枚まで
			break
		}

		x := float32(i * 40)
		y := float32(300)

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
func (c *CardDeckView) drawCard(screen *ebiten.Image, card interface{}, x, y float32, selected bool) {
	// カード背景色を決定
	var colorR, colorG, colorB float32
	switch card.(type) {
	case *core.BattleCard:
		colorR, colorG, colorB = 0.8, 0.4, 0.2 // オレンジ系
	case *core.StructureCard:
		colorR, colorG, colorB = 0.2, 0.8, 0.4 // 緑系
	default:
		colorR, colorG, colorB = 0.5, 0.5, 0.5 // グレー
	}

	// 選択中の場合は明るくする
	if selected {
		colorR = min(colorR*1.5, 1.0)
		colorG = min(colorG*1.5, 1.0)
		colorB = min(colorB*1.5, 1.0)
	}

	// カード背景描画 (40x60)
	vertices := []ebiten.Vertex{
		{DstX: x, DstY: y, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: x + 40, DstY: y, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: x + 40, DstY: y + 60, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: x, DstY: y + 60, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// 選択中の場合は枠を描画
	if selected {
		c.drawCardBorder(screen, x, y)
	}

	// カード情報描画
	switch cardData := card.(type) {
	case *core.BattleCard:
		c.drawBattleCardInfo(screen, cardData, x, y)
	case *core.StructureCard:
		c.drawStructureCardInfo(screen, cardData, x, y)
	}
}

// drawCardBorder カード枠を描画
func (c *CardDeckView) drawCardBorder(screen *ebiten.Image, x, y float32) {
	// 上枠
	vertices := []ebiten.Vertex{
		{DstX: x, DstY: y, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
		{DstX: x + 40, DstY: y, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
		{DstX: x + 40, DstY: y + 2, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
		{DstX: x, DstY: y + 2, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// 下枠
	vertices = []ebiten.Vertex{
		{DstX: x, DstY: y + 58, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
		{DstX: x + 40, DstY: y + 58, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
		{DstX: x + 40, DstY: y + 60, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
		{DstX: x, DstY: y + 60, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
	}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// 左枠
	vertices = []ebiten.Vertex{
		{DstX: x, DstY: y, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
		{DstX: x + 2, DstY: y, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
		{DstX: x + 2, DstY: y + 60, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
		{DstX: x, DstY: y + 60, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
	}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// 右枠
	vertices = []ebiten.Vertex{
		{DstX: x + 38, DstY: y, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
		{DstX: x + 40, DstY: y, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
		{DstX: x + 40, DstY: y + 60, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
		{DstX: x + 38, DstY: y + 60, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 0, ColorA: 1},
	}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
}

// drawBattleCardInfo BattleCardの詳細を描画
func (c *CardDeckView) drawBattleCardInfo(screen *ebiten.Image, card *core.BattleCard, x, y float32) {
	// カードID (上部、8pt)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(x+2), float64(y+8))
	cardID := string(card.CardID)
	if len(cardID) > 6 {
		cardID = cardID[:5] + "..."
	}
	drawing.DrawText(screen, cardID, 8, opt)

	// Power (中央、12pt)
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(x+12), float64(y+25))
	drawing.DrawText(screen, fmt.Sprintf("%.1f", card.Power), 12, opt)

	// Type (下部、8pt)
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(x+2), float64(y+45))
	cardType := string(card.Type)
	if len(cardType) > 6 {
		cardType = cardType[:5] + "..."
	}
	drawing.DrawText(screen, cardType, 8, opt)
}

// drawStructureCardInfo StructureCardの詳細を描画
func (c *CardDeckView) drawStructureCardInfo(screen *ebiten.Image, card *core.StructureCard, x, y float32) {
	// カードID (上部、8pt)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(x+2), float64(y+8))
	cardID := string(card.CardID)
	if len(cardID) > 6 {
		cardID = cardID[:5] + "..."
	}
	drawing.DrawText(screen, cardID, 8, opt)

	// 効果マーク (中央、14pt)
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(x+15), float64(y+25))
	if card.YieldModifier != nil {
		drawing.DrawText(screen, "⚡", 14, opt) // 効果ありマーク
	} else {
		drawing.DrawText(screen, "○", 14, opt) // 効果なしマーク
	}

	// "STR" (下部、8pt)
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(float64(x+10), float64(y+45))
	drawing.DrawText(screen, "STR", 8, opt)
}

// min 最小値を返すヘルパー関数
func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
