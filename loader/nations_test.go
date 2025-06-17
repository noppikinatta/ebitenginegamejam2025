package loader_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/loader"
	"github.com/noppikinatta/ebitenginegamejam2025/system"
)

// T8.3: Verify LoadNations builds map with correct initial relationship values.
func TestLoadNations_InitialRelationship(t *testing.T) {
    nations, err := loader.LoadNations("../testdata/nations.csv")
    if err != nil {
        t.Fatalf("LoadNations error: %v", err)
    }
    if rel := nations["N1"].InitialRelationship; rel != 50 {
        t.Fatalf("relationship mismatch: got %d want 50", rel)
    }
}

// T8.4: Verify AllianceManager reflects CSV ally bonuses after load.
func TestAllianceManager_BonusesFromCSV(t *testing.T) {
    nations, _ := loader.LoadNations("../testdata/nations.csv")
    am := system.NewAllianceManagerFromData(nations)
    // Increase relationship to allow alliance
    am.ImproveRelationship("Iron Republic", 20)
    if !am.FormAlliance("Iron Republic") {
        t.Fatalf("failed to form alliance")
    }
    bonus := am.GetAllianceBonuses()
    if g := bonus.ResourceBonus["Gold"]; g != 5 {
        t.Fatalf("gold bonus mismatch: got %d want 5", g)
    }
    if bonus.MilitaryBonus != 10 {
        t.Fatalf("military bonus mismatch: got %d want 10", bonus.MilitaryBonus)
    }
} 