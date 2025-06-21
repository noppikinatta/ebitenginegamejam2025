package core

// string CardID
// string CardPackID

// struct CardPack カードパックは、Marketで購入できます。
// - CardPackID CardPackID
// - Ratios map[CardID]int カードIDとそのカードが出現する確率。合計は幾つでもいい。
// - NumPerOpen int 一度のOpenで返すCardIDの数
// - func Open(intner Intner) []CardID カードパックを開きます。Ratiosの値を足してIntnでそれ未満の乱数を取得し、割合に応じたカードを得る抽選を行います。抽選はNumPerOpen回行います。
// interface Intner : rand.Intnのインタフェース。テスタビリティを考えた物

// struct Cards Cardsは各種類のカードの集合。これは単純なデータコンテナとして使います。
// - BattleCards []*BattleCard
// - StructureCards []*StructureCard
// - ResourceCards []*ResourceCard

// struct BattleCard は、戦闘でBattlefieldに出すカードです。この構造体はimmutableです。
// - CardID CardID
// - Power BattleCardPower はカードの戦闘力です。
// - Skill BattleCardSkill はカードが持っているスキルです。
// - Type BattleCardType は戦士、魔法使い、動物などのカードタイプです。スキルの効果対象の判定に使います。

// struct StructureCard は、Territoryに配置するカードです。この構造体はimmutableです。
// - CardID CardID
// - YieldModifier YieldModifier はTerritoryのYieldを変更するスキルです。
// - BattlefieldModifier BattlefieldModifier はBattlefieldの状態を変更するスキルです。

// struct ResourceCard は、Resourceを獲得するカードです。この構造体はimmutableです。
// - CardID CardID
// - ResourceQuantity ResourceQuantity は獲得するResourceの量です。

// struct CardDatabase
// - BattleCards map[CardID]*BattleCard
// - StructureCards map[CardID]*StructureCard
// - ResourceCards map[CardID]*ResourceCard
// - func GetCards(cardIDs []CardID) (*Cards, bool) 引数cardIDsに対応するカードを返す。1枚でも対応するカードがなければfalseを返す。
//   正しくデータを作っていれば、GetCardsは常にtrueを返す。

// struct CardDeck
// Cards (embedded struct)
// ここにドメインの機能を追加するかもしれない。

// float64 BattleCardPower 本来intにしたいが、計算上float64の方が都合が良い。表示時には小数点以下第一位までにするかもしれない。

// string BattleCardSkillID
// struct BattleCardSkill
// - BattleCardSkillID BattleCardSkillID
// - EnemySkillBlocker EnemySkillBlocker
// - BattleCardPowerModifier BattleCardPowerModifier

// interface EnemySkillBlocker は、Enemyのスキルをブロックするスキルです。
// - CanBlock(enemySkill EnemySkill) bool 引数enemySkillがブロックできるかどうかを返す。
//   引数に情報を足す必要があるかもしれない。

// struct BattleCardPowerModifier は、BattleCardの戦闘力を変更するスキルです。
// - func CanAffect(battleCard *BattleCard) bool 引数battleCardの戦闘力を変更できるかどうかを返す。
// - func Modify(battleCard *BattleCard) float64 引数battleCardの戦闘力を変更する。

// string BattleCardType これは単純な文字列

// struct BattleCardState はBattleCardの戦闘中の状態を表します。
// - BattleCard *BattleCard
// - AffectedEnemySkills []bool EnemySkillの対象になっているかのフラグ。EnemySkillは複数あるのでそれぞれの判定結果を格納する。
// - AffectedCardSkills []bool BattleCardSkillの対象になっているかのフラグ。場に出ている全てのBattleCardのSkillが影響するかもしれない。このスライスの長さは場に出ているBattleCardの枚数に等しい。

// interface YieldModifier
// - func Modify(quantity ResourceQuantity) ResourceQuantity 引数quantityを変更する。

// interface BattlefieldModifier
// - func Modify(battlefield *Battlefield) *Battlefield 引数battlefieldを変更する。
