package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// BattleView is a Widget for displaying battles.
// Position: Drawn within MainView.
type BattleView struct {
	BattlePoint core.BattlePoint  // The point to be battled.
	PointName   string            // The name of the battle point.
	Battlefield *core.Battlefield // Battlefield information.
	GameState   *core.GameState   // Game state.
	HoveredCard interface{}
	MouseX      int
	MouseY      int

	// Callback for view switching.
	OnBackClicked func() // Return to MapGridView.
}

// NewBattleView creates a BattleView.
func NewBattleView(onBackClicked func()) *BattleView {
	return &BattleView{
		OnBackClicked: onBackClicked,
	}
}

// SetBattlePoint sets the battle point to be displayed.
func (bv *BattleView) SetBattlePoint(point core.BattlePoint) {
	bv.BattlePoint = point
	bv.Battlefield = bv.createBattlefield(point)
}

// SetPointName sets the name of the battle point.
func (bv *BattleView) SetPointName(pointName string) {
	bv.PointName = pointName
}

// SetGameState sets the game state.
func (bv *BattleView) SetGameState(gameState *core.GameState) {
	bv.GameState = gameState
}

// GetTotalPower calculates the total Power value of the placed BattleCards.
func (bv *BattleView) GetTotalPower() float64 {
	return bv.Battlefield.CalculateTotalPower()
}

// CanDefeatEnemy determines if the enemy can be defeated.
func (bv *BattleView) CanDefeatEnemy() bool {
	if bv.Battlefield != nil {
		return bv.Battlefield.CanBeat()
	}
	// Keep existing logic for backward compatibility.
	if bv.BattlePoint == nil {
		return false
	}
	return bv.GetTotalPower() >= bv.BattlePoint.GetEnemy().Power
}

// HandleInput handles input.
func (bv *BattleView) HandleInput(input *Input) error {
	cursorX, cursorY := input.Mouse.CursorPosition()
	cardIndex := bv.cardIndex(cursorX, cursorY)
	bv.MouseX = cursorX
	bv.MouseY = cursorY

	if cardIndex != -1 {
		bv.HoveredCard = bv.Battlefield.BattleCards[cardIndex]
	} else {
		bv.HoveredCard = nil
	}

	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		if cardIndex != -1 {
			bv.handleBattleCardClick(cursorX, cursorY)
		}

		// Click detection for the back button (480,20,40,40).
		if cursorX >= 480 && cursorX < 520 && cursorY >= 20 && cursorY < 60 {
			// Return all placed BattleCards to the CardDeck.
			bv.GameState.CardDeck.Add(&core.Cards{BattleCards: bv.Battlefield.BattleCards})
			bv.Battlefield.BattleCards = make([]*core.BattleCard, 0)
			if bv.OnBackClicked != nil {
				bv.OnBackClicked()
				return nil
			}
		}

		// Click detection for the conquer button (200,280,120,40).
		if cursorX >= 200 && cursorX < 320 && cursorY >= 280 && cursorY < 300 {
			if bv.CanDefeatEnemy() {
				bv.Conquer()
			}
			if bv.OnBackClicked != nil {
				bv.GameState.CardDeck.Add(&core.Cards{BattleCards: bv.Battlefield.BattleCards})
				bv.Battlefield.BattleCards = make([]*core.BattleCard, 0)
				bv.OnBackClicked()
			}
			return nil
		}

		// Click detection for the enemy image (victory process).
		if bv.CanDefeatEnemy() && cursorX >= 180 && cursorX < 340 && cursorY >= 60 && cursorY < 220 {
			if bv.Conquer() {
				// On successful conquest, return to MapGridView.
				if bv.OnBackClicked != nil {
					bv.OnBackClicked()
				}
			}
			return nil
		}

		// Click detection for BattleCard (return to CardDeck).
		bv.handleBattleCardClick(cursorX, cursorY)
	}
	return nil
}

