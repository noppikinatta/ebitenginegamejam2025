package ui

import "github.com/hajimehoshi/ebiten/v2"

type Widget interface {
	HandleInput(input *Input) error
	Draw(screen *ebiten.Image)
}
