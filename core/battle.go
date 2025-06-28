package core

// BattlefieldEffect は戦場の効果（将来的に使用予定）
type BattlefieldEffect interface {
	// 将来的に戦場効果のロジックを追加
}

// Battlefield は未制圧のWildernessで戦闘開始時に生成される戦場オブジェクトです。
type Battlefield struct {
	Enemy            *Enemy              // 戦闘相手のEnemy。
	Effects          []BattlefieldEffect // 戦場効果の配列
	BaseSupportPower float64             // 周囲のTerritoryに置いたStructureCardの影響で増加したPower。
	BattleCards      []*BattleCard       // 戦闘中に出すBattleCardの集合。
	CardSlot         int                 // BattleCardを置くことができる枚数。
}

// CanBeat 戦闘を勝利できるかどうかを返す。
func (b *Battlefield) CanBeat() bool {
	totalPower := b.CalculateTotalPower()
	return totalPower >= b.Enemy.Power
}

// Beat 戦闘を勝利する。
func (b *Battlefield) Beat() {
	// 現在は勝利処理のみ実装
	// 将来的に戦闘結果の処理、報酬の付与等を追加するかもしれない
}

// NewBattlefield Battlefieldのインスタンスを作成する。
func NewBattlefield(enemy *Enemy, supportPower float64) *Battlefield {
	return &Battlefield{
		Enemy:            enemy,
		Effects:          make([]BattlefieldEffect, 0),
		BaseSupportPower: supportPower,
		BattleCards:      make([]*BattleCard, 0, enemy.BattleCardSlot),
		CardSlot:         enemy.BattleCardSlot,
	}
}

// AddBattleCard 戦場にBattleCardを追加する。
func (b *Battlefield) AddBattleCard(card *BattleCard) bool {
	if len(b.BattleCards) >= b.CardSlot {
		return false // スロット上限に達している
	}
	b.BattleCards = append(b.BattleCards, card)
	return true
}

// RemoveBattleCard 戦場からBattleCardを除去する。
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
		power := float64(card.Power)
		power = modifiers[i].Calculate(power)
		totalPower += power
	}
	return float64(int(totalPower))
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
