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
	Territory *core.Territory // 表示するTerritory
	PointName string          // 地点名
	GameState *core.GameState // ゲーム状態

	// View切り替えのコールバック
	OnBackClicked func()                         // MapGridViewに戻る
	OnCardClicked func(card *core.StructureCard) // カードクリック時（CardDeckに戻す）
}

// NewTerritoryView TerritoryViewを作成する
func NewTerritoryView(onBackClicked func()) *TerritoryView {
	return &TerritoryView{
		OnBackClicked: onBackClicked,
	}
}

// SetTerritory 表示するTerritoryを設定
func (tv *TerritoryView) SetTerritory(territory *core.Territory) {
	tv.Territory = territory
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
			if tv.OnBackClicked != nil {
				tv.OnBackClicked()
				return nil
			}
		}

		// TODO: StructureCardのクリック判定（CardDeckに戻す）
	}
	return nil
}

// Draw 描画処理
func (tv *TerritoryView) Draw(screen *ebiten.Image) {
	// ヘッダ描画 (0,20,520,40)
	tv.drawHeader(screen)

	// 戻るボタン描画 (480,20,40,40)
	tv.drawBackButton(screen)

	// 産出量表示 (0,60,60,100)
	tv.drawYield(screen)

	// 効果説明 (60,60,460,100)
	tv.drawEffectDescription(screen)

	// StructureCard置き場 (0,160,520,60)
	tv.drawStructureCards(screen)
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

	// 配置されたStructureCardを描画
	if tv.Territory != nil {
		for i, card := range tv.Territory.Cards {
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
