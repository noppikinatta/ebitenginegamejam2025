package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Ebitengine Game Jam 2025")

	seq := scene.CreateSequence()
	ebiten.RunGame(seq)
}
