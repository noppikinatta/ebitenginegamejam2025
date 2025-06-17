package system_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

func TestResourceValuesUpdateCorrectly(t *testing.T) {
    g := scene.CreateSequence(); g.SetCurrentScene("ingame")
    rm := g.GetInGameScene().GetResourceManager()
    init := rm.GetResource("Gold")
    rm.AddResource("Gold", 50)
    if rm.GetResource("Gold") != init+50 { t.Error("add resource failed") }
    if !rm.ConsumeResources(map[string]int{"Gold":25}) { t.Error("consume failed") }
    if rm.GetResource("Gold") != init+25 { t.Error("consume incorrect") }
} 