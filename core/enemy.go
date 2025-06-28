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
}

type EnemySkill interface {
	Calculate(options *EnemySkillCalculationOptions)
}

type EnemySkillCalculationOptions struct {
	SupportPowerMultiplier   float64
	BattleCards              []*BattleCard
	BattleCardPowerModifiers []*BattleCardPowerModifier
	Enemy                    *Enemy
}

type EnemySkillAdditiveDebuff float64

func (s EnemySkillAdditiveDebuff) Calculate(options *EnemySkillCalculationOptions) {
	for i := range options.BattleCards {
		options.BattleCardPowerModifiers[i].AdditiveDebuff += float64(s)
	}
}

type EnemySkillCardTypeAdditiveDebuff struct {
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
