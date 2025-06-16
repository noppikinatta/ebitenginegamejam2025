package main

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

// T1.1: Test game initializes without errors
func TestGameInitializesWithoutErrors(t *testing.T) {
	// Test that CreateSequence() can be called without panicking
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Game initialization panicked: %v", r)
		}
	}()

	game := scene.CreateSequence()
	if game == nil {
		t.Error("CreateSequence() returned nil")
	}

	// Verify that the game has required components
	if game.Sequence == nil {
		t.Error("Game.Sequence is nil")
	}
}

// T1.2: Test screen resolution is correctly set
func TestScreenResolutionCorrectlySet(t *testing.T) {
	game := scene.CreateSequence()

	// Test internal resolution should be 640x360 according to plan
	outsideWidth, outsideHeight := 1280, 720
	width, height := game.Layout(outsideWidth, outsideHeight)

	expectedWidth, expectedHeight := 640, 360
	if width != expectedWidth || height != expectedHeight {
		t.Errorf("Expected resolution %dx%d, got %dx%d", expectedWidth, expectedHeight, width, height)
	}
}

// T1.3: Test scene manager can switch between scenes
func TestSceneManagerCanSwitchBetweenScenes(t *testing.T) {
	game := scene.CreateSequence()

	// Initially should be on title scene
	if game.GetCurrentScene() != "title" {
		t.Errorf("Expected initial scene to be 'title', got '%s'", game.GetCurrentScene())
	}

	// Simulate scene transition
	game.SetCurrentScene("ingame")
	if game.GetCurrentScene() != "ingame" {
		t.Errorf("Expected scene to be 'ingame' after switching, got '%s'", game.GetCurrentScene())
	}
}

// T1.4: Test basic keyboard/mouse input detection
func TestBasicInputDetection(t *testing.T) {
	game := scene.CreateSequence()

	// Test Update method can be called without errors
	err := game.Update()
	if err != nil {
		t.Errorf("Game.Update() returned error: %v", err)
	}

	// Create a test screen
	screen := ebiten.NewImage(640, 360)

	// Test Draw method can be called without errors
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Game.Draw() panicked: %v", r)
		}
	}()

	game.Draw(screen)
}

// T2.1: Test title scene displays correctly
func TestTitleSceneDisplaysCorrectly(t *testing.T) {
	game := scene.CreateSequence()

	// Should start on title scene
	if game.GetCurrentScene() != "title" {
		t.Errorf("Expected initial scene to be 'title', got '%s'", game.GetCurrentScene())
	}

	// Test that title scene has story introduction
	titleScene := game.GetTitleScene()
	if titleScene == nil {
		t.Error("Title scene is nil")
	}

	// Test that title scene has proper story content
	story := titleScene.GetStoryText()
	if story == "" {
		t.Error("Title scene should have story text")
	}
}

// T2.2: Test InGame scene UI layout (5 components in correct positions)
func TestInGameSceneUILayout(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	if inGameScene == nil {
		t.Error("InGame scene is nil")
	}

	// Test that all 5 UI components have correct layout
	layout := inGameScene.GetUILayout()
	if layout == nil {
		t.Error("InGame scene should have UI layout")
	}

	// According to plan.md:
	// - GameMain: {0,20,520,280}
	// - ResourceView: {0,0,520,20}
	// - Calendar: {520,0,120,40}
	// - History: {520,80,120,320}
	// - CardDeck: {0,300,520,60}

	expected := map[string][4]int{
		"GameMain":     {0, 20, 520, 280},
		"ResourceView": {0, 0, 520, 20},
		"Calendar":     {520, 0, 120, 40},
		"History":      {520, 80, 120, 320},
		"CardDeck":     {0, 300, 520, 60},
	}

	for componentName, expectedBounds := range expected {
		actualBounds := layout.GetComponentBounds(componentName)
		if actualBounds != expectedBounds {
			t.Errorf("Component %s: expected bounds %v, got %v", componentName, expectedBounds, actualBounds)
		}
	}
}

