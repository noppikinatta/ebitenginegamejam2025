package core

// MarketLevel はMarketのレベル。MarketItemがプレイヤーから見える状態になる判定に使う。
type MarketLevel float64

// Market カードパックを購入するためのMarket。
type Market struct {
	Level MarketLevel     // このMarketのレベル。
	Items []*MarketItem   // カードパックの一覧。
}

// VisibleCardPacks カードパックの一覧を返す。
func (m *Market) VisibleCardPacks() []*CardPack {
	var visiblePacks []*CardPack
	
	for _, item := range m.Items {
		if m.Level >= item.RequiredLevel {
			visiblePacks = append(visiblePacks, item.CardPack)
		}
	}
	
	return visiblePacks
}

// CanPurchase 引数indexのカードパックを購入できるかどうかを返す。
func (m *Market) CanPurchase(index int, treasury *Treasury) bool {
	if index < 0 || index >= len(m.Items) {
		return false
	}
	
	item := m.Items[index]
	
	// マーケットレベルチェック
	if m.Level < item.RequiredLevel {
		return false
	}
	
	// アイテム自体の購入可能性チェック
	return item.CanPurchase(treasury)
}

// Purchase 引数indexのカードパックを購入する。国庫が不足していればfalseを返す。
func (m *Market) Purchase(index int, treasury *Treasury) (*CardPack, bool) {
	if !m.CanPurchase(index, treasury) {
		return nil, false
	}
	
	item := m.Items[index]
	
	// 国庫から価格を引く
	if !treasury.Sub(item.Price) {
		return nil, false
	}
	
	return item.CardPack, true
}

// MarketItem カードパックを表す。
type MarketItem struct {
	CardPack      *CardPack
	Price         ResourceQuantity // カードパックの価格。
	RequiredLevel MarketLevel      // カードパックを購入するために必要なMarketのレベル。
}

// CanPurchase 引数treasuryの国庫がカードパックを購入できるかどうかを返す。
func (mi *MarketItem) CanPurchase(treasury *Treasury) bool {
	return treasury.Resources.CanPurchase(mi.Price)
}

// Treasury 国庫を表す。
type Treasury struct {
	Resources ResourceQuantity // 国庫のResourceの量。
}

// Add 引数otherを国庫に加える。
func (t *Treasury) Add(other ResourceQuantity) {
	t.Resources = t.Resources.Add(other)
}

// Sub 引数otherを国庫から引き出す。国庫が不足していれば引き算を実行せず、falseを返す。
func (t *Treasury) Sub(other ResourceQuantity) bool {
	if !t.Resources.CanPurchase(other) {
		return false
	}
	t.Resources = t.Resources.Sub(other)
	return true
}
