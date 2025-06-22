package core

// 旧仕様コメントを削除し、実装に置き換え済み

// NationID は国家の一意識別子
type NationID string

// Nation 国家を表す。
type Nation struct {
	NationID NationID
	Market   *Market // カードパックを購入するためのMarket。
}

// VisibleCardPacks Marketで可視化されているカードパックの一覧を返す。
func (n *Nation) VisibleCardPacks() []*CardPack {
	return n.Market.VisibleCardPacks()
}

// CanPurchase 引数indexのカードパックを購入できるかどうかを返す。
func (n *Nation) CanPurchase(index int, treasury *Treasury) bool {
	return n.Market.CanPurchase(index, treasury)
}

// Purchase 引数indexのカードパックを購入する。国庫が不足していればfalseを返す。
func (n *Nation) Purchase(index int, treasury *Treasury) (*CardPack, bool) {
	return n.Market.Purchase(index, treasury)
}

// MyNation プレイヤー国家を表す。
type MyNation struct {
	Nation                      // 埋め込み構造体
	BasicYield ResourceQuantity // 基本Yield。
}

// AppendMarketItem 引数itemをMarket.Itemsに追加する。これは、Enemy撃破時の報酬として自国で買えるカードパックを増やす機能。
func (mn *MyNation) AppendMarketItem(item *MarketItem) {
	mn.Market.Items = append(mn.Market.Items, item)
}

// AppendLevel 引数marketLevelをMarket.Levelに加算する。これは、Enemy撃破時の報酬として自国で買えるカードパックを増やす機能。
func (mn *MyNation) AppendLevel(marketLevel MarketLevel) {
	mn.Market.Level += marketLevel
}

// OtherNation NPC(カードの取引相手)の国家を表す。
type OtherNation struct {
	Nation // 埋め込み構造体
}

// Purchase Nation.Purchaseを呼び出し、Market.Levelに0.2を加算する。
func (on *OtherNation) Purchase(index int, treasury *Treasury) (*CardPack, bool) {
	cardPack, ok := on.Nation.Purchase(index, treasury)
	if ok {
		on.Market.Level += 0.2
	}
	return cardPack, ok
}