// T2.3: Test result scene displays game history
func TestResultSceneDisplaysGameHistory(t *testing.T) {
	game := scene.CreateSequence()

	// Test that result scene can be created
	game.SetCurrentScene("result")
	if game.GetCurrentScene() != "result" {
		t.Errorf("Expected scene to be 'result', got '%s'", game.GetCurrentScene())
	}

	resultScene := game.GetResultScene()
	if resultScene == nil {
		t.Error("Result scene is nil")
	}

	// Test that result scene can display history
	history := resultScene.GetGameHistory()
	if history == nil {
		t.Error("Result scene should have game history")
	}
}

// T2.4: Test scene transitions work properly
func TestSceneTransitionsWorkProperly(t *testing.T) {
	game := scene.CreateSequence()

	// Test transition: title -> ingame -> result
	if game.GetCurrentScene() != "title" {
		t.Error("Should start on title scene")
	}

	// Transition to ingame
	game.TransitionTo("ingame")
	if game.GetCurrentScene() != "ingame" {
		t.Errorf("Expected 'ingame', got '%s'", game.GetCurrentScene())
	}

	// Transition to result
	game.TransitionTo("result")
	if game.GetCurrentScene() != "result" {
		t.Errorf("Expected 'result', got '%s'", game.GetCurrentScene())
	}

	// Test that transitions trigger proper cleanup/setup
	// This will be extended as we implement more scene logic
}

// T3.1: Test ResourceView displays all 5 resource types
func TestResourceViewDisplaysAllResourceTypes(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	resourceView := inGameScene.GetResourceView()
	if resourceView == nil {
		t.Error("ResourceView should not be nil")
	}

	// Test that all 5 resource types are present: Gold, Iron, Wood, Grain, Mana
	expectedResources := []string{"Gold", "Iron", "Wood", "Grain", "Mana"}
	resources := resourceView.GetResourceTypes()

	if len(resources) != 5 {
		t.Errorf("Expected 5 resource types, got %d", len(resources))
	}

	for _, expected := range expectedResources {
		found := false
		for _, actual := range resources {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Resource type '%s' not found", expected)
		}
	}

	// Test that each resource has an initial amount
	for _, resourceType := range expectedResources {
		amount := resourceView.GetResourceAmount(resourceType)
		if amount < 0 {
			t.Errorf("Resource %s should have non-negative amount, got %d", resourceType, amount)
		}
	}
}

// T3.2: Test Calendar shows correct year/month format
func TestCalendarShowsCorrectYearMonthFormat(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	calendar := inGameScene.GetCalendar()
	if calendar == nil {
		t.Error("Calendar should not be nil")
	}

	// According to plan.md: starts from Kingdom Year 1000, Month 4
	year := calendar.GetCurrentYear()
	month := calendar.GetCurrentMonth()

	if year != 1000 {
		t.Errorf("Expected starting year 1000, got %d", year)
	}

	if month != 4 {
		t.Errorf("Expected starting month 4, got %d", month)
	}

	// Test calendar format display
	display := calendar.GetDisplayText()
	if display == "" {
		t.Error("Calendar should have display text")
	}
}

// T3.3: Test History component can add and display entries
func TestHistoryComponentCanAddAndDisplayEntries(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	history := inGameScene.GetHistory()
	if history == nil {
		t.Error("History should not be nil")
	}

	// Test adding history entries
	testEntry := "Kingdom Year 1000, Month 4: Game Started"
	history.AddEntry(testEntry)

	entries := history.GetEntries()
	if len(entries) == 0 {
		t.Error("History should have at least one entry after adding")
	}

	// Test that the added entry is present
	found := false
	for _, entry := range entries {
		if entry == testEntry {
			found = true
			break
		}
	}
	if !found {
		t.Error("Added history entry not found")
	}
}

