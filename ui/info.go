package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// InfoViewMode is the display mode of InfoView.
type InfoViewMode int

const (
	InfoModeHistory InfoViewMode = iota
	InfoModeCardInfo
	InfoModeNationPoint
	InfoModeWildernessPoint
	InfoModeEnemySkill
)

// InfoView is a widget for displaying information.
// Position: (520,20,120,280).
// Changes the content of the information displayed according to the situation.
type InfoView struct {
	CurrentMode InfoViewMode

	// Display data.
	SelectedCard  interface{}     // Selected card (BattleCard or StructureCard).
	SelectedPoint core.Point      // Selected Point.
	SelectedEnemy *core.Enemy     // Selected Enemy.
	History       []string        // Event history.
	GameState     *core.GameState // Game state.

	// Mouse cursor position (set externally).
	MouseX, MouseY int
}

// NewInfoView creates an InfoView.
func NewInfoView() *InfoView {
	return &InfoView{
		CurrentMode: InfoModeHistory, // The default is HistoryView.
		History:     make([]string, 0),
	}
}

// SetGameState sets the game state.
func (iv *InfoView) SetGameState(gameState *core.GameState) {
	iv.GameState = gameState
}

// AddHistoryEvent adds an event to the history.
func (iv *InfoView) AddHistoryEvent(event string) {
	iv.History = append(iv.History, event)
	// Holds up to 14 lines.
	if len(iv.History) > 14 {
		iv.History = iv.History[1:]
	}
}

// SetSelectedCard sets the selected card.
func (iv *InfoView) SetSelectedCard(card interface{}) {
	iv.SelectedCard = card
	if card != nil {
		iv.CurrentMode = InfoModeCardInfo
	} else {
		iv.CurrentMode = InfoModeHistory
	}
}

// SetSelectedPoint sets the selected Point.
func (iv *InfoView) SetSelectedPoint(point core.Point) {
	iv.SelectedPoint = point
	if point != nil {
		switch point.(type) {
		case *core.MyNationPoint, *core.OtherNationPoint:
			iv.CurrentMode = InfoModeNationPoint
		case *core.WildernessPoint, *core.BossPoint:
			iv.CurrentMode = InfoModeWildernessPoint
		}
	} else {
		iv.CurrentMode = InfoModeHistory
	}
}

// SetEnemySkillMode sets the mode to EnemySkillView.
func (iv *InfoView) SetEnemySkillMode(enemy *core.Enemy) {
	iv.SelectedEnemy = enemy
	iv.CurrentMode = InfoModeEnemySkill
}

// HandleInput handles input.
func (iv *InfoView) HandleInput(input *Input) error {
	// InfoView basically does not accept input (display only).
	return nil
}

// Draw handles drawing.
func (iv *InfoView) Draw(screen *ebiten.Image) {
	// Draw background.
	iv.drawBackground(screen)

	// Draw the content according to the current mode.
	switch iv.CurrentMode {
	case InfoModeHistory:
		iv.drawHistoryView(screen)
	case InfoModeCardInfo:
		iv.drawCardInfoView(screen)
	case InfoModeNationPoint:
		iv.drawNationPointView(screen)
	case InfoModeWildernessPoint:
		iv.drawWildernessPointView(screen)
	case InfoModeEnemySkill:
		iv.drawEnemySkillView(screen)
	}
}

// drawBackground draws the background.
func (iv *InfoView) drawBackground(screen *ebiten.Image) {
	// InfoView background (520,20,120,280).
	vertices := []ebiten.Vertex{
		{DstX: 520, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.15, ColorG: 0.15, ColorB: 0.2, ColorA: 1},
		{DstX: 640, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.15, ColorG: 0.15, ColorB: 0.2, ColorA: 1},
		{DstX: 640, DstY: 300, SrcX: 0, SrcY: 0, ColorR: 0.15, ColorG: 0.15, ColorB: 0.2, ColorA: 1},
		{DstX: 520, DstY: 300, SrcX: 0, SrcY: 0, ColorR: 0.15, ColorG: 0.15, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
}

// drawHistoryView draws the HistoryView.
func (iv *InfoView) drawHistoryView(screen *ebiten.Image) {
	// Title.
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(525, 25)
	drawing.DrawText(screen, lang.Text("ui-history"), 12, opt)

	// Event history display (120x20 x 14 lines).
	for i, event := range iv.History {
		if i >= 14 { // Maximum 14 lines.
			break
		}

		y := 45.0 + float64(i)*18
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)

		// Omit if the text is long.
		displayText := event
		if len(displayText) > 15 {
			displayText = displayText[:12] + "..."
		}

		drawing.DrawText(screen, displayText, 9, opt)
	}

	// Display when there is no history.
	if len(iv.History) == 0 {
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, 50)
		drawing.DrawText(screen, lang.Text("ui-no-events"), 9, opt)
	}
}

