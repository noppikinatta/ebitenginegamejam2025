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
	totalPower := b.BaseSupportPower

	// 各BattleCardのパワーを計算（敵スキルの影響を考慮）
	for _, card := range b.BattleCards {
		cardPower := float64(card.Power)

		// 敵のスキルの影響を適用
		for _, skill := range b.Enemy.Skills {
			if skill.CanAffect(card) {
				cardPower = skill.Modify(card)
			}
		}

		totalPower += cardPower
	}

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
	totalPower := b.BaseSupportPower
	for _, card := range b.BattleCards {
		totalPower += float64(card.Power)
	}
	return totalPower
}

type BattleCardPowerModifier struct {
	MultiplicativeBuff   float64 // 乗算バフ（1.0以上、1.2なら20%増加）
	MultiplicativeDebuff float64 // 乗算デバフ（1.0未満、0.8なら20%減少、ProtectedFromDebuffがtrueなら無効）
	AdditiveBuff         float64 // 加算バフ（常に適用される正の値）
	AdditiveDebuff       float64 // 加算デバフ（ProtectedFromDebuffがtrueなら無効、正の値として格納し負として扱う）
	ProtectionFromDebuff float64 // デバフ耐性 (0.0~1.0 0.0が基準、1.0でデバフ無効)
}