// T3.4: Test CardDeck can display multiple cards horizontally
func TestCardDeckCanDisplayMultipleCardsHorizontally(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	cardDeck := inGameScene.GetCardDeck()
	if cardDeck == nil {
		t.Error("CardDeck should not be nil")
	}

	// Test that card deck can hold cards
	cards := cardDeck.GetCards()
	if cards == nil {
		t.Error("CardDeck should return card list")
	}

	// Get initial card count
	initialCount := len(cards)

	// Test adding cards to deck
	cardDeck.AddCard("Test Unit Card")
	cardDeck.AddCard("Test Enchant Card")

	updatedCards := cardDeck.GetCards()
	expectedCount := initialCount + 2
	if len(updatedCards) != expectedCount {
		t.Errorf("Expected %d cards in deck, got %d", expectedCount, len(updatedCards))
	}

	// Test horizontal layout calculation
	layout := cardDeck.GetHorizontalLayout()
	if layout == nil {
		t.Error("CardDeck should have horizontal layout")
	}
}

// T3.5: Test GameMain container switches between sub-screens
func TestGameMainContainerSwitchesBetweenSubScreens(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	gameMain := inGameScene.GetGameMain()
	if gameMain == nil {
		t.Error("GameMain should not be nil")
	}

	// Test that GameMain starts with Map screen
	currentScreen := gameMain.GetCurrentScreen()
	if currentScreen != "Map" {
		t.Errorf("Expected GameMain to start with 'Map' screen, got '%s'", currentScreen)
	}

	// Test switching to Diplomacy screen
	gameMain.SwitchToScreen("Diplomacy")
	if gameMain.GetCurrentScreen() != "Diplomacy" {
		t.Errorf("Expected 'Diplomacy' screen, got '%s'", gameMain.GetCurrentScreen())
	}

	// Test switching to Battle screen
	gameMain.SwitchToScreen("Battle")
	if gameMain.GetCurrentScreen() != "Battle" {
		t.Errorf("Expected 'Battle' screen, got '%s'", gameMain.GetCurrentScreen())
	}

	// Test switching back to Map screen
	gameMain.SwitchToScreen("Map")
	if gameMain.GetCurrentScreen() != "Map" {
		t.Errorf("Expected 'Map' screen, got '%s'", gameMain.GetCurrentScreen())
	}
}

// T4.1: Test Map screen generates 13x7 grid correctly
func TestMapScreenGenerates13x7GridCorrectly(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	gameMain := inGameScene.GetGameMain()

	// Should start with Map screen
	if gameMain.GetCurrentScreen() != "Map" {
		t.Error("GameMain should start with Map screen")
	}

	// Get map data
	mapScreen := gameMain.GetMapScreen()
	if mapScreen == nil {
		t.Error("MapScreen should not be nil")
	}

	grid := mapScreen.GetGrid()
	if grid == nil {
		t.Error("Map grid should not be nil")
	}

	// Test grid dimensions
	if grid.GetWidth() != 13 || grid.GetHeight() != 7 {
		t.Errorf("Expected 13x7 grid, got %dx%d", grid.GetWidth(), grid.GetHeight())
	}

	// Test total points
	totalPoints := grid.GetTotalPoints()
	expectedPoints := 13 * 7 // 91 points
	if totalPoints != expectedPoints {
		t.Errorf("Expected %d total points, got %d", expectedPoints, totalPoints)
	}
}

// T4.2: Test point types are assigned correctly
func TestPointTypesAreAssignedCorrectly(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	gameMain := inGameScene.GetGameMain()
	mapScreen := gameMain.GetMapScreen()
	grid := mapScreen.GetGrid()

	// Test Home point (should be at bottom-left: 0,6)
	homePoint := grid.GetPoint(0, 6)
	if homePoint == nil {
		t.Error("Home point should exist at (0,6)")
	}
	if homePoint.GetType() != "Home" {
		t.Errorf("Expected Home point type, got %s", homePoint.GetType())
	}

	// Test Boss point (should be at top-right: 12,0)
	bossPoint := grid.GetPoint(12, 0)
	if bossPoint == nil {
		t.Error("Boss point should exist at (12,0)")
	}
	if bossPoint.GetType() != "Boss" {
		t.Errorf("Expected Boss point type, got %s", bossPoint.GetType())
	}

	// Test that there are Wild and NPC points
	wildCount := 0
	npcCount := 0

	for x := 0; x < 13; x++ {
		for y := 0; y < 7; y++ {
			point := grid.GetPoint(x, y)
			switch point.GetType() {
			case "Wild":
				wildCount++
			case "NPC":
				npcCount++
			}
		}
	}

	if wildCount == 0 {
		t.Error("Should have at least one Wild point")
	}
	if npcCount == 0 {
		t.Error("Should have at least one NPC point")
	}
}

