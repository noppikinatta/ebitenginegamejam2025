package component_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

// TestGameMainContainerSwitchesBetweenSubScreens checks screen switching
func TestGameMainContainerSwitchesBetweenSubScreens(t *testing.T) {
    game := scene.CreateSequence()
    game.SetCurrentScene("ingame")
    gm := game.GetInGameScene().GetGameMain()
    if gm == nil { t.Fatal("game main nil") }
    if gm.GetCurrentScreen() != "Map" {
        t.Errorf("expected Map, got %s", gm.GetCurrentScreen())
    }
    gm.SwitchToScreen("Diplomacy")
    if gm.GetCurrentScreen() != "Diplomacy" {
        t.Errorf("expected Diplomacy, got %s", gm.GetCurrentScreen())
    }
    gm.SwitchToScreen("Battle")
    if gm.GetCurrentScreen() != "Battle" {
        t.Errorf("expected Battle, got %s", gm.GetCurrentScreen())
    }
    gm.SwitchToScreen("Map")
    if gm.GetCurrentScreen() != "Map" {
        t.Errorf("expected Map again, got %s", gm.GetCurrentScreen())
    }
} 