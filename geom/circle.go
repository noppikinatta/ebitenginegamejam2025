package geom

type Circle struct {
	Center PointF
	Radius float64
}

func (c Circle) Left() float64 {
	return c.Center.X - c.Radius
}

func (c Circle) Right() float64 {
	return c.Center.X + c.Radius
}

func (c Circle) Top() float64 {
	return c.Center.Y - c.Radius
}

func (c Circle) Bottom() float64 {
	return c.Center.Y + c.Radius
}

func (c Circle) IntersectsWith(other Circle) bool {
	distance := c.Center.Distance(other.Center)
	return distance < (c.Radius + other.Radius)
}
