package core

// EnemyID は敵の一意識別子
type EnemyID string

// EnemyType は敵のタイプ
type EnemyType string

// EnemySkillID は敵スキルの一意識別子
type EnemySkillID string

// Enemy は敵を表す構造体
type Enemy struct {
	EnemyID       EnemyID
	EnemyType     EnemyType
	Power         float64
	Skills        []EnemySkill
	BattleCardSlot int // このEnemyとの戦闘で、プレイヤーが出せるBattleCardの枚数。
}

// EnemySkill は敵のスキルを表すインターフェース
type EnemySkill interface {
	CanAffect(battleCard *BattleCard) bool // 引数battleCardの戦闘力を変更できるかどうかを返す。
	Modify(battleCard *BattleCard) float64 // 引数battleCardの戦闘力を変更する。
}

// EnemySkillCondition は敵スキルの発動条件（将来的に使用予定）
type EnemySkillCondition interface {
	// 将来的にスキル発動条件のロジックを追加
}
