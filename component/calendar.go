package component

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Calendar struct {
	currentYear  int
	currentMonth int
	x, y         int
	width        int
	height       int
}

func NewCalendar(x, y, width, height int) *Calendar {
	return &Calendar{
		currentYear:  1000, // Kingdom Year 1000
		currentMonth: 4,    // Month 4
		x:            x,
		y:            y,
		width:        width,
		height:       height,
	}
}

func (c *Calendar) GetCurrentYear() int {
	return c.currentYear
}

func (c *Calendar) GetCurrentMonth() int {
	return c.currentMonth
}

func (c *Calendar) GetDisplayText() string {
	return fmt.Sprintf("Kingdom Year %d, Month %d", c.currentYear, c.currentMonth)
}

func (c *Calendar) AdvanceMonth() {
	c.currentMonth++
	if c.currentMonth > 12 {
		c.currentMonth = 1
		c.currentYear++
	}
}

func (c *Calendar) Draw(screen *ebiten.Image) {
	// Draw background
	calendarBg := ebiten.NewImage(c.width, c.height)
	calendarBg.Fill(color.RGBA{60, 60, 80, 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.x), float64(c.y))
	screen.DrawImage(calendarBg, op)

	// Draw calendar text
	ebitenutil.DebugPrintAt(screen, "Calendar", c.x+5, c.y+5)
	ebitenutil.DebugPrintAt(screen, c.GetDisplayText(), c.x+5, c.y+20)
}
