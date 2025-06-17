package system_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

func TestTurnCounterIncrementsMonthly(t *testing.T) {
    g := scene.CreateSequence(); g.SetCurrentScene("ingame")
    tm := g.GetInGameScene().GetTurnManager()
    if tm.GetCurrentYear()!=1000 || tm.GetCurrentMonth()!=4 { t.Fatalf("start date incorrect") }
    tm.AdvanceTurn()
    if tm.GetCurrentMonth()!=5 { t.Error("advance turn failed") }
    for i:=0; i<8; i++ { tm.AdvanceTurn() }
    if tm.GetCurrentYear()!=1001 || tm.GetCurrentMonth()!=1 { t.Error("year rollover failed") }
} 