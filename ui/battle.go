package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/flow"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// BattleView is a Widget for displaying battles.
// Position: Drawn within MainView.
type BattleView struct {
	BattlePoint core.BattlePoint  // The point to be battled.
	PointName   string            // The name of the battle point.
	Battlefield *core.Battlefield // Battlefield information.
	GameState   *core.GameState   // Game state.

	// ViewModels and Flows
	BattleViewModel *viewmodel.BattleViewModel
	BattleFlow      *flow.BattleFlow

	HoveredCard interface{}
	MouseX      int
	MouseY      int

	// Callback for view switching.
	OnBackClicked func() // Return to MapGridView.
}

// NewBattleView creates a BattleView.
func NewBattleView(onBackClicked func()) *BattleView {
	return &BattleView{
		OnBackClicked: onBackClicked,
	}
}

// SetBattlePoint sets the battle point to be displayed.
func (bv *BattleView) SetBattlePoint(point core.BattlePoint) {
	bv.BattlePoint = point
	bv.Battlefield = bv.createBattlefield(point)
	
	// Create viewmodel and flow with new battlefield
	if bv.GameState != nil {
		bv.BattleViewModel = viewmodel.NewBattleViewModel(bv.GameState, bv.Battlefield, point)
		bv.BattleFlow = flow.NewBattleFlow(bv.GameState, bv.Battlefield)
	}
}

// SetPointName sets the name of the battle point.
func (bv *BattleView) SetPointName(pointName string) {
	bv.PointName = pointName
}

// SetGameState sets the game state.
func (bv *BattleView) SetGameState(gameState *core.GameState) {
	bv.GameState = gameState
	
	// Recreate viewmodel and flow if battlefield exists
	if bv.Battlefield != nil && bv.BattlePoint != nil {
		bv.BattleViewModel = viewmodel.NewBattleViewModel(gameState, bv.Battlefield, bv.BattlePoint)
		bv.BattleFlow = flow.NewBattleFlow(gameState, bv.Battlefield)
	}
}

// GetTotalPower calculates the total Power value of the placed BattleCards.
func (bv *BattleView) GetTotalPower() float64 {
	if bv.BattleViewModel != nil {
		return bv.BattleViewModel.TotalPower()
	}
	// Fallback to direct calculation
	if bv.Battlefield != nil {
		return bv.Battlefield.CalculateTotalPower()
	}
	return 0.0
}

// CanDefeatEnemy determines if the enemy can be defeated.
func (bv *BattleView) CanDefeatEnemy() bool {
	if bv.BattleViewModel != nil {
		return bv.BattleViewModel.CanBeat()
	}
	// Fallback logic
	if bv.Battlefield != nil {
		return bv.Battlefield.CanBeat()
	}
	if bv.BattlePoint == nil {
		return false
	}
	return bv.GetTotalPower() >= bv.BattlePoint.Enemy().Power()
}

// CanPlaceCard checks if a card can be placed
func (bv *BattleView) CanPlaceCard() bool {
	if bv.BattleFlow != nil {
		return bv.BattleFlow.CanPlaceCard()
	}
	// Fallback logic
	if bv.Battlefield == nil {
		return false
	}
	return len(bv.Battlefield.BattleCards) < bv.Battlefield.CardSlot
}

// PlaceCard places a battle card on the battlefield
func (bv *BattleView) PlaceCard(card *core.BattleCard) bool {
	if bv.BattleFlow != nil {
		return bv.BattleFlow.PlaceCard(card)
	}
	// Fallback logic
	return false
}

// HandleInput handles input.
func (bv *BattleView) HandleInput(input *Input) error {
	cursorX, cursorY := input.Mouse.CursorPosition()
	cardIndex := bv.cardIndex(cursorX, cursorY)
	bv.MouseX = cursorX
	bv.MouseY = cursorY

	if cardIndex != -1 && bv.BattleViewModel != nil {
		cardVM := bv.BattleViewModel.Card(cardIndex)
		if cardVM != nil {
			bv.HoveredCard = cardVM
		}
	} else {
		bv.HoveredCard = nil
	}

	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		if cardIndex != -1 {
			bv.handleBattleCardClick(cursorX, cursorY)
		}

		// Click detection for the back button (960,40,80,80).
		if cursorX >= 960 && cursorX < 1040 && cursorY >= 40 && cursorY < 120 {
			// Use flow to rollback cards
			if bv.BattleFlow != nil {
				bv.BattleFlow.Rollback()
			}
			if bv.OnBackClicked != nil {
				bv.OnBackClicked()
				return nil
			}
		}

		// Click detection for the conquer button (400,560,240,40).
		if cursorX >= 400 && cursorX < 640 && cursorY >= 560 && cursorY < 600 {
			if bv.CanDefeatEnemy() {
				// Use flow to attempt conquest
				if bv.BattleFlow != nil && bv.BattleFlow.Conquer() {
					// Conquest successful, return to map view
					if bv.OnBackClicked != nil {
						bv.OnBackClicked()
						return nil
					}
				}
			}
		}
	}

	return nil
}

