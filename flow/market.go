package flow

import (
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// MarketFlow handles market purchase operations
type MarketFlow struct {
	gameState *core.GameState
	market    *core.Market
	nation    core.Nation
	intner    core.Intner
}

// NewMarketFlow creates a new MarketFlow
func NewMarketFlow(gameState *core.GameState, intner core.Intner) *MarketFlow {
	return &MarketFlow{
		gameState: gameState,
		intner:    intner,
	}
}

func (mf *MarketFlow) SelectMarket(x, y int) {
	point, ok := mf.gameState.MapGrid.GetPoint(x, y)
	if !ok {
		return
	}
	marketPoint, ok := point.AsMarketPoint()
	if !ok {
		return
	}

	market, ok := mf.gameState.Markets[marketPoint.Nation().ID()]
	if !ok {
		return
	}

	mf.nation = marketPoint.Nation()
	mf.market = market
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

	// Open card pack
	cardIDs := cardPack.Open(mf.intner)

	// Add cards to deck
	for _, cardID := range cardIDs {
		mf.gameState.CardDeck.Add(cardID)
	}
}
