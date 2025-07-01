package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// TerritoryView is a widget for displaying a Territory.
// Position: Drawn within MainView
type TerritoryView struct {
	Territory   *core.Territory       // Territory to display
	OldCards    []*core.StructureCard // Cards before changes
	TerrainType string                // Terrain name
	GameState   *core.GameState       // Game state
	HoveredCard interface{}
	MouseX      int
	MouseY      int

	// View switching callbacks
	OnBackClicked func()                         // Return to MapGridView
	OnCardClicked func(card *core.StructureCard) // When a card is clicked (return to CardDeck)
}

// NewTerritoryView creates a TerritoryView
func NewTerritoryView(onBackClicked func()) *TerritoryView {
	return &TerritoryView{
		OnBackClicked: onBackClicked,
	}
}

// SetTerritory sets the Territory to display
func (tv *TerritoryView) SetTerritory(territory *core.Territory, terrainType string) {
	tv.Territory = territory
	tv.OldCards = make([]*core.StructureCard, len(territory.Cards))
	copy(tv.OldCards, territory.Cards)
	tv.TerrainType = terrainType
}

// SetGameState sets the game state
func (tv *TerritoryView) SetGameState(gameState *core.GameState) {
	tv.GameState = gameState
}

// AddStructureCard places a StructureCard
func (tv *TerritoryView) AddStructureCard(card *core.StructureCard) bool {
	if tv.Territory == nil {
		return false
	}
	return tv.Territory.AppendCard(card)
}

// RemoveStructureCard removes a StructureCard
func (tv *TerritoryView) RemoveStructureCard(index int) *core.StructureCard {
	if tv.Territory == nil {
		return nil
	}
	card, ok := tv.Territory.RemoveCard(index)
	if !ok {
		return nil
	}
	return card
}

// GetCurrentYield gets the current yield
func (tv *TerritoryView) GetCurrentYield() core.ResourceQuantity {
	if tv.Territory == nil {
		return core.ResourceQuantity{}
	}
	return tv.Territory.Yield()
}

// HandleInput processes input
func (tv *TerritoryView) HandleInput(input *Input) error {
	cursorX, cursorY := input.Mouse.CursorPosition()
	tv.MouseX = cursorX
	tv.MouseY = cursorY
	cardIndex := tv.cardIndex(cursorX, cursorY)
	if cardIndex != -1 {
		tv.HoveredCard = tv.Territory.Cards[cardIndex]
	} else {
		tv.HoveredCard = nil
	}

	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {

		// Back button click detection (960,40,80,80)
		if cursorX >= 960 && cursorX < 1040 && cursorY >= 40 && cursorY < 120 {
			// Execute construction confirmation if there are changes
			if tv.IsChanged() {
				tv.ConfirmConstruction()
			}

			if tv.OnBackClicked != nil {
				tv.OnBackClicked()
				return nil
			}
		}

		// Construction confirmation button click detection (400,440,240,80)
		if cursorX >= 400 && cursorX < 640 && cursorY >= 440 && cursorY < 520 {
			// Execute construction confirmation if there are changes
			if tv.IsChanged() {
				tv.ConfirmConstruction()
			}
			if tv.OnBackClicked != nil {
				tv.OnBackClicked()
				return nil
			}
		}

		// StructureCard click detection (return to CardDeck)
		tv.handleStructureCardClick(cursorX, cursorY)
	}
	return nil
}

func (tv *TerritoryView) cardIndex(cursorX, cursorY int) int {
	if cursorX < 0 || cursorX >= 1040 || cursorY < 320 || cursorY >= 440 {
		return -1
	}
	cardIndex := cursorX / 80
	if cardIndex < 0 || cardIndex >= len(tv.Territory.Cards) {
		return -1
	}
	return cardIndex
}

// Draw handles the drawing process
func (tv *TerritoryView) Draw(screen *ebiten.Image) {
	// Draw header (0,40,1040,80)
	tv.drawHeader(screen)

	// Draw back button (960,40,80,80)
	tv.drawBackButton(screen)

	// Draw change indicator (880,40,80,80)
	tv.drawChangeIndicator(screen)

	// Draw yield display (0,120,120,200)
	tv.drawYield(screen)

	// Draw effect description (120,120,920,200)
	tv.drawEffectDescription(screen)

	// Draw StructureCard slots (0,320,1040,120)
	tv.drawStructureCards(screen)

	// Draw construction confirmation button (400,440,240,80)
	tv.drawConstructionButton(screen)

	tv.drawHoveredCardTooltip(screen)
}

// drawHeader draws the header
func (tv *TerritoryView) drawHeader(screen *ebiten.Image) {
	// Header background
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 40, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: 1040, DstY: 40, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: 1040, DstY: 120, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: 0, DstY: 120, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Title text
	terrainName := lang.Text(tv.TerrainType)

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(20, 60)
	drawing.DrawText(screen, terrainName, 32, opt)
}

