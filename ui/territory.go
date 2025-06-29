package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// TerritoryView Territory display Widget
// Position: Drawn within MainView
type TerritoryView struct {
	Territory   *core.Territory       // Territory to display
	OldCards    []*core.StructureCard // Cards before changes
	TerrainType string                // Terrain name
	GameState   *core.GameState       // Game state
	HoveredCard interface{}
	MouseX      int
	MouseY      int

	// View switching callbacks
	OnBackClicked func()                         // Return to MapGridView
	OnCardClicked func(card *core.StructureCard) // When card is clicked (return to CardDeck)
}

// NewTerritoryView creates a TerritoryView
func NewTerritoryView(onBackClicked func()) *TerritoryView {
	return &TerritoryView{
		OnBackClicked: onBackClicked,
	}
}

// SetTerritory sets the Territory to display
func (tv *TerritoryView) SetTerritory(territory *core.Territory, terrainType string) {
	tv.Territory = territory
	tv.OldCards = make([]*core.StructureCard, len(territory.Cards))
	copy(tv.OldCards, territory.Cards)
	tv.TerrainType = terrainType
}

// SetGameState sets the game state
func (tv *TerritoryView) SetGameState(gameState *core.GameState) {
	tv.GameState = gameState
}

// AddStructureCard places a StructureCard
func (tv *TerritoryView) AddStructureCard(card *core.StructureCard) bool {
	if tv.Territory == nil {
		return false
	}
	return tv.Territory.AppendCard(card)
}

// RemoveStructureCard removes a StructureCard
func (tv *TerritoryView) RemoveStructureCard(index int) *core.StructureCard {
	if tv.Territory == nil {
		return nil
	}
	card, ok := tv.Territory.RemoveCard(index)
	if !ok {
		return nil
	}
	return card
}

// GetCurrentYield gets the current yield
func (tv *TerritoryView) GetCurrentYield() core.ResourceQuantity {
	if tv.Territory == nil {
		return core.ResourceQuantity{}
	}
	return tv.Territory.Yield()
}

// HandleInput processes input
func (tv *TerritoryView) HandleInput(input *Input) error {
	cursorX, cursorY := input.Mouse.CursorPosition()
	tv.MouseX = cursorX
	tv.MouseY = cursorY
	cardIndex := tv.cardIndex(cursorX, cursorY)
	if cardIndex != -1 {
		tv.HoveredCard = tv.Territory.Cards[cardIndex]
	} else {
		tv.HoveredCard = nil
	}

	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {

		// Back button click detection (480,20,40,40)
		if cursorX >= 480 && cursorX < 520 && cursorY >= 20 && cursorY < 60 {
			// Execute construction confirmation if there are changes
			if tv.IsChanged() {
				tv.ConfirmConstruction()
			}

			if tv.OnBackClicked != nil {
				tv.OnBackClicked()
				return nil
			}
		}

		// Construction confirmation button click detection (200,220,120,40)
		if cursorX >= 200 && cursorX < 320 && cursorY >= 220 && cursorY < 260 {
			// Execute construction confirmation if there are changes
			if tv.IsChanged() {
				tv.ConfirmConstruction()
			}
			if tv.OnBackClicked != nil {
				tv.OnBackClicked()
				return nil
			}
		}

		// StructureCard click detection (return to CardDeck)
		tv.handleStructureCardClick(cursorX, cursorY)
	}
	return nil
}

func (tv *TerritoryView) cardIndex(cursorX, cursorY int) int {
	if cursorX < 0 || cursorX >= 520 || cursorY < 160 || cursorY >= 220 {
		return -1
	}
	cardIndex := cursorX / 40
	if cardIndex < 0 || cardIndex >= len(tv.Territory.Cards) {
		return -1
	}
	return cardIndex
}

// Draw drawing process
func (tv *TerritoryView) Draw(screen *ebiten.Image) {
	// Draw header (0,20,520,40)
	tv.drawHeader(screen)

	// Draw back button (480,20,40,40)
	tv.drawBackButton(screen)

	// Draw change indicator (440,20,40,40)
	tv.drawChangeIndicator(screen)

	// Draw yield display (0,60,60,100)
	tv.drawYield(screen)

	// Draw effect description (60,60,460,100)
	tv.drawEffectDescription(screen)

	// Draw StructureCard slots (0,160,520,60)
	tv.drawStructureCards(screen)

	// Draw construction confirmation button (200,220,120,40)
	tv.drawConstructionButton(screen)

	tv.drawHoveredCardTooltip(screen)
}

