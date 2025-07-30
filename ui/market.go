package ui

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// MarketView is a widget for displaying the Market.
// Position: Drawn within MainView
type MarketView struct {
	Nation    core.Nation     // MyNation or OtherNation
	GameState *core.GameState // Game state

	// View switching callbacks
	OnBackClicked func()                        // Return to MapGridView
	OnPurchase    func(cardPack *core.CardPack) // When a CardPack is purchased
}

// NewMarketView creates a MarketView
func NewMarketView(onBackClicked func()) *MarketView {
	return &MarketView{
		OnBackClicked: onBackClicked,
	}
}

// SetNation sets the nation to be displayed
func (mv *MarketView) SetNation(nation core.Nation) {
	mv.Nation = nation
}

// SetGameState sets the game state
func (mv *MarketView) SetGameState(gameState *core.GameState) {
	mv.GameState = gameState
}

// HandleInput processes input
func (mv *MarketView) HandleInput(input *Input) error {
	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		cursorX, cursorY := input.Mouse.CursorPosition()

		// Back button click detection (960,40,80,80)
		if cursorX >= 960 && cursorX < 1040 && cursorY >= 40 && cursorY < 120 {
			if mv.OnBackClicked != nil {
				mv.OnBackClicked()
				return nil
			}
		}

		// CardPack click detection and purchase processing
		mv.handleMarketItemClick(cursorX, cursorY)
	}
	return nil
}

// Draw handles the drawing process
func (mv *MarketView) Draw(screen *ebiten.Image) {
	if mv.Nation == nil {
		return
	}

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
	nationName := lang.Text(string(mv.Nation.ID()))
	drawing.DrawText(screen, nationName, 32, opt)

	// Market level text
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(600, 60)
	market := mv.getMarket()
	if market != nil {
		marketLevel := lang.ExecuteTemplate("ui-market-level", map[string]any{"level": market.Level})
		drawing.DrawText(screen, marketLevel, 28, opt)
	}
}

// drawBackButton draws the back button
func (mv *MarketView) drawBackButton(screen *ebiten.Image) {
	DrawButton(screen, 960, 40, 80, 80, "ui-close")
}

// drawMarketItems draws the list of MarketItems
func (mv *MarketView) drawMarketItems(screen *ebiten.Image) {
	marketItems := mv.getAllMarketItems()

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

	for i, item := range marketItems {
		if i >= 6 { // Display up to 6 items
			break
		}

		pos := positions[i]
		mv.drawMarketItem(screen, item, i, pos[0], pos[1], pos[2], pos[3])
	}
}

