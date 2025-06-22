package core

// ResourceQuantityは、5種類のResourceの量を表します。

// ResourceQuantity は5種類のリソースの量を表す構造体です
type ResourceQuantity struct {
	Money int
	Food  int
	Wood  int
	Iron  int
	Mana  int
}

// Add は単純な足し算を行います。結果がマイナスになってもよい。
func (rq ResourceQuantity) Add(other ResourceQuantity) ResourceQuantity {
	return ResourceQuantity{
		Money: rq.Money + other.Money,
		Food:  rq.Food + other.Food,
		Wood:  rq.Wood + other.Wood,
		Iron:  rq.Iron + other.Iron,
		Mana:  rq.Mana + other.Mana,
	}
}

// Sub は単純な引き算を行います。結果がマイナスになってもよい。
func (rq ResourceQuantity) Sub(other ResourceQuantity) ResourceQuantity {
	return ResourceQuantity{
		Money: rq.Money - other.Money,
		Food:  rq.Food - other.Food,
		Wood:  rq.Wood - other.Wood,
		Iron:  rq.Iron - other.Iron,
		Mana:  rq.Mana - other.Mana,
	}
}

// CanPurchase は引数priceを充足していればtrueを返します
func (rq ResourceQuantity) CanPurchase(price ResourceQuantity) bool {
	return rq.Money >= price.Money &&
		rq.Food >= price.Food &&
		rq.Wood >= price.Wood &&
		rq.Iron >= price.Iron &&
		rq.Mana >= price.Mana
}
