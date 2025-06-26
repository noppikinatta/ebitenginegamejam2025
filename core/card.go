package core

// CardID はカードの一意識別子
type CardID string

// CardPackID はカードパックの一意識別子
type CardPackID string

// Intner は乱数生成のインターフェース。テスタビリティを考慮
type Intner interface {
	Intn(n int) int
}

// EnemySkill はenemy.goで定義されます

// Battlefield はbattle.goで定義されます

// CardPack カードパックは、Marketで購入できます。
type CardPack struct {
	CardPackID CardPackID
	Ratios     map[CardID]int // カードIDとそのカードが出現する確率。合計は幾つでもいい。
	NumPerOpen int            // 一度のOpenで返すCardIDの数
}

// Open カードパックを開きます。Ratiosの値を足してIntnでそれ未満の乱数を取得し、割合に応じたカードを得る抽選を行います。抽選はNumPerOpen回行います。
func (c *CardPack) Open(intner Intner) []CardID {
	if len(c.Ratios) == 0 {
		return []CardID{}
	}

	result := make([]CardID, 0, c.NumPerOpen)
	for i := 0; i < c.NumPerOpen; i++ {
		// 合計重みを計算
		totalWeight := 0
		for _, weight := range c.Ratios {
			totalWeight += weight
		}

		// 乱数を生成
		rand := intner.Intn(totalWeight)

		// 累積確率からカードを選択
		current := 0
		for cardID, weight := range c.Ratios {
			current += weight
			if rand < current {
				result = append(result, cardID)
				break
			}
		}
	}

	return result
}

// Cards Cardsは各種類のカードの集合。これは単純なデータコンテナとして使います。
type Cards struct {
	BattleCards    []*BattleCard
	StructureCards []*StructureCard
	ResourceCards  []*ResourceCard
}

// BattleCardPower 本来intにしたいが、計算上float64の方が都合が良い。表示時には小数点以下第一位までにするかもしれない。
type BattleCardPower float64

// BattleCardType これは単純な文字列
type BattleCardType string

// BattleCard は、戦闘でBattlefieldに出すカードです。この構造体はimmutableです。
type BattleCard struct {
	CardID CardID
	Power  BattleCardPower  // はカードの戦闘力です。
	Skill  *BattleCardSkill // はカードが持っているスキルです。
	Type   BattleCardType   // は戦士、魔法使い、動物などのカードタイプです。スキルの効果対象の判定に使います。
}

// StructureCard は、Territoryに配置するカードです。この構造体はimmutableです。
type StructureCard struct {
	CardID              CardID
	YieldModifier       YieldModifier       // はTerritoryのYieldを変更するスキルです。
	BattlefieldModifier BattlefieldModifier // はBattlefieldの状態を変更するスキルです。
}

// ResourceCard は、Resourceを獲得するカードです。この構造体はimmutableです。
type ResourceCard struct {
	CardID           CardID
	ResourceQuantity ResourceQuantity // は獲得するResourceの量です。
}

// CardDatabase はカードの検索データベース
type CardDatabase struct {
	BattleCards    map[CardID]*BattleCard
	StructureCards map[CardID]*StructureCard
	ResourceCards  map[CardID]*ResourceCard
}

// GetCards 引数cardIDsに対応するカードを返す。1枚でも対応するカードがなければfalseを返す。
// 正しくデータを作っていれば、GetCardsは常にtrueを返す。
func (d *CardDatabase) GetCards(cardIDs []CardID) (*Cards, bool) {
	cards := &Cards{
		BattleCards:    make([]*BattleCard, 0),
		StructureCards: make([]*StructureCard, 0),
		ResourceCards:  make([]*ResourceCard, 0),
	}

	for _, cardID := range cardIDs {
		// BattleCardとして存在するかチェック
		if battleCard, exists := d.BattleCards[cardID]; exists {
			cards.BattleCards = append(cards.BattleCards, battleCard)
			continue
		}

		// StructureCardとして存在するかチェック
		if structureCard, exists := d.StructureCards[cardID]; exists {
			cards.StructureCards = append(cards.StructureCards, structureCard)
			continue
		}

		// ResourceCardとして存在するかチェック
		if resourceCard, exists := d.ResourceCards[cardID]; exists {
			cards.ResourceCards = append(cards.ResourceCards, resourceCard)
			continue
		}

		// どのタイプにも存在しない場合はfalseを返す
		return nil, false
	}

	return cards, true
}

// CardDeck はプレイヤーが持つカードデッキ
type CardDeck struct {
	Cards // 埋め込み構造体
	// ここにドメインの機能を追加するかもしれない。
}

// Add は引数のCardsをCardDeckに追加します
func (cd *CardDeck) Add(cards *Cards) {
	if cards == nil {
		return
	}

	cd.BattleCards = append(cd.BattleCards, cards.BattleCards...)
	cd.StructureCards = append(cd.StructureCards, cards.StructureCards...)
	cd.ResourceCards = append(cd.ResourceCards, cards.ResourceCards...)
}

// BattleCardSkillID はバトルカードスキルの識別子
type BattleCardSkillID string

// BattleCardSkill はバトルカードのスキル
type BattleCardSkill struct {
	BattleCardSkillID       BattleCardSkillID
	EnemySkillBlocker       EnemySkillBlocker
	BattleCardPowerModifier BattleCardPowerModifier
}

// EnemySkillBlocker は、Enemyのスキルをブロックするスキルです。
type EnemySkillBlocker interface {
	CanBlock(enemySkill *EnemySkill) bool // 引数enemySkillがブロックできるかどうかを返す。
}

// BattleCardPowerModifier は、BattleCardの戦闘力を変更するスキルです。
type BattleCardPowerModifier interface {
	CanAffect(battleCard *BattleCard) bool // 引数battleCardの戦闘力を変更できるかどうかを返す。
	Modify(battleCard *BattleCard) float64 // 引数battleCardの戦闘力を変更する。
}

// BattleCardState はBattleCardの戦闘中の状態を表します。
type BattleCardState struct {
	BattleCard          *BattleCard
	AffectedEnemySkills []bool // EnemySkillの対象になっているかのフラグ。EnemySkillは複数あるのでそれぞれの判定結果を格納する。
	AffectedCardSkills  []bool // BattleCardSkillの対象になっているかのフラグ。場に出ている全てのBattleCardのSkillが影響するかもしれない。このスライスの長さは場に出ているBattleCardの枚数に等しい。
}

// YieldModifier はTerritoryの収益を変更するスキル
type YieldModifier interface {
	Modify(quantity ResourceQuantity) ResourceQuantity // 引数quantityを変更する。
}

// BattlefieldModifier はBattlefieldの状態を変更するスキル
type BattlefieldModifier interface {
	Modify(battlefield *Battlefield) *Battlefield // 引数battlefieldを変更する。
}
