package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/flow"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// CardDeckView is a Widget for the card deck.
// Position: (0,600,1280,120).
// Displays up to 16 cards at 80x120.
type CardDeckView struct {
	ViewModel *viewmodel.CardDeckViewModel // ViewModel for display information
	Flow      *flow.CardDeckFlow           // Flow for operations

	OnCardSelected func(interface{}) // Callback when a card is selected.

	// New callbacks.
	OnBattleCardClicked    func(*core.BattleCard) bool    // Callback when a BattleCard is clicked.
	OnStructureCardClicked func(*core.StructureCard) bool // Callback when a StructureCard is clicked.

	// Mouse cursor position (set externally).
	MouseX, MouseY int

	HoveredCard interface{}
}

// NewCardDeckView creates a CardDeckView.
func NewCardDeckView(viewModel *viewmodel.CardDeckViewModel, flow *flow.CardDeckFlow) *CardDeckView {
	return &CardDeckView{
		ViewModel: viewModel,
		Flow:      flow,
	}
}

// SetGameState sets the GameState reference
func (c *CardDeckView) SetGameState(gameState *core.GameState) {
	// This function is no longer needed as GameState is managed by ViewModel
}

// GetSelectedCard gets the selected card.
func (c *CardDeckView) GetSelectedCard() interface{} {
	return c.Flow.GetSelectedCard()
}

// SetSelectedIndex sets the selected card index
func (c *CardDeckView) SetSelectedIndex(index int) {
	c.Flow.Select(index)
}

// getAllCards gets all cards in a single slice, sorted by display order.
func (c *CardDeckView) getAllCards() []interface{} {
	return c.Flow.GetAllCards()
}

// SetDisplayOrder sets the display order of cards (legacy method for compatibility)
func (c *CardDeckView) SetDisplayOrder(order []core.CardID) {
	// This functionality should be moved to flow or viewmodel
	// For now, we'll ignore this as it's handled by the viewmodel
}

// SelectCard selects a card at the specified index
func (c *CardDeckView) SelectCard(index int) {
	c.Flow.Select(index)
}

// RemoveSelectedCard removes the currently selected card (legacy method)
func (c *CardDeckView) RemoveSelectedCard() {
	// This operation should be handled by flows in other components
	// For now, we'll implement a simple clear selection
	c.Flow.ClearSelection()
}

// AddCard adds a card back to the deck (legacy method)
func (c *CardDeckView) AddCard(card interface{}) {
	// This operation should be handled by the appropriate flow
	// For now, this is a placeholder
}

// HandleInput handles input for card selection and clicking.
func (c *CardDeckView) HandleInput(input *Input) error {
	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		cursorX, cursorY := input.Mouse.CursorPosition()

		// Determine card click area
		index := c.cardIndex(cursorX, cursorY)
		if index != -1 {
			// Select the card
			c.Flow.Select(index)

			// Notify selection
			if c.OnCardSelected != nil {
				selectedCard := c.Flow.GetSelectedCard()
				c.OnCardSelected(selectedCard)
			}

			// Handle card type specific clicks
			selectedCard := c.Flow.GetSelectedCard()
			if battleCard, ok := selectedCard.(*core.BattleCard); ok {
				if c.OnBattleCardClicked != nil {
					c.OnBattleCardClicked(battleCard)
				}
			} else if structureCard, ok := selectedCard.(*core.StructureCard); ok {
				if c.OnStructureCardClicked != nil {
					c.OnStructureCardClicked(structureCard)
				}
			}
		}
	}

	// Update hovered card
	cursorX, cursorY := input.Mouse.CursorPosition()
	index := c.cardIndex(cursorX, cursorY)
	if index != -1 {
		allCards := c.getAllCards()
		if index < len(allCards) {
			c.HoveredCard = allCards[index]
		}
	} else {
		c.HoveredCard = nil
	}

	return nil
}

// cardIndex calculates which card index the cursor is over
func (c *CardDeckView) cardIndex(cursorX, cursorY int) int {
	// Card deck area: (0,600,1280,120)
	if cursorY < 600 || cursorY >= 720 {
		return -1
	}

	// Each card is 80px wide
	cardX := cursorX / 80
	if cardX >= 16 { // Maximum 16 cards displayed
		return -1
	}

	allCards := c.getAllCards()
	if cardX >= len(allCards) {
		return -1
	}

	return cardX
}

// Draw draws all cards in the deck.
func (c *CardDeckView) Draw(screen *ebiten.Image) {
	allCards := c.getAllCards()
	selectedIndex := c.Flow.GetSelectedIndex()

	for i, card := range allCards {
		if i >= 16 { // Display only the first 16 cards
			break
		}

		x := float64(i * 80)
		y := float64(600)

		// Highlight if selected
		alpha := float32(1.0)
		if i == selectedIndex {
			alpha = 0.7
		}
		if card == c.HoveredCard {
			alpha = 0.8
		}

		// Draw card background
		DrawCardBackground(screen, x, y, alpha)

		// Draw card
		switch typedCard := card.(type) {
		case *core.BattleCard:
			DrawBattleCard(screen, x, y, typedCard)
		case *core.StructureCard:
			DrawCard(screen, x, y, string(typedCard.ID()))
		}

		// Display duplicates count
		duplicates := c.getDuplicateCount(card)
		if duplicates > 1 {
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(x+60, y+100)
			drawing.DrawText(screen, fmt.Sprintf("x%d", duplicates), 12, opt)
		}
	}
}

// getDuplicateCount returns the number of duplicates for a card
func (c *CardDeckView) getDuplicateCount(card interface{}) int {
	return c.ViewModel.GetDuplicateCount(card)
}
