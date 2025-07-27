package core

// EnemyID is a unique identifier for an enemy.
type EnemyID string

// EnemyType is the type of an enemy.
type EnemyType string

// EnemySkillID is a unique identifier for an enemy skill.
type EnemySkillID string

// Enemy represents an enemy.
type Enemy struct {
	id             EnemyID
	enemyType      EnemyType
	power          float64
	skills         []*EnemySkill
	battleCardSlot int // The number of BattleCards a player can play in a battle against this Enemy.
}

// NewEnemy creates a new Enemy instance.
func NewEnemy(id EnemyID, enemyType EnemyType, power float64, skills []*EnemySkill, battleCardSlot int) *Enemy {
	return &Enemy{
		id:             id,
		enemyType:      enemyType,
		power:          power,
		skills:         skills,
		battleCardSlot: battleCardSlot,
	}
}

// ID returns the enemy ID.
func (e *Enemy) ID() EnemyID {
	return e.id
}

// Type returns the enemy type.
func (e *Enemy) Type() EnemyType {
	return e.enemyType
}

// Power returns the enemy power.
func (e *Enemy) Power() float64 {
	return e.power
}

// Skills returns a copy of the enemy skills.
func (e *Enemy) Skills() []*EnemySkill {
	result := make([]*EnemySkill, len(e.skills))
	copy(result, e.skills)
	return result
}

// BattleCardSlot returns the number of BattleCards a player can play in a battle against this Enemy.
func (e *Enemy) BattleCardSlot() int {
	return e.battleCardSlot
}

// EnemySkill represents an enemy skill.
type EnemySkill struct {
	id        EnemySkillID
	condition func(idx int, options *EnemySkillCalculationOptions) bool
	modifier  *BattleCardPowerModifier
}

// NewEnemySkill creates a new EnemySkill instance.
func NewEnemySkill(id EnemySkillID, condition func(idx int, options *EnemySkillCalculationOptions) bool, modifier *BattleCardPowerModifier) *EnemySkill {
	return &EnemySkill{
		id:        id,
		condition: condition,
		modifier:  modifier,
	}
}

// ID returns the enemy skill ID.
func (s *EnemySkill) ID() EnemySkillID {
	return s.id
}

// Calculate applies the enemy skill effects.
func (s *EnemySkill) Calculate(options *EnemySkillCalculationOptions) {
	for i := range options.BattleCards {
		if !s.condition(i, options) {
			continue
		}
		modifier := options.BattleCardPowerModifiers[i]
		modifier.Union(s.modifier)
	}
}

type EnemySkillCalculationOptions struct {
	BattleCards              []*BattleCard
	BattleCardPowerModifiers []*BattleCardPowerModifier
	Enemy                    *Enemy
}