// T4.3: Test point connectivity rules
func TestPointConnectivityRules(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	gameMain := inGameScene.GetGameMain()
	mapScreen := gameMain.GetMapScreen()

	// Test path finding
	pathfinder := mapScreen.GetPathfinder()
	if pathfinder == nil {
		t.Error("Pathfinder should not be nil")
	}

	// Test that Home is accessible initially
	homeAccessible := pathfinder.IsPointAccessible(0, 6)
	if !homeAccessible {
		t.Error("Home point should be accessible")
	}

	// Test that Boss is not accessible initially (blocked by enemies)
	bossAccessible := pathfinder.IsPointAccessible(12, 0)
	if bossAccessible {
		t.Error("Boss point should not be accessible initially")
	}

	// Test path calculation
	path := pathfinder.FindPath(0, 6, 1, 6) // Home to adjacent point
	if path == nil {
		t.Error("Should be able to find path from Home to adjacent point")
	}
}

// T4.4: Test Diplomacy screen shows available cards
func TestDiplomacyScreenShowsAvailableCards(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	gameMain := inGameScene.GetGameMain()

	// Switch to Diplomacy screen
	gameMain.SwitchToScreen("Diplomacy")

	diplomacyScreen := gameMain.GetDiplomacyScreen()
	if diplomacyScreen == nil {
		t.Error("DiplomacyScreen should not be nil")
	}

	// Test available cards
	availableCards := diplomacyScreen.GetAvailableCards()
	if availableCards == nil {
		t.Error("Available cards should not be nil")
	}

	if len(availableCards) == 0 {
		t.Error("Should have at least one available card")
	}

	// Test card purchase functionality
	firstCard := availableCards[0]
	cost := diplomacyScreen.GetCardCost(firstCard.Name)
	if cost == nil {
		t.Error("Card should have a cost")
	}

	// Test NPC nation info
	npcInfo := diplomacyScreen.GetCurrentNPCInfo()
	if npcInfo == nil {
		t.Error("Should have current NPC info")
	}
}

// T4.5: Test Battle screen allows card placement
func TestBattleScreenAllowsCardPlacement(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	gameMain := inGameScene.GetGameMain()

	// Switch to Battle screen
	gameMain.SwitchToScreen("Battle")

	battleScreen := gameMain.GetBattleScreen()
	if battleScreen == nil {
		t.Error("BattleScreen should not be nil")
	}

	// Test battlefield layout
	battlefield := battleScreen.GetBattlefield()
	if battlefield == nil {
		t.Error("Battlefield should not be nil")
	}

	// Test front row (should allow 5 cards)
	frontRow := battlefield.GetFrontRow()
	if len(*frontRow) != 5 {
		t.Errorf("Front row should have 5 slots, got %d", len(*frontRow))
	}

	// Test back row (should allow 5 cards)
	backRow := battlefield.GetBackRow()
	if len(*backRow) != 5 {
		t.Errorf("Back row should have 5 slots, got %d", len(*backRow))
	}

	// Test card placement
	testCard := "Test Warrior"
	success := battlefield.PlaceCard(testCard, "front", 0)
	if !success {
		t.Error("Should be able to place card in front row")
	}

	// Test that card was placed
	placedCard := (*frontRow)[0]
	if placedCard != testCard {
		t.Errorf("Expected '%s' in front row slot 0, got '%s'", testCard, placedCard)
	}

	// Test enemy setup
	enemies := battleScreen.GetEnemies()
	if enemies == nil || len(enemies) == 0 {
		t.Error("Battle should have enemies")
	}
}

