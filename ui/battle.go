package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/flow"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// BattleView is a Widget for displaying battles.
// Position: Drawn within MainView.
type BattleView struct {
	BattleFlow      *flow.BattleFlow
	BattleViewModel *viewmodel.BattleViewModel

	HoveredCardIndex int
}

// NewBattleView creates a BattleView.
func NewBattleView(battleFlow *flow.BattleFlow) *BattleView {
	return &BattleView{
		BattleFlow: battleFlow,
	}
}

func (bv *BattleView) Select(x, y int) {
	vm, ok := bv.BattleFlow.Select(x, y)
	if !ok {
		return
	}
	bv.BattleViewModel = vm
}

// HandleInput handles input.
func (bv *BattleView) HandleInput(input *Input) (bool, error) {
	cursorX, cursorY := input.Mouse.CursorPosition()
	cardIndex := bv.cardIndex(cursorX, cursorY)
	bv.HoveredCardIndex = cardIndex

	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		if cardIndex != -1 {
			bv.handleBattleCardClick(cursorX, cursorY)
		}

		// Click detection for the back button (960,40,80,80).
		if cursorX >= 960 && cursorX < 1040 && cursorY >= 40 && cursorY < 120 {
			bv.BattleFlow.Rollback()
			return true, nil
		}

		// Click detection for the conquer button (400,560,240,40).
		if cursorX >= 400 && cursorX < 640 && cursorY >= 560 && cursorY < 600 {
			ok := bv.BattleFlow.Conquer()
			if !ok {
				bv.BattleFlow.Rollback()
			}
			return true, nil
		}
	}

	return false, nil
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
		x := float64(100 + i*80)
		y := float64(400)

		cardVM, ok := bv.BattleViewModel.Card(i)
		if !ok {
			DrawCardBackground(screen, x, y, 0.5)
			continue
		}

		DrawCard(screen, x, y, cardVM, i == bv.HoveredCardIndex)

		maxCards := bv.BattleViewModel.CardSlot()
		for i := numCards; i < maxCards; i++ {
			x := float64(100 + i*80)
			y := float64(400)
			DrawCardBackground(screen, x, y, 0.3)
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
