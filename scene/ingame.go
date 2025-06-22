package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InGame struct {
}

func NewInGame() *InGame {
	layout := NewUILayout()

	// Create system managers
	resourceManager := system.NewResourceManager()
	turnManager := system.NewTurnManager()
	// Try to load card templates from CSV; fall back to defaults on error.
	var cardManager *system.CardManager
	if cards, err := loader.LoadCards("data/cards.csv"); err == nil {
		cardManager = system.NewCardManagerFromData(cards)
	} else {
		cardManager = system.NewCardManager()
	}
	territoryManager := system.NewTerritoryManager(resourceManager)
	var allianceManager *system.AllianceManager
	if nations, err := loader.LoadNations("data/nations.csv"); err == nil {
		allianceManager = system.NewAllianceManagerFromData(nations)
	} else {
		allianceManager = system.NewAllianceManager()
	}
	var combatManager *system.CombatManager
	if enemyMap, errE := loader.LoadEnemies("data/enemies.csv"); errE == nil {
		bossMap, _ := loader.LoadBosses("data/bosses.csv")
		combatManager = system.NewCombatManagerWithTemplates(enemyMap, bossMap)
	} else {
		combatManager = system.NewCombatManager()
	}

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
	return nil
}

func (g *InGame) Draw(screen *ebiten.Image) {

}

func (g *InGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}