// drawHeader draws the header
func (tv *TerritoryView) drawHeader(screen *ebiten.Image) {
	// Header background
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: 520, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: 520, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: 0, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Title text
	terrainName := lang.Text(tv.TerrainType)

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(10, 30)
	drawing.DrawText(screen, terrainName, 16, opt)
}

// drawBackButton draws the back button
func (tv *TerritoryView) drawBackButton(screen *ebiten.Image) {
	DrawButton(screen, 480, 20, 40, 40, "ui-close")
}

// drawYield draws the yield display
func (tv *TerritoryView) drawYield(screen *ebiten.Image) {
	// Yield display background (0,60,60,100)
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 60, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 60, DstY: 160, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 0, DstY: 160, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Get current yield
	currentYield := tv.GetCurrentYield()

	// Display 5 resource types at 60x20 each
	resourceTypes := []struct {
		name  string
		value int
	}{
		{"resource-money", currentYield.Money},
		{"resource-food", currentYield.Food},
		{"resource-wood", currentYield.Wood},
		{"resource-iron", currentYield.Iron},
		{"resource-mana", currentYield.Mana},
	}

	for i, resource := range resourceTypes {
		y := 60.0 + float64(i)*20

		// Resource image (20x20)
		icon := drawing.Image(resource.name)
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(5, y)
		screen.DrawImage(icon, opt)

		// Yield number (40x20)
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(25, y)
		yieldText := fmt.Sprintf("%d", resource.value)
		drawing.DrawText(screen, yieldText, 12, opt)
	}
}

