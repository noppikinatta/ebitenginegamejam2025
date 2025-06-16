package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Title struct {
	nextScene string
}

func (t *Title) Update() error {
	// Enterキーでゲーム画面に遷移
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		t.nextScene = "ingame"
	}
	return nil
}

func (t *Title) Draw(screen *ebiten.Image) {
	// 背景色を設定
	screen.Fill(color.RGBA{20, 20, 40, 255})

	// タイトルテキストを表示
	ebitenutil.DebugPrintAt(screen, "UNION TOWER DEFENSE", 100, 100)
	ebitenutil.DebugPrintAt(screen, "Press ENTER to Start", 120, 150)
}

func (t *Title) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}

func (t *Title) GetNextScene() string {
	return t.nextScene
}
