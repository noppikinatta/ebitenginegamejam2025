package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Title struct {
	nextScene string
	story     string
}

func NewTitle() *Title {
	return &Title{
		story: "In the Kingdom Year 1000, nations stood divided.\nYou must unite them under one banner.\nForge alliances, defeat enemies, and bring peace to the land.\n\nThe Union awaits your leadership.",
	}
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
	ebitenutil.DebugPrintAt(screen, "UNION STRATEGY GAME", 200, 50)

	// ストーリーテキストを表示
	ebitenutil.DebugPrintAt(screen, "Story:", 50, 100)
	lines := []string{
		"In the Kingdom Year 1000, nations stood divided.",
		"You must unite them under one banner.",
		"Forge alliances, defeat enemies, and bring peace to the land.",
		"",
		"The Union awaits your leadership.",
	}
	for i, line := range lines {
		ebitenutil.DebugPrintAt(screen, line, 60, 120+i*20)
	}

	ebitenutil.DebugPrintAt(screen, "Press ENTER to Start", 220, 300)
}

func (t *Title) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}

func (t *Title) GetNextScene() string {
	return t.nextScene
}

func (t *Title) GetStoryText() string {
	return t.story
}
