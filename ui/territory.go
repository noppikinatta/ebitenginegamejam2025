package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
)

// TerritoryView Territory表示Widget
// 位置: MainView内で描画
type TerritoryView struct {
	Territory *core.Territory       // 表示するTerritory
	Cards     []*core.StructureCard // 一時的カード置き場（Territory.Cardsとメモリ共有しない）
	PointName string                // 地点名
	GameState *core.GameState       // ゲーム状態

	// View切り替えのコールバック
	OnBackClicked func()                         // MapGridViewに戻る
	OnCardClicked func(card *core.StructureCard) // カードクリック時（CardDeckに戻す）
}

// NewTerritoryView TerritoryViewを作成する
func NewTerritoryView(onBackClicked func()) *TerritoryView {
	return &TerritoryView{
		Cards:         make([]*core.StructureCard, 0),
		OnBackClicked: onBackClicked,
	}
}

// SetTerritory 表示するTerritoryを設定
func (tv *TerritoryView) SetTerritory(territory *core.Territory) {
	tv.Territory = territory

	// Territory.Cardsの内容を一時置き場にコピー（メモリ共有を避ける）
	tv.Cards = make([]*core.StructureCard, len(territory.Cards))
	copy(tv.Cards, territory.Cards)
}

// SetPointName 地点名を設定
func (tv *TerritoryView) SetPointName(pointName string) {
	tv.PointName = pointName
}

// SetGameState ゲーム状態を設定
func (tv *TerritoryView) SetGameState(gameState *core.GameState) {
	tv.GameState = gameState
}

// AddStructureCard StructureCardを配置する
func (tv *TerritoryView) AddStructureCard(card *core.StructureCard) bool {
	if tv.Territory == nil {
		return false
	}
	return tv.Territory.AppendCard(card)
}

// RemoveStructureCard StructureCardを除去する
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

// GetCurrentYield 現在の産出量を取得
func (tv *TerritoryView) GetCurrentYield() core.ResourceQuantity {
	if tv.Territory == nil {
		return core.ResourceQuantity{}
	}
	return tv.Territory.Yield()
}

// HandleInput 入力処理
func (tv *TerritoryView) HandleInput(input *Input) error {
	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		cursorX, cursorY := input.Mouse.CursorPosition()

		// 戻るボタンのクリック判定 (480,20,40,40)
		if cursorX >= 480 && cursorX < 520 && cursorY >= 20 && cursorY < 60 {
			// 変更があった場合は建設決定処理を実行
			if tv.IsChanged() {
				tv.ConfirmConstruction()
			}

			if tv.OnBackClicked != nil {
				tv.OnBackClicked()
				return nil
			}
		}

		// StructureCardのクリック判定（CardDeckに戻す）
		tv.handleStructureCardClick(cursorX, cursorY)
	}
	return nil
}

// Draw 描画処理
func (tv *TerritoryView) Draw(screen *ebiten.Image) {
	// ヘッダ描画 (0,20,520,40)
	tv.drawHeader(screen)

	// 戻るボタン描画 (480,20,40,40)
	tv.drawBackButton(screen)

	// 変更状態表示 (440,20,40,40)
	tv.drawChangeIndicator(screen)

	// 産出量表示 (0,60,60,100)
	tv.drawYield(screen)

	// 効果説明 (60,60,460,100)
	tv.drawEffectDescription(screen)

	// StructureCard置き場 (0,160,520,60)
	tv.drawStructureCards(screen)

	// 建設決定ボタン (200,220,120,40)
	tv.drawConstructionButton(screen)
}

// drawHeader ヘッダを描画
func (tv *TerritoryView) drawHeader(screen *ebiten.Image) {
	// ヘッダ背景
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: 520, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: 520, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: 0, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// タイトルテキスト
	pointName := tv.PointName
	if pointName == "" && tv.Territory != nil {
		pointName = string(tv.Territory.TerritoryID)
	}
	if pointName == "" {
		pointName = "Unknown Territory"
	}

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(10, 30)
	drawing.DrawText(screen, pointName, 16, opt)
}

// drawBackButton 戻るボタンを描画
func (tv *TerritoryView) drawBackButton(screen *ebiten.Image) {
	DrawButton(screen, 480, 20, 40, 40, "X")
}

// drawYield 産出量表示を描画
func (tv *TerritoryView) drawYield(screen *ebiten.Image) {
	// 産出量表示の背景 (0,60,60,100)
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 60, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 60, DstY: 160, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 0, DstY: 160, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// 現在の産出量を取得
	currentYield := tv.GetCurrentYield()

	// 5種類の資源を60x20ずつで表示
	resourceTypes := []struct {
		name  string
		value int
	}{
		{"Money", currentYield.Money},
		{"Food", currentYield.Food},
		{"Wood", currentYield.Wood},
		{"Iron", currentYield.Iron},
		{"Mana", currentYield.Mana},
	}

	for i, resource := range resourceTypes {
		y := 60.0 + float64(i)*20

		// Resource画像(20x20)
		icon := GetResourceIcon(resource.name)
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(5, y)
		drawing.DrawText(screen, icon, 12, opt)

		// 産出量数字(40x20)
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(25, y)
		yieldText := fmt.Sprintf("%d", resource.value)
		drawing.DrawText(screen, yieldText, 12, opt)
	}
}

