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

// CardGenerator はカードを生成するための構造体です。
type CardGenerator struct {
	BattleCards    map[CardID]*BattleCard
	StructureCards map[CardID]*StructureCard
	ResourceCards  map[CardID]*ResourceCard
}

// Generate は引数のCardIDの配列に対応するカードを生成します。
// 1枚でも対応するカードがなければfalseを返します。
// 正しくデータを作っていれば、Generateは常にtrueを返します。
func (g *CardGenerator) Generate(cardIDs []CardID) (*Cards, bool) {
	cards := &Cards{
		BattleCards:    make([]*BattleCard, 0),
		StructureCards: make([]*StructureCard, 0),
		ResourceCards:  make([]*ResourceCard, 0),
	}

	for _, cardID := range cardIDs {
		// BattleCardとして存在するかチェック
		if battleCard, exists := g.BattleCards[cardID]; exists {
			newBattleCard := *battleCard
			cards.BattleCards = append(cards.BattleCards, &newBattleCard)
			continue
		}

		// StructureCardとして存在するかチェック
		if structureCard, exists := g.StructureCards[cardID]; exists {
			newStructureCard := *structureCard
			cards.StructureCards = append(cards.StructureCards, &newStructureCard)
			continue
		}

		// ResourceCardとして存在するかチェック
		if resourceCard, exists := g.ResourceCards[cardID]; exists {
			newResourceCard := *resourceCard
			cards.ResourceCards = append(cards.ResourceCards, &newResourceCard)
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
	BattleCardSkillID BattleCardSkillID
	DescriptionKey    string
	Calculator        BattleCardSkillCalculator
}

func (s *BattleCardSkill) Calculate(options *BattleCardSkillCalculationOptions) {
	s.Calculator.Calculate(options)
}

type BattleCardSkillCalculator interface {
	Calculate(options *BattleCardSkillCalculationOptions)
}

type BattleCardSkillCalculationOptions struct {
	SupportPowerMultiplier   float64
	BattleCardIndex          int
	BattleCards              []*BattleCard
	BattleCardPowerModifiers []*BattleCardPowerModifier
	Enemy                    *Enemy
}

type BattleCardSkillCalculationFunc func(options *BattleCardSkillCalculationOptions)

func (f BattleCardSkillCalculationFunc) Calculate(options *BattleCardSkillCalculationOptions) {
	f(options)
}

var NopBattleCardSkillCalculation = BattleCardSkillCalculationFunc(func(options *BattleCardSkillCalculationOptions) {})

type BattleCardSkillCalculatorComposite struct {
	Calculators []BattleCardSkillCalculator
}

func (c *BattleCardSkillCalculatorComposite) Calculate(options *BattleCardSkillCalculationOptions) {
	for _, calculator := range c.Calculators {
		calculator.Calculate(options)
	}
}

type BattleCardSkillCalculatorSupportPowerMultiplier struct {
	Multiplier float64
}

func (c *BattleCardSkillCalculatorSupportPowerMultiplier) Calculate(options *BattleCardSkillCalculationOptions) {
	options.SupportPowerMultiplier += c.Multiplier
}

type BattleCardSkillCalculatorEnemyType struct {
	EnemyType  EnemyType
	Multiplier float64
}

func (c *BattleCardSkillCalculatorEnemyType) Calculate(options *BattleCardSkillCalculationOptions) {
	if options.Enemy.EnemyType == c.EnemyType {
		options.BattleCardPowerModifiers[options.BattleCardIndex].MultiplicativeBuff += c.Multiplier
	}
}

type BattleCardSkillCalculatorBoostBuff struct {
	BoostBuff float64
}

func (c *BattleCardSkillCalculatorBoostBuff) Calculate(options *BattleCardSkillCalculationOptions) {
	modifier := options.BattleCardPowerModifiers[options.BattleCardIndex]
	modifier.MultiplicativeBuff *= c.BoostBuff
	modifier.AdditiveBuff *= c.BoostBuff
}

type BattleCardSkillCalculatorTrailings struct {
	CardType   BattleCardType
	Multiplier float64
}

func (c *BattleCardSkillCalculatorTrailings) Calculate(options *BattleCardSkillCalculationOptions) {
	for i, card := range options.BattleCards {
		if i <= options.BattleCardIndex {
			continue
		}
		if c.CardType == "" || card.Type == c.CardType {
			options.BattleCardPowerModifiers[i].MultiplicativeBuff += c.Multiplier
		}
	}
}

type BattleCardSkillCalculatorAll struct {
	ModifierFunc func(modifier *BattleCardPowerModifier)
}

func (c *BattleCardSkillCalculatorAll) Calculate(options *BattleCardSkillCalculationOptions) {
	for _, modifier := range options.BattleCardPowerModifiers {
		c.ModifierFunc(modifier)
	}
}

var AddingByIndexBattleCardSkillCalculator = BattleCardSkillCalculationFunc(func(options *BattleCardSkillCalculationOptions) {
	modifier := options.BattleCardPowerModifiers[options.BattleCardIndex]
	modifier.AdditiveBuff += float64(options.BattleCardIndex)
})

type BattleCardSkillCalculatorAllByCardType struct {
	CardType   BattleCardType
	Multiplier float64
}

func (c *BattleCardSkillCalculatorAllByCardType) Calculate(options *BattleCardSkillCalculationOptions) {
	for i, card := range options.BattleCards {
		if card.Type == c.CardType {
			options.BattleCardPowerModifiers[i].MultiplicativeBuff += c.Multiplier
		}
	}
}

type BattleCardSkillCalculatorByIdx struct {
	Index      int
	Multiplier float64
}

func (c *BattleCardSkillCalculatorByIdx) Calculate(options *BattleCardSkillCalculationOptions) {
	if options.BattleCardIndex == c.Index {
		options.BattleCardPowerModifiers[c.Index].MultiplicativeBuff += c.Multiplier
	}
}

type BattleCardSkillCalculatorProofBuff struct {
	Value float64
}

func (c *BattleCardSkillCalculatorProofBuff) Calculate(options *BattleCardSkillCalculationOptions) {
	options.BattleCardPowerModifiers[options.BattleCardIndex].ProtectionFromDebuff += c.Value
}

type BattleCardSkillCalculatorProofDebufNeighboring struct {
	Value float64
}

func (c *BattleCardSkillCalculatorProofDebufNeighboring) Calculate(options *BattleCardSkillCalculationOptions) {
	leftIdx := options.BattleCardIndex - 1
	rightIdx := options.BattleCardIndex + 1

	if leftIdx >= 0 {
		options.BattleCardPowerModifiers[leftIdx].ProtectionFromDebuff += c.Value
	}
	if rightIdx < len(options.BattleCards) {
		options.BattleCardPowerModifiers[rightIdx].ProtectionFromDebuff += c.Value
	}
}

type BattleCardSkillCalculatorTwoPlatoon struct {
	Multiplier float64
	CardType   BattleCardType
}

func (c *BattleCardSkillCalculatorTwoPlatoon) Calculate(options *BattleCardSkillCalculationOptions) {
	rightIdx := options.BattleCardIndex + 1
	if rightIdx >= len(options.BattleCards) {
		return
	}

	if options.BattleCards[rightIdx].Type == c.CardType {
		options.BattleCardPowerModifiers[options.BattleCardIndex].MultiplicativeBuff += c.Multiplier
		options.BattleCardPowerModifiers[rightIdx].MultiplicativeBuff += c.Multiplier
	}
}

// YieldModifier はTerritoryの収益を変更するスキル
type YieldModifier interface {
	Modify(quantity ResourceQuantity) ResourceQuantity // 引数quantityを変更する。
}

type MultiplyYieldModifier struct {
	Multiply float64
}

func (m *MultiplyYieldModifier) Modify(quantity ResourceQuantity) ResourceQuantity {
	return ResourceQuantity{
		Food:  int(float64(quantity.Food) * m.Multiply),
		Money: int(float64(quantity.Money) * m.Multiply),
		Wood:  int(float64(quantity.Wood) * m.Multiply),
		Iron:  int(float64(quantity.Iron) * m.Multiply),
		Mana:  int(float64(quantity.Mana) * m.Multiply),
	}
}

// BattlefieldModifier はBattlefieldの状態を変更するスキル
type BattlefieldModifier interface {
	Modify(battlefield *Battlefield) *Battlefield // 引数battlefieldを変更する。
}

type CardSlotBattlefieldModifier struct {
	Value int
}

func (m *CardSlotBattlefieldModifier) Modify(battlefield *Battlefield) *Battlefield {
	battlefield.CardSlot += m.Value
	return battlefield
}

type SupportPowerBattlefieldModifier struct {
	Value float64
}

func (m *SupportPowerBattlefieldModifier) Modify(battlefield *Battlefield) *Battlefield {
	battlefield.BaseSupportPower += m.Value
	return battlefield
}
