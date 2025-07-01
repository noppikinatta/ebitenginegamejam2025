package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/bamenn"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
	"github.com/noppikinatta/ebitenginegamejam2025/ui"
)

type GameHistory struct {
	entries []string
}

type Result struct {
	history    *GameHistory
	input      *ui.Input
	nextScene  ebiten.Game
	sequence   *bamenn.Sequence
	transition bamenn.Transition
}

func NewResult(input *ui.Input) *Result {
	return &Result{
		history: &GameHistory{
			entries: []string{
				"Kingdom Year 1000, Month 4: Game Started",
				"Kingdom Year 1000, Month 5: First Alliance Formed",
				"Kingdom Year 1000, Month 8: Boss Defeated",
			},
		},
		input: input,
	}
}

func (r *Result) Init(nextScene ebiten.Game, sequence *bamenn.Sequence, transition bamenn.Transition) {
	r.nextScene = nextScene
	r.sequence = sequence
	r.transition = transition
}

func (r *Result) Update() error {

	return nil
}

func (r *Result) Draw(screen *ebiten.Image) {
	// Background color
	screen.Fill(color.RGBA{60, 40, 80, 255})

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(40, 360)
	drawing.DrawText(screen, lang.Text("story-2"), 24, opt)
}

func (r *Result) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1280, 720
}

func (r *Result) GetGameHistory() *GameHistory {
	return r.history
}
