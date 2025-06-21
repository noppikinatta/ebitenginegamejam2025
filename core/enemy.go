package core

// string EnemyID
// struct Enemy
// - EnemyID EnemyID
// - EnemyType EnemyType
// - Power float64
// - Skills []*EnemySkill
// - BattleCardSlot int このEnemyとの戦闘で、プレイヤーが出せるBattleCardの枚数。

// string EnemyType
// struct EnemySkill
// - EnemySkillID EnemySkillID
// - func CanAffect(battleCard *BattleCard) bool 引数battleCardの戦闘力を変更できるかどうかを返す。
// - func Modify(battleCard *BattleCard) float64 引数battleCardの戦闘力を変更する。

// interface EnemySkillCondition
