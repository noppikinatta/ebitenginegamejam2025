package core

// string CardID
// string CardPackID
// struct CardPack
// - CardPackID CardPackID
// - Ratios map[CardID]int
// - NumPerOpen int
// - func Open(intner Intner) []CardID
// interface Intner : rand.Intnのインタフェース。テスタビリティを考えた物

// struct Cards
// - BattleCards []*BattleCard
// - StructureCards []*StructureCard
// - ResourceCards []*ResourceCard

// struct BattleCard
// - CardID CardID
// - Power BattleCardPower
// - Skills []*BattleCardSkill
// - Type BattleCardType

// struct StructureCard
// - CardID CardID
// - YieldModifier YieldModifier
// - BattlefieldModifier BattlefieldModifier

// struct ResourceCard
// - CardID CardID
// struct CardDatabase

// struct CardDeck
// Cards (embedded struct)

// float64 BattleCardPower
// string BattleCardSkillID
// struct BattleCardSkill
// struct EnemySkillBlocker
// struct BattleCardPowerModifier
// string BattleCardType
// struct BattleCardState

// interface YieldModifier
// interface BattlefieldModifier
