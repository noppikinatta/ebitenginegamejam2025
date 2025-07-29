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
	Territory   *core.Territory       // Territory to display
	TerrainType string                // Terrain name
	GameState   *core.GameState       // Game state

	// ViewModels and Flows
	TerritoryViewModel *viewmodel.TerritoryViewModel
	TerritoryFlow      *flow.TerritoryFlow

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
	tv.TerrainType = terrainType
	
	// Create viewmodel and flow with new territory
	if tv.GameState != nil {
		tv.TerritoryViewModel = viewmodel.NewTerritoryViewModel(tv.GameState, territory)
		tv.TerritoryFlow = flow.NewTerritoryFlow(tv.GameState, territory)
	}
}

// SetGameState sets the game state
func (tv *TerritoryView) SetGameState(gameState *core.GameState) {
	tv.GameState = gameState
	
	// Recreate viewmodel and flow if territory exists
	if tv.Territory != nil {
		tv.TerritoryViewModel = viewmodel.NewTerritoryViewModel(gameState, tv.Territory)
		tv.TerritoryFlow = flow.NewTerritoryFlow(gameState, tv.Territory)
	}
}

// CanPlaceCard checks if a card can be placed
func (tv *TerritoryView) CanPlaceCard() bool {
	if tv.TerritoryFlow != nil {
		return tv.TerritoryFlow.CanPlaceCard()
	}
	// Fallback logic
	if tv.Territory == nil {
		return false
	}
	// TODO: GetCardSlot method not available, use simple check
	return len(tv.Territory.Cards()) < 10 // Default limit
}

// PlaceCard places a StructureCard
func (tv *TerritoryView) PlaceCard(card *core.StructureCard) bool {
	if tv.TerritoryFlow != nil {
		return tv.TerritoryFlow.PlaceCard(card)
	}
	// Fallback logic
	if tv.Territory == nil {
		return false
	}
	return tv.Territory.AppendCard(card)
}

// GetCurrentYield gets the current yield
func (tv *TerritoryView) GetCurrentYield() core.ResourceQuantity {
	// TODO: Implement proper yield calculation once Territory API is available
	if tv.Territory == nil {
		return core.ResourceQuantity{}
	}
	// Simple calculation based on cards
	var total core.ResourceQuantity
	for _, card := range tv.Territory.Cards() {
		// TODO: StructureCard.Yield() method not available yet
		// For now, return default values
		_ = card // Use card to avoid unused variable warning
	}
	return total
}

// GetNewYield gets the predicted yield after changes
func (tv *TerritoryView) GetNewYield() core.ResourceQuantity {
	// For now, return the same as current yield
	return tv.GetCurrentYield()
}

// HandleInput handles input
func (tv *TerritoryView) HandleInput(input *Input) error {
	cursorX, cursorY := input.Mouse.CursorPosition()
	cardIndex := tv.cardIndex(cursorX, cursorY)
	tv.MouseX = cursorX
	tv.MouseY = cursorY

	if cardIndex != -1 && tv.TerritoryViewModel != nil {
		cardVM := tv.TerritoryViewModel.Card(cardIndex)
		if cardVM != nil {
			tv.HoveredCard = cardVM
		}
	} else {
		tv.HoveredCard = nil
	}

	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		if cardIndex != -1 {
			tv.handleCardClick(cardIndex)
		}

		// Click detection for back button (960,40,80,80)
		if cursorX >= 960 && cursorX < 1040 && cursorY >= 40 && cursorY < 120 {
			// Use flow to rollback cards
			if tv.TerritoryFlow != nil {
				tv.TerritoryFlow.Rollback()
			}
			if tv.OnBackClicked != nil {
				tv.OnBackClicked()
				return nil
			}
		}

		// Click detection for confirm button (400,560,240,40)
		if cursorX >= 400 && cursorX < 640 && cursorY >= 560 && cursorY < 600 {
			// Use flow to commit construction changes
			// TODO: Implement Confirm method in TerritoryFlow
			if tv.OnBackClicked != nil {
				tv.OnBackClicked()
				return nil
			}
		}
	}

	return nil
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
	if tv.Territory != nil && cardIndex >= 0 && cardIndex < len(tv.Territory.Cards()) {
		// Remove card from territory
		card, ok := tv.Territory.RemoveCard(cardIndex)
		if ok && tv.OnCardClicked != nil {
			tv.OnCardClicked(card)
		}
	}
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
	// Draw terrain type
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(50, 120)
	drawing.DrawText(screen, fmt.Sprintf("Terrain: %s", tv.TerrainType), 24, opt)

	// Draw card slots
	currentCards := tv.TerritoryViewModel.NumCards()
	maxCards := 10 // Default max cards, TODO: get from terrain
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(50, 150)
	drawing.DrawText(screen, fmt.Sprintf("Cards: %d/%d", currentCards, maxCards), 20, opt)
}

// drawStructureCards draws the placed structure cards
func (tv *TerritoryView) drawStructureCards(screen *ebiten.Image) {
	numCards := tv.TerritoryViewModel.NumCards()
	
	for i := 0; i < numCards; i++ {
		cardVM := tv.TerritoryViewModel.Card(i)
		if cardVM != nil {
			x := float64(100 + i*80)
			y := float64(400)
			
			// Draw card background
			alpha := float32(1.0)
			if cardVM == tv.HoveredCard {
				alpha = 0.8
			}
			DrawCardBackground(screen, x, y, alpha)
			
			// Draw card content (simplified)
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(x+10, y+10)
			drawing.DrawText(screen, cardVM.Name(), 12, opt)
			
			// Draw card name only for now
			// TODO: Add yield display once StructureCardViewModel.Yield() is implemented
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(x+10, y+90)
			drawing.DrawText(screen, "Structure", 8, opt)
		}
	}
	
	// Draw empty card slots
	maxCards := 10 // Default max cards
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
	currentYield := tv.GetCurrentYield()
	predictedYield := tv.GetNewYield()
	
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
