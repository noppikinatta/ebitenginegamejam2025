package core

// ResourceQuantityは、5種類のResourceの量を表します。

// struct ResourceQuantity
// - Money int
// - Food int
// - Wood int
// - Iron int
// - Mana int
// - func Add(other ResourceQuantity) ResourceQuantity 単純な足し算。結果がマイナスになってもいい。
// - func Sub(other ResourceQuantity) ResourceQuantity 単純な引き算。結果がマイナスになってもいい。
// - func CanPurchase(price ResourceQuantity) bool 引数priceを充足していればtrueを返す
