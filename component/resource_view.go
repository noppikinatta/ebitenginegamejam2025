package component

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ResourceView struct {
	resources map[string]int
	x, y      int
	width     int
	height    int
}

func NewResourceView(x, y, width, height int) *ResourceView {
	return &ResourceView{
		resources: map[string]int{
			"Gold":  100,
			"Iron":  50,
			"Wood":  75,
			"Grain": 60,
			"Mana":  25,
		},
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}

func (rv *ResourceView) GetResourceTypes() []string {
	types := make([]string, 0, len(rv.resources))
	for resourceType := range rv.resources {
		types = append(types, resourceType)
	}
	return types
}

func (rv *ResourceView) GetResourceAmount(resourceType string) int {
	return rv.resources[resourceType]
}

func (rv *ResourceView) SetResourceAmount(resourceType string, amount int) {
	rv.resources[resourceType] = amount
}

func (rv *ResourceView) Draw(screen *ebiten.Image) {
	// Draw background
	resourceViewBg := ebiten.NewImage(rv.width, rv.height)
	resourceViewBg.Fill(color.RGBA{40, 40, 40, 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(rv.x), float64(rv.y))
	screen.DrawImage(resourceViewBg, op)

	// Draw resource information
	resourceOrder := []string{"Gold", "Iron", "Wood", "Grain", "Mana"}
	for i, resourceType := range resourceOrder {
		amount := rv.resources[resourceType]
		text := fmt.Sprintf("%s: %d", resourceType, amount)
		ebitenutil.DebugPrintAt(screen, text, rv.x+i*100+5, rv.y+5)
	}
}
