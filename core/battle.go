package core

// BattlefieldEffect は戦場の効果（将来的に使用予定）
type BattlefieldEffect interface {
	// 将来的に戦場効果のロジックを追加
}

// Battlefield は未制圧のWildernessで戦闘開始時に生成される戦場オブジェクトです。
type Battlefield struct {
	Enemy        *Enemy              // 戦闘相手のEnemy。
	Effects      []BattlefieldEffect // 戦場効果の配列
	SupportPower float64             // 周囲のTerritoryに置いたStructureCardの影響で増加したPower。
	BattleCards  []*BattleCard       // 戦闘中に出すBattleCardの集合。
	CardSlot     int                 // BattleCardを置くことができる枚数。
}

// CanBeat 戦闘を勝利できるかどうかを返す。
func (bf *Battlefield) CanBeat() bool {
	totalPower := bf.SupportPower

	// 各BattleCardのパワーを計算（敵スキルの影響を考慮）
	for _, card := range bf.BattleCards {
		cardPower := float64(card.Power)

		// 敵のスキルの影響を適用
		for _, skill := range bf.Enemy.Skills {
			if skill.CanAffect(card) {
				cardPower = skill.Modify(card)
			}
		}

		totalPower += cardPower
	}

	return totalPower >= bf.Enemy.Power
}

// Beat 戦闘を勝利する。
func (bf *Battlefield) Beat() {
	// 現在は勝利処理のみ実装
	// 将来的に戦闘結果の処理、報酬の付与等を追加するかもしれない
}

// NewBattlefield Battlefieldのインスタンスを作成する。
func NewBattlefield(enemy *Enemy, supportPower float64) *Battlefield {
	return &Battlefield{
		Enemy:        enemy,
		Effects:      make([]BattlefieldEffect, 0),
		SupportPower: supportPower,
		BattleCards:  make([]*BattleCard, 0, enemy.BattleCardSlot),
		CardSlot:     enemy.BattleCardSlot,
	}
}

// AddBattleCard 戦場にBattleCardを追加する。
func (bf *Battlefield) AddBattleCard(card *BattleCard) bool {
	if len(bf.BattleCards) >= bf.CardSlot {
		return false // スロット上限に達している
	}
	bf.BattleCards = append(bf.BattleCards, card)
	return true
}

// RemoveBattleCard 戦場からBattleCardを除去する。
func (bf *Battlefield) RemoveBattleCard(index int) (*BattleCard, bool) {
	if index < 0 || index >= len(bf.BattleCards) {
		return nil, false
	}

	card := bf.BattleCards[index]
	bf.BattleCards = append(bf.BattleCards[:index], bf.BattleCards[index+1:]...)
	return card, true
}
