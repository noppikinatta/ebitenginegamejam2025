package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InGame struct {
}

func NewInGame() *InGame {
	return &InGame{}
}

func (g *InGame) Update() error {
	return nil
}

func (g *InGame) Draw(screen *ebiten.Image) {

}

func (g *InGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}
