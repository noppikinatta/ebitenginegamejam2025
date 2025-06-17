package scene

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// TestGameInitializesWithoutErrors checks CreateSequence and internal setup
func TestGameInitializesWithoutErrors(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("Game initialization panicked: %v", r)
        }
    }()

    game := CreateSequence()
    if game == nil {
        t.Error("CreateSequence() returned nil")
    }
    if game.Sequence == nil {
        t.Error("Game.Sequence is nil")
    }
}

// TestScreenResolutionCorrectlySet verifies Layout returns fixed 640x360
func TestScreenResolutionCorrectlySet(t *testing.T) {
    game := CreateSequence()
    w, h := game.Layout(1280, 720)
    if w != 640 || h != 360 {
        t.Errorf("expected 640x360, got %dx%d", w, h)
    }
}

// TestSceneManagerCanSwitchBetweenScenes validates SetCurrentScene logic
func TestSceneManagerCanSwitchBetweenScenes(t *testing.T) {
    game := CreateSequence()
    if game.GetCurrentScene() != "title" {
        t.Errorf("initial scene should be title, got %s", game.GetCurrentScene())
    }
    game.SetCurrentScene("ingame")
    if game.GetCurrentScene() != "ingame" {
        t.Errorf("scene should be ingame after switch, got %s", game.GetCurrentScene())
    }
}

// TestBasicInputDetection makes sure Update/Draw run without panic
func TestBasicInputDetection(t *testing.T) {
    game := CreateSequence()
    if err := game.Update(); err != nil {
        t.Errorf("Update returned error: %v", err)
    }
    screen := ebiten.NewImage(640, 360)
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("Draw panicked: %v", r)
        }
    }()
    game.Draw(screen)
}

// TestSceneTransitionsWorkProperly ensures TransitionTo handles switches
func TestSceneTransitionsWorkProperly(t *testing.T) {
    game := CreateSequence()
    game.TransitionTo("ingame")
    if game.GetCurrentScene() != "ingame" {
        t.Errorf("expected ingame, got %s", game.GetCurrentScene())
    }
    game.TransitionTo("result")
    if game.GetCurrentScene() != "result" {
        t.Errorf("expected result, got %s", game.GetCurrentScene())
    }
} 