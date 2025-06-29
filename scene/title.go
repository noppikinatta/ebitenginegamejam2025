package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/bamenn"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
	"github.com/noppikinatta/ebitenginegamejam2025/ui"
)

type Title struct {
	input      *ui.Input
	nextScene  ebiten.Game
	sequence   *bamenn.Sequence
	transition bamenn.Transition
}

func NewTitle(input *ui.Input) *Title {
	return &Title{
		input: input,
	}
}

func (t *Title) Init(nextScene ebiten.Game, sequence *bamenn.Sequence, transition bamenn.Transition) {
	t.nextScene = nextScene
	t.sequence = sequence
	t.transition = transition
}

func (t *Title) Update() error {
	if t.input.Mouse.IsJustPressed(ebiten.MouseButtonLeft) {
		t.sequence.SwitchWithTransition(t.nextScene, t.transition)
	}

	return nil
}

func (t *Title) Draw(screen *ebiten.Image) {
	// 背景色を設定
	screen.Fill(color.RGBA{20, 20, 40, 255})

	titleImg := drawing.Image("title")
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(160, 0)
	screen.DrawImage(titleImg, opt)

	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(20, 200)
	drawing.DrawText(screen, lang.Text("story-1"), 12, opt)

	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(220, 320)
	drawing.DrawText(screen, "Click to Start", 14, opt)
}

func (t *Title) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}
