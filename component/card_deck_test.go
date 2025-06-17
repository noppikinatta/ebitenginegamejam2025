package component_test

import (
	"testing"

	"github.com/noppikinatta/ebitenginegamejam2025/scene"
)

// TestCardDeckCanDisplayMultipleCardsHorizontally validates AddCard and layout
func TestCardDeckCanDisplayMultipleCardsHorizontally(t *testing.T) {
    g := scene.CreateSequence()
    g.SetCurrentScene("ingame")
    deck := g.GetInGameScene().GetCardDeck()
    if deck == nil { t.Fatal("card deck nil") }
    initLen := len(deck.GetCards())
    deck.AddCard("Unit A")
    deck.AddCard("Unit B")
    if len(deck.GetCards()) != initLen+2 {
        t.Errorf("expected %d cards, got %d", initLen+2, len(deck.GetCards()))
    }
    if deck.GetHorizontalLayout() == nil {
        t.Error("layout nil")
    }
} 