package loader_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/loader"
	"github.com/noppikinatta/ebitenginegamejam2025/system"
)

// T8.1: Verify LoadCards parses headers and row count matches CSV.
func TestLoadCards_ParseAndCount(t *testing.T) {
	cards, err := loader.LoadCards("../testdata/cards.csv")
	if err != nil {
		t.Fatalf("LoadCards returned error: %v", err)
	}
	if got, want := len(cards), 2; got != want {
		t.Fatalf("row count mismatch: got %d want %d", got, want)
	}
}

// T8.2: Verify CardManager returns stats identical to CSV values.
func TestCardManager_StatsFromCSV(t *testing.T) {
	cards, _ := loader.LoadCards("../testdata/cards.csv")
	cm := system.NewCardManagerFromData(cards)

	c := cm.GetCardTemplate("C1")
	if c == nil {
		t.Fatalf("expected card template for C1")
	}
	if c.Attack != 4 || c.Defense != 3 {
		t.Fatalf("card stats mismatch: got atk=%d def=%d", c.Attack, c.Defense)
	}
	if gold := c.Cost["Gold"]; gold != 50 {
		t.Fatalf("card cost mismatch: gold=%d want 50", gold)
	}
}