// drawMarketItem draws an individual MarketItem
func (mv *MarketView) drawMarketItem(screen *ebiten.Image, item *core.MarketItem, index int, x, y, width, height float64) {
	isAvailable := mv.isMarketItemAvailable(item)

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
	cardPack := item.CardPack()
	var cardPackName string
	if cardPack != nil {
		cardPackName = lang.Text(string(cardPack.CardPackID))
	}
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
func (mv *MarketView) drawCardPackPrice(screen *ebiten.Image, item *core.MarketItem, index int, x, y, width, height float64) {
	// Get price information
	_, canPurchase := mv.getCardPackPrice(index)
	price := item.Price()
	subtracted := mv.GameState.Treasury.Resources.Sub(price)

	// Display each resource type in 120x40
	resourceTypes := []struct {
		name  string
		value int
		red   bool
	}{
		{"resource-money", price.Money, subtracted.Money < 0},
		{"resource-food", price.Food, subtracted.Food < 0},
		{"resource-wood", price.Wood, subtracted.Wood < 0},
		{"resource-iron", price.Iron, subtracted.Iron < 0},
		{"resource-mana", price.Mana, subtracted.Mana < 0},
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

// getMarket gets the market for the current nation
func (mv *MarketView) getMarket() *core.Market {
	if mv.Nation == nil || mv.GameState == nil {
		return nil
	}
	return mv.GameState.Markets[mv.Nation.ID()]
}

// getAllMarketItems gets a list of all MarketItems (including insufficient level)
func (mv *MarketView) getAllMarketItems() []*core.MarketItem {
	market := mv.getMarket()
	if market == nil {
		return []*core.MarketItem{}
	}
	return market.Items
}

// getCardPackPrice gets the price and purchasability of a CardPack
func (mv *MarketView) getCardPackPrice(index int) (*core.ResourceQuantity, bool) {
	if mv.GameState.Treasury == nil {
		return nil, false
	}

	market := mv.getMarket()
	if market != nil && index < len(market.Items) {
		item := market.Items[index]
		canPurchase := market.CanPurchase(index, mv.GameState.Treasury)
		price := item.Price()
		return &price, canPurchase
	}

	return nil, false
}

// handleMarketItemClick handles MarketItem clicks
func (mv *MarketView) handleMarketItemClick(cursorX, cursorY int) {
	positions := [][4]int{
		{0, 120, 520, 160},   // Top left
		{520, 120, 520, 160}, // Top right
		{0, 280, 520, 160},   // Middle left
		{520, 280, 520, 160}, // Middle right
		{0, 440, 520, 160},   // Bottom left
		{520, 440, 520, 160}, // Bottom right
	}

	marketItems := mv.getAllMarketItems()

	for i, pos := range positions {
		if i >= len(marketItems) {
			break
		}

		if cursorX >= pos[0] && cursorX < pos[0]+pos[2] &&
			cursorY >= pos[1] && cursorY < pos[1]+pos[3] {
			// MarketItem was clicked
			item := marketItems[i]

			// Cannot purchase if level is insufficient
			if !mv.isMarketItemAvailable(item) {
				return // Do nothing
			}

			if err := mv.PurchaseCardPack(item); err == nil {
				// Return to MapGridView on successful purchase
				if mv.OnBackClicked != nil {
					mv.OnBackClicked()
				}
			}
			break
		}
	}
}

type simpleRand struct {
	*rand.Rand
}

func newSimpleRand() *simpleRand {
	return &simpleRand{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
func (sr *simpleRand) Intn(n int) int {
	return sr.Rand.Intn(n)
}

// PurchaseCardPack purchases a CardPack
func (mv *MarketView) PurchaseCardPack(item *core.MarketItem) error {
	if mv.GameState == nil || mv.Nation == nil {
		return fmt.Errorf("GameState or Nation is nil")
	}

	// Find item index
	itemIndex := -1
	for i, marketItem := range mv.getMarket().Items {
		if marketItem == item {
			itemIndex = i
			break
		}
	}

	if itemIndex == -1 {
		return fmt.Errorf("item not found in market")
	}

	// Purchase processing
	market := mv.getMarket()
	if market == nil {
		return fmt.Errorf("market is nil")
	}

	oldLevel := market.Level
	cardPack, ok := market.Purchase(itemIndex, mv.GameState.Treasury)
	if !ok {
		return fmt.Errorf("purchase failed")
	}

	if int(market.Level) > int(oldLevel) {
		mv.GameState.AddHistory(core.History{
			Turn: mv.GameState.CurrentTurn,
			Key:  "history-market",
			Data: map[string]any{
				"nation": string(mv.Nation.ID()),
				"level":  int(market.Level),
			},
		})
	}

	// Open CardPack and get Cards
	if cardPack != nil {
		rng := newSimpleRand()
		cardIDs := cardPack.Open(rng)

		// Add CardIDs to GameState.CardDeck
		for _, cardID := range cardIDs {
			mv.GameState.CardDeck.Add(cardID)
		}
	}

	mv.GameState.NextTurn()

	return nil
}

// isMarketItemAvailable Determines if a MarketItem is available
func (mv *MarketView) isMarketItemAvailable(item *core.MarketItem) bool {
	if mv.Nation == nil {
		return false
	}
	market := mv.getMarket()
	if market == nil {
		return false
	}
	// Currently, only level is checked
	return market.Level >= item.RequiredLevel()
}