// drawBackButton draws the back button
func (tv *TerritoryView) drawBackButton(screen *ebiten.Image) {
	DrawButton(screen, 960, 40, 80, 80, "ui-close")
}

// drawYield draws the yield display
func (tv *TerritoryView) drawYield(screen *ebiten.Image) {
	// Yield display background (0,120,120,200)
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 120, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 120, DstY: 120, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 120, DstY: 320, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 0, DstY: 320, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Get current yield
	currentYield := tv.GetCurrentYield()

	// Display 5 resource types in 120x40
	resourceTypes := []struct {
		name  string
		value int
	}{
		{"resource-money", currentYield.Money},
		{"resource-food", currentYield.Food},
		{"resource-wood", currentYield.Wood},
		{"resource-iron", currentYield.Iron},
		{"resource-mana", currentYield.Mana},
	}

	for i, resource := range resourceTypes {
		y := 120.0 + float64(i)*40

		// Resource image (40x40)
		icon := drawing.Image(resource.name)
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Scale(2.0, 2.0)
		opt.GeoM.Translate(10, y)
		screen.DrawImage(icon, opt)

		// Yield number (80x40)
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(50, y)
		yieldText := fmt.Sprintf("%d", resource.value)
		drawing.DrawText(screen, yieldText, 24, opt)
	}
}

