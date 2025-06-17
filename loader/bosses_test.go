package loader_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/loader"
)

// T8.5 (part): Verify LoadBosses produce expected stats.
func TestLoadBossesStats(t *testing.T) {
    bosses, err := loader.LoadBosses("../testdata/bosses.csv")
    if err != nil {
        t.Fatalf("LoadBosses error: %v", err)
    }
    if hp := bosses["B1"].Health; hp != 50 {
        t.Fatalf("boss health mismatch: got %d want 50", hp)
    }
} 