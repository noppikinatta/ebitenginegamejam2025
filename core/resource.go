package core

// ResourceQuantity represents the amount of 5 types of Resources.

// ResourceQuantity is a struct that represents the amount of 5 types of resources.
type ResourceQuantity struct {
	Money int
	Food  int
	Wood  int
	Iron  int
	Mana  int
}

// Add performs a simple addition. The result can be negative.
func (rq ResourceQuantity) Add(other ResourceQuantity) ResourceQuantity {
	return ResourceQuantity{
		Money: rq.Money + other.Money,
		Food:  rq.Food + other.Food,
		Wood:  rq.Wood + other.Wood,
		Iron:  rq.Iron + other.Iron,
		Mana:  rq.Mana + other.Mana,
	}
}

// Sub performs a simple subtraction. The result can be negative.
func (rq ResourceQuantity) Sub(other ResourceQuantity) ResourceQuantity {
	return ResourceQuantity{
		Money: rq.Money - other.Money,
		Food:  rq.Food - other.Food,
		Wood:  rq.Wood - other.Wood,
		Iron:  rq.Iron - other.Iron,
		Mana:  rq.Mana - other.Mana,
	}
}

// CanPurchase returns true if the price can be afforded.
func (rq ResourceQuantity) CanPurchase(price ResourceQuantity) bool {
	return rq.Money >= price.Money &&
		rq.Food >= price.Food &&
		rq.Wood >= price.Wood &&
		rq.Iron >= price.Iron &&
		rq.Mana >= price.Mana
}