// drawCardInfoView draws the CardInfoView.
func (iv *InfoView) drawCardInfoView(screen *ebiten.Image) {
	if iv.SelectedCard == nil {
		iv.drawHistoryView(screen)
		return
	}

	startY := 25.0

	switch card := iv.SelectedCard.(type) {
	case *core.BattleCard:
		iv.drawBattleCardInfo(screen, card, startY)
	case *core.StructureCard:
		iv.drawStructureCardInfo(screen, card, startY)
	default:
		iv.drawHistoryView(screen)
	}
}

// drawBattleCardInfo draws the detailed information of a BattleCard.
func (iv *InfoView) drawBattleCardInfo(screen *ebiten.Image, card *core.BattleCard, y float64) {
	// Card name (20).
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(525, y)
	drawing.DrawText(screen, fmt.Sprintf("Card: %s", card.CardID), 10, opt)
	y += 20

	// Illustration (60) - dummy rectangle.
	vertices := []ebiten.Vertex{
		{DstX: 525, DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: 585, DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: 585, DstY: float32(y + 60), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: 525, DstY: float32(y + 60), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
	y += 60

	// Card type (20).
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(525, y)
	drawing.DrawText(screen, "Type: Battle", 9, opt)
	y += 20

	// Card class (20).
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(525, y)
	drawing.DrawText(screen, fmt.Sprintf("Class: %s", card.Type), 9, opt)
	y += 20

	// Power (20).
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(525, y)
	drawing.DrawText(screen, fmt.Sprintf("Power: %.1f", card.BasePower), 9, opt)
	y += 20

	// Skill name (20).
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(525, y)
	if card.Skill != nil {
		drawing.DrawText(screen, "Skill: Active", 9, opt) // Dummy text.
	} else {
		drawing.DrawText(screen, "Skill: None", 9, opt)
	}
	y += 20

	// Skill description (40).
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(525, y)
	if card.Skill != nil {
		drawing.DrawText(screen, "Special battle", 8, opt)
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y+12)
		drawing.DrawText(screen, "effect active", 8, opt)
	} else {
		drawing.DrawText(screen, "No special effect", 8, opt)
	}
}

// drawStructureCardInfo draws the detailed information of a StructureCard.
func (iv *InfoView) drawStructureCardInfo(screen *ebiten.Image, card *core.StructureCard, y float64) {
	// Card name (20).
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(525, y)
	drawing.DrawText(screen, fmt.Sprintf("Card: %s", card.CardID), 10, opt)
	y += 20

	// Illustration (60).
	vertices := []ebiten.Vertex{
		{DstX: 525, DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.6, ColorB: 0.4, ColorA: 1},
		{DstX: 585, DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.6, ColorB: 0.4, ColorA: 1},
		{DstX: 585, DstY: float32(y + 60), SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.6, ColorB: 0.4, ColorA: 1},
		{DstX: 525, DstY: float32(y + 60), SrcX: 0, SrcY: 0, ColorR: 0.2, ColorG: 0.6, ColorB: 0.4, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
	y += 60

	// Card type.
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(525, y)
	drawing.DrawText(screen, "Type: Structure", 9, opt)
	y += 20

	// YieldModifier効果 (20×9)
	if card.YieldModifier != nil {
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "Yield Effect:", 9, opt)
		y += 15

	opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "Boosts resource", 8, opt)
		y += 12

		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "production", 8, opt)
	} else {
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "No yield effect", 9, opt)
	}
}

// drawNationPointView draws the NationPointView.
func (iv *InfoView) drawNationPointView(screen *ebiten.Image) {
	if iv.SelectedPoint == nil {
		iv.drawHistoryView(screen)
		return
	}

	y := 25.0

	// Nation name.
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(525, y)

	switch point := iv.SelectedPoint.(type) {
	case *core.MyNationPoint:
		drawing.DrawText(screen, "My Nation", 12, opt)
		y += 20

		// Card Packs (20)
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "Card Packs:", 10, opt)
		y += 20

		// CardPack一覧 (20×12)
		marketItems := point.MyNation.VisibleMarketItems()
		for i, item := range marketItems {
			if i >= 12 {
				break
			}

			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(525, y)
			packName := string(item.CardPack.CardPackID)
			if len(packName) > 12 {
				packName = packName[:9] + "..."
			}
			drawing.DrawText(screen, packName, 9, opt)
			y += 18
		}

	case *core.OtherNationPoint:
		drawing.DrawText(screen, fmt.Sprintf("Nation %s", point.OtherNation.NationID), 10, opt)
		y += 20

		// Card Packs (20)
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "Card Packs:", 10, opt)
		y += 20

		// CardPack一覧 (20×12)
		marketItems := point.OtherNation.VisibleMarketItems()
		for i, item := range marketItems {
			if i >= 12 {
				break
			}

			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(525, y)
			packName := string(item.CardPack.CardPackID)
			if len(packName) > 12 {
				packName = packName[:9] + "..."
			}
			drawing.DrawText(screen, packName, 9, opt)
			y += 18
		}
	}
}

