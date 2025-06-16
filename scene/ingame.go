package scene

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/component"
	"github.com/noppikinatta/ebitenginegamejam2025/system"
)

type UILayout struct {
	components map[string][4]int // [x, y, width, height]
}

func NewUILayout() *UILayout {
	return &UILayout{
		components: map[string][4]int{
			"GameMain":     {0, 20, 520, 280},
			"ResourceView": {0, 0, 520, 20},
			"Calendar":     {520, 0, 120, 40},
			"History":      {520, 80, 120, 320},
			"CardDeck":     {0, 300, 520, 60},
		},
	}
}

func (ul *UILayout) GetComponentBounds(componentName string) [4]int {
	return ul.components[componentName]
}

type InGame struct {
	layout       *UILayout
	resourceView *component.ResourceView
	calendar     *component.Calendar
	history      *component.History
	cardDeck     *component.CardDeck
	gameMain     *component.GameMain

	// System managers
	resourceManager  *system.ResourceManager
	turnManager      *system.TurnManager
	cardManager      *system.CardManager
	territoryManager *system.TerritoryManager
	allianceManager  *system.AllianceManager
	combatManager    *system.CombatManager
}

func NewInGame() *InGame {
	layout := NewUILayout()

	// Create system managers
	resourceManager := system.NewResourceManager()
	turnManager := system.NewTurnManager()
	cardManager := system.NewCardManager()
	territoryManager := system.NewTerritoryManager(resourceManager)
	allianceManager := system.NewAllianceManager()
	combatManager := system.NewCombatManager()

	// Create components with their layout positions
	resourceViewBounds := layout.GetComponentBounds("ResourceView")
	calendarBounds := layout.GetComponentBounds("Calendar")
	historyBounds := layout.GetComponentBounds("History")
	cardDeckBounds := layout.GetComponentBounds("CardDeck")
	gameMainBounds := layout.GetComponentBounds("GameMain")

	return &InGame{
		layout:       layout,
		resourceView: component.NewResourceView(resourceManager, resourceViewBounds[0], resourceViewBounds[1], resourceViewBounds[2], resourceViewBounds[3]),
		calendar:     component.NewCalendar(calendarBounds[0], calendarBounds[1], calendarBounds[2], calendarBounds[3]),
		history:      component.NewHistory(historyBounds[0], historyBounds[1], historyBounds[2], historyBounds[3]),
		cardDeck:     component.NewCardDeck(cardManager, cardDeckBounds[0], cardDeckBounds[1], cardDeckBounds[2], cardDeckBounds[3]),
		gameMain:     component.NewGameMain(combatManager, gameMainBounds[0], gameMainBounds[1], gameMainBounds[2], gameMainBounds[3]),

		// System managers
		resourceManager:  resourceManager,
		turnManager:      turnManager,
		cardManager:      cardManager,
		territoryManager: territoryManager,
		allianceManager:  allianceManager,
		combatManager:    combatManager,
	}
}

func (g *InGame) Update() error {
	// 1. ターン進行
	g.turnManager.AdvanceTurn()
	// Calendar は TurnManager の進捗に合わせて AdvanceMonth せず、表示だけ同期する (実装簡易)

	// 2. 資源生成 (制圧地から取得)
	g.territoryManager.GenerateResources()

	// 2.5 同盟ボーナスの適用
	bonuses := g.allianceManager.GetAllianceBonuses()
	// 資源追加
	for rtype, amt := range bonuses.ResourceBonus {
		g.resourceManager.AddResource(rtype, amt)
	}
	// 戦闘ボーナス
	g.combatManager.SetPlayerAttackBonus(bonuses.MilitaryBonus)

	// 3. 歴史ログ追加 (簡易)
	entry := fmt.Sprintf("Kingdom Year %d, Month %d: Turn %d", g.turnManager.GetCurrentYear(), g.turnManager.GetCurrentMonth(), g.turnManager.GetTurnCount())
	g.history.AddEntry(entry)

	return nil
}

func (g *InGame) Draw(screen *ebiten.Image) {
	// Background color
	screen.Fill(color.RGBA{40, 40, 60, 255})

	// Draw all UI components
	g.resourceView.Draw(screen)
	g.calendar.Draw(screen)
	g.history.Draw(screen)
	g.cardDeck.Draw(screen)
	g.gameMain.Draw(screen)
}

func (g *InGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}

func (g *InGame) GetUILayout() *UILayout {
	return g.layout
}

func (g *InGame) GetResourceView() *component.ResourceView {
	return g.resourceView
}

func (g *InGame) GetCalendar() *component.Calendar {
	return g.calendar
}

func (g *InGame) GetHistory() *component.History {
	return g.history
}

func (g *InGame) GetCardDeck() *component.CardDeck {
	return g.cardDeck
}

func (g *InGame) GetGameMain() *component.GameMain {
	return g.gameMain
}

// System manager getters
func (g *InGame) GetResourceManager() *system.ResourceManager {
	return g.resourceManager
}

func (g *InGame) GetTurnManager() *system.TurnManager {
	return g.turnManager
}

func (g *InGame) GetCardManager() *system.CardManager {
	return g.cardManager
}

func (g *InGame) GetTerritoryManager() *system.TerritoryManager {
	return g.territoryManager
}

func (g *InGame) GetAllianceManager() *system.AllianceManager {
	return g.allianceManager
}

func (g *InGame) GetCombatManager() *system.CombatManager {
	return g.combatManager
}
