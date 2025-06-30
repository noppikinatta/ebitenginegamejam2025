package core

// BattlefieldEffect represents an effect on the battlefield.
type BattlefieldEffect interface {
	// TODO: Add logic for battlefield effects.
}

// Battlefield represents a battle instance, created when starting a battle in an unconquered Wilderness.
type Battlefield struct {
	Enemy            *Enemy              // Enemy is the opponent in the battle.
	Effects          []BattlefieldEffect // Effects is a slice of battlefield effects.
	BaseSupportPower float64             // BaseSupportPower is the power gained from StructureCards in adjacent Territories.
	BattleCards      []*BattleCard       // BattleCards is a collection of BattleCards played during the battle.
	CardSlot         int                 // CardSlot is the maximum number of BattleCards that can be placed.
}

// CanBeat returns true if the player's power is enough to defeat the enemy.
func (b *Battlefield) CanBeat() bool {
	totalPower := b.CalculateTotalPower()
	return totalPower >= b.Enemy.Power
}

// Beat handles the logic for winning a battle.
func (b *Battlefield) Beat() {
	// Currently, only the victory process is implemented.
	// In the future, this might include processing battle results, granting rewards, etc.
}

// NewBattlefield creates a new Battlefield instance.
func NewBattlefield(enemy *Enemy, supportPower float64) *Battlefield {
	return &Battlefield{
		Enemy:            enemy,
		Effects:          make([]BattlefieldEffect, 0),
		BaseSupportPower: supportPower,
		BattleCards:      make([]*BattleCard, 0, enemy.BattleCardSlot),
		CardSlot:         enemy.BattleCardSlot,
	}
}

// AddBattleCard adds a BattleCard to the battlefield.
func (b *Battlefield) AddBattleCard(card *BattleCard) bool {
	if len(b.BattleCards) >= b.CardSlot {
		return false // Slot limit reached
	}
	b.BattleCards = append(b.BattleCards, card)
	return true
}

// RemoveBattleCard removes a BattleCard from the battlefield.
func (b *Battlefield) RemoveBattleCard(index int) (*BattleCard, bool) {
	if index < 0 || index >= len(b.BattleCards) {
		return nil, false
	}

	card := b.BattleCards[index]
	b.BattleCards = append(b.BattleCards[:index], b.BattleCards[index+1:]...)
	return card, true
}

func (b *Battlefield) CalculateTotalPower() float64 {
	modifiers := make([]*BattleCardPowerModifier, len(b.BattleCards))
	for i := range modifiers {
		modifiers[i] = &BattleCardPowerModifier{}
	}

	cardCalcOptions := &BattleCardSkillCalculationOptions{
		BattleCards:              b.BattleCards,
		BattleCardPowerModifiers: modifiers,
		Enemy:                    b.Enemy,
	}

	for i, card := range b.BattleCards {
		cardCalcOptions.BattleCardIndex = i
		if card.Skill != nil {
			card.Skill.Calculate(cardCalcOptions)
		}
	}

	enemyCalcOptions := &EnemySkillCalculationOptions{
		BattleCards:              b.BattleCards,
		BattleCardPowerModifiers: modifiers,
		Enemy:                    b.Enemy,
	}

	for _, skill := range b.Enemy.Skills {
		skill.Calculate(enemyCalcOptions)
	}

	totalPower := b.BaseSupportPower * (cardCalcOptions.SupportPowerMultiplier + 1.0)
	for i, card := range b.BattleCards {
		power := float64(card.Power())
		power = modifiers[i].Calculate(power)
		totalPower += power
	}
	return totalPower
}

type BattleCardPowerModifier struct {
	MultiplicativeBuff   float64
	MultiplicativeDebuff float64
	BuffBoostedPower     float64
	AdditiveBuff         float64
	AdditiveDebuff       float64
	ProtectionFromDebuff float64
}

func (m *BattleCardPowerModifier) Calculate(power float64) float64 {
	power += m.additiveBuffValue()
	power *= m.multiplicativeBuffValue()
	power *= m.multiplicativeDebuffValue()
	power += m.additiveDebuffValue()
	if power < 0.0 {
		return 0.0
	}
	return power
}

func (m *BattleCardPowerModifier) buffBoostedPowerValue() float64 {
	return m.BuffBoostedPower + 1.0
}

func (m *BattleCardPowerModifier) additiveBuffValue() float64 {
	return m.AdditiveBuff * m.buffBoostedPowerValue()
}

func (m *BattleCardPowerModifier) multiplicativeBuffValue() float64 {
	return m.MultiplicativeBuff*m.buffBoostedPowerValue() + 1.0
}

func (m *BattleCardPowerModifier) protectionFromDebuffValue() float64 {
	v := 1.0 - m.ProtectionFromDebuff
	if v < 0.0 {
		return 0.0
	}
	return v
}

func (m *BattleCardPowerModifier) multiplicativeDebuffValue() float64 {
	v := 1.0 - m.MultiplicativeDebuff*m.protectionFromDebuffValue()
	if v < 0.0 {
		return 0.0
	}
	return v
}

func (m *BattleCardPowerModifier) additiveDebuffValue() float64 {
	return -m.AdditiveDebuff * m.protectionFromDebuffValue()
}
