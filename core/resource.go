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

type ResourceModifier struct {
	Money float64
	Food  float64
	Wood  float64
	Iron  float64
	Mana  float64
}

// NewResourceModifier creates a new ResourceModifier with no modification (all values 0.0).
func NewResourceModifier() ResourceModifier {
	return ResourceModifier{
		Money: 0.0,
		Food:  0.0,
		Wood:  0.0,
		Iron:  0.0,
		Mana:  0.0,
	}
}

func (m ResourceModifier) Modify(quantity ResourceQuantity) ResourceQuantity {
	return ResourceQuantity{
		Money: m.multiply(quantity.Money, m.Money),
		Food:  m.multiply(quantity.Food, m.Food),
		Wood:  m.multiply(quantity.Wood, m.Wood),
		Iron:  m.multiply(quantity.Iron, m.Iron),
		Mana:  m.multiply(quantity.Mana, m.Mana),
	}
}

func (m ResourceModifier) multiply(value int, modifier float64) int {
	return int(float64(value) * (modifier + 1))
}
