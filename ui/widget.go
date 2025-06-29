package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

type Widget interface {
	HandleInput(input *Input) error
	Draw(screen *ebiten.Image)
}

// 共通描画機能

// DrawResource 資源量を描画する（60x20領域）
// x, y: 描画位置
// resourceType: 資源の種類（表示用）
// value: 表示する値
// increment: 増分（+2のように表示）
func DrawResource(screen *ebiten.Image, x, y float64, resourceType string, value int, increment int) {
	// 左20x20に資源アイコン（今は文字で代用）
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x, y)
	resourceIcon := drawing.Image(resourceType)
	screen.DrawImage(resourceIcon, opt)

	// 右40x20に数値表示
	var text string
	if increment != 0 {
		text = fmt.Sprintf("%d(+%d)", value, increment)
	} else {
		text = fmt.Sprintf("%d", value)
	}

	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+20, y+2)
	drawing.DrawText(screen, text, 10, opt)
}

func DrawCardBackground(screen *ebiten.Image, x, y float64, alpha float32) {
	var r, g, b float32
	r = 1.9
	g = 0.8
	b = 0.7
	// 枠を描画（四角形）
	drawing.DrawRect(screen, x+1, y+1, 38, 58, r*alpha, g*alpha, b*alpha, alpha)
	r = 0.9
	g = 0.7
	b = 0.5
	drawing.DrawRect(screen, x+3, y+3, 34, 54, r*alpha, g*alpha, b*alpha, alpha)
}

// DrawCard カードを描画する（40x60領域）
func DrawCard(screen *ebiten.Image, x, y float64, cardID string) {
	DrawCardBackground(screen, x, y, 1)

	// カード名を描画
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+4, y+4)
	opt.ColorScale.Scale(1, 1, 1, 0.8)
	cardImage := drawing.Image(cardID)
	screen.DrawImage(cardImage, opt)
	cardName := lang.Text(cardID)
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+1, y+48)
	drawing.DrawText(screen, cardName, 9, opt)
}

func DrawBattleCard(screen *ebiten.Image, x, y float64, battleCard *core.BattleCard) {
	DrawCard(screen, x, y, string(battleCard.CardID))
	// カードタイプを描画
	var rectR, rectG, rectB, rectA float32
	rectR, rectG, rectB, rectA = 1, 1, 1, 1
	optText := &ebiten.DrawImageOptions{}
	optText.GeoM.Translate(x+3, y+1)
	cardTypeText := ""
	switch battleCard.Type {
	case "cardtype-str":
		cardTypeText = "S"
		rectR, rectG, rectB, rectA = 1, 0.2, 0.2, 1
	case "cardtype-agi":
		cardTypeText = "A"
		rectR, rectG, rectB, rectA = 0.2, 1, 0.2, 1
	case "cardtype-mag":
		cardTypeText = "M"
		rectR, rectG, rectB, rectA = 0.2, 0.2, 1, 1
	}
	drawing.DrawRect(screen, x+3, y+3, 8, 12, rectR, rectG, rectB, rectA)
	drawing.DrawText(screen, cardTypeText, 10, optText)

	// パワーを描画
	drawing.DrawRect(screen, x+3, y+34, 34, 14, 0.2, 0.2, 0.2, 0.75)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+1, y+32)
	powerIcon := drawing.Image("ui-power")
	screen.DrawImage(powerIcon, opt)
	opt.GeoM.Translate(16, 0)
	drawing.DrawText(screen, fmt.Sprintf("%.1f", battleCard.Power()), 12, opt)
}

// DrawButton ボタンを描画する（クリック判定は後で実装）
func DrawButton(screen *ebiten.Image, x, y, width, height float64, imageKey string) {
	// ボタンの枠を描画
	vertices := []ebiten.Vertex{
		{DstX: float32(x), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.8, ColorB: 0.8, ColorA: 1},
		{DstX: float32(x + width), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.8, ColorB: 0.8, ColorA: 1},
		{DstX: float32(x + width), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.8, ColorB: 0.8, ColorA: 1},
		{DstX: float32(x), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.8, ColorB: 0.8, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// ボタン画像を描画
	image := drawing.Image(imageKey)
	imageBounds := image.Bounds()
	imageWidth := float64(imageBounds.Dx())
	imageHeight := float64(imageBounds.Dy())
	imageX := x + width/2 - imageWidth/2
	imageY := y + height/2 - imageHeight/2
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(imageX, imageY)
	screen.DrawImage(image, opt)
}

func DrawCardDescriptionTooltip(screen *ebiten.Image, card interface{}, mouseX, mouseY int) {
	left, top := float64(mouseX-20), float64(mouseY-80)

	switch typedCard := card.(type) {
	case *core.BattleCard:
		drawing.DrawRect(screen, left, top, 240, 40, 0, 0, 0, 0.5)
		skillName := lang.Text(string(typedCard.Skill.BattleCardSkillID))
		skillDescription := lang.Text(string(typedCard.Skill.DescriptionKey))
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(left, top)
		drawing.DrawText(screen, skillName, 12, opt)
		opt.GeoM.Translate(0, 16)
		drawing.DrawText(screen, skillDescription, 9, opt)
	case *core.StructureCard:
		drawing.DrawRect(screen, left, top, 400, 60, 0, 0, 0, 0.5)
		description := lang.Text(string(typedCard.DescriptionKey))
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(left, top)
		drawing.DrawText(screen, description, 9, opt)
	}
}
