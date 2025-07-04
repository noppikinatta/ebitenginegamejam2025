package core

// GameState manages the overall state of the game.
type GameState struct {
	MyNation      *MyNation      // Player's nation
	CardDeck      *CardDeck      // Player's card deck
	MapGrid       *MapGrid       // Map grid
	Treasury      *Treasury      // Player's treasury
	CurrentTurn   int            // Current turn number
	CardGenerator *CardGenerator // Card generator
}

func (g *GameState) GetYield() ResourceQuantity {
	totalYield := g.MyNation.BasicYield

	// Add the Yield of the Territory of the controlled WildernessPoint
	for _, point := range g.MapGrid.Points {
		if wilderness, ok := point.(*WildernessPoint); ok && wilderness.Controlled {
			totalYield = totalYield.Add(wilderness.Territory.Yield())
		}
	}

	return totalYield
}

// AddYield adds the BasicYield of the controlled Territory and MyNation to the Treasury.
func (g *GameState) AddYield() {
	g.Treasury.Add(g.GetYield())
}

// NextTurn advances the turn and adds Yield.
func (g *GameState) NextTurn() {
	g.CurrentTurn++
	g.AddYield()
}

// IsVictory determines the victory condition (whether all BossPoints have been defeated).
func (g *GameState) IsVictory() bool {
	for _, point := range g.MapGrid.Points {
		if bossPoint, ok := point.(*BossPoint); ok {
			if !bossPoint.Defeated {
				return false // A non-defeated Boss exists
			}
		}
	}
	return true // All Bosses have been defeated (or no Bosses exist)
}

// CanInteract determines whether the Point at the specified coordinates can be interacted with.
func (g *GameState) CanInteract(x, y int) bool {
	return g.MapGrid.CanInteract(x, y)
}

// GetPoint gets the Point at the specified coordinates.
func (g *GameState) GetPoint(x, y int) Point {
	return g.MapGrid.GetPoint(x, y)
}
