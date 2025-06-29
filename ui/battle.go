package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// BattleView Battle表示Widget
// 位置: MainView内で描画
type BattleView struct {
	BattlePoint core.BattlePoint  // 戦闘対象の地点
	PointName   string            // 戦闘地点名
	Battlefield *core.Battlefield // 戦場情報
	GameState   *core.GameState   // ゲーム状態

	// View切り替えのコールバック
	OnBackClicked func() // MapGridViewに戻る
}

// NewBattleView BattleViewを作成する
func NewBattleView(onBackClicked func()) *BattleView {
	return &BattleView{
		OnBackClicked: onBackClicked,
	}
}

// SetBattlePoint 表示する戦闘地点を設定
func (bv *BattleView) SetBattlePoint(point core.BattlePoint) {
	bv.BattlePoint = point
	bv.Battlefield = bv.createBattlefield(point)
}

// SetPointName 戦闘地点名を設定
func (bv *BattleView) SetPointName(pointName string) {
	bv.PointName = pointName
}

// SetGameState ゲーム状態を設定
func (bv *BattleView) SetGameState(gameState *core.GameState) {
	bv.GameState = gameState
}

// GetTotalPower 配置されたBattleCardの総Power値を計算
func (bv *BattleView) GetTotalPower() float64 {
	return bv.Battlefield.CalculateTotalPower()
}

// CanDefeatEnemy 敵を倒せるかどうかを判定
func (bv *BattleView) CanDefeatEnemy() bool {
	if bv.Battlefield != nil {
		return bv.Battlefield.CanBeat()
	}
	// 後方互換性のため既存ロジックも残す
	if bv.BattlePoint == nil {
		return false
	}
	return bv.GetTotalPower() >= bv.BattlePoint.GetEnemy().Power
}

// HandleInput 入力処理
func (bv *BattleView) HandleInput(input *Input) error {
	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		cursorX, cursorY := input.Mouse.CursorPosition()

		// 戻るボタンのクリック判定 (480,20,40,40)
		if cursorX >= 480 && cursorX < 520 && cursorY >= 20 && cursorY < 60 {
			// 置いたBattleCardをすべてCardDeckに戻す
			bv.GameState.CardDeck.Add(&core.Cards{BattleCards: bv.Battlefield.BattleCards})
			bv.Battlefield.BattleCards = make([]*core.BattleCard, 0)
			if bv.OnBackClicked != nil {
				bv.OnBackClicked()
				return nil
			}
		}

		// 制圧ボタンのクリック判定 (200,280,120,40)
		if bv.CanDefeatEnemy() && cursorX >= 200 && cursorX < 320 && cursorY >= 280 && cursorY < 320 {
			if bv.Conquer() {
				// 制圧成功時、MapGridViewに戻る
				if bv.OnBackClicked != nil {
					bv.OnBackClicked()
				}
			}
			return nil
		}

		// 敵画像のクリック判定（勝利処理）
		if bv.CanDefeatEnemy() && cursorX >= 180 && cursorX < 340 && cursorY >= 60 && cursorY < 220 {
			if bv.Conquer() {
				// 制圧成功時、MapGridViewに戻る
				if bv.OnBackClicked != nil {
					bv.OnBackClicked()
				}
			}
			return nil
		}

		// BattleCardのクリック判定（CardDeckに戻す）
		bv.handleBattleCardClick(cursorX, cursorY)
	}
	return nil
}

// handleBattleCardClick BattleCardのクリック処理
func (bv *BattleView) handleBattleCardClick(cursorX, cursorY int) {
	// BattleCard置き場 (0,220,480,60) 内の各カード (40x60)
	if cursorY >= 220 && cursorY < 280 {
		cardIndex := cursorX / 40

		var targetCard *core.BattleCard
		if bv.Battlefield != nil && cardIndex >= 0 && cardIndex < len(bv.Battlefield.BattleCards) {
			targetCard = bv.Battlefield.BattleCards[cardIndex]
		}

		if targetCard != nil {
			// カードをCardDeckに戻す
			bv.RemoveCard(targetCard)
		}
	}
}

// Draw 描画処理
func (bv *BattleView) Draw(screen *ebiten.Image) {
	// ヘッダ描画 (0,20,520,40)
	bv.drawHeader(screen)

	// 戻るボタン描画 (480,20,40,40)
	bv.drawBackButton(screen)

	// 敵画像描画 (180,60,160,160)
	bv.drawEnemy(screen)

	// BattleCard置き場描画 (0,220,480,60)
	bv.drawBattleCards(screen)

	// Power表示 (480,220,40,60)
	bv.drawPowerDisplay(screen)

	// 制圧ボタン描画 (200,280,120,40)
	bv.drawConquerButton(screen)
}