// drawWildernessPointView draws the WildernessPointView.
func (iv *InfoView) drawWildernessPointView(screen *ebiten.Image) {
	if iv.SelectedPoint == nil {
		iv.drawHistoryView(screen)
		return
	}

	y := 25.0

	switch point := iv.SelectedPoint.(type) {
	case *core.WildernessPoint:
		// Point名 (20)
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "Wilderness", 12, opt)
		y += 20

		// Enemy (20)
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "Enemy:", 10, opt)
		y += 20

		// Enemy情報 (40)
		if point.Enemy != nil {
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(525, y)
			enemyName := string(point.Enemy.EnemyID)
			if len(enemyName) > 12 {
				enemyName = enemyName[:9] + "..."
			}
			if point.Controlled {
				enemyName += " (X)" // 制圧済み
			}
			drawing.DrawText(screen, enemyName, 9, opt)
			y += 20

			// EnemyのPower (20)
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(525, y)
			drawing.DrawText(screen, fmt.Sprintf("Power: %.1f", point.Enemy.Power), 9, opt)
			y += 20
		}

		// Yields (20)
		if point.Territory != nil {
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(525, y)
			drawing.DrawText(screen, "Yields:", 10, opt)
			y += 20

			// Resource種類ごとのYield (20×3 - 1行に2種類表示)
			yield := point.Territory.BaseYield
			resources := []struct {
				name  string
				value int
			}{
				{"Money", yield.Money},
				{"Food", yield.Food},
				{"Wood", yield.Wood},
				{"Iron", yield.Iron},
				{"Mana", yield.Mana},
			}

			for i := 0; i < len(resources); i += 2 {
				opt = &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(525, y)

				text := fmt.Sprintf("%s:%d", resources[i].name[:3], resources[i].value)
				if i+1 < len(resources) {
					text += fmt.Sprintf(" %s:%d", resources[i+1].name[:3], resources[i+1].value)
				}
				drawing.DrawText(screen, text, 8, opt)
				y += 16
			}

			// Structure Cards (20)
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(525, y)
			drawing.DrawText(screen, "Structures:", 9, opt)
			y += 18

			// 配置されているStructureCard (20×4)
			for i, card := range point.Territory.Cards {
				if i >= 4 {
					break
				}

				opt = &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(525, y)
				cardName := string(card.CardID)
				if len(cardName) > 12 {
					cardName = cardName[:9] + "..."
				}
				drawing.DrawText(screen, cardName, 8, opt)
				y += 16
			}
		}

	case *core.BossPoint:
		// Point名 (20)
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "Boss Point", 12, opt)
		y += 20

		// Enemy (20)
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "Boss:", 10, opt)
		y += 20

		// Boss情報 (40)
		if point.Boss != nil {
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(525, y)
			bossName := string(point.Boss.EnemyID)
			if len(bossName) > 12 {
				bossName = bossName[:9] + "..."
			}
			if point.Defeated {
				bossName += " (X)" // 撃破済み
			}
			drawing.DrawText(screen, bossName, 9, opt)
			y += 20

			// BossのPower (20)
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(525, y)
			drawing.DrawText(screen, fmt.Sprintf("Power: %.1f", point.Boss.Power), 9, opt)
		}
	}
}

// drawEnemySkillView draws the EnemySkillView.
func (iv *InfoView) drawEnemySkillView(screen *ebiten.Image) {
	if iv.SelectedEnemy == nil {
		iv.drawHistoryView(screen)
		return
	}

	y := 25.0

	// Enemy name.
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(525, y)
	drawing.DrawText(screen, "Enemy Skills", 12, opt)
	y += 20

	// Enemy Skills (60×4)
	if len(iv.SelectedEnemy.Skills) == 0 {
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "No special skills", 9, opt)
		return
	}

	for i, _ := range iv.SelectedEnemy.Skills {
		if i >= 4 { // 最大4個のスキル
			break
		}

		// スキル名 (20)
			opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, fmt.Sprintf("Skill %d", i+1), 10, opt) // ダミーテキスト
		y += 20

		// スキル説明 (40)
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "Special enemy", 8, opt)
		y += 12

		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "ability effect", 8, opt)
		y += 12

		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(525, y)
		drawing.DrawText(screen, "in battle", 8, opt)
		y += 16
	}
}
