package component

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/noppikinatta/ebitenginegamejam2025/screen"
)

type GameMain struct {
	currentScreen   string
	x, y            int
	width           int
	height          int
	mapScreen       *screen.MapScreen
	diplomacyScreen *screen.DiplomacyScreen
	battleScreen    *screen.BattleScreen
}

func NewGameMain(x, y, width, height int) *GameMain {
	return &GameMain{
		currentScreen:   "Map", // Start with Map screen
		x:               x,
		y:               y,
		width:           width,
		height:          height,
		mapScreen:       screen.NewMapScreen(),
		diplomacyScreen: screen.NewDiplomacyScreen(),
		battleScreen:    screen.NewBattleScreen(),
	}
}

func (gm *GameMain) GetCurrentScreen() string {
	return gm.currentScreen
}

func (gm *GameMain) SwitchToScreen(screenName string) {
	validScreens := map[string]bool{
		"Map":       true,
		"Diplomacy": true,
		"Battle":    true,
	}

	if validScreens[screenName] {
		gm.currentScreen = screenName
	}
}

func (gm *GameMain) GetMapScreen() *screen.MapScreen {
	return gm.mapScreen
}

func (gm *GameMain) GetDiplomacyScreen() *screen.DiplomacyScreen {
	return gm.diplomacyScreen
}

func (gm *GameMain) GetBattleScreen() *screen.BattleScreen {
	return gm.battleScreen
}

func (gm *GameMain) Draw(screen *ebiten.Image) {
	// Draw background
	mainBg := ebiten.NewImage(gm.width, gm.height)
	mainBg.Fill(color.RGBA{20, 40, 60, 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(gm.x), float64(gm.y))
	screen.DrawImage(mainBg, op)

	// Draw current screen content
	switch gm.currentScreen {
	case "Map":
		gm.drawMapScreen(screen)
	case "Diplomacy":
		gm.drawDiplomacyScreen(screen)
	case "Battle":
		gm.drawBattleScreen(screen)
	}
}

func (gm *GameMain) drawMapScreen(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "MAP SCREEN", gm.x+10, gm.y+10)
	ebitenutil.DebugPrintAt(screen, "13x7 Grid World Map", gm.x+10, gm.y+30)

	// Draw a simple grid representation
	grid := gm.mapScreen.GetGrid()
	for row := 0; row < 7; row++ {
		for col := 0; col < 13; col++ {
			pointX := gm.x + 20 + col*35
			pointY := gm.y + 50 + row*35

			point := grid.GetPoint(col, row)
			if point != nil {
				// Draw point based on type
				symbol := "."
				switch point.GetType() {
				case "Home":
					symbol = "H"
				case "Boss":
					symbol = "B"
				case "NPC":
					symbol = "N"
				case "Wild":
					symbol = "W"
				}
				ebitenutil.DebugPrintAt(screen, symbol, pointX, pointY)
			}
		}
	}
}

func (gm *GameMain) drawDiplomacyScreen(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "DIPLOMACY SCREEN", gm.x+10, gm.y+10)

	npc := gm.diplomacyScreen.GetCurrentNPCInfo()
	if npc != nil {
		ebitenutil.DebugPrintAt(screen, "Trading with: "+npc.Name, gm.x+10, gm.y+30)
		ebitenutil.DebugPrintAt(screen, "Specialty: "+npc.Specialty, gm.x+10, gm.y+45)
	}

	ebitenutil.DebugPrintAt(screen, "Available Cards:", gm.x+10, gm.y+70)

	// Show available cards
	cards := gm.diplomacyScreen.GetAvailableCards()
	for i, card := range cards {
		if i < 5 { // Limit display
			text := card.Name + " (" + card.Type + ")"
			ebitenutil.DebugPrintAt(screen, text, gm.x+20, gm.y+90+i*20)
		}
	}
}

func (gm *GameMain) drawBattleScreen(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "BATTLE SCREEN", gm.x+10, gm.y+10)
	ebitenutil.DebugPrintAt(screen, "Place your cards to fight", gm.x+10, gm.y+30)

	battlefield := gm.battleScreen.GetBattlefield()

	// Draw battle field layout
	ebitenutil.DebugPrintAt(screen, "Your Front Row:", gm.x+10, gm.y+60)
	frontRow := battlefield.GetFrontRow()
	for i, card := range *frontRow {
		symbol := "_"
		if card != "" {
			symbol = "C"
		}
		ebitenutil.DebugPrintAt(screen, "["+symbol+"]", gm.x+20+i*30, gm.y+80)
	}

	ebitenutil.DebugPrintAt(screen, "Your Back Row:", gm.x+10, gm.y+100)
	backRow := battlefield.GetBackRow()
	for i, card := range *backRow {
		symbol := "_"
		if card != "" {
			symbol = "C"
		}
		ebitenutil.DebugPrintAt(screen, "["+symbol+"]", gm.x+20+i*30, gm.y+120)
	}

	ebitenutil.DebugPrintAt(screen, "Enemies:", gm.x+10, gm.y+160)
	enemies := gm.battleScreen.GetEnemies()
	for i, enemy := range enemies {
		if i < 3 { // Limit display
			ebitenutil.DebugPrintAt(screen, enemy.Name, gm.x+20+i*80, gm.y+180)
		}
	}
}
