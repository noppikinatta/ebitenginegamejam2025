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
