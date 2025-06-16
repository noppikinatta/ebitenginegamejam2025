package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type InGame struct {
}

func (g *InGame) Update() error {
	return nil
}

func (g *InGame) Draw(screen *ebiten.Image) {
	// 背景色を設定
	screen.Fill(color.RGBA{40, 40, 60, 255})

	// ゲーム画面のテキストを表示
	ebitenutil.DebugPrintAt(screen, "IN GAME", 280, 100)
	ebitenutil.DebugPrintAt(screen, "Game scene - coming soon!", 200, 200)
	ebitenutil.DebugPrintAt(screen, "Press ESC to return to Title", 200, 400)
}

func (g *InGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}
