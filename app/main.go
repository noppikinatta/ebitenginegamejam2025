package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/scene"
	"github.com/noppikinatta/ebitenginegamejam2025/ui"
	"github.com/noppikinatta/nyuuryoku"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Ebitengine Game Jam 2025")

	input := ui.Input{Mouse: nyuuryoku.NewMouse()}
	seq := scene.CreateSequence(&input)
	ebiten.RunGame(seq)
}
