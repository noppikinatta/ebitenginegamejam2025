package scene

import "testing"

// TestTitleSceneDisplaysCorrectly validates story text exists
func TestTitleSceneDisplaysCorrectly(t *testing.T) {
    game := CreateSequence()
    if game.GetCurrentScene() != "title" {
        t.Fatalf("expected initial scene title, got %s", game.GetCurrentScene())
    }
    titleScene := game.GetTitleScene()
    if titleScene == nil {
        t.Fatal("title scene nil")
    }
    if titleScene.GetStoryText() == "" {
        t.Error("title scene should have story text")
    }
} 