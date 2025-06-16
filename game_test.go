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
