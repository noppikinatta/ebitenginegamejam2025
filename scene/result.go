package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameHistory struct {
	entries []string
}

type Result struct {
	history *GameHistory
}

func NewResult() *Result {
	return &Result{
		history: &GameHistory{
			entries: []string{
				"Kingdom Year 1000, Month 4: Game Started",
				"Kingdom Year 1000, Month 5: First Alliance Formed",
				"Kingdom Year 1000, Month 8: Boss Defeated",
			},
		},
	}
}

func (r *Result) Update() error {
	return nil
}

func (r *Result) Draw(screen *ebiten.Image) {
	// Background color
	screen.Fill(color.RGBA{60, 40, 80, 255})

	// Title
	ebitenutil.DebugPrintAt(screen, "VICTORY!", 280, 50)
	ebitenutil.DebugPrintAt(screen, "Game History:", 50, 100)

	// Display history entries
	for i, entry := range r.history.entries {
		ebitenutil.DebugPrintAt(screen, entry, 60, 130+i*20)
	}
}

func (r *Result) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}

func (r *Result) GetGameHistory() *GameHistory {
	return r.history
}
