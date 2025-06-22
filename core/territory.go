package core

// TerritoryID は領土の一意識別子
type TerritoryID string

// Territory は、制圧したWildernessPointです。
// Territory は、ターンごとにYield分のResourceを獲得します。
// Territory には、StructureCardを配置できます。
type Territory struct {
	TerritoryID TerritoryID
	Cards       []*StructureCard
	CardSlot    int
	BaseYield   ResourceQuantity
}

// AppendCard StructureCardを領土に配置します
func (t *Territory) AppendCard(card *StructureCard) bool {
	if len(t.Cards) >= t.CardSlot {
		return false // スロット上限に達している
	}
	t.Cards = append(t.Cards, card)
	return true
}

// RemoveCard 指定インデックスのStructureCardを領土から除去します
func (t *Territory) RemoveCard(index int) (*StructureCard, bool) {
	if index < 0 || index >= len(t.Cards) {
		return nil, false
	}

	card := t.Cards[index]
	t.Cards = append(t.Cards[:index], t.Cards[index+1:]...)
	return card, true
}

// Yield BaseYieldを置かれているStructureCardのYieldModifierに通して結果を返す。
func (t *Territory) Yield() ResourceQuantity {
	yield := t.BaseYield

	// 配置されているStructureCardのYieldModifierを順次適用
	for _, card := range t.Cards {
		if card.YieldModifier != nil {
			yield = card.YieldModifier.Modify(yield)
		}
	}

	return yield
}
