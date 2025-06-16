package entity

type Point struct {
	X        int
	Y        int
	Type     string // "Home", "Wild", "NPC", "Boss"
	Defeated bool   // For Wild/Boss points
	Name     string // For NPC points
}

func NewPoint(x, y int, pointType string) *Point {
	return &Point{
		X:        x,
		Y:        y,
		Type:     pointType,
		Defeated: false,
		Name:     "",
	}
}

func (p *Point) GetType() string {
	return p.Type
}

func (p *Point) IsDefeated() bool {
	return p.Defeated
}

func (p *Point) SetDefeated(defeated bool) {
	p.Defeated = defeated
}

func (p *Point) GetName() string {
	return p.Name
}

func (p *Point) SetName(name string) {
	p.Name = name
}
