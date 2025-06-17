package component_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

// TestHistoryComponentCanAddAndDisplayEntries checks AddEntry/ GetEntries
func TestHistoryComponentCanAddAndDisplayEntries(t *testing.T) {
    game := scene.CreateSequence()
    game.SetCurrentScene("ingame")
    hist := game.GetInGameScene().GetHistory()
    if hist == nil {
        t.Fatal("history nil")
    }
    entry := "Test entry"
    hist.AddEntry(entry)
    found := false
    for _, e := range hist.GetEntries() {
        if e == entry { found = true; break }
    }
    if !found {
        t.Error("added entry not found")
    }
} 