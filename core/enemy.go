package core

// EnemyID is a unique identifier for an enemy.
type EnemyID string

// EnemyType is the type of an enemy.
type EnemyType string

// EnemySkillID is a unique identifier for an enemy skill.
type EnemySkillID string

// Enemy represents an enemy.
type Enemy struct {
	EnemyID        EnemyID
	EnemyType      EnemyType
	Power          float64
	Skills         []EnemySkill
	BattleCardSlot int // The number of BattleCards a player can play in a battle against this Enemy.
	Question       string
}

type EnemySkill interface {
	ID() EnemySkillID
	Calculate(options *EnemySkillCalculationOptions)
}

type EnemySkillImpl struct {
	IDField   EnemySkillID
	Condition func(idx int, options *EnemySkillCalculationOptions) bool
	Modifier  *BattleCardPowerModifier
}

func (s *EnemySkillImpl) ID() EnemySkillID {
	return s.IDField
}

func (s *EnemySkillImpl) Calculate(options *EnemySkillCalculationOptions) {
	for i := range options.BattleCards {
		if !s.Condition(i, options) {
			continue
		}
		modifier := options.BattleCardPowerModifiers[i]
		modifier.Union(s.Modifier)
	}
}

type EnemySkillCalculationOptions struct {
	BattleCards              []*BattleCard
	BattleCardPowerModifiers []*BattleCardPowerModifier
	Enemy                    *Enemy
}