func (bv *BattleView) cardIndex(cursorX, cursorY int) int {
	if cursorX < 0 || cursorX >= 480 || cursorY < 220 || cursorY >= 280 {
		return -1
	}
	cardIndex := cursorX / 40
	if cardIndex < 0 || cardIndex >= len(bv.Battlefield.BattleCards) {
		return -1
	}
	return cardIndex
}

// handleBattleCardClick handles BattleCard clicks.
func (bv *BattleView) handleBattleCardClick(cursorX, cursorY int) {
	// Each card (40x60) in the BattleCard area (0,220,480,60).
	if cursorY >= 220 && cursorY < 280 {
		cardIndex := cursorX / 40

		var targetCard *core.BattleCard
		if bv.Battlefield != nil && cardIndex >= 0 && cardIndex < len(bv.Battlefield.BattleCards) {
			targetCard = bv.Battlefield.BattleCards[cardIndex]
		}

		if targetCard != nil {
			// Return the card to the CardDeck.
			bv.RemoveCard(targetCard)
		}
	}
}

// Draw handles drawing.
func (bv *BattleView) Draw(screen *ebiten.Image) {
	// Draw header (0,20,520,40).
	bv.drawHeader(screen)

	// Draw back button (480,20,40,40).
	bv.drawBackButton(screen)

	// Draw enemy image (180,60,160,160).
	bv.drawEnemy(screen)

	// Draw BattleCard area (0,220,480,60).
	bv.drawBattleCards(screen)

	// Draw power display (480,220,40,60).
	bv.drawPowerDisplay(screen)

	// Draw conquer button (200,280,120,40).
	bv.drawConquerButton(screen)

	bv.drawHoveredCardTooltip(screen)
}

