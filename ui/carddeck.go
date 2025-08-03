package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/flow"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// CardDeckView is a Widget for the card deck.
// Position: (0,600,1280,120).
// Displays up to 16 cards at 80x120.
type CardDeckView struct {
	centerViewModer CenterViewModer
	ViewModel       *viewmodel.CardDeckViewModel // ViewModel for display information
	Flow            *flow.CardDeckFlow           // Flow for operations

	HoveredCardIndex int
}

// NewCardDeckView creates a CardDeckView.
func NewCardDeckView(centerViewModer CenterViewModer, viewModel *viewmodel.CardDeckViewModel, flow *flow.CardDeckFlow) *CardDeckView {
	return &CardDeckView{
		centerViewModer: centerViewModer,
		ViewModel:       viewModel,
		Flow:            flow,
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
		c.clickCard(index)
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

func (c *CardDeckView) clickCard(idx int) {
	cardID, ok := c.ViewModel.CardID(idx)
	if !ok {
		return
	}

	switch c.centerViewModer.CurrentViewMode() {
	case ViewTypeBattle:
		c.clickBattleCard(cardID)
	case ViewTypeTerritory:
		c.clickStructureCard(cardID)
	}
}

func (c *CardDeckView) clickBattleCard(cardID core.CardID) {
	c.Flow.PlayBattleCardInBattle(cardID)
}

func (c *CardDeckView) clickStructureCard(cardID core.CardID) {
	c.Flow.PlayStructureCardInTerritory(cardID)
}

// Draw draws all cards in the deck.
func (c *CardDeckView) Draw(screen *ebiten.Image) {
	length := c.ViewModel.CountTypesInHand()
	for i := range length {
		c.drawACard(screen, i)
	}
}

func (c *CardDeckView) drawACard(screen *ebiten.Image, idx int) {
	x, y := c.locationForCard(idx)

	card, ok := c.ViewModel.Card(idx)
	if !ok {
		DrawCardBackground(screen, float64(x), float64(y), 0.5)
		return
	}

	DrawCard(screen, x, y, card, idx == c.HoveredCardIndex)

}

func (c *CardDeckView) locationForCard(index int) (x, y float64) {
	return float64(index) * 80, 600
}
