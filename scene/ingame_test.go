package scene

import "testing"

// TestInGameSceneUILayout ensures layout matches plan
func TestInGameSceneUILayout(t *testing.T) {
    game := CreateSequence()
    game.SetCurrentScene("ingame")

    inGame := game.GetInGameScene()
    if inGame == nil {
        t.Fatal("inGame scene nil")
    }
    layout := inGame.GetUILayout()
    if layout == nil {
        t.Fatal("layout nil")
    }
    expected := map[string][4]int{
        "GameMain":     {0, 20, 520, 280},
        "ResourceView": {0, 0, 520, 20},
        "Calendar":     {520, 0, 120, 40},
        "History":      {520, 80, 120, 320},
        "CardDeck":     {0, 300, 520, 60},
    }
    for name, want := range expected {
        got := layout.GetComponentBounds(name)
        if got != want {
            t.Errorf("%s bounds expected %v got %v", name, want, got)
        }
    }
} 