// drawHeader draws the header.
func (bv *BattleView) drawHeader(screen *ebiten.Image) {
	// Header background.
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		{DstX: 520, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		{DstX: 520, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		{DstX: 0, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Title text.
	pointName := ""
	if p, ok := bv.BattlePoint.(*core.WildernessPoint); ok {
		pointName = p.TerrainType
	}
	if _, ok := bv.BattlePoint.(*core.BossPoint); ok {
		pointName = "point-boss"
	}
	title := lang.ExecuteTemplate("battle-title", map[string]any{"location": lang.Text(pointName)})

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(10, 30)
	drawing.DrawText(screen, title, 16, opt)
}

// drawBackButton draws the back button.
func (bv *BattleView) drawBackButton(screen *ebiten.Image) {
	DrawButton(screen, 480, 20, 40, 40, "ui-close")
}

// drawEnemy draws the enemy image.
func (bv *BattleView) drawEnemy(screen *ebiten.Image) {
	// Enemy image background (180,60,160,160).
	var color [4]float32
	color = [4]float32{0.8, 0.8, 0.8, 1}

	vertices := []ebiten.Vertex{
		{DstX: 180, DstY: 60, SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
		{DstX: 340, DstY: 60, SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
		{DstX: 340, DstY: 220, SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
		{DstX: 180, DstY: 220, SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Draw enemy information.
	if bv.BattlePoint != nil {
		enemy := bv.BattlePoint.GetEnemy()

		enemyImage := drawing.Image(string(enemy.EnemyID))
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Scale(2, 2)
		opt.GeoM.Translate(180, 60)
		screen.DrawImage(enemyImage, opt)

		// Enemy type.
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(185, 70)
		enemyType := lang.ExecuteTemplate("battle-enemy-type", map[string]any{"type": lang.Text(string(enemy.EnemyType))})
		drawing.DrawText(screen, enemyType, 12, opt)

		// Enemy's Power.
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(185, 90)
		powerIcon := drawing.Image("ui-power")
		screen.DrawImage(powerIcon, opt)
		opt.GeoM.Translate(16, 0)
		powerText := fmt.Sprintf("%s: %.1f", lang.Text("battle-power"), enemy.Power)
		drawing.DrawText(screen, powerText, 12, opt)

		// Enemy's quote.
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(16, 160)
		enemyTalk := lang.ExecuteTemplate("battle-enemy-talk", map[string]any{"name": lang.Text(string(enemy.EnemyID)), "text": lang.Text(string(enemy.Question))})
		drawing.DrawText(screen, enemyTalk, 12, opt)

		// 敵のスキル
		for i, skill := range enemy.Skills {
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(350, 60+float64(i*60))
			skillName := lang.Text(string(skill.ID()))
			drawing.DrawText(screen, skillName, 12, opt)
			opt.GeoM.Translate(0, 16)
			skillDescription := lang.Text(string(skill.ID()) + "-desc")
			drawing.DrawText(screen, skillDescription, 9, opt)
		}
	}
}

// drawBattleCards BattleCard置き場を描画
func (bv *BattleView) drawBattleCards(screen *ebiten.Image) {
	// BattleCard置き場の背景 (0,220,480,60)
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 220, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.2, ColorB: 0.3, ColorA: 1},
		{DstX: 480, DstY: 220, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.2, ColorB: 0.3, ColorA: 1},
		{DstX: 480, DstY: 280, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.2, ColorB: 0.3, ColorA: 1},
		{DstX: 0, DstY: 280, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.2, ColorB: 0.3, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Draw each card.
	// 配置されたBattleCardを描画（40x60 × 12枚）
		for i, card := range bv.Battlefield.BattleCards {
		if i >= 12 { // 最大12枚まで
			break
		}

		cardX := float64(i * 40)
		cardY := 220.0

		// カード描画
		DrawBattleCard(screen, cardX, cardY, card)
	}

	for i := len(bv.Battlefield.BattleCards); i < bv.Battlefield.CardSlot; i++ {
		cardX := float64(i * 40)
		cardY := 220.0
		DrawCardBackground(screen, cardX, cardY, 0.5)
	}
}

// drawPowerDisplay draws the power display.
func (bv *BattleView) drawPowerDisplay(screen *ebiten.Image) {
	// Power表示の背景 (480,220,40,60)
	vertices := []ebiten.Vertex{
		{DstX: 480, DstY: 220, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.4, ColorA: 1},
		{DstX: 520, DstY: 220, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.4, ColorA: 1},
		{DstX: 520, DstY: 280, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.4, ColorA: 1},
		{DstX: 480, DstY: 280, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.4, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// 総Power値を描画
	totalPower := bv.GetTotalPower()

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(485, 225)
	powerIcon := drawing.Image("ui-power")
	screen.DrawImage(powerIcon, opt)

	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(485, 240)
	powerText := fmt.Sprintf("%.1f", totalPower)
	drawing.DrawText(screen, powerText, 14, opt)

	// 必要Power（敵のPower）を表示
	if bv.BattlePoint != nil {
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(485, 255)
		requiredText := fmt.Sprintf("/%.1f", bv.BattlePoint.GetEnemy().Power)
		drawing.DrawText(screen, requiredText, 10, opt)
	}
}

// drawConquerButton draws the conquer button.
func (bv *BattleView) drawConquerButton(screen *ebiten.Image) {
	canConquer := bv.CanDefeatEnemy()

	// ボタンの色を決定
	var colorR, colorG, colorB float32 = 0.5, 0.5, 0.5 // 無効時は灰色
	if canConquer {
		colorR, colorG, colorB = 0.2, 0.8, 0.2 // 有効時は緑
	}

	// ボタン背景
	vertices := []ebiten.Vertex{
		{DstX: 200, DstY: 280, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: 320, DstY: 280, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: 320, DstY: 320, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: 200, DstY: 320, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// ボタンテキスト
	if canConquer {
	opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(210, 280)
		drawing.DrawText(screen, lang.Text("ui-conquer"), 14, opt)
	} else {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(210, 280)
		drawing.DrawText(screen, lang.Text("ui-need-power"), 12, opt)
	}
}

func (bv *BattleView) drawHoveredCardTooltip(screen *ebiten.Image) {
	if bv.HoveredCard == nil {
		return
	}

	DrawCardDescriptionTooltip(screen, bv.HoveredCard, bv.MouseX, bv.MouseY)
}

// createBattlefield 戦場を作成する
func (bv *BattleView) createBattlefield(point core.BattlePoint) *core.Battlefield {
	if bv.GameState == nil {
		return core.NewBattlefield(point.GetEnemy(), 0.0)
	}

	x, y, ok := bv.GameState.MapGrid.XYOfPoint(point)
	if !ok {
		panic("BattleView.createBattlefield: 戦闘地点がマップグリッドに存在しません")
	}
	enemy := point.GetEnemy()
	battlefield := core.NewBattlefield(enemy, 0.0)

	// x,yをもとに上下左右のPointを調査
	mapGrid := bv.GameState.MapGrid
	directions := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} // 上下左右

	for _, dir := range directions {
		checkX := x + dir[0]
		checkY := y + dir[1]

		// マップ範囲内かチェック
		if checkX >= 0 && checkX < mapGrid.Size.X && checkY >= 0 && checkY < mapGrid.Size.Y {
			p := mapGrid.GetPoint(checkX, checkY)

			// ControlledなWildernessPointからTerritoryを取得
			if wildernessPoint, ok := p.(*core.WildernessPoint); ok {
				if wildernessPoint.Controlled && wildernessPoint.Territory != nil {
					territory := wildernessPoint.Territory

					// Territory.CardsのStructureCard.BattlefieldModifierを適用
					for _, card := range territory.Cards {
				if card.BattlefieldModifier != nil {
							card.BattlefieldModifier.Modify(battlefield)
						}
					}
				}
			}
		}
	}

	return battlefield
}

// CanPlaceCard カードを配置できるかどうかを判定
func (bv *BattleView) CanPlaceCard() bool {
	if bv.Battlefield == nil {
		return false
	}
	return len(bv.Battlefield.BattleCards) < bv.Battlefield.CardSlot
}

// PlaceCard カードを配置する
func (bv *BattleView) PlaceCard(card *core.BattleCard) bool {
	if bv.Battlefield == nil {
	return false
}

	return bv.Battlefield.AddBattleCard(card)
}

// RemoveCard カードを除去する
func (bv *BattleView) RemoveCard(card *core.BattleCard) bool {
	if bv.Battlefield == nil {
		return false
	}

	// カードのインデックスを見つける
	cardIndex := -1
	for i, battleCard := range bv.Battlefield.BattleCards {
		if battleCard == card {
			cardIndex = i
				break
			}
		}

	if cardIndex == -1 {
		return false
	}

	// Battlefieldから除去
	removedCard, success := bv.Battlefield.RemoveBattleCard(cardIndex)
	if success && removedCard != nil {
		// GameState.CardDeckに追加
		if bv.GameState != nil {
			cards := &core.Cards{BattleCards: []*core.BattleCard{removedCard}}
			bv.GameState.CardDeck.Add(cards)
			}
		}

	return success
}

// Conquer 制圧処理を実行する
func (bv *BattleView) Conquer() bool {
	if bv.Battlefield == nil || bv.GameState == nil || bv.BattlePoint == nil {
		return false
	}

	// 勝利可能かチェック
	if !bv.Battlefield.CanBeat() {
		return false
	}

	// 戦闘勝利処理
	bv.Battlefield.Beat()

	// 置いたBattleCardをすべてCardDeckに戻す
	bv.GameState.CardDeck.Add(&core.Cards{BattleCards: bv.Battlefield.BattleCards})

	// 対象BattlePointのControlledをtrueに変更
	bv.BattlePoint.SetControlled(true)
	bv.GameState.MapGrid.UpdateAccesibles()

	bv.GameState.MyNation.AppendLevel(0.5)
	bv.GameState.NextTurn()

	return true
}
