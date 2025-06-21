package core

// string NationID

// struct Nation
// - NationID NationID
// - Market *Market
// - func CardPacks() []*CardPack
// - func CanPurchase(index int, treasury *Treasury) bool
// - func Purchase(index int, treasury *Treasury) (*CardPack, bool)

// struct MyNation
// - Nation (embedded struct)

// struct OtherNation
// - Nation (embedded struct)

// struct Treasury
// - Resoruces ResourceQuantity
// - func Add(other ResourceQuantity)
// - func Sub(other ResourceQuantity) bool
