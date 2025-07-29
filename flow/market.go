package flow

import (
	"math/rand"
	"time"

	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// MarketFlow handles market purchase operations
type MarketFlow struct {
	gameState *core.GameState
	market    *core.Market
	nation    core.Nation
}

// NewMarketFlow creates a new MarketFlow
func NewMarketFlow(gameState *core.GameState, market *core.Market, nation core.Nation) *MarketFlow {
	return &MarketFlow{
		gameState: gameState,
		market:    market,
		nation:    nation,
	}
}

// Purchase attempts to purchase a market item at the specified index
func (mf *MarketFlow) Purchase(marketItemIdx int) bool {
	if mf.market == nil || mf.gameState == nil {
		return false
	}
	
	if marketItemIdx < 0 || marketItemIdx >= len(mf.market.Items) {
		return false
	}
	
	item := mf.market.Items[marketItemIdx]
	
	// Check if item is available
	if !mf.isMarketItemAvailable(item) {
		return false
	}
	
	// Attempt purchase
	oldLevel := mf.market.Level
	cardPack, ok := mf.market.Purchase(marketItemIdx, mf.gameState.Treasury)
	if !ok {
		return false
	}
	
	// Add history if market level increased
	if int(mf.market.Level) > int(oldLevel) {
		mf.addMarketLevelHistory()
	}
	
	// Process card pack if purchased
	if cardPack != nil {
		mf.processCardPack(cardPack)
	}
	
	// Advance turn
	mf.gameState.NextTurn()
	
	return true
}

// CanPurchase checks if a market item can be purchased
func (mf *MarketFlow) CanPurchase(marketItemIdx int) bool {
	if mf.market == nil || mf.gameState == nil {
		return false
	}
	
	if marketItemIdx < 0 || marketItemIdx >= len(mf.market.Items) {
		return false
	}
	
	item := mf.market.Items[marketItemIdx]
	
	// Check availability and affordability
	return mf.isMarketItemAvailable(item) && item.CanPurchase(mf.gameState.Treasury)
}

// isMarketItemAvailable checks if a market item is available for purchase
func (mf *MarketFlow) isMarketItemAvailable(item *core.MarketItem) bool {
	if mf.market == nil {
		return false
	}
	// Check if market level meets requirement
	return mf.market.Level >= item.RequiredLevel()
}

// addMarketLevelHistory adds a history entry for market level increase
func (mf *MarketFlow) addMarketLevelHistory() {
	if mf.gameState == nil || mf.nation == nil {
		return
	}
	
	mf.gameState.AddHistory(core.History{
		Turn: mf.gameState.CurrentTurn,
		Key:  "history-market",
		Data: map[string]any{
			"nation": string(mf.nation.ID()),
			"level":  int(mf.market.Level),
		},
	})
}

// processCardPack opens a card pack and adds cards to the deck
func (mf *MarketFlow) processCardPack(cardPack *core.CardPack) {
	if cardPack == nil || mf.gameState == nil {
		return
	}
	
	// Create random number generator
	rng := &simpleRand{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	
	// Open card pack
	cardIDs := cardPack.Open(rng)
	
	// Add cards to deck
	for _, cardID := range cardIDs {
		mf.gameState.CardDeck.Add(cardID)
	}
}

// simpleRand implements core.RNG interface
type simpleRand struct {
	*rand.Rand
}

func (sr *simpleRand) Intn(n int) int {
	return sr.Rand.Intn(n)
} 