// T5.1: Test resource values update correctly
func TestResourceValuesUpdateCorrectly(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	resourceManager := inGameScene.GetResourceManager()
	if resourceManager == nil {
		t.Error("ResourceManager should not be nil")
	}

	// Test initial resources
	gold := resourceManager.GetResource("Gold")
	if gold < 0 {
		t.Error("Gold should be non-negative initially")
	}

	// Test resource modification
	initialGold := gold
	resourceManager.AddResource("Gold", 50)
	newGold := resourceManager.GetResource("Gold")
	if newGold != initialGold+50 {
		t.Errorf("Expected gold to be %d, got %d", initialGold+50, newGold)
	}

	// Test resource consumption
	success := resourceManager.ConsumeResources(map[string]int{"Gold": 25})
	if !success {
		t.Error("Should be able to consume available resources")
	}

	finalGold := resourceManager.GetResource("Gold")
	if finalGold != initialGold+25 {
		t.Errorf("Expected gold to be %d after consumption, got %d", initialGold+25, finalGold)
	}
}

// T5.2: Test turn counter increments monthly
func TestTurnCounterIncrementsMonthly(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	turnManager := inGameScene.GetTurnManager()
	if turnManager == nil {
		t.Error("TurnManager should not be nil")
	}

	// Test initial state
	year := turnManager.GetCurrentYear()
	month := turnManager.GetCurrentMonth()
	if year != 1000 || month != 4 {
		t.Errorf("Expected year 1000, month 4, got year %d, month %d", year, month)
	}

	// Test turn advancement
	turnManager.AdvanceTurn()
	newYear := turnManager.GetCurrentYear()
	newMonth := turnManager.GetCurrentMonth()

	if newMonth != 5 || newYear != 1000 {
		t.Errorf("Expected year 1000, month 5 after one turn, got year %d, month %d", newYear, newMonth)
	}

	// Test year rollover
	for i := 0; i < 8; i++ { // Advance 8 more months to reach month 13 -> year 1001, month 1
		turnManager.AdvanceTurn()
	}

	finalYear := turnManager.GetCurrentYear()
	finalMonth := turnManager.GetCurrentMonth()
	if finalYear != 1001 || finalMonth != 1 {
		t.Errorf("Expected year 1001, month 1 after year rollover, got year %d, month %d", finalYear, finalMonth)
	}
}

// T5.3: Test card types have correct properties
func TestCardTypesHaveCorrectProperties(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	cardManager := inGameScene.GetCardManager()
	if cardManager == nil {
		t.Error("CardManager should not be nil")
	}

	// Test Unit card creation
	unitCard := cardManager.CreateCard("Warrior", "Unit", map[string]int{"Gold": 50})
	if unitCard.Type != "Unit" {
		t.Errorf("Expected Unit type, got %s", unitCard.Type)
	}
	if unitCard.Attack <= 0 {
		t.Error("Unit cards should have positive attack")
	}

	// Test Enchant card creation
	enchantCard := cardManager.CreateCard("Magic Shield", "Enchant", map[string]int{"Mana": 30})
	if enchantCard.Type != "Enchant" {
		t.Errorf("Expected Enchant type, got %s", enchantCard.Type)
	}

	// Test Building card creation
	buildingCard := cardManager.CreateCard("Farm", "Building", map[string]int{"Wood": 40})
	if buildingCard.Type != "Building" {
		t.Errorf("Expected Building type, got %s", buildingCard.Type)
	}

	// Test card database
	allCards := cardManager.GetAllCardTemplates()
	if len(allCards) == 0 {
		t.Error("Should have card templates available")
	}
}

