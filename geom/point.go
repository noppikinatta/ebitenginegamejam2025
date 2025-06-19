package geom

import (
	"image"
	"math"
)

type PointF struct {
	X, Y float64
}

func (p PointF) Add(other PointF) PointF {
	return PointF{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

func (p PointF) Subtract(other PointF) PointF {
	return PointF{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}

func (p PointF) Multiply(value float64) PointF {
	return PointF{
		X: p.X * value,
		Y: p.Y * value,
	}
}

func (p PointF) Angle() float64 {
	return math.Atan2(p.Y, p.X)
}

func (p PointF) Abs() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p PointF) Distance(other PointF) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (p PointF) InnerProduct(other PointF) float64 {
	return p.X*other.X + p.Y*other.Y
}

func PointFFromPolar(abs float64, angleRad float64) PointF {
	x := abs * math.Cos(angleRad)
	y := abs * math.Sin(angleRad)

	return PointF{X: x, Y: y}
}

func PointFFromPoint(p image.Point) PointF {
	return PointF{
		X: float64(p.X),
		Y: float64(p.Y),
	}
}