// drawHeader ヘッダを描画
func (bv *BattleView) drawHeader(screen *ebiten.Image) {
	// ヘッダ背景
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		{DstX: 520, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		{DstX: 520, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		{DstX: 0, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// タイトルテキスト
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

// drawBackButton 戻るボタンを描画
func (bv *BattleView) drawBackButton(screen *ebiten.Image) {
	DrawButton(screen, 480, 20, 40, 40, "ui-close")
}

// drawEnemy 敵画像を描画
func (bv *BattleView) drawEnemy(screen *ebiten.Image) {
	// 敵画像の背景 (180,60,160,160)
	var color [4]float32
	if bv.CanDefeatEnemy() {
		color = [4]float32{0.2, 0.8, 0.2, 1} // 緑（勝てる）
	} else {
		color = [4]float32{0.8, 0.2, 0.2, 1} // 赤（勝てない）
	}

	vertices := []ebiten.Vertex{
		{DstX: 180, DstY: 60, SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
		{DstX: 340, DstY: 60, SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
		{DstX: 340, DstY: 220, SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
		{DstX: 180, DstY: 220, SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// 敵の情報を描画
	if bv.BattlePoint != nil {
		enemy := bv.BattlePoint.GetEnemy()

		enemyImage := drawing.Image(string(enemy.EnemyID))
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Scale(2, 2)
		opt.GeoM.Translate(180, 60)
		screen.DrawImage(enemyImage, opt)

		// 敵のタイプ
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(185, 70)
		enemyType := lang.ExecuteTemplate("battle-enemy-type", map[string]any{"type": lang.Text(string(enemy.EnemyType))})
		drawing.DrawText(screen, enemyType, 12, opt)

		// 敵のPower
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(185, 90)
		powerIcon := drawing.Image("ui-power")
		screen.DrawImage(powerIcon, opt)
		opt.GeoM.Translate(16, 0)
		powerText := fmt.Sprintf("%s: %.1f", lang.Text("battle-power"), enemy.Power)
		drawing.DrawText(screen, powerText, 12, opt)

		// 敵のセリフ
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(16, 160)
		enemyTalk := lang.ExecuteTemplate("battle-enemy-talk", map[string]any{"name": lang.Text(string(enemy.EnemyID)), "text": lang.Text(string(enemy.Question))})
		drawing.DrawText(screen, enemyTalk, 12, opt)

		// 勝利可能性
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(185, 130)
		if bv.CanDefeatEnemy() {
			drawing.DrawText(screen, lang.Text("ui-click-to-win"), 14, opt)
		} else {
			drawing.DrawText(screen, lang.Text("ui-need-more-power"), 10, opt)
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

	// 空きスロットを表示
	for i := len(bv.Battlefield.BattleCards); i < 12; i++ {
		cardX := float64(i * 40)
		cardY := 220.0

		// 空きスロットの枠線
		vertices := []ebiten.Vertex{
			{DstX: float32(cardX), DstY: float32(cardY), SrcX: 0, SrcY: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5, ColorA: 0.5},
			{DstX: float32(cardX + 40), DstY: float32(cardY), SrcX: 0, SrcY: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5, ColorA: 0.5},
			{DstX: float32(cardX + 40), DstY: float32(cardY + 60), SrcX: 0, SrcY: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5, ColorA: 0.5},
			{DstX: float32(cardX), DstY: float32(cardY + 60), SrcX: 0, SrcY: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5, ColorA: 0.5},
		}
		indices := []uint16{0, 1, 2, 0, 2, 3}
		screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

		// 敵のBattleCardSlot制限で使用不可の場合
		if bv.BattlePoint != nil && i >= bv.BattlePoint.GetEnemy().BattleCardSlot {
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(cardX+10, cardY+25)
			drawing.DrawText(screen, "X", 16, opt)
		}
	}
}

// drawPowerDisplay Power表示を描画
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
	drawing.DrawText(screen, "Power", 10, opt)

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

// drawConquerButton 制圧ボタンを描画
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
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(215, 295)
	if canConquer {
		drawing.DrawText(screen, lang.Text("ui-conquer"), 14, opt)
	} else {
		drawing.DrawText(screen, lang.Text("ui-need-power"), 12, opt)
	}
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
	supportPower := 0.0

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
							// BattlefieldModifierを適用（現在は簡単な実装）
							// TODO: 実際のModifierの適用ロジック
						}
					}

					// 簡単な支援力計算（隣接する制圧済み領土ごとに+1）
					supportPower += 1.0
				}
			}
		}
	}

	return core.NewBattlefield(enemy, supportPower)
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