// cardIndex calculates which card index the cursor is over
func (bv *BattleView) cardIndex(cursorX, cursorY int) int {
	// Battle card area calculation
	// This is a simplified implementation
	if cursorY < 400 || cursorY >= 520 { // Card area roughly
		return -1
	}
	
	cardX := (cursorX - 100) / 80 // Assuming cards start at x=100 and are 80px wide
	if cardX < 0 {
		return -1
	}
	
	if bv.BattleViewModel != nil {
		numCards := bv.BattleViewModel.NumCards()
		if cardX >= numCards {
			return -1
		}
	}
	
	return cardX
}

// handleBattleCardClick handles clicking on battle cards
func (bv *BattleView) handleBattleCardClick(cursorX, cursorY int) {
	cardIndex := bv.cardIndex(cursorX, cursorY)
	if cardIndex != -1 && bv.BattleFlow != nil {
		// Remove card from battlefield
		bv.BattleFlow.RemoveFromBattle(cardIndex)
	}
}

// createBattlefield creates a battlefield from a battle point
func (bv *BattleView) createBattlefield(point core.BattlePoint) *core.Battlefield {
	if point == nil || point.Enemy() == nil {
		return nil
	}
	
	enemy := point.Enemy()
	supportPower := 0.0 // TODO: Calculate support power from adjacent territories
	
	return core.NewBattlefield(enemy, supportPower)
}

// Draw handles drawing.
func (bv *BattleView) Draw(screen *ebiten.Image) {
	if bv.BattleViewModel == nil {
		// Draw a message if no battle is set up
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(400, 300)
		drawing.DrawText(screen, "No battle in progress", 24, opt)
		return
	}

	// Draw battle title
	title := bv.BattleViewModel.Title()
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(20, 60)
	drawing.DrawText(screen, title, 32, opt)

	// Draw enemy information
	bv.drawEnemyInfo(screen)

	// Draw battle cards
	bv.drawBattleCards(screen)

	// Draw UI buttons
	bv.drawButtons(screen)

	// Draw total power and battle result
	bv.drawBattleStatus(screen)
}

// drawEnemyInfo draws enemy information
func (bv *BattleView) drawEnemyInfo(screen *ebiten.Image) {
	// Draw enemy image
	enemyImage := bv.BattleViewModel.EnemyImage()
	if enemyImage != nil {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Scale(2.0, 2.0)
		opt.GeoM.Translate(50, 120)
		screen.DrawImage(enemyImage, opt)
	}

	// Draw enemy power
	power := bv.BattleViewModel.EnemyPower()
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(200, 150)
	drawing.DrawText(screen, fmt.Sprintf("Power: %.1f", power), 24, opt)

	// Draw enemy type
	enemyType := bv.BattleViewModel.EnemyType()
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(200, 180)
	drawing.DrawText(screen, enemyType, 20, opt)
}

// drawBattleCards draws the placed battle cards
func (bv *BattleView) drawBattleCards(screen *ebiten.Image) {
	numCards := bv.BattleViewModel.NumCards()
	
	for i := 0; i < numCards; i++ {
		cardVM := bv.BattleViewModel.Card(i)
		if cardVM != nil {
			x := float64(100 + i*80)
			y := float64(400)
			
			// Draw card background
			alpha := float32(1.0)
			if cardVM == bv.HoveredCard {
				alpha = 0.8
			}
			DrawCardBackground(screen, x, y, alpha)
			
			// Draw card content (simplified)
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(x+10, y+10)
			drawing.DrawText(screen, cardVM.Name(), 12, opt)
			
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(x+10, y+90)
			drawing.DrawText(screen, fmt.Sprintf("%.1f", cardVM.Power()), 16, opt)
		}
	}
}

// drawButtons draws UI buttons
func (bv *BattleView) drawButtons(screen *ebiten.Image) {
	// Back button (960,40,80,80)
	drawing.DrawRect(screen, 960, 40, 80, 80, 0.3, 0.3, 0.3, 1.0)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(980, 65)
	drawing.DrawText(screen, "Back", 20, opt)

	// Conquer button (400,560,240,40)
	canWin := bv.BattleViewModel.CanBeat()
	if canWin {
		drawing.DrawRect(screen, 400, 560, 240, 40, 0.2, 0.6, 0.2, 1.0)
	} else {
		drawing.DrawRect(screen, 400, 560, 240, 40, 0.6, 0.2, 0.2, 1.0)
	}
	
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(480, 575)
	buttonText := "Conquer"
	if !canWin {
		buttonText = "Cannot Win"
	}
	drawing.DrawText(screen, buttonText, 20, opt)
}

// drawBattleStatus draws battle status information
func (bv *BattleView) drawBattleStatus(screen *ebiten.Image) {
	totalPower := bv.BattleViewModel.TotalPower()
	enemyPower := bv.BattleViewModel.EnemyPower()
	
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(400, 350)
	drawing.DrawText(screen, fmt.Sprintf("Your Power: %.1f", totalPower), 24, opt)
	
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(600, 350)
	drawing.DrawText(screen, fmt.Sprintf("Enemy: %.1f", enemyPower), 24, opt)
	
	canWin := bv.BattleViewModel.CanBeat()
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(500, 380)
	if canWin {
		drawing.DrawText(screen, "Victory!", 20, opt)
	} else {
		drawing.DrawText(screen, "Need more power", 20, opt)
	}
}
