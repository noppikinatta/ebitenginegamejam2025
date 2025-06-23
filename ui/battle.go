package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
)

// BattleView Battle表示Widget
// 位置: MainView内で描画
type BattleView struct {
	Enemy       *core.Enemy        // 戦闘対象の敵
	PointName   string             // 戦闘地点名
	BattleCards []*core.BattleCard // 配置されたBattleCard（最大12枚）
	GameState   *core.GameState    // ゲーム状態

	// View切り替えのコールバック
	OnBackClicked  func()                      // MapGridViewに戻る
	OnEnemyClicked func(enemy *core.Enemy)     // 敵クリック時の勝利処理
	OnCardClicked  func(card *core.BattleCard) // カードクリック時（CardDeckに戻す）
}

// NewBattleView BattleViewを作成する
func NewBattleView(onBackClicked func()) *BattleView {
	return &BattleView{
		BattleCards:   make([]*core.BattleCard, 0, 12), // 最大12枚
		OnBackClicked: onBackClicked,
	}
}

// SetEnemy 表示する敵を設定
func (bv *BattleView) SetEnemy(enemy *core.Enemy) {
	bv.Enemy = enemy
}

// SetPointName 戦闘地点名を設定
func (bv *BattleView) SetPointName(pointName string) {
	bv.PointName = pointName
}

// SetGameState ゲーム状態を設定
func (bv *BattleView) SetGameState(gameState *core.GameState) {
	bv.GameState = gameState
}

// AddBattleCard BattleCardを配置する
func (bv *BattleView) AddBattleCard(card *core.BattleCard) bool {
	if len(bv.BattleCards) >= 12 {
		return false // 最大12枚まで
	}

	// 敵のBattleCardSlot制限をチェック
	if bv.Enemy != nil && len(bv.BattleCards) >= bv.Enemy.BattleCardSlot {
		return false
	}

	bv.BattleCards = append(bv.BattleCards, card)
	return true
}

// RemoveBattleCard BattleCardを除去する
func (bv *BattleView) RemoveBattleCard(index int) *core.BattleCard {
	if index < 0 || index >= len(bv.BattleCards) {
		return nil
	}

	card := bv.BattleCards[index]
	bv.BattleCards = append(bv.BattleCards[:index], bv.BattleCards[index+1:]...)
	return card
}

// GetTotalPower 配置されたBattleCardの総Power値を計算
func (bv *BattleView) GetTotalPower() float64 {
	var totalPower float64
	for _, card := range bv.BattleCards {
		totalPower += float64(card.Power)
	}
	return totalPower
}

// CanDefeatEnemy 敵を倒せるかどうかを判定
func (bv *BattleView) CanDefeatEnemy() bool {
	if bv.Enemy == nil {
		return false
	}
	return bv.GetTotalPower() >= bv.Enemy.Power
}

// HandleInput 入力処理
func (bv *BattleView) HandleInput(input *Input) error {
	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		cursorX, cursorY := input.Mouse.CursorPosition()

		// 戻るボタンのクリック判定 (480,20,40,40)
		if cursorX >= 480 && cursorX < 520 && cursorY >= 20 && cursorY < 60 {
			if bv.OnBackClicked != nil {
				bv.OnBackClicked()
				return nil
			}
		}

		// TODO: 敵画像のクリック判定（勝利処理）
		// TODO: BattleCardのクリック判定（CardDeckに戻す）
	}
	return nil
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
	pointName := bv.PointName
	if pointName == "" {
		pointName = "Unknown Point"
	}
	title := fmt.Sprintf("Battle of %s", pointName)

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(10, 30)
	drawing.DrawText(screen, title, 16, opt)
}

// drawBackButton 戻るボタンを描画
func (bv *BattleView) drawBackButton(screen *ebiten.Image) {
	DrawButton(screen, 480, 20, 40, 40, "X")
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
	if bv.Enemy != nil {
		// 敵の名前
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(185, 70)
		enemyName := fmt.Sprintf("Enemy: %s", bv.Enemy.EnemyID)
		drawing.DrawText(screen, enemyName, 12, opt)

		// 敵のPower
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(185, 90)
		powerText := fmt.Sprintf("Power: %.1f", bv.Enemy.Power)
		drawing.DrawText(screen, powerText, 12, opt)

		// CardSlot制限
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(185, 110)
		slotText := fmt.Sprintf("Card Limit: %d", bv.Enemy.BattleCardSlot)
		drawing.DrawText(screen, slotText, 10, opt)

		// 勝利可能性
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(185, 130)
		if bv.CanDefeatEnemy() {
			drawing.DrawText(screen, "CLICK TO WIN!", 14, opt)
		} else {
			drawing.DrawText(screen, "Need more power", 10, opt)
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
	for i, card := range bv.BattleCards {
		if i >= 12 { // 最大12枚まで
			break
		}

		cardX := float64(i * 40)
		cardY := 220.0

		// カード描画
		DrawCard(screen, cardX, cardY, fmt.Sprintf("%.1f", card.Power))
	}

	// 空きスロットを表示
	for i := len(bv.BattleCards); i < 12; i++ {
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
		if bv.Enemy != nil && i >= bv.Enemy.BattleCardSlot {
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
	if bv.Enemy != nil {
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(485, 255)
		requiredText := fmt.Sprintf("/%.1f", bv.Enemy.Power)
		drawing.DrawText(screen, requiredText, 10, opt)
	}
}
