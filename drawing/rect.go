package drawing

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func DrawRect(screen *ebiten.Image, x, y, width, height float64, r, g, b, a float32) {
	vertices := []ebiten.Vertex{
		{DstX: float32(x), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(x + width), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(x + width), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: r, ColorG: g, ColorB: b, ColorA: a},
		{DstX: float32(x), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: r, ColorG: g, ColorB: b, ColorA: a},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, WhitePixel, &ebiten.DrawTrianglesOptions{})
}

type ColorF32 struct {
	r, g, b, a_inv float32
}

func NewColorF32(r, g, b, a float32) ColorF32 {
	return ColorF32{
		r:     r,
		g:     g,
		b:     b,
		a_inv: 1 - a,
	}
}

func (c ColorF32) RGBA() (r, g, b, a float32) {
	return c.r, c.g, c.b, 1 - c.a_inv
}
