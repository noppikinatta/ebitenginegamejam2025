package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
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
	drawing.DrawText(screen, resourceType, 10, opt)

	// 右40x20に数値表示
	var text string
	if increment != 0 {
		text = fmt.Sprintf("%d(+%d)", value, increment)
	} else {
		text = fmt.Sprintf("%d", value)
	}

	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+20, y)
	drawing.DrawText(screen, text, 12, opt)
}

// DrawCard カードを描画する（40x60領域）
func DrawCard(screen *ebiten.Image, x, y float64, cardName string) {
	// 枠を描画（四角形）
	vertices := []ebiten.Vertex{
		{DstX: float32(x), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x + 40), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x + 40), DstY: float32(y + 60), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: float32(x), DstY: float32(y + 60), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// カード名を描画
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+2, y+2)
	drawing.DrawText(screen, cardName, 8, opt)
}

// DrawButton ボタンを描画する（クリック判定は後で実装）
func DrawButton(screen *ebiten.Image, x, y, width, height float64, text string) {
	// ボタンの枠を描画
	vertices := []ebiten.Vertex{
		{DstX: float32(x), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.8, ColorB: 0.8, ColorA: 1},
		{DstX: float32(x + width), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.8, ColorB: 0.8, ColorA: 1},
		{DstX: float32(x + width), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.8, ColorB: 0.8, ColorA: 1},
		{DstX: float32(x), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.8, ColorB: 0.8, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// ボタンテキストを描画
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+width/2-10, y+height/2-6) // 中央寄せ（概算）
	drawing.DrawText(screen, text, 12, opt)
}

// GetResourceIcon 資源タイプに応じたアイコン文字を返す
func GetResourceIcon(resourceType string) string {
	switch resourceType {
	case "Money":
		return "¥"
	case "Food":
		return "🍞"
	case "Wood":
		return "🌲"
	case "Iron":
		return "⚒"
	case "Mana":
		return "✨"
	default:
		return "?"
	}
}
