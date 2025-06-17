package screen_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

// TestMapScreenGenerates13x7GridCorrectly verifies grid dimensions and totals
func TestMapScreenGenerates13x7GridCorrectly(t *testing.T) {
    game := scene.CreateSequence()
    game.SetCurrentScene("ingame")
    ms := game.GetInGameScene().GetGameMain().GetMapScreen()
    if ms == nil { t.Fatal("map screen nil") }
    grid := ms.GetGrid()
    if grid == nil { t.Fatal("grid nil") }
    if grid.GetWidth() != 13 || grid.GetHeight() != 7 {
        t.Errorf("expected 13x7 grid, got %dx%d", grid.GetWidth(), grid.GetHeight())
    }
    if grid.GetTotalPoints() != 91 {
        t.Errorf("expected 91 points, got %d", grid.GetTotalPoints())
    }
}

// TestPointTypesAreAssignedCorrectly checks Home/Boss and counts NPC/Wild
func TestPointTypesAreAssignedCorrectly(t *testing.T) {
    g := scene.CreateSequence(); g.SetCurrentScene("ingame")
    grid := g.GetInGameScene().GetGameMain().GetMapScreen().GetGrid()
    if grid.GetPoint(0,6).GetType() != "Home" { t.Error("home not Home type") }
    if grid.GetPoint(12,0).GetType() != "Boss" { t.Error("boss not Boss type") }
    wild, npc := 0,0
    for x:=0;x<13;x++{ for y:=0;y<7;y++{ switch grid.GetPoint(x,y).GetType(){case "Wild":wild++;case "NPC":npc++}}}
    if wild==0 {t.Error("no wild points")}
    if npc==0 {t.Error("no npc points")}
}

// TestPointConnectivityRules validates pathfinder access and path
func TestPointConnectivityRules(t *testing.T) {
    g := scene.CreateSequence(); g.SetCurrentScene("ingame")
    pf := g.GetInGameScene().GetGameMain().GetMapScreen().GetPathfinder()
    if !pf.IsPointAccessible(0,6) { t.Error("home should accessible") }
    if pf.IsPointAccessible(12,0) { t.Error("boss should not accessible initially") }
    if pf.FindPath(0,6,1,6) == nil { t.Error("path from home to adjacent should exist") }
} 