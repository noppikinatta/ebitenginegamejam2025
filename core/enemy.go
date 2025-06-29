package core

// EnemyID は敵の一意識別子
type EnemyID string

// EnemyType は敵のタイプ
type EnemyType string

// EnemySkillID は敵スキルの一意識別子
type EnemySkillID string

// Enemy は敵を表す構造体
type Enemy struct {
	EnemyID        EnemyID
	EnemyType      EnemyType
	Power          float64
	Skills         []EnemySkill
	BattleCardSlot int // このEnemyとの戦闘で、プレイヤーが出せるBattleCardの枚数。
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