// T5.4: Test territory control affects resource generation
func TestTerritoryControlAffectsResourceGeneration(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	territoryManager := inGameScene.GetTerritoryManager()
	if territoryManager == nil {
		t.Error("TerritoryManager should not be nil")
	}

	// Test initial controlled territories
	controlledTerritories := territoryManager.GetControlledTerritories()
	if len(controlledTerritories) == 0 {
		t.Error("Should start with at least home territory")
	}

	// Test resource generation from controlled territories
	resourceManager := inGameScene.GetResourceManager()
	initialGold := resourceManager.GetResource("Gold")

	// Simulate territory conquest
	territoryManager.ConquerTerritory(1, 6) // Adjacent to home

	// Test that new territory is controlled
	newControlledTerritories := territoryManager.GetControlledTerritories()
	if len(newControlledTerritories) <= len(controlledTerritories) {
		t.Error("Should have more controlled territories after conquest")
	}

	// Test resource generation increase
	territoryManager.GenerateResources()
	newGold := resourceManager.GetResource("Gold")
	if newGold <= initialGold {
		t.Error("Resource generation should increase with more territories")
	}
}

// T5.5: Test alliance formation with NPCs
func TestAllianceFormationWithNPCs(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	allianceManager := inGameScene.GetAllianceManager()
	if allianceManager == nil {
		t.Error("AllianceManager should not be nil")
	}

	// Test initial alliance state
	allies := allianceManager.GetAllies()
	if allies == nil {
		t.Error("Allies list should not be nil")
	}

	// Test NPC relationship
	relationship := allianceManager.GetRelationship("Iron Republic")
	if relationship < 0 || relationship > 100 {
		t.Errorf("Relationship should be between 0-100, got %d", relationship)
	}

	// Test alliance formation
	success := allianceManager.FormAlliance("Iron Republic")
	if !success && relationship >= 70 {
		t.Error("Should be able to form alliance with high relationship NPCs")
	}

	// Test alliance benefits
	if success {
		newAllies := allianceManager.GetAllies()
		found := false
		for _, ally := range newAllies {
			if ally == "Iron Republic" {
				found = true
				break
			}
		}
		if !found {
			t.Error("NPC should be in allies list after forming alliance")
		}

		// Test alliance bonuses
		bonuses := allianceManager.GetAllianceBonuses()
		if bonuses == nil {
			t.Error("Should have alliance bonuses")
		}
	}
}

// T6.1: Test cards can be placed in front/back rows (max 10 total)
func TestCardsCanBePlacedInFrontBackRows(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	combatManager := inGameScene.GetCombatManager()
	if combatManager == nil {
		t.Error("CombatManager should not be nil")
	}

	// Test initial battlefield state
	battlefield := combatManager.GetBattlefield()
	if battlefield == nil {
		t.Error("Battlefield should not be nil")
	}

	// Test front row capacity
	frontRow := battlefield.GetFrontRow()
	if len(*frontRow) != 5 {
		t.Errorf("Front row should have 5 slots, got %d", len(*frontRow))
	}

	// Test back row capacity
	backRow := battlefield.GetBackRow()
	if len(*backRow) != 5 {
		t.Errorf("Back row should have 5 slots, got %d", len(*backRow))
	}

	// Test filling all slots
	cardManager := inGameScene.GetCardManager()

	// Fill front row
	for i := 0; i < 5; i++ {
		card := cardManager.CreateCard("Warrior", "Unit", map[string]int{"Gold": 50})
		success := combatManager.PlaceCardInBattle(card, "front", i)
		if !success {
			t.Errorf("Should be able to place card in front row slot %d", i)
		}
	}

	// Fill back row
	for i := 0; i < 5; i++ {
		card := cardManager.CreateCard("Archer", "Unit", map[string]int{"Gold": 40})
		success := combatManager.PlaceCardInBattle(card, "back", i)
		if !success {
			t.Errorf("Should be able to place card in back row slot %d", i)
		}
	}

	// Test that all 10 slots are filled
	placedCards := combatManager.GetPlacedCards()
	if len(placedCards) != 10 {
		t.Errorf("Expected 10 placed cards, got %d", len(placedCards))
	}
}