// drawEffectDescription draws the effect description
func (tv *TerritoryView) drawEffectDescription(screen *ebiten.Image) {
	// Effect description background (60,60,460,100)
	vertices := []ebiten.Vertex{
		{DstX: 60, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.25, ColorG: 0.25, ColorB: 0.25, ColorA: 1},
		{DstX: 520, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.25, ColorG: 0.25, ColorB: 0.25, ColorA: 1},
		{DstX: 520, DstY: 160, SrcX: 0, SrcY: 0, ColorR: 0.25, ColorG: 0.25, ColorB: 0.25, ColorA: 1},
		{DstX: 60, DstY: 160, SrcX: 0, SrcY: 0, ColorR: 0.25, ColorG: 0.25, ColorB: 0.25, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Title
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(65, 65)
	drawing.DrawText(screen, lang.Text("territory-structure-effects"), 12, opt)

	if tv.Territory == nil || len(tv.Territory.Cards) == 0 {
		// When no cards are placed
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(65, 85)
		drawing.DrawText(screen, lang.Text("territory-no-structures"), 10, opt)

		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(65, 105)
		drawing.DrawText(screen, lang.Text("territory-place-cards"), 10, opt)
		return
	}

	// Display effects of placed StructureCards
	startY := 85.0
	for i, card := range tv.Territory.Cards {
		if i >= 4 { // Display up to 4 cards maximum
			break
		}

		y := startY + float64(i)*18

		// カード名
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(65, y)
		cardName := lang.Text(string(card.CardID))
		drawing.DrawText(screen, cardName, 10, opt)

		// 効果説明
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(200, y)
		effect := lang.Text(string(card.DescriptionKey))
		drawing.DrawText(screen, effect, 9, opt)
	}
}

// drawStructureCards StructureCard置き場を描画
func (tv *TerritoryView) drawStructureCards(screen *ebiten.Image) {
	// StructureCard置き場の背景 (0,160,520,60)
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 160, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 520, DstY: 160, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 520, DstY: 220, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 0, DstY: 220, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// 配置されたStructureCardを描画（一時置き場tv.Cardsを使用）
	for i, card := range tv.Territory.Cards {
		cardX := float64(i * 40)
		cardY := 160.0

		DrawCard(screen, cardX, cardY, string(card.CardID))
	}

	// 空きスロットを表示
	if tv.Territory != nil {
		maxSlots := tv.Territory.CardSlot
		for i := len(tv.Territory.Cards); i < maxSlots && i < 13; i++ { // 最大13枚まで表示（520÷40=13）
			cardX := float64(i * 40)
			cardY := 160.0

			// 空きスロットの枠線
			vertices := []ebiten.Vertex{
				{DstX: float32(cardX), DstY: float32(cardY), SrcX: 0, SrcY: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5, ColorA: 0.5},
				{DstX: float32(cardX + 40), DstY: float32(cardY), SrcX: 0, SrcY: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5, ColorA: 0.5},
				{DstX: float32(cardX + 40), DstY: float32(cardY + 60), SrcX: 0, SrcY: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5, ColorA: 0.5},
				{DstX: float32(cardX), DstY: float32(cardY + 60), SrcX: 0, SrcY: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5, ColorA: 0.5},
			}
			indices := []uint16{0, 1, 2, 0, 2, 3}
			screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
		}
	}
}

func (tv *TerritoryView) drawHoveredCardTooltip(screen *ebiten.Image) {
	if tv.HoveredCard == nil {
		return
	}

	DrawCardDescriptionTooltip(screen, tv.HoveredCard, tv.MouseX, tv.MouseY)
}

// handleStructureCardClick StructureCardのクリック処理
func (tv *TerritoryView) handleStructureCardClick(cursorX, cursorY int) {
	// StructureCard置き場 (0,160,520,60) 内の各カード (40x60)
	if cursorY >= 160 && cursorY < 220 {
		cardIndex := cursorX / 40

		if cardIndex >= 0 && cardIndex < len(tv.Territory.Cards) {
			targetCard := tv.Territory.Cards[cardIndex]

			// カードをCardDeckに戻す
			tv.RemoveCard(targetCard)
			if tv.OnCardClicked != nil {
				tv.OnCardClicked(targetCard)
			}
		}
	}
}

// drawChangeIndicator 変更状態表示を描画
func (tv *TerritoryView) drawChangeIndicator(screen *ebiten.Image) {
	if !tv.IsChanged() {
		return // 変更がなければ何も表示しない
	}

	// 変更状態表示の背景 (440,20,40,40)
	vertices := []ebiten.Vertex{
		{DstX: 440, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.6, ColorB: 0.2, ColorA: 1},
		{DstX: 480, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.6, ColorB: 0.2, ColorA: 1},
		{DstX: 480, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.6, ColorB: 0.2, ColorA: 1},
		{DstX: 440, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.6, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// "*"マークを表示
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(455, 35)
	drawing.DrawText(screen, "*", 20, opt)
}

// drawConstructionButton 建設決定ボタンを描画
func (tv *TerritoryView) drawConstructionButton(screen *ebiten.Image) {
	isChanged := tv.IsChanged()

	// ボタンの色を決定
	var colorR, colorG, colorB float32 = 0.4, 0.4, 0.4 // 変更なしは灰色
	var buttonText string = lang.Text("ui-no-changes")

	if isChanged {
		colorR, colorG, colorB = 0.2, 0.6, 0.8 // 変更ありは青
		buttonText = lang.Text("ui-confirm")
	}

	// ボタン背景 (200,220,120,40)
	vertices := []ebiten.Vertex{
		{DstX: 200, DstY: 220, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: 320, DstY: 220, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: 320, DstY: 260, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: 200, DstY: 260, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// ボタンテキスト
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(220, 235)
	drawing.DrawText(screen, buttonText, 12, opt)
}

// IsChanged 構成が変わっているかどうかを判定
func (tv *TerritoryView) IsChanged() bool {
	if tv.Territory == nil {
		return false
	}

	// 長さが異なれば変更有り
	if len(tv.Territory.Cards) != len(tv.OldCards) {
		return true
	}

	// 両スライスのStructureCardポインタが全て同じかチェック（順番は問わない）
	for _, oldCard := range tv.OldCards {
		found := false
		for _, territoryCard := range tv.Territory.Cards {
			if oldCard == territoryCard {
				found = true
				break
			}
		}
		if !found {
			return true
		}
	}

	return false
}

// ConfirmConstruction 建設決定処理
func (tv *TerritoryView) ConfirmConstruction() {
	if tv.Territory == nil {
		return
	}
}

// CanPlaceCard カードを配置できるかどうかを判定
func (tv *TerritoryView) CanPlaceCard() bool {
	if tv.Territory == nil {
		return false
	}
	return len(tv.Territory.Cards) < tv.Territory.CardSlot
}

// PlaceCard カードを配置する
func (tv *TerritoryView) PlaceCard(card *core.StructureCard) bool {
	if !tv.CanPlaceCard() {
		return false
	}

	tv.Territory.Cards = append(tv.Territory.Cards, card)
	return true
}

// RemoveCard カードを除去する
func (tv *TerritoryView) RemoveCard(card *core.StructureCard) bool {
	// カードのインデックスを見つける
	cardIndex := -1
	for i, structureCard := range tv.Territory.Cards {
		if structureCard == card {
			cardIndex = i
			break
		}
	}

	if cardIndex == -1 {
		return false
	}

	// Cardsから除去
	tv.Territory.Cards = append(tv.Territory.Cards[:cardIndex], tv.Territory.Cards[cardIndex+1:]...)

	// GameState.CardDeckに追加
	if tv.GameState != nil {
		cards := &core.Cards{StructureCards: []*core.StructureCard{card}}
		tv.GameState.CardDeck.Add(cards)
	}

	return true
}
