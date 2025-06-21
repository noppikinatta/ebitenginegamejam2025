package core

// struct Market カードパックを購入するためのMarket。
// - Level MarketLevel このMarketのレベル。
// - Items []*MarketItem カードパックの一覧。
// - func VisibleCardPacks() []*CardPack カードパックの一覧を返す。
// - func CanPurchase(index int, treasury *Treasury) bool 引数indexのカードパックを購入できるかどうかを返す。
// - func Purchase(index int, treasury *Treasury) (*CardPack, bool) 引数indexのカードパックを購入する。国庫が不足していればfalseを返す。

// struct MarketItem カードパックを表す。
// - CardPack *CardPack
// - Price ResourceQuantity カードパックの価格。
// - RequiredLevel MarketLevel カードパックを購入するために必要なMarketのレベル。
// - func CanPurchase(treasury *Treasury) bool 引数treasuryの国庫がカードパックを購入できるかどうかを返す。

// float63 MarketLevel Marketのレベル。MarketItemがプレイヤーから見える状態になる判定に使う。
