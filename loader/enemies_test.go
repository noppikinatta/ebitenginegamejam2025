package loader_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/loader"
)

// T8.5 (part): Verify LoadEnemies produce expected stats.
func TestLoadEnemiesStats(t *testing.T) {
    enemies, err := loader.LoadEnemies("../testdata/enemies.csv")
    if err != nil {
        t.Fatalf("LoadEnemies error: %v", err)
    }
    if e := enemies["E1"].Attack; e != 3 {
        t.Fatalf("enemy attack mismatch: got %d want 3", e)
    }
}

// T8.6: Verify malformed CSV row is reported as error but does not crash loader.
func TestLoadEnemies_MalformedRowError(t *testing.T) {
    if _, err := loader.LoadEnemies("../testdata/enemies_bad.csv"); err == nil {
        t.Fatalf("expected error for malformed CSV, got nil")
    }
} 