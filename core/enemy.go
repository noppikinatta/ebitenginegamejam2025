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

type BaseEnemySkill struct {
	IDField EnemySkillID
}

func (s *BaseEnemySkill) ID() EnemySkillID {
	return s.IDField
}

type EnemySkillCalculationOptions struct {
	BaseEnemySkill
	SupportPowerMultiplier   float64
	BattleCards              []*BattleCard
	BattleCardPowerModifiers []*BattleCardPowerModifier
	Enemy                    *Enemy
}

type EnemySkillAdditiveDebuff struct {
	BaseEnemySkill
	Value float64
}

func (s *EnemySkillAdditiveDebuff) Calculate(options *EnemySkillCalculationOptions) {
	for i := range options.BattleCards {
		options.BattleCardPowerModifiers[i].AdditiveDebuff += s.Value
	}
}

type EnemySkillCardTypeAdditiveDebuff struct {
	BaseEnemySkill
	CardType BattleCardType
	Value    float64
}

func (s *EnemySkillCardTypeAdditiveDebuff) Calculate(options *EnemySkillCalculationOptions) {
	for i, card := range options.BattleCards {
		if card.Type == s.CardType {
			options.BattleCardPowerModifiers[i].AdditiveDebuff += s.Value
		}
	}
}

type EnemySkillCardTypeMultiplicativeDebuff struct {
	BaseEnemySkill
	CardType BattleCardType
	Value    float64
}

func (s *EnemySkillCardTypeMultiplicativeDebuff) Calculate(options *EnemySkillCalculationOptions) {
	for i, card := range options.BattleCards {
		if card.Type == s.CardType {
			options.BattleCardPowerModifiers[i].MultiplicativeDebuff += s.Value
		}
	}
}

type EnemySkillCardTypeExceptMultiplicativeDebuff struct {
	BaseEnemySkill
	CardType BattleCardType
	Value    float64
}

func (s *EnemySkillCardTypeExceptMultiplicativeDebuff) Calculate(options *EnemySkillCalculationOptions) {
	for i, card := range options.BattleCards {
		if card.Type != s.CardType {
			options.BattleCardPowerModifiers[i].MultiplicativeDebuff += s.Value
		}
	}
}

type EnemySkillIndexForwardMultiplicativeDebuff struct {
	BaseEnemySkill
	NumOfCards int
	Value      float64
}

func (s *EnemySkillIndexForwardMultiplicativeDebuff) Calculate(options *EnemySkillCalculationOptions) {
	for i := range options.BattleCards {
		if i < s.NumOfCards {
			options.BattleCardPowerModifiers[i].MultiplicativeDebuff += s.Value
		}
	}
}

type EnemySkillIndexBackwardMultiplicativeDebuff struct {
	BaseEnemySkill
	NumOfCards int
	Value      float64
}

func (s *EnemySkillIndexBackwardMultiplicativeDebuff) Calculate(options *EnemySkillCalculationOptions) {
	for i := range options.BattleCards {
		if i >= len(options.BattleCards)-s.NumOfCards {
			options.BattleCardPowerModifiers[i].MultiplicativeDebuff += s.Value
		}
	}
}