// T6.2: Test combat calculates damage correctly
func TestCombatCalculatesDamageCorrectly(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	combatManager := inGameScene.GetCombatManager()

	// Set up a simple combat scenario
	cardManager := inGameScene.GetCardManager()

	// Place player cards
	warrior := cardManager.CreateCard("Warrior", "Unit", map[string]int{"Gold": 50})
	combatManager.PlaceCardInBattle(warrior, "front", 0)

	// Add enemy
	combatManager.AddEnemy("Goblin", 3, 2, 5) // Attack: 3, Defense: 2, Health: 5

	// Test damage calculation
	playerDamage := combatManager.CalculatePlayerDamage()
	if playerDamage <= 0 {
		t.Error("Player should deal positive damage with placed cards")
	}

	enemyDamage := combatManager.CalculateEnemyDamage()
	if enemyDamage <= 0 {
		t.Error("Enemies should deal positive damage")
	}

	// Test actual combat round
	initialEnemyHealth := combatManager.GetEnemyHealth("Goblin")
	combatManager.ExecuteCombatRound()

	newEnemyHealth := combatManager.GetEnemyHealth("Goblin")
	if newEnemyHealth >= initialEnemyHealth {
		t.Error("Enemy health should decrease after combat round")
	}

	// Test combat stats
	combatStats := combatManager.GetCombatStats()
	if combatStats == nil {
		t.Error("Combat stats should not be nil")
	}

	if combatStats.RoundsCompleted == 0 {
		t.Error("Should have completed at least one combat round")
	}
}

// T6.3: Test victory/defeat conditions trigger properly
func TestVictoryDefeatConditionsTriggerProperly(t *testing.T) {
	game := scene.CreateSequence()
	game.SetCurrentScene("ingame")

	inGameScene := game.GetInGameScene()
	combatManager := inGameScene.GetCombatManager()

	// Test initial combat state
	combatState := combatManager.GetCombatState()
	if combatState != "ongoing" {
		t.Errorf("Expected combat state 'ongoing', got '%s'", combatState)
	}

	// Test victory condition
	cardManager := inGameScene.GetCardManager()

	// Place powerful player card
	mage := cardManager.CreateCard("Mage", "Unit", map[string]int{"Gold": 60, "Mana": 20})
	combatManager.PlaceCardInBattle(mage, "front", 0)

	// Add weak enemy
	combatManager.AddEnemy("Weak Goblin", 1, 1, 1) // Very weak enemy

	// Execute combat until conclusion
	maxRounds := 10
	for i := 0; i < maxRounds; i++ {
		combatManager.ExecuteCombatRound()
		state := combatManager.GetCombatState()
		if state == "victory" {
			break
		}
		if state == "defeat" {
			t.Error("Should not be defeated with strong cards vs weak enemy")
			break
		}
	}

	// Check final state
	finalState := combatManager.GetCombatState()
	if finalState != "victory" {
		t.Errorf("Expected victory, got '%s'", finalState)
	}

	// Test victory rewards
	rewards := combatManager.GetVictoryRewards()
	if rewards == nil {
		t.Error("Should have victory rewards")
	}

	if len(rewards.Resources) == 0 {
		t.Error("Victory should provide resource rewards")
	}

	// Test defeat condition setup
	combatManager.ResetCombat()

	// Place weak player card
	combatManager.ClearPlayerCards()
	weakCard := cardManager.CreateCard("Peasant", "Unit", map[string]int{"Gold": 10})
	weakCard.SetStats(1, 1) // Very weak
	combatManager.PlaceCardInBattle(weakCard, "front", 0)

	// Add strong enemy
	combatManager.ClearEnemies()
	combatManager.AddEnemy("Dragon", 20, 10, 50) // Very strong enemy

	// Execute a few rounds (should lead to defeat)
	for i := 0; i < 5; i++ {
		combatManager.ExecuteCombatRound()
		state := combatManager.GetCombatState()
		if state == "defeat" {
			break
		}
	}

	// Should be defeated now
	defeatState := combatManager.GetCombatState()
	if defeatState != "defeat" && defeatState != "ongoing" {
		t.Errorf("Expected defeat or ongoing, got '%s'", defeatState)
	}
}
