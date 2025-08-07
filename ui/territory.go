package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/flow"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// TerritoryView is a widget for displaying a Territory.
// Position: Drawn within MainView
type TerritoryView struct {
	TerritoryViewModel *viewmodel.TerritoryViewModel
	TerritoryFlow      *flow.TerritoryFlow

	hoveredCardIndex int
}

// NewTerritoryView creates a TerritoryView
func NewTerritoryView(territoryFlow *flow.TerritoryFlow) *TerritoryView {
	return &TerritoryView{
		TerritoryFlow: territoryFlow,
	}
}

func (tv *TerritoryView) Select(x, y int) {
	vm, ok := tv.TerritoryFlow.SelectTerritory(x, y)
	if ok {
		tv.TerritoryViewModel = vm
	}
}

// HandleInput handles input
func (tv *TerritoryView) HandleInput(input *Input) (back bool, err error) {
	cursorX, cursorY := input.Mouse.CursorPosition()
	cardIndex := tv.cardIndex(cursorX, cursorY)
	tv.hoveredCardIndex = cardIndex

	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		if cardIndex != -1 {
			tv.handleCardClick(cardIndex)
		}

		// Click detection for back button (960,40,80,80)
		if cursorX >= 960 && cursorX < 1040 && cursorY >= 40 && cursorY < 120 {
			tv.TerritoryFlow.Rollback()
			return true, nil
		}

		// Click detection for confirm button (400,560,240,40)
		if cursorX >= 400 && cursorX < 640 && cursorY >= 560 && cursorY < 600 {
			tv.TerritoryFlow.Commit()
			return true, nil
		}
	}

	return false, nil
}

// cardIndex calculates which card index the cursor is over
func (tv *TerritoryView) cardIndex(cursorX, cursorY int) int {
	// Territory card area calculation (simplified)
	if cursorY < 400 || cursorY >= 520 { // Card area roughly
		return -1
	}

	cardX := (cursorX - 100) / 80 // Assuming cards start at x=100 and are 80px wide
	if cardX < 0 {
		return -1
	}

	if tv.TerritoryViewModel != nil {
		numCards := tv.TerritoryViewModel.NumCards()
		if cardX >= numCards {
			return -1
		}
	}

	return cardX
}

// handleCardClick handles clicking on territory cards
func (tv *TerritoryView) handleCardClick(cardIndex int) {
	tv.TerritoryFlow.RemoveFromPlan(cardIndex)
}

// Draw handles drawing
func (tv *TerritoryView) Draw(screen *ebiten.Image) {
	if tv.TerritoryViewModel == nil {
		// Draw a message if no territory is set up
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(400, 300)
		drawing.DrawText(screen, "No territory selected", 24, opt)
		return
	}

	// Draw territory title
	title := tv.TerritoryViewModel.Title()
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(20, 60)
	drawing.DrawText(screen, title, 32, opt)

	// Draw territory information
	tv.drawTerritoryInfo(screen)

	// Draw structure cards
	tv.drawStructureCards(screen)

	// Draw UI buttons
	tv.drawButtons(screen)

	// Draw yield information
	tv.drawYieldInfo(screen)
}

// drawTerritoryInfo draws territory information
func (tv *TerritoryView) drawTerritoryInfo(screen *ebiten.Image) {
	// Draw card slots
	currentCards := tv.TerritoryViewModel.NumCards()
	maxCards := tv.TerritoryViewModel.CardSlot()
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(50, 120)
	drawing.DrawText(screen, fmt.Sprintf("Cards: %d/%d", currentCards, maxCards), 20, opt)
}

// drawStructureCards draws the placed structure cards
func (tv *TerritoryView) drawStructureCards(screen *ebiten.Image) {
	numCards := tv.TerritoryViewModel.NumCards()

	for i := 0; i < numCards; i++ {
		x := float64(100 + i*80)
		y := float64(400)

		card, ok := tv.TerritoryViewModel.Card(i)
		if !ok {
			DrawCardBackground(screen, x, y, 0.5)
			continue
		}

		DrawCard(screen, x, y, card, i == tv.hoveredCardIndex)
	}

	// Draw empty card slots
	maxCards := tv.TerritoryViewModel.CardSlot()
	for i := numCards; i < maxCards; i++ {
		x := float64(100 + i*80)
		y := float64(400)
		DrawCardBackground(screen, x, y, 0.3)
	}
}

// drawButtons draws UI buttons
func (tv *TerritoryView) drawButtons(screen *ebiten.Image) {
	// Back button (960,40,80,80)
	drawing.DrawRect(screen, 960, 40, 80, 80, 0.3, 0.3, 0.3, 1.0)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(980, 65)
	drawing.DrawText(screen, "Back", 20, opt)

	// Confirm button (400,560,240,40)
	drawing.DrawRect(screen, 400, 560, 240, 40, 0.2, 0.6, 0.2, 1.0)
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(480, 575)
	drawing.DrawText(screen, "Confirm", 20, opt)
}

// drawYieldInfo draws yield information
func (tv *TerritoryView) drawYieldInfo(screen *ebiten.Image) {
	currentYield := tv.TerritoryViewModel.CurrentYield()
	predictedYield := tv.TerritoryViewModel.PredictedYield()

	// Current yield
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(400, 200)
	drawing.DrawText(screen, "Current Yield:", 20, opt)

	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(400, 220)
	currentText := fmt.Sprintf("M:%.0f F:%.0f W:%.0f I:%.0f A:%.0f",
		currentYield.Money, currentYield.Food, currentYield.Wood, currentYield.Iron, currentYield.Mana)
	drawing.DrawText(screen, currentText, 16, opt)

	// Predicted yield
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(400, 260)
	drawing.DrawText(screen, "Predicted Yield:", 20, opt)

	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(400, 280)
	predictedText := fmt.Sprintf("M:%.0f F:%.0f W:%.0f I:%.0f A:%.0f",
		predictedYield.Money, predictedYield.Food, predictedYield.Wood, predictedYield.Iron, predictedYield.Mana)
	drawing.DrawText(screen, predictedText, 16, opt)

	// Calculate difference
	diff := core.ResourceQuantity{
		Money: predictedYield.Money - currentYield.Money,
		Food:  predictedYield.Food - currentYield.Food,
		Wood:  predictedYield.Wood - currentYield.Wood,
		Iron:  predictedYield.Iron - currentYield.Iron,
		Mana:  predictedYield.Mana - currentYield.Mana,
	}

	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(400, 320)
	drawing.DrawText(screen, "Difference:", 20, opt)

	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(400, 340)
	diffText := fmt.Sprintf("M:%+.0f F:%+.0f W:%+.0f I:%+.0f A:%+.0f",
		diff.Money, diff.Food, diff.Wood, diff.Iron, diff.Mana)
	drawing.DrawText(screen, diffText, 16, opt)
}
