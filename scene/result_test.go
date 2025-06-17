package scene

import "testing"

// TestResultSceneDisplaysGameHistory checks history accessibility
func TestResultSceneDisplaysGameHistory(t *testing.T) {
	game := CreateSequence()
	game.SetCurrentScene("result")
	if game.GetCurrentScene() != "result" {
		t.Fatalf("expected result scene, got %s", game.GetCurrentScene())
	}
	res := game.GetResultScene()
	if res == nil {
		t.Fatal("result scene nil")
	}
	if res.GetGameHistory() == nil {
		t.Error("game history should not be nil")
	}
}
