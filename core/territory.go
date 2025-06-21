package core

// Terriroryは、制圧したWildernessPointです。
// Territoryは、ターンごとにYield分のResourceを獲得します。
// Terriroryには、StructureCardを配置できます。

// string TerritoryID
// struct Territory
// - TerritoryID TerritoryID
// - Cards []*StructureCard
// - CardSlot int
// - BaseYield ResourceQuantity
// - func AppendCard(card *StructureCard) bool
// - func RemoveCard(index int) (*StructureCard, bool)
// - func Yield() ResourceQuantity
