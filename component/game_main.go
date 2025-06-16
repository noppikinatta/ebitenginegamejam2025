package component

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameMain struct {
	currentScreen string
	x, y          int
	width         int
	height        int
}

func NewGameMain(x, y, width, height int) *GameMain {
	return &GameMain{
		currentScreen: "Map", // Start with Map screen
		x:             x,
		y:             y,
		width:         width,
		height:        height,
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
	for row := 0; row < 7; row++ {
		for col := 0; col < 13; col++ {
			pointX := gm.x + 20 + col*35
			pointY := gm.y + 50 + row*35

			// Draw point based on type
			if row == 6 && col == 0 {
				// Home (bottom-left)
				ebitenutil.DebugPrintAt(screen, "H", pointX, pointY)
			} else if row == 0 && col == 12 {
				// Boss (top-right)
				ebitenutil.DebugPrintAt(screen, "B", pointX, pointY)
			} else {
				// Wild/NPC points
				ebitenutil.DebugPrintAt(screen, ".", pointX, pointY)
			}
		}
	}
}

func (gm *GameMain) drawDiplomacyScreen(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "DIPLOMACY SCREEN", gm.x+10, gm.y+10)
	ebitenutil.DebugPrintAt(screen, "NPC Nation Trade", gm.x+10, gm.y+30)
	ebitenutil.DebugPrintAt(screen, "Available Cards:", gm.x+10, gm.y+60)

	// Show sample cards for purchase
	cards := []string{
		"Iron Warrior - Cost: 50 Gold",
		"Wooden Spear - Cost: 25 Wood",
		"Magic Shield - Cost: 30 Mana",
	}

	for i, card := range cards {
		ebitenutil.DebugPrintAt(screen, card, gm.x+20, gm.y+80+i*20)
	}
}

func (gm *GameMain) drawBattleScreen(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "BATTLE SCREEN", gm.x+10, gm.y+10)
	ebitenutil.DebugPrintAt(screen, "Place your cards to fight", gm.x+10, gm.y+30)

	// Draw battle field layout
	ebitenutil.DebugPrintAt(screen, "Your Front Row:", gm.x+10, gm.y+60)
	ebitenutil.DebugPrintAt(screen, "[_] [_] [_] [_] [_]", gm.x+20, gm.y+80)

	ebitenutil.DebugPrintAt(screen, "Your Back Row:", gm.x+10, gm.y+100)
	ebitenutil.DebugPrintAt(screen, "[_] [_] [_] [_] [_]", gm.x+20, gm.y+120)

	ebitenutil.DebugPrintAt(screen, "Enemy Forces:", gm.x+10, gm.y+160)
	ebitenutil.DebugPrintAt(screen, "[E] [E] [E]", gm.x+20, gm.y+180)
}