// drawEffectDescription 効果説明を描画
func (tv *TerritoryView) drawEffectDescription(screen *ebiten.Image) {
	// 効果説明の背景 (60,60,460,100)
	vertices := []ebiten.Vertex{
		{DstX: 60, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.25, ColorG: 0.25, ColorB: 0.25, ColorA: 1},
		{DstX: 520, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.25, ColorG: 0.25, ColorB: 0.25, ColorA: 1},
		{DstX: 520, DstY: 160, SrcX: 0, SrcY: 0, ColorR: 0.25, ColorG: 0.25, ColorB: 0.25, ColorA: 1},
		{DstX: 60, DstY: 160, SrcX: 0, SrcY: 0, ColorR: 0.25, ColorG: 0.25, ColorB: 0.25, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// タイトル
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(65, 65)
	drawing.DrawText(screen, "Structure Effects:", 12, opt)

	if tv.Territory == nil || len(tv.Territory.Cards) == 0 {
		// カードが配置されていない場合
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(65, 85)
		drawing.DrawText(screen, "No structure cards placed.", 10, opt)

		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(65, 105)
		drawing.DrawText(screen, "Place cards to get bonuses!", 10, opt)
		return
	}

	// 配置されたStructureCardの効果を表示
	startY := 85.0
	for i, card := range tv.Territory.Cards {
		if i >= 4 { // 最大4枚まで表示
			break
		}

		y := startY + float64(i)*18

		// カード名
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(65, y)
		cardName := fmt.Sprintf("Card: %s", card.CardID)
		drawing.DrawText(screen, cardName, 10, opt)

		// 効果説明（ダミーテキスト）
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(200, y)
		effect := "Boosts resource production" // 実際の効果説明は後で実装
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
	for i, card := range tv.Cards {
		cardX := float64(i * 40)
		cardY := 160.0

		// カード描画
		cardName := fmt.Sprintf("%s", card.CardID)
		if len(cardName) > 6 {
			cardName = cardName[:6] // 表示用に短縮
		}
		DrawCard(screen, cardX, cardY, cardName)
	}

	// 空きスロットを表示
	if tv.Territory != nil {
		maxSlots := tv.Territory.CardSlot
		for i := len(tv.Cards); i < maxSlots && i < 13; i++ { // 最大13枚まで表示（520÷40=13）
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

// handleStructureCardClick StructureCardのクリック処理
func (tv *TerritoryView) handleStructureCardClick(cursorX, cursorY int) {
	// StructureCard置き場 (0,160,520,60) 内の各カード (40x60)
	if cursorY >= 160 && cursorY < 220 {
		cardIndex := cursorX / 40

		if cardIndex >= 0 && cardIndex < len(tv.Cards) {
			targetCard := tv.Cards[cardIndex]

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

	// 「*」マークを表示
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(455, 35)
	drawing.DrawText(screen, "*", 20, opt)
}

// drawConstructionButton 建設決定ボタンを描画
func (tv *TerritoryView) drawConstructionButton(screen *ebiten.Image) {
	isChanged := tv.IsChanged()

	// ボタンの色を決定
	var colorR, colorG, colorB float32 = 0.4, 0.4, 0.4 // 変更なしは灰色
	var buttonText string = "No Changes"

	if isChanged {
		colorR, colorG, colorB = 0.2, 0.6, 0.8 // 変更ありは青
		buttonText = "CONFIRM"
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
	if len(tv.Cards) != len(tv.Territory.Cards) {
		return true
	}

	// 両スライスのStructureCardポインタが全て同じかチェック（順番は問わない）
	for _, viewCard := range tv.Cards {
		found := false
		for _, territoryCard := range tv.Territory.Cards {
			if viewCard == territoryCard {
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

	// TerritoryView.Cardsの内容をTerritory.Cardsにコピー
	tv.Territory.Cards = make([]*core.StructureCard, len(tv.Cards))
	copy(tv.Territory.Cards, tv.Cards)
}

// CanPlaceCard カードを配置できるかどうかを判定
func (tv *TerritoryView) CanPlaceCard() bool {
	if tv.Territory == nil {
		return false
	}
	return len(tv.Cards) < tv.Territory.CardSlot
}

// PlaceCard カードを配置する
func (tv *TerritoryView) PlaceCard(card *core.StructureCard) bool {
	if !tv.CanPlaceCard() {
		return false
	}

	tv.Cards = append(tv.Cards, card)
	return true
}

// RemoveCard カードを除去する
func (tv *TerritoryView) RemoveCard(card *core.StructureCard) bool {
	// カードのインデックスを見つける
	cardIndex := -1
	for i, structureCard := range tv.Cards {
		if structureCard == card {
			cardIndex = i
			break
		}
	}

	if cardIndex == -1 {
		return false
	}

	// Cardsから除去
	tv.Cards = append(tv.Cards[:cardIndex], tv.Cards[cardIndex+1:]...)

	// GameState.CardDeckに追加
	if tv.GameState != nil {
		cards := &core.Cards{StructureCards: []*core.StructureCard{card}}
		tv.GameState.CardDeck.Add(cards)
	}

	return true
}
