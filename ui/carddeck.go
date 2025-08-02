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

	HoveredCardIndex int
}

// NewCardDeckView creates a CardDeckView.
func NewCardDeckView(viewModel *viewmodel.CardDeckViewModel, flow *flow.CardDeckFlow) *CardDeckView {
	return &CardDeckView{
		ViewModel: viewModel,
		Flow:      flow,
	}
}

// HandleInput handles input for card selection and clicking.
func (c *CardDeckView) HandleInput(input *Input) error {
	cursorX, cursorY := input.Mouse.CursorPosition()
	index := c.cardIndex(cursorX, cursorY)

	if index == -1 {
		c.HoveredCardIndex = -1
		return nil
	}

	c.HoveredCardIndex = index

	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		c.Flow.Select(index)
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
	idx := cursorX / 80
	length := c.ViewModel.CountTypesInHand()
	if idx >= length {
		return -1
	}

	return idx
}

// Draw draws all cards in the deck.
func (c *CardDeckView) Draw(screen *ebiten.Image) {
	length := c.ViewModel.CountTypesInHand()
	for i := range length {
		x, y, width, height := c.boundsForCard(i)
	}

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

func (c *CardDeckView) boundsForCard(index int) (x, y, width, height int) {
	return index * 80, 600, 80, 120
}
