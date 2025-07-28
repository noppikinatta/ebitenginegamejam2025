package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// CardDeckView is a Widget for the card deck.
// Position: (0,600,1280,120).
// Displays up to 16 cards at 80x120.
type CardDeckView struct {
	CardDeck       *core.CardDeck      // The card deck to display.
	CardGenerator  *core.CardGenerator // Card generator to create card objects from IDs.
	SelectedIndex  int                 // The index of the selected card (-1 for none).
	OnCardSelected func(interface{})   // Callback when a card is selected.

	// New callbacks.
	OnBattleCardClicked    func(*core.BattleCard) bool    // Callback when a BattleCard is clicked.
	OnStructureCardClicked func(*core.StructureCard) bool // Callback when a StructureCard is clicked.

	// Mouse cursor position (set externally).
	MouseX, MouseY int

	HoveredCard interface{}
}

// NewCardDeckView creates a CardDeckView.
func NewCardDeckView(cardDeck *core.CardDeck, cardGenerator *core.CardGenerator) *CardDeckView {
	return &CardDeckView{
		CardDeck:      cardDeck,
		CardGenerator: cardGenerator,
		SelectedIndex: -1, // Nothing is selected initially.
	}
}

// GetSelectedCard gets the selected card.
func (c *CardDeckView) GetSelectedCard() interface{} {
	if c.CardDeck == nil || c.SelectedIndex < 0 {
		return nil
	}

	allCards := c.getAllCards()
	if c.SelectedIndex >= len(allCards) {
		return nil
	}

	return allCards[c.SelectedIndex]
}

// getAllCards gets all cards in a single slice.
func (c *CardDeckView) getAllCards() []interface{} {
	if c.CardDeck == nil || c.CardGenerator == nil {
		return []interface{}{}
	}

	allCards := make([]interface{}, 0)
	cardIDs := c.CardDeck.GetAllCardIDs()

	// Generate cards from CardIDs
	if len(cardIDs) > 0 {
		cards, ok := c.CardGenerator.Generate(cardIDs)
		if ok {
			// Add BattleCards
			for _, card := range cards.BattleCards {
				allCards = append(allCards, card)
			}
			// Add StructureCards
			for _, card := range cards.StructureCards {
				allCards = append(allCards, card)
			}
		}
	}

	return allCards
}

// SelectCard selects a card.
func (c *CardDeckView) SelectCard(index int) {
	if c.CardDeck == nil {
		return
	}

	allCards := c.getAllCards()
	if index < 0 || index >= len(allCards) {
		c.SelectedIndex = -1
		if c.OnCardSelected != nil {
			c.OnCardSelected(nil)
		}
		return
	}

	c.SelectedIndex = index
	if c.OnCardSelected != nil {
		c.OnCardSelected(allCards[index])
	}
}

// ClearSelection clears the selection.
func (c *CardDeckView) ClearSelection() {
	c.SelectedIndex = -1
	if c.OnCardSelected != nil {
		c.OnCardSelected(nil)
	}
}

// RemoveSelectedCard removes the selected card from the deck.
func (c *CardDeckView) RemoveSelectedCard() interface{} {
	if c.CardDeck == nil || c.SelectedIndex < 0 {
		return nil
	}

	allCards := c.getAllCards()
	if c.SelectedIndex >= len(allCards) {
		return nil
	}

	// Get the selected card.
	selectedCard := allCards[c.SelectedIndex]

	// Remove the card from the deck.
	switch card := selectedCard.(type) {
	case *core.BattleCard:
		c.removeBattleCard(card)
	case *core.StructureCard:
		c.removeStructureCard(card)
	}

	// Clear the selection.
	c.ClearSelection()

	return selectedCard
}

// removeBattleCard removes a BattleCard from the deck.
func (c *CardDeckView) removeBattleCard(card *core.BattleCard) {
	if c.CardDeck != nil {
		c.CardDeck.Remove(card.CardID)
	}
}

// removeStructureCard removes a StructureCard from the deck.
func (c *CardDeckView) removeStructureCard(card *core.StructureCard) {
	if c.CardDeck != nil {
		c.CardDeck.Remove(card.ID())
	}
}

// AddCard adds a card to the deck.
func (c *CardDeckView) AddCard(card interface{}) {
	if c.CardDeck == nil {
		return
	}

	switch newCard := card.(type) {
	case *core.BattleCard:
		c.CardDeck.Add(newCard.CardID)
	case *core.StructureCard:
		c.CardDeck.Add(newCard.ID())
	}
}

