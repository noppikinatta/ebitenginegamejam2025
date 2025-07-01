package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// BattleView is a Widget for displaying battles.
// Position: Drawn within MainView.
type BattleView struct {
	BattlePoint core.BattlePoint  // The point to be battled.
	PointName   string            // The name of the battle point.
	Battlefield *core.Battlefield // Battlefield information.
	GameState   *core.GameState   // Game state.
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
}

// SetPointName sets the name of the battle point.
func (bv *BattleView) SetPointName(pointName string) {
	bv.PointName = pointName
}

// SetGameState sets the game state.
func (bv *BattleView) SetGameState(gameState *core.GameState) {
	bv.GameState = gameState
}

// GetTotalPower calculates the total Power value of the placed BattleCards.
func (bv *BattleView) GetTotalPower() float64 {
	return bv.Battlefield.CalculateTotalPower()
}

// CanDefeatEnemy determines if the enemy can be defeated.
func (bv *BattleView) CanDefeatEnemy() bool {
	if bv.Battlefield != nil {
		return bv.Battlefield.CanBeat()
	}
	// Keep existing logic for backward compatibility.
	if bv.BattlePoint == nil {
		return false
	}
	return bv.GetTotalPower() >= bv.BattlePoint.GetEnemy().Power
}

// HandleInput handles input.
func (bv *BattleView) HandleInput(input *Input) error {
	cursorX, cursorY := input.Mouse.CursorPosition()
	cardIndex := bv.cardIndex(cursorX, cursorY)
	bv.MouseX = cursorX
	bv.MouseY = cursorY

	if cardIndex != -1 {
		bv.HoveredCard = bv.Battlefield.BattleCards[cardIndex]
	} else {
		bv.HoveredCard = nil
	}

	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		if cardIndex != -1 {
			bv.handleBattleCardClick(cursorX, cursorY)
		}

		// Click detection for the back button (960,40,80,80).
		if cursorX >= 960 && cursorX < 1040 && cursorY >= 40 && cursorY < 120 {
			// Return all placed BattleCards to the CardDeck.
			bv.GameState.CardDeck.Add(&core.Cards{BattleCards: bv.Battlefield.BattleCards})
			bv.Battlefield.BattleCards = make([]*core.BattleCard, 0)
			if bv.OnBackClicked != nil {
				bv.OnBackClicked()
				return nil
			}
		}

		// Click detection for the conquer button (400,560,240,40).
		if cursorX >= 400 && cursorX < 640 && cursorY >= 560 && cursorY < 600 {
			if bv.CanDefeatEnemy() {
				bv.Conquer()
			}
			if bv.OnBackClicked != nil {
				bv.GameState.CardDeck.Add(&core.Cards{BattleCards: bv.Battlefield.BattleCards})
				bv.Battlefield.BattleCards = make([]*core.BattleCard, 0)
				bv.OnBackClicked()
			}
			return nil
		}

		// Click detection for the enemy image (victory process).
		if bv.CanDefeatEnemy() && cursorX >= 360 && cursorX < 680 && cursorY >= 120 && cursorY < 440 {
			if bv.Conquer() {
				// On successful conquest, return to MapGridView.
				if bv.OnBackClicked != nil {
					bv.OnBackClicked()
				}
			}
			return nil
		}

		// Click detection for BattleCard (return to CardDeck).
		bv.handleBattleCardClick(cursorX, cursorY)
	}
	return nil
}

func (bv *BattleView) cardIndex(cursorX, cursorY int) int {
	if cursorX < 0 || cursorX >= 960 || cursorY < 440 || cursorY >= 560 {
		return -1
	}
	cardIndex := cursorX / 80
	if cardIndex < 0 || cardIndex >= len(bv.Battlefield.BattleCards) {
		return -1
	}
	return cardIndex
}

// handleBattleCardClick handles BattleCard clicks.
func (bv *BattleView) handleBattleCardClick(cursorX, cursorY int) {
	// Each card (80x120) in the BattleCard area (0,440,960,120).
	if cursorY >= 440 && cursorY < 560 {
		cardIndex := cursorX / 80

		var targetCard *core.BattleCard
		if bv.Battlefield != nil && cardIndex >= 0 && cardIndex < len(bv.Battlefield.BattleCards) {
			targetCard = bv.Battlefield.BattleCards[cardIndex]
		}

		if targetCard != nil {
			// Return the card to the CardDeck.
			bv.RemoveCard(targetCard)
		}
	}
}

// Draw handles drawing.
func (bv *BattleView) Draw(screen *ebiten.Image) {
	// Draw header (0,40,1040,80).
	bv.drawHeader(screen)

	// Draw back button (960,40,80,80).
	bv.drawBackButton(screen)

	// Draw enemy image (360,120,320,320).
	bv.drawEnemy(screen)

	// Draw BattleCard area (0,440,960,120).
	bv.drawBattleCards(screen)

	// Draw power display (960,440,80,120).
	bv.drawPowerDisplay(screen)

	// Draw conquer button (400,560,240,80).
	bv.drawConquerButton(screen)

	bv.drawHoveredCardTooltip(screen)
}

