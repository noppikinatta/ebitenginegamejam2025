package component

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/noppikinatta/ebitenginegamejam2025/system"
)

type ResourceView struct {
	resourceManager *system.ResourceManager
	x, y            int
	width           int
	height          int
}

func NewResourceView(rm *system.ResourceManager, x, y, width, height int) *ResourceView {
	return &ResourceView{
		resourceManager: rm,
		x:               x,
		y:               y,
		width:           width,
		height:          height,
	}
}

func (rv *ResourceView) GetResourceTypes() []string {
	resources := rv.resourceManager.GetAllResources()
	types := make([]string, 0, len(resources))
	for resourceType := range resources {
		types = append(types, resourceType)
	}
	return types
}

func (rv *ResourceView) GetResourceAmount(resourceType string) int {
	return rv.resourceManager.GetResource(resourceType)
}

func (rv *ResourceView) SetResourceAmount(resourceType string, amount int) {
	rv.resourceManager.SetResource(resourceType, amount)
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
		amount := rv.resourceManager.GetResource(resourceType)
		text := fmt.Sprintf("%s: %d", resourceType, amount)
		ebitenutil.DebugPrintAt(screen, text, rv.x+i*100+5, rv.y+5)
	}
}
