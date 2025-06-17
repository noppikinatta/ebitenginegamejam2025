package screen_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

// TestBattleScreenAllowsCardPlacement checks battlefield slots and card placement
func TestBattleScreenAllowsCardPlacement(t *testing.T) {
    g := scene.CreateSequence(); g.SetCurrentScene("ingame")
    gm := g.GetInGameScene().GetGameMain()
    gm.SwitchToScreen("Battle")
    bs := gm.GetBattleScreen()
    if bs == nil { t.Fatal("battle screen nil") }
    bf := bs.GetBattlefield()
    front, back := bf.GetFrontRow(), bf.GetBackRow()
    if len(*front)!=5 || len(*back)!=5 { t.Fatalf("row size mismatch") }
    if !bf.PlaceCard("Test C", "front", 0) { t.Error("place card failed") }
    if (*front)[0] != "Test C" { t.Error("card not placed correctly") }
    if len(bs.GetEnemies())==0 { t.Error("enemies list empty") }
} 