package component_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

// TestResourceViewDisplaysAllResourceTypes checks resource list and amounts
func TestResourceViewDisplaysAllResourceTypes(t *testing.T) {
    game := scene.CreateSequence()
    game.SetCurrentScene("ingame")
    rv := game.GetInGameScene().GetResourceView()
    if rv == nil {
        t.Fatal("resource view nil")
    }
    expected := []string{"Gold", "Iron", "Wood", "Grain", "Mana"}
    got := rv.GetResourceTypes()
    if len(got) != 5 {
        t.Fatalf("expected 5 resources, got %d", len(got))
    }
    for _, e := range expected {
        found := false
        for _, g := range got {
            if e == g { found = true; break }
        }
        if !found {
            t.Errorf("resource %s missing", e)
        }
        if rv.GetResourceAmount(e) < 0 {
            t.Errorf("resource %s negative amount", e)
        }
    }
} 