package system_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

func TestTerritoryControlAffectsResourceGeneration(t *testing.T) {
    g := scene.CreateSequence(); g.SetCurrentScene("ingame")
    tm := g.GetInGameScene().GetTerritoryManager()
    rm := g.GetInGameScene().GetResourceManager()
    initialGold := rm.GetResource("Gold")
    tm.ConquerTerritory(1,6)
    tm.GenerateResources()
    if rm.GetResource("Gold") <= initialGold { t.Error("resource not increased") }
} 