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

		// Back button click detection (480,20,40,40)
		if cursorX >= 480 && cursorX < 520 && cursorY >= 20 && cursorY < 60 {
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

	// Draw header (0,20,480,40)
	mv.drawHeader(screen)

	// Draw back button (480,20,40,40)
	mv.drawBackButton(screen)

	// Draw CardPack list
	mv.drawMarketItems(screen)
}

// drawHeader draws the Nation name header
func (mv *MarketView) drawHeader(screen *ebiten.Image) {
	nationName := mv.Nation.Name()

	// Header background
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.3, ColorA: 1},
		{DstX: 480, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.3, ColorA: 1},
		{DstX: 480, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.3, ColorA: 1},
		{DstX: 0, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.3, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Nation name text
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(10, 30)
	drawing.DrawText(screen, nationName, 16, opt)
}

// drawBackButton draws the back button
func (mv *MarketView) drawBackButton(screen *ebiten.Image) {
	DrawButton(screen, 480, 20, 40, 40, "ui-close")
}

// drawMarketItems draws the list of MarketItems
func (mv *MarketView) drawMarketItems(screen *ebiten.Image) {
	marketItems := mv.getAllMarketItems()

	// CardPack display area: 260x80 x 6
	// Positions: (0,60,260,80), (260,60,260,80), (0,140,260,80), (260,140,260,80), (0,220,260,80), (260,220,260,80)
	positions := [][4]float64{
		{0, 60, 260, 80},    // Top left
		{260, 60, 260, 80},  // Top right
		{0, 140, 260, 80},   // Middle left
		{260, 140, 260, 80}, // Middle right
		{0, 220, 260, 80},   // Bottom left
		{260, 220, 260, 80}, // Bottom right
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

	// CardPack image (0,60,40,40) -> relative position (0,0,40,40)
	mv.drawCardPackImage(screen, x, y, 40, 40)

	// CardPack name (40,60,220,20) -> relative position (40,0,220,20)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+40, y)
	cardPackName := lang.Text(string(item.CardPack.CardPackID))
	drawing.DrawText(screen, cardPackName, 14, opt)

	// CardPack description (40,80,220,40) -> relative position (40,20,220,40)
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+40, y+20)
	var description string
	if !isAvailable {
		description = lang.ExecuteTemplate("market-required-level", map[string]any{"level": item.RequiredLevel})
		drawing.DrawText(screen, description, 10, opt)
	}

	// CardPack price (0,120,260,20) -> relative position (0,60,260,20)
	mv.drawCardPackPrice(screen, item, index, x, y+60, 260, 20)
}

// drawCardPackImage draws the CardPack image
func (mv *MarketView) drawCardPackImage(screen *ebiten.Image, x, y, width, height float64) {
	// 24x32 image (rectangle as a dummy)
	imageX := x + (width-24)/2
	imageY := y + (height-32)/2

	vertices := []ebiten.Vertex{
		{DstX: float32(imageX), DstY: float32(imageY), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: float32(imageX + 24), DstY: float32(imageY), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: float32(imageX + 24), DstY: float32(imageY + 32), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: float32(imageX), DstY: float32(imageY + 32), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
}

// drawCardPackPrice draws the CardPack price
func (mv *MarketView) drawCardPackPrice(screen *ebiten.Image, item *core.MarketItem, index int, x, y, width, height float64) {
	// Get price information
	_, canPurchase := mv.getCardPackPrice(index)
	price := item.Price
	subtracted := mv.GameState.Treasury.Resources.Sub(price)

	// Display each resource type in 60x20
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
		if resource.value > 0 && currentX < x+width-60 {
			// Resource image (20x20) and Price number (40x20)
			icon := drawing.Image(resource.name)

			// Resource icon
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(currentX, y)
			screen.DrawImage(icon, opt)

			// Price number (red if not purchasable)
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(currentX+20, y)
			if !canPurchase || resource.red {
				opt.ColorScale.Scale(1, 0, 0, 1)
			}
			priceText := fmt.Sprintf("%d", resource.value)
			drawing.DrawText(screen, priceText, 12, opt)

			currentX += 60
		}
	}
}

// getAllMarketItems gets a list of all MarketItems (including insufficient level)
func (mv *MarketView) getAllMarketItems() []*core.MarketItem {
	if mv.Nation == nil {
		return []*core.MarketItem{}
	}
	market := mv.Nation.GetMarket()
	if market == nil {
		return []*core.MarketItem{}
	}
	return market.Items
}

// getCardPackPrice gets the price and purchasability of a CardPack
func (mv *MarketView) getCardPackPrice(index int) (*core.ResourceQuantity, bool) {
	if mv.GameState.Treasury == nil || mv.Nation == nil {
		return nil, false
	}

	market := mv.Nation.GetMarket()
	if market != nil && index < len(market.Items) {
		item := market.Items[index]
		canPurchase := mv.Nation.CanPurchase(index, mv.GameState.Treasury)
		return &item.Price, canPurchase
	}

	return nil, false
}

// handleMarketItemClick handles MarketItem clicks
func (mv *MarketView) handleMarketItemClick(cursorX, cursorY int) {
	positions := [][4]int{
		{0, 60, 260, 80},    // 左上
		{260, 60, 260, 80},  // 右上
		{0, 140, 260, 80},   // 左中
		{260, 140, 260, 80}, // 右中
		{0, 220, 260, 80},   // 左下
		{260, 220, 260, 80}, // 右下
	}

	marketItems := mv.getAllMarketItems()

	for i, pos := range positions {
		if i >= len(marketItems) {
			break
		}

		if cursorX >= pos[0] && cursorX < pos[0]+pos[2] &&
			cursorY >= pos[1] && cursorY < pos[1]+pos[3] {
			// MarketItemがクリックされた
			item := marketItems[i]

			// レベル不足の場合は購入できない
			if !mv.isMarketItemAvailable(item) {
				return // 何もしない
			}

			if err := mv.PurchaseCardPack(item); err == nil {
				// 購入成功時、MapGridViewに戻る
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

	market := mv.Nation.GetMarket()
	if market == nil {
		return fmt.Errorf("Market is nil")
	}

	// アイテムのインデックスを見つける
	itemIndex := -1
	for i, marketItem := range market.Items {
		if marketItem == item {
			itemIndex = i
			break
		}
	}

	if itemIndex == -1 {
		return fmt.Errorf("Item not found in market")
	}

	// 購入処理
	cardPack, ok := mv.Nation.Purchase(itemIndex, mv.GameState.Treasury)
	if !ok {
		return fmt.Errorf("Purchase failed")
	}

	// CardPackを開いてCardsを取得
	rng := newSimpleRand()
	cardIDs := cardPack.Open(rng)

	cards, ok := mv.GameState.CardGenerator.Generate(cardIDs)
	if !ok {
		return fmt.Errorf("Card generation failed")
	}

	// GameState.CardDeckに追加
	mv.GameState.CardDeck.Add(cards)

	mv.GameState.NextTurn()

	return nil
}

// isMarketItemAvailable MarketItemが利用可能かどうかを判定
func (mv *MarketView) isMarketItemAvailable(item *core.MarketItem) bool {
	if mv.Nation == nil {
		return false
	}
	market := mv.Nation.GetMarket()
	if market == nil {
		return false
	}
	return market.Level >= item.RequiredLevel
}
