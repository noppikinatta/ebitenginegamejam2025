package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

func main() {
	seq := scene.CreateSequence()
	ebiten.RunGame(seq)
}
