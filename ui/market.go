package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/flow"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// MarketView is a widget for displaying the Market.
type MarketView struct {
	flow      *flow.MarketFlow
	viewModel *viewmodel.MarketViewModel
}

// NewMarketView creates a MarketView
func NewMarketView(flow *flow.MarketFlow, viewModel *viewmodel.MarketViewModel) *MarketView {
	return &MarketView{
		flow:      flow,
		viewModel: viewModel,
	}
}

// HandleInput processes input
func (mv *MarketView) HandleInput(input *Input) (back bool, err error) {
	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		cursorX, cursorY := input.Mouse.CursorPosition()

		// Back button click detection (960,40,80,80)
		if cursorX >= 960 && cursorX < 1040 && cursorY >= 40 && cursorY < 120 {
			return true, nil
		}

		// CardPack click detection and purchase processing
		if mv.handleMarketItemClick(cursorX, cursorY) {
			return true, nil
		}
	}
	return false, nil
}

// handleMarketItemClick handles MarketItem clicks
func (mv *MarketView) handleMarketItemClick(cursorX, cursorY int) (purchased bool) {
	positions := [][4]int{
		{0, 120, 520, 160},   // Top left
		{520, 120, 520, 160}, // Top right
		{0, 280, 520, 160},   // Middle left
		{520, 280, 520, 160}, // Middle right
		{0, 440, 520, 160},   // Bottom left
		{520, 440, 520, 160}, // Bottom right
	}

	for i, pos := range positions {
		if cursorX >= pos[0] && cursorX < pos[0]+pos[2] &&
			cursorY >= pos[1] && cursorY < pos[1]+pos[3] {

			return mv.flow.Purchase(i)
		}
	}

	return false
}

// Draw handles the drawing process
func (mv *MarketView) Draw(screen *ebiten.Image) {
	// Draw header (0,40,960,80)
	mv.drawHeader(screen)

	// Draw back button (960,40,80,80)
	mv.drawBackButton(screen)

	// Draw CardPack list
	mv.drawMarketItems(screen)
}

// drawHeader draws the Nation name header
func (mv *MarketView) drawHeader(screen *ebiten.Image) {
	// Header background
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 40, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.3, ColorA: 1},
		{DstX: 960, DstY: 40, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.3, ColorA: 1},
		{DstX: 960, DstY: 120, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.3, ColorA: 1},
		{DstX: 0, DstY: 120, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.3, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Nation name text
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(20, 60)
	drawing.DrawText(screen, mv.viewModel.Title(), 32, opt)

	// Market level text
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(600, 60)
	marketLevel := lang.ExecuteTemplate("ui-market-level", map[string]any{"level": mv.viewModel.Level()})
	drawing.DrawText(screen, marketLevel, 28, opt)
}

// drawBackButton draws the back button
func (mv *MarketView) drawBackButton(screen *ebiten.Image) {
	DrawButton(screen, 960, 40, 80, 80, "ui-close")
}

// drawMarketItems draws the list of MarketItems
func (mv *MarketView) drawMarketItems(screen *ebiten.Image) {
	numItems := mv.viewModel.NumItems()

	// CardPack display area: 520x160 x 6
	// Positions: (0,120,520,160), (520,120,520,160), (0,280,520,160), (520,280,520,160), (0,440,520,160), (520,440,520,160)
	positions := [][4]float64{
		{0, 120, 520, 160},   // Top left
		{520, 120, 520, 160}, // Top right
		{0, 280, 520, 160},   // Middle left
		{520, 280, 520, 160}, // Middle right
		{0, 440, 520, 160},   // Bottom left
		{520, 440, 520, 160}, // Bottom right
	}

	for i := 0; i < numItems; i++ {
		item, ok := mv.viewModel.Item(i)
		if !ok {
			continue
		}

		pos := positions[i]
		mv.drawMarketItem(screen, item, i, pos[0], pos[1], pos[2], pos[3])
	}
}

// drawMarketItem draws an individual MarketItem
func (mv *MarketView) drawMarketItem(screen *ebiten.Image, item *viewmodel.MarketItemViewModel, index int, x, y, width, height float64) {
	isAvailable := item.Unlocked()

	// Draw CardPack frame (dim if level is insufficient)
	var colorR, colorG, colorB float32 = 0.9, 0.9, 0.9
	if !isAvailable {
		colorR, colorG, colorB = 0.5, 0.5, 0.5 // Dim
	}

	vertices := []ebiten.Vertex{
		{DstX: float32(x), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: float32(x + width), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: float32(x + width), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: float32(x), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// CardPack image (0,120,80,80) -> relative position (0,0,80,80)
	mv.drawCardPackImage(screen, x, y, 80, 80)

	// CardPack name (80,120,440,40) -> relative position (80,0,440,40)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+80, y)
	cardPackName := item.ItemName()
	drawing.DrawText(screen, cardPackName, 28, opt)

	// CardPack description (80,160,440,80) -> relative position (80,40,440,80)
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+80, y+40)
	var description string
	if !isAvailable {
		description = lang.ExecuteTemplate("market-required-level", map[string]any{"level": item.RequiredLevel()})
		drawing.DrawText(screen, description, 20, opt)
	}

	// CardPack price (0,240,520,40) -> relative position (0,120,520,40)
	mv.drawCardPackPrice(screen, item, index, x, y+120, 520, 40)
}

// drawCardPackImage draws the CardPack image
func (mv *MarketView) drawCardPackImage(screen *ebiten.Image, x, y, width, height float64) {
	// 48x64 image (rectangle as a dummy)
	imageX := x + (width-48)/2
	imageY := y + (height-64)/2

	vertices := []ebiten.Vertex{
		{DstX: float32(imageX), DstY: float32(imageY), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: float32(imageX + 48), DstY: float32(imageY), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: float32(imageX + 48), DstY: float32(imageY + 64), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: float32(imageX), DstY: float32(imageY + 64), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
}

// drawCardPackPrice draws the CardPack price
func (mv *MarketView) drawCardPackPrice(screen *ebiten.Image, item *viewmodel.MarketItemViewModel, index int, x, y, width, height float64) {
	// Get price information
	price := item.Price()
	canPurchase := item.CanPurchase()
	sufficiency := item.ResourceSufficiency()

	// Display each resource type in 120x40
	resourceTypes := []struct {
		name  string
		value int
		red   bool
	}{
		{"resource-money", price.Money, !sufficiency.Money},
		{"resource-food", price.Food, !sufficiency.Food},
		{"resource-wood", price.Wood, !sufficiency.Wood},
		{"resource-iron", price.Iron, !sufficiency.Iron},
		{"resource-mana", price.Mana, !sufficiency.Mana},
	}

	currentX := x
	for _, resource := range resourceTypes {
		if resource.value > 0 && currentX < x+width-120 {
			// Resource image (40x40) and Price number (80x40)
			icon := drawing.Image(resource.name)

			// Resource icon
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Scale(2.0, 2.0)
			opt.GeoM.Translate(currentX, y)
			screen.DrawImage(icon, opt)

			// Price number (red if not purchasable)
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(currentX+40, y)
			if !canPurchase || resource.red {
				opt.ColorScale.Scale(1, 0, 0, 1)
			}
			priceText := fmt.Sprintf("%d", resource.value)
			drawing.DrawText(screen, priceText, 24, opt)

			currentX += 120
		}
	}
}
