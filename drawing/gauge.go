package drawing

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/geom"
)

type GaugeDrawer struct {
	Max           int
	Current       int
	TopLeft       geom.PointF
	BottomRight   geom.PointF
	TextOffset    geom.PointF
	FontSize      float64
	ColorScaleMax ebiten.ColorScale
	ColorScaleMin ebiten.ColorScale
}

func (d *GaugeDrawer) Draw(screen *ebiten.Image) {
	d.drawRect(screen)

	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(d.TopLeft.X+d.TextOffset.X, d.TopLeft.Y+d.TextOffset.Y)

	DrawText(screen, fmt.Sprint(d.Current), d.FontSize, &opt)
}

func (d *GaugeDrawer) drawRect(screen *ebiten.Image) {
	idxs := []uint16{0, 1, 2, 0, 2, 3}

	v := d.colorVert()

	width := (d.BottomRight.X - d.TopLeft.X) * float64(d.Proportion())
	bottomRight := geom.PointF{
		X: d.TopLeft.X + width,
		Y: d.BottomRight.Y,
	}

	vertices := make([]ebiten.Vertex, 4)
	v.DstX = float32(d.TopLeft.X)
	v.DstY = float32(d.TopLeft.Y)
	vertices[0] = v
	v.DstX = float32(d.TopLeft.X)
	v.DstY = float32(bottomRight.Y)
	vertices[1] = v
	v.DstX = float32(bottomRight.X)
	v.DstY = float32(bottomRight.Y)
	vertices[2] = v
	v.DstX = float32(bottomRight.X)
	v.DstY = float32(d.TopLeft.Y)
	vertices[3] = v

	opt := ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(vertices, idxs, WhitePixel, &opt)
}

func (d *GaugeDrawer) Proportion() float32 {
	if d.Max == 0 {
		return 0 // avoid panic
	}

	return float32(d.Current) / float32(d.Max)
}

func (d *GaugeDrawer) colorVert() ebiten.Vertex {
	p := d.Proportion()

	v := func(min, max float32) float32 {
		return (1-p)*min + p*max
	}

	return ebiten.Vertex{
		ColorR: v(d.ColorScaleMin.R(), d.ColorScaleMax.R()),
		ColorG: v(d.ColorScaleMin.G(), d.ColorScaleMax.G()),
		ColorB: v(d.ColorScaleMin.B(), d.ColorScaleMax.B()),
		ColorA: v(d.ColorScaleMin.A(), d.ColorScaleMax.A()),
	}
}
