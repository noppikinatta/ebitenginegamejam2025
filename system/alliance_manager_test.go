package system_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

func TestAllianceFormationWithNPCs(t *testing.T) {
    g := scene.CreateSequence(); g.SetCurrentScene("ingame")
    am := g.GetInGameScene().GetAllianceManager()
    rel := am.GetRelationship("Iron Republic")
    if rel<0 || rel>100 { t.Fatalf("relationship out of range") }
    am.FormAlliance("Iron Republic")
    if am.GetAllianceBonuses()==nil { t.Error("bonuses nil") }
} 