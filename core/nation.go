package core

// string NationID

// struct Nation 国家を表す。
// - NationID NationID
// - Market *Market カードパックを購入するためのMarket。
// - func VisibleCardPacks() []*CardPack Marketで可視化されているカードパックの一覧を返す。
// - func CanPurchase(index int, treasury *Treasury) bool 引数indexのカードパックを購入できるかどうかを返す。
// - func Purchase(index int, treasury *Treasury) (*CardPack, bool) 引数indexのカードパックを購入する。国庫が不足していればfalseを返す。

// struct MyNation プレイヤー国家を表す。
// - Nation (embedded struct)
// - func AppendMarketItem(item *MarketItem) 引数itemをMarket.Itemsに追加する。これは、Enemy撃破時の報酬として自国で買えるカードパックを増やす機能。
// - func AppendLevel(marketLevel float64) 引数marketLevelをMarket.Levelに加算する。これは、Enemy撃破時の報酬として自国で買えるカードパックを増やす機能。

// struct OtherNation NPC(カードの取引相手)の国家を表す。
// - Nation (embedded struct)
// - func Purchase(index int, treasury *Treasury) (*CardPack, bool) Nation.Purchaseを呼び出し、Market.Levelに0.2を加算する。

// struct Treasury 国庫を表す。
// - Resoruces ResourceQuantity 国庫のResourceの量。
// - func Add(other ResourceQuantity) 引数otherを国庫に加える。
// - func Sub(other ResourceQuantity) bool 引数otherを国庫から引き出す。国庫が不足していれば引き算を実行せず、falseを返す。
