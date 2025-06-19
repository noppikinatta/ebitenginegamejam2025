package geom

import (
	"math"
)

// LinearFunc: Ax + By + C = 0.
type LinearFunc struct {
	A, B, C float64
}

func LinearFuncFromPt(pt1, pt2 PointF) LinearFunc {
	a := (pt2.Y - pt1.Y)
	b := (pt1.X - pt2.X)
	c := -1 * (pt1.X*a + pt1.Y*b)

	return LinearFunc{a, b, c}
}

func (f LinearFunc) X(y float64) (float64, bool) {
	if f.A == 0 {
		return 0, false
	}
	return (f.B*y + f.C) / (-1 * f.A), true
}

func (f LinearFunc) Y(x float64) (float64, bool) {
	if f.B == 0 {
		return 0, false
	}

	return (f.A*x + f.C) / (-1 * f.B), true
}

func (f LinearFunc) Distance(pt PointF) float64 {
	numerator := math.Abs(f.A*pt.X + f.B*pt.Y + f.C)
	denominator := math.Sqrt(f.A*f.A + f.B*f.B)

	if denominator == 0 {
		// Not mathematically correct, but I don't want the program to stop in an impossible case.
		return 0
	}

	return numerator / denominator
}

type LineSegment struct {
	Pt1, Pt2 PointF
}

func (l LineSegment) Left() float64 {
	return math.Min(l.Pt1.X, l.Pt2.X)
}

func (l LineSegment) Right() float64 {
	return math.Max(l.Pt1.X, l.Pt2.X)
}

func (l LineSegment) Top() float64 {
	return math.Min(l.Pt1.Y, l.Pt2.Y)
}

func (l LineSegment) Bottom() float64 {
	return math.Max(l.Pt1.Y, l.Pt2.Y)
}

func (l LineSegment) Center() PointF {
	return PointF{
		X: (l.Pt1.X + l.Pt2.X) * 0.5,
		Y: (l.Pt1.Y + l.Pt2.Y) * 0.5,
	}
}

func (l LineSegment) Length() float64 {
	dx := l.Pt1.X - l.Pt2.X
	dy := l.Pt1.Y - l.Pt2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (l LineSegment) CrossesWith(other LineSegment) bool {
	if l.Right() < other.Left() {
		return false
	}
	if other.Right() < l.Left() {
		return false
	}
	if l.Bottom() < other.Top() {
		return false
	}
	if other.Bottom() < l.Top() {
		return false
	}

	fn1 := LinearFuncFromPt(l.Pt1, l.Pt2)
	fn2 := LinearFuncFromPt(other.Pt1, other.Pt2)

	if fn1.A*fn2.B == fn2.A*fn1.B {
		return false // parallel
	}

	crossX := (fn2.C*fn1.B - fn1.C*fn2.B) / (fn1.A*fn2.B - fn2.A*fn1.B)
	if crossX < l.Left() || crossX > l.Right() {
		return false
	}
	if crossX < other.Left() || crossX > other.Right() {
		return false
	}

	return true
}
