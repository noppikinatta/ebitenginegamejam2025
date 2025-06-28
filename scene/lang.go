package scene

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/geom"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

type langSwitcher struct {
	alpha       float32
	keys        []ebiten.Key
	currentLang string
}

func (s *langSwitcher) Update() {
	s.keys = s.keys[:0]
	s.keys = inpututil.AppendJustPressedKeys(s.keys)
	for _, k := range s.keys {
		if k == ebiten.KeyL {
			s.currentLang = lang.Switch()
			s.alpha = 1
			break
		}
	}

	s.alpha -= 0.01
	if s.alpha < 0 {
		s.alpha = 0
	}
}

func (s *langSwitcher) Draw(screen *ebiten.Image) {
	if s.alpha == 0 {
		return
	}

	idxs := []uint16{0, 1, 2, 0, 2, 3}

	topLeft := geom.PointF{
		X: 0,
		Y: 0,
	}

	bottomRight := geom.PointF{
		X: 200,
		Y: 40,
	}

	v := ebiten.Vertex{
		ColorA: s.alpha,
	}
	vertices := make([]ebiten.Vertex, 4)
	v.DstX = float32(topLeft.X)
	v.DstY = float32(topLeft.Y)
	vertices[0] = v
	v.DstX = float32(topLeft.X)
	v.DstY = float32(bottomRight.Y)
	vertices[1] = v
	v.DstX = float32(bottomRight.X)
	v.DstY = float32(bottomRight.Y)
	v.ColorA = 0
	vertices[2] = v
	v.DstX = float32(bottomRight.X)
	v.DstY = float32(topLeft.Y)
	v.ColorA = 0
	vertices[3] = v

	ropt := ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(vertices, idxs, drawing.WhitePixel, &ropt)

	topt := ebiten.DrawImageOptions{}
	topt.GeoM.Translate(0, 8)
	topt.ColorScale.Scale(s.alpha, s.alpha, s.alpha, s.alpha)
	drawing.DrawText(screen, fmt.Sprintf("Current Language: %s", s.currentLang), 16, &topt)
}
