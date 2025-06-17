package screen_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

// TestDiplomacyScreenShowsAvailableCards ensures cards, cost and npc info present
func TestDiplomacyScreenShowsAvailableCards(t *testing.T) {
    g := scene.CreateSequence(); g.SetCurrentScene("ingame")
    ds := g.GetInGameScene().GetGameMain()
    ds.SwitchToScreen("Diplomacy")
    screen := ds.GetDiplomacyScreen()
    if screen == nil { t.Fatal("diplomacy screen nil") }
    cards := screen.GetAvailableCards()
    if len(cards) == 0 { t.Fatal("no available cards") }
    if screen.GetCardCost(cards[0].Name) == nil {
        t.Error("card cost nil")
    }
    if screen.GetCurrentNPCInfo() == nil { t.Error("npc info nil") }
} 