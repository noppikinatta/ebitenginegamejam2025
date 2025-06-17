package system_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

func TestCardTypesHaveCorrectProperties(t *testing.T) {
    g := scene.CreateSequence(); g.SetCurrentScene("ingame")
    cm := g.GetInGameScene().GetCardManager()
    unit := cm.CreateCard("Warrior","Unit",map[string]int{"Gold":50})
    if unit.Type!="Unit" || unit.Attack<=0 { t.Error("unit card properties wrong") }
    enchant := cm.CreateCard("Magic Shield","Enchant",map[string]int{"Mana":30})
    if enchant.Type!="Enchant" { t.Error("enchant type wrong") }
    build := cm.CreateCard("Farm","Building",map[string]int{"Wood":40})
    if build.Type!="Building" { t.Error("building type wrong") }
    if len(cm.GetAllCardTemplates())==0 { t.Error("template empty") }
} 