// drawHeader draws the header.
func (bv *BattleView) drawHeader(screen *ebiten.Image) {
	// Header background.
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 40, SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		{DstX: 1040, DstY: 40, SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		{DstX: 1040, DstY: 120, SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
		{DstX: 0, DstY: 120, SrcX: 0, SrcY: 0, ColorR: 0.4, ColorG: 0.2, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Title text.
	pointName := ""
	if p, ok := bv.BattlePoint.(*core.WildernessPoint); ok {
		pointName = p.TerrainType
	}
	if _, ok := bv.BattlePoint.(*core.BossPoint); ok {
		pointName = "point-boss"
	}
	title := lang.ExecuteTemplate("battle-title", map[string]any{"location": lang.Text(pointName)})

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(20, 60)
	drawing.DrawText(screen, title, 32, opt)
}

// drawBackButton draws the back button.
func (bv *BattleView) drawBackButton(screen *ebiten.Image) {
	DrawButton(screen, 960, 40, 80, 80, "ui-close")
}

// drawEnemy draws the enemy image.
func (bv *BattleView) drawEnemy(screen *ebiten.Image) {
	// Enemy image background (360,120,320,320).
	var color [4]float32
	color = [4]float32{0.8, 0.8, 0.8, 1}

	vertices := []ebiten.Vertex{
		{DstX: 360, DstY: 120, SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
		{DstX: 680, DstY: 120, SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
		{DstX: 680, DstY: 440, SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
		{DstX: 360, DstY: 440, SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Draw enemy information.
	if bv.BattlePoint != nil {
		enemy := bv.BattlePoint.GetEnemy()

		enemyImage := drawing.Image(string(enemy.EnemyID))
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Scale(4, 4)
		opt.GeoM.Translate(360, 120)
		screen.DrawImage(enemyImage, opt)

		// Enemy type.
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(370, 140)
		enemyType := lang.ExecuteTemplate("battle-enemy-type", map[string]any{"type": lang.Text(string(enemy.EnemyType))})
		drawing.DrawText(screen, enemyType, 24, opt)

		// Enemy's Power.
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Scale(2.0, 2.0)
		opt.GeoM.Translate(370, 180)
		powerIcon := drawing.Image("ui-power")
		screen.DrawImage(powerIcon, opt)
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(402, 180)
		powerText := fmt.Sprintf("%s: %.1f", lang.Text("battle-power"), enemy.Power)
		drawing.DrawText(screen, powerText, 24, opt)

		// Enemy's quote.
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(32, 320)
		enemyTalk := lang.ExecuteTemplate("battle-enemy-talk", map[string]any{"name": lang.Text(string(enemy.EnemyID)), "text": lang.Text(string(enemy.Question))})
		drawing.DrawText(screen, enemyTalk, 24, opt)

		// Enemy skills
		for i, skill := range enemy.Skills {
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(700, 120+float64(i*120))
			skillName := lang.Text(string(skill.ID()))
			drawing.DrawText(screen, skillName, 24, opt)
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(700, 152+float64(i*120))
			skillDescription := lang.Text(string(skill.ID()) + "-desc")
			drawing.DrawText(screen, skillDescription, 18, opt)
		}
	}
}

// drawBattleCards draws the BattleCard area
func (bv *BattleView) drawBattleCards(screen *ebiten.Image) {
	// Background of the BattleCard area (0,440,960,120)
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 440, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.2, ColorB: 0.3, ColorA: 1},
		{DstX: 960, DstY: 440, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.2, ColorB: 0.3, ColorA: 1},
		{DstX: 960, DstY: 560, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.2, ColorB: 0.3, ColorA: 1},
		{DstX: 0, DstY: 560, SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.2, ColorB: 0.3, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Draw each card.
	// Draw deployed BattleCards (80x120 × 12 cards)
	for i, card := range bv.Battlefield.BattleCards {
		if i >= 12 { // Maximum of 12 cards
			break
		}

		cardX := float64(i * 80)
		cardY := 440.0

		// Draw card
		DrawBattleCard(screen, cardX, cardY, card)
	}

	for i := len(bv.Battlefield.BattleCards); i < bv.Battlefield.CardSlot; i++ {
		cardX := float64(i * 80)
		cardY := 440.0
		DrawCardBackground(screen, cardX, cardY, 0.5)
	}
}

// drawPowerDisplay draws the Power display.
func (bv *BattleView) drawPowerDisplay(screen *ebiten.Image) {
	// Background of the Power display (960,440,80,120)
	vertices := []ebiten.Vertex{
		{DstX: 960, DstY: 440, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.4, ColorA: 1},
		{DstX: 1040, DstY: 440, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.4, ColorA: 1},
		{DstX: 1040, DstY: 560, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.4, ColorA: 1},
		{DstX: 960, DstY: 560, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.4, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Draw total Power value
	totalPower := bv.GetTotalPower()

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(2.0, 2.0)
	opt.GeoM.Translate(970, 450)
	powerIcon := drawing.Image("ui-power")
	screen.DrawImage(powerIcon, opt)

	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(970, 480)
	powerText := fmt.Sprintf("%.1f", totalPower)
	drawing.DrawText(screen, powerText, 28, opt)

	// Display required Power (enemy's Power)
	if bv.BattlePoint != nil {
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(970, 510)
		requiredText := fmt.Sprintf("/%.1f", bv.BattlePoint.GetEnemy().Power)
		drawing.DrawText(screen, requiredText, 20, opt)
	}
}

// drawConquerButton draws the conquer button.
func (bv *BattleView) drawConquerButton(screen *ebiten.Image) {
	canConquer := bv.CanDefeatEnemy()

	// Determine button color
	var colorR, colorG, colorB float32 = 0.5, 0.5, 0.5 // Gray when invalid
	if canConquer {
		colorR, colorG, colorB = 0.2, 0.8, 0.2 // Green when valid
	}

	// Button background
	vertices := []ebiten.Vertex{
		{DstX: 400, DstY: 560, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: 640, DstY: 560, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: 640, DstY: 600, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: 400, DstY: 600, SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Button text
	if canConquer {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(420, 560)
		drawing.DrawText(screen, lang.Text("ui-conquer"), 28, opt)
	} else {
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(420, 560)
		drawing.DrawText(screen, lang.Text("ui-need-power"), 24, opt)
	}
}

func (bv *BattleView) drawHoveredCardTooltip(screen *ebiten.Image) {
	if bv.HoveredCard == nil {
		return
	}

	DrawCardDescriptionTooltip(screen, bv.HoveredCard, bv.MouseX, bv.MouseY)
}

// createBattlefield creates a battlefield
func (bv *BattleView) createBattlefield(point core.BattlePoint) *core.Battlefield {
	if bv.GameState == nil {
		return core.NewBattlefield(point.GetEnemy(), 0.0)
	}

	x, y, ok := bv.GameState.MapGrid.XYOfPoint(point)
	if !ok {
		panic("BattleView.createBattlefield: battle point does not exist in map grid")
	}
	enemy := point.GetEnemy()
	battlefield := core.NewBattlefield(enemy, 0.0)

	// Investigate Points in four directions based on x,y
	mapGrid := bv.GameState.MapGrid
	directions := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} // Up, down, left, right

	for _, dir := range directions {
		checkX := x + dir[0]
		checkY := y + dir[1]

		// Check if within map range
		if checkX >= 0 && checkX < mapGrid.Size.X && checkY >= 0 && checkY < mapGrid.Size.Y {
			p := mapGrid.GetPoint(checkX, checkY)

			// Get Territory from controlled WildernessPoint
			if wildernessPoint, ok := p.(*core.WildernessPoint); ok {
				if wildernessPoint.Controlled && wildernessPoint.Territory != nil {
					territory := wildernessPoint.Territory

					// Apply StructureCard.BattlefieldModifier from Territory.Cards
					for _, card := range territory.Cards {
						if card.BattlefieldModifier != nil {
							card.BattlefieldModifier.Modify(battlefield)
						}
					}
				}
			}
		}
	}

	return battlefield
}

// CanPlaceCard determines if a card can be placed
func (bv *BattleView) CanPlaceCard() bool {
	if bv.Battlefield == nil {
		return false
	}
	return len(bv.Battlefield.BattleCards) < bv.Battlefield.CardSlot
}

// PlaceCard places a card
func (bv *BattleView) PlaceCard(card *core.BattleCard) bool {
	if bv.Battlefield == nil {
		return false
	}

	return bv.Battlefield.AddBattleCard(card)
}

// RemoveCard removes a card
func (bv *BattleView) RemoveCard(card *core.BattleCard) bool {
	if bv.Battlefield == nil {
		return false
	}

	// Find card index
	cardIndex := -1
	for i, battleCard := range bv.Battlefield.BattleCards {
		if battleCard == card {
			cardIndex = i
			break
		}
	}

	if cardIndex == -1 {
		return false
	}

	// Remove from Battlefield
	removedCard, success := bv.Battlefield.RemoveBattleCard(cardIndex)
	if success && removedCard != nil {
		// Add to GameState.CardDeck
		if bv.GameState != nil {
			cards := &core.Cards{BattleCards: []*core.BattleCard{removedCard}}
			bv.GameState.CardDeck.Add(cards)
		}
	}

	return success
}

// Conquer executes the conquest process
func (bv *BattleView) Conquer() bool {
	if bv.Battlefield == nil || bv.GameState == nil || bv.BattlePoint == nil {
		return false
	}
	// Check if victory is possible
	if !bv.Battlefield.CanBeat() {
		return false
	}

	// Battle victory process
	bv.Battlefield.Beat()

	// Return all placed BattleCards to the CardDeck
	bv.GameState.CardDeck.Add(&core.Cards{BattleCards: bv.Battlefield.BattleCards})

	// Change the target BattlePoint's Controlled to true
	bv.BattlePoint.SetControlled(true)
	bv.GameState.MapGrid.UpdateAccesibles()

	bv.GameState.MyNation.AppendLevel(0.5)
	bv.GameState.NextTurn()

	return true
}