// drawEffectDescription draws the effect description
func (tv *TerritoryView) drawEffectDescription(screen *ebiten.Image) {
	// Effect description background (120,120,920,200)
	vertices := []ebiten.Vertex{
		{DstX: 120, DstY: 120, SrcX: 0, SrcY: 0, ColorR: 0.25, ColorG: 0.25, ColorB: 0.25, ColorA: 1},
		{DstX: 1040, DstY: 120, SrcX: 0, SrcY: 0, ColorR: 0.25, ColorG: 0.25, ColorB: 0.25, ColorA: 1},
		{DstX: 1040, DstY: 320, SrcX: 0, SrcY: 0, ColorR: 0.25, ColorG: 0.25, ColorB: 0.25, ColorA: 1},
		{DstX: 120, DstY: 320, SrcX: 0, SrcY: 0, ColorR: 0.25, ColorG: 0.25, ColorB: 0.25, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Title
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(130, 130)
	drawing.DrawText(screen, lang.Text("territory-structure-effects"), 24, opt)

	if tv.Territory == nil || len(tv.Territory.Cards) == 0 {
		// When no cards are placed
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(130, 170)
		drawing.DrawText(screen, lang.Text("territory-no-structures"), 20, opt)

		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(130, 210)
		drawing.DrawText(screen, lang.Text("territory-place-cards"), 20, opt)
		return
	}

	// Display effects of placed StructureCards
	startY := 170.0
	for i, card := range tv.Territory.Cards {
		if i >= 4 { // Display up to 4 cards maximum
			break
		}

		y := startY + float64(i)*36

		// Card name
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(130, y)
		cardName := lang.Text(string(card.CardID))
		drawing.DrawText(screen, cardName, 20, opt)

		// Effect description
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(400, y)
		effect := lang.Text(string(card.DescriptionKey))
		drawing.DrawText(screen, effect, 18, opt)
	}
}

// drawStructureCards draws the StructureCard slots
func (tv *TerritoryView) drawStructureCards(screen *ebiten.Image) {
	// Slot background
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 320, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 1040, DstY: 320, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 1040, DstY: 440, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
		{DstX: 0, DstY: 440, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.3, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Draw deployed StructureCards (using temporary storage tv.Cards)
	for i, card := range tv.Territory.Cards {
		cardX := float64(i * 80)
		cardY := 320.0

		DrawCard(screen, cardX, cardY, string(card.CardID))
	}

	// Display empty slots
	if tv.Territory != nil {
		maxSlots := tv.Territory.CardSlot
		for i := len(tv.Territory.Cards); i < maxSlots && i < 13; i++ { // Display up to 13 cards max (1040÷80=13)
			cardX := float64(i * 80)
			cardY := 320.0

			// Empty slot border
			vertices := []ebiten.Vertex{
				{DstX: float32(cardX), DstY: float32(cardY), SrcX: 0, SrcY: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5, ColorA: 0.5},
				{DstX: float32(cardX + 80), DstY: float32(cardY), SrcX: 0, SrcY: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5, ColorA: 0.5},
				{DstX: float32(cardX + 80), DstY: float32(cardY + 120), SrcX: 0, SrcY: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5, ColorA: 0.5},
				{DstX: float32(cardX), DstY: float32(cardY + 120), SrcX: 0, SrcY: 0, ColorR: 0.5, ColorG: 0.5, ColorB: 0.5, ColorA: 0.5},
			}
			indices := []uint16{0, 1, 2, 0, 2, 3}
			screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
		}
	}
}

func (tv *TerritoryView) drawHoveredCardTooltip(screen *ebiten.Image) {
	if tv.HoveredCard == nil {
		return
	}

	DrawCardDescriptionTooltip(screen, tv.HoveredCard, tv.MouseX, tv.MouseY)
}

// handleStructureCardClick handles StructureCard clicks
func (tv *TerritoryView) handleStructureCardClick(cursorX, cursorY int) {
	// Each card (80x120) in the StructureCard area (0,320,1040,120)
	if cursorY >= 320 && cursorY < 440 {
		cardIndex := cursorX / 80

		if cardIndex >= 0 && cardIndex < len(tv.Territory.Cards) {
			targetCard := tv.Territory.Cards[cardIndex]

			// Return card to CardDeck
			tv.RemoveCard(targetCard)
			if tv.OnCardClicked != nil {
				tv.OnCardClicked(targetCard)
			}
		}
	}
}

// drawChangeIndicator draws the change indicator
func (tv *TerritoryView) drawChangeIndicator(screen *ebiten.Image) {
	// Don't display anything if there are no changes
	if !tv.IsChanged() {
		return
	}
	// Background for change status display (880,40,80,80)
	vertices := []ebiten.Vertex{
		{DstX: 880, DstY: 40, SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.6, ColorB: 0.2, ColorA: 1},
		{DstX: 960, DstY: 40, SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.6, ColorB: 0.2, ColorA: 1},
		{DstX: 960, DstY: 120, SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.6, ColorB: 0.2, ColorA: 1},
		{DstX: 880, DstY: 120, SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.6, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Display "*" mark
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(910, 70)
	drawing.DrawText(screen, "*", 40, opt)
}

// drawConstructionButton draws the construction confirmation button
func (tv *TerritoryView) drawConstructionButton(screen *ebiten.Image) {
	isChanged := tv.IsChanged()

	// Determine button color
	var colorR, colorG, colorB float32 = 0.4, 0.4, 0.4 // Gray when no changes
	var buttonText string = lang.Text("ui-no-changes")

	if isChanged {
		colorR, colorG, colorB = 0.2, 0.6, 0.8 // Blue when there are changes
		buttonText = lang.Text("ui-confirm")
	}

	// Button background (400,440,240,80)
	vertices := []ebiten.Vertex{
		{DstX: 400, DstY: 440, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: 640, DstY: 440, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: 640, DstY: 520, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: 400, DstY: 520, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Button text
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(440, 470)
	drawing.DrawText(screen, buttonText, 24, opt)
}

// IsChanged checks if the cards have been changed
func (tv *TerritoryView) IsChanged() bool {
	if tv.Territory == nil {
		return false
	}
	if len(tv.Territory.Cards) != len(tv.OldCards) {
		return true
	}

	// Check if all StructureCard pointers in both slices are the same (order doesn't matter)
	for _, oldCard := range tv.OldCards {
		found := false
		for _, territoryCard := range tv.Territory.Cards {
			if oldCard == territoryCard {
				found = true
				break
			}
		}
		if !found {
			return true
		}
	}

	return false
}

// ConfirmConstruction Construction confirmation process
func (tv *TerritoryView) ConfirmConstruction() {
	if tv.Territory == nil {
		return
	}
}

// CanPlaceCard checks if a card can be placed
func (tv *TerritoryView) CanPlaceCard() bool {
	if tv.Territory == nil {
		return false
	}
	return len(tv.Territory.Cards) < tv.Territory.CardSlot
}

// PlaceCard places a card
func (tv *TerritoryView) PlaceCard(card *core.StructureCard) bool {
	if !tv.CanPlaceCard() {
		return false
	}
	tv.Territory.Cards = append(tv.Territory.Cards, card)
	return true
}

// RemoveCard removes a card
func (tv *TerritoryView) RemoveCard(card *core.StructureCard) bool {
	// Find card index
	cardIndex := -1
	for i, structureCard := range tv.Territory.Cards {
		if structureCard == card {
			cardIndex = i
			break
		}
	}

	if cardIndex == -1 {
		return false
	}

	// Remove from Cards
	tv.Territory.Cards = append(tv.Territory.Cards[:cardIndex], tv.Territory.Cards[cardIndex+1:]...)

	// Add to GameState.CardDeck
	if tv.GameState != nil {
		cards := &core.Cards{StructureCards: []*core.StructureCard{card}}
		tv.GameState.CardDeck.Add(cards)
	}

	return true
}