// HandleInput handles input.
func (c *CardDeckView) HandleInput(input *Input) error {
	cursorX, cursorY := input.Mouse.CursorPosition()
	c.MouseX = cursorX
	c.MouseY = cursorY

	cardIndex := c.CardIndex(cursorX, cursorY)

	if cardIndex != -1 {
		c.HoveredCard = c.getAllCards()[cardIndex]
	} else {
		c.HoveredCard = nil
	}

	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		// Check if within the CardDeckView area (0,600,1280,120).
		if cursorY >= 600 && cursorY < 720 && cursorX >= 0 && cursorX < 1280 {
			c.handleCardClick(cursorX, cursorY)
		}
	}
	return nil
}

func (c *CardDeckView) CardIndex(cursorX, cursorY int) int {
	if cursorX < 0 || cursorX >= 1280 || cursorY < 600 || cursorY >= 720 {
		return -1
	}

	allCards := c.getAllCards()
	cardIndex := cursorX / 80

	if cardIndex < 0 || cardIndex >= len(allCards) || cardIndex >= 16 {
		return -1
	}

	return cardIndex
}

// handleCardClick handles card clicks.
func (c *CardDeckView) handleCardClick(cursorX, cursorY int) {
	if c.CardDeck == nil {
		return
	}

	allCards := c.getAllCards()
	cardIndex := c.CardIndex(cursorX, cursorY)
	if cardIndex == -1 {
		return
	}

	card := allCards[cardIndex]

	// Call the callback according to the card type.
	switch cardData := card.(type) {
	case *core.BattleCard:
		if c.OnBattleCardClicked != nil {
			if c.OnBattleCardClicked(cardData) {
				// If true is returned, remove the card from the deck.
				c.removeBattleCard(cardData)
			}
		}
	case *core.StructureCard:
		if c.OnStructureCardClicked != nil {
			if c.OnStructureCardClicked(cardData) {
				// If true is returned, remove the card from the deck.
				c.removeStructureCard(cardData)
			}
		}
	}
}

// Draw handles drawing.
func (c *CardDeckView) Draw(screen *ebiten.Image) {
	// Draw background.
	c.drawBackground(screen)

	// Draw cards.
	c.drawCards(screen)

	c.drawHoveredCardTooltip(screen)
}

// drawBackground draws the background.
func (c *CardDeckView) drawBackground(screen *ebiten.Image) {
	// CardDeckView background (0,600,1280,120).
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 600, SrcX: 0, SrcY: 0, ColorR: 0.1, ColorG: 0.1, ColorB: 0.15, ColorA: 1},
		{DstX: 1280, DstY: 600, SrcX: 0, SrcY: 0, ColorR: 0.1, ColorG: 0.1, ColorB: 0.15, ColorA: 1},
		{DstX: 1280, DstY: 720, SrcX: 0, SrcY: 0, ColorR: 0.1, ColorG: 0.1, ColorB: 0.15, ColorA: 1},
		{DstX: 0, DstY: 720, SrcX: 0, SrcY: 0, ColorR: 0.1, ColorG: 0.1, ColorB: 0.15, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
}

// drawCards draws the cards.
func (c *CardDeckView) drawCards(screen *ebiten.Image) {
	if c.CardDeck == nil {
		// Message when card deck is empty
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(640, 640)
		drawing.DrawText(screen, lang.Text("card-no-deck"), 24, opt)
		return
	}

	allCards := c.getAllCards()
	if len(allCards) == 0 {
		// Message when there are no cards
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(600, 640)
		drawing.DrawText(screen, lang.Text("card-no-cards"), 24, opt)
		return
	}

	// Draw cards in 80x120 size (max 16)
	for i, card := range allCards {
		if i >= 16 { // Up to 16 cards
			break
		}

		x := float64(i * 80)
		y := float64(600)

		c.drawCard(screen, card, x, y, i == c.SelectedIndex)
	}

	// Abbreviated display if over 16 cards
	if len(allCards) > 16 {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(1160, 690)
		drawing.DrawText(screen, fmt.Sprintf("+%d", len(allCards)-16), 18, opt)
	}
}

// drawCard draws an individual card.
func (c *CardDeckView) drawCard(screen *ebiten.Image, card interface{}, x, y float64, selected bool) {
	switch typedCard := card.(type) {
	case *core.BattleCard:
		DrawBattleCard(screen, x, y, typedCard)
	case *core.StructureCard:
		DrawCard(screen, x, y, string(typedCard.ID()))
	}
}

func (c *CardDeckView) drawHoveredCardTooltip(screen *ebiten.Image) {
	if c.HoveredCard == nil {
		return
	}

	DrawCardDescriptionTooltip(screen, c.HoveredCard, c.MouseX, c.MouseY)
}
