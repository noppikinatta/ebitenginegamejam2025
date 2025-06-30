package core

// CardID is a unique identifier for a card.
type CardID string

// CardPackID is a unique identifier for a card pack.
type CardPackID string

// Intner is an interface for random number generation, for testability.
type Intner interface {
	Intn(n int) int
}

// EnemySkill is defined in enemy.go

// Battlefield is defined in battle.go

// CardPack can be purchased at the Market.
type CardPack struct {
	CardPackID CardPackID
	Ratios     map[CardID]int // Card ID and the probability of that card appearing. The sum can be anything.
	NumPerOpen int            // The number of CardIDs returned in a single Open.
}

// Open opens a card pack. It sums the values in Ratios, gets a random number less than the sum using Intn,
// and draws cards according to the ratio. The draw is performed NumPerOpen times.
func (c *CardPack) Open(intner Intner) []CardID {
	if len(c.Ratios) == 0 {
		return []CardID{}
	}

	result := make([]CardID, 0, c.NumPerOpen)
	for i := 0; i < c.NumPerOpen; i++ {
		// Calculate total weight
		totalWeight := 0
		for _, weight := range c.Ratios {
			totalWeight += weight
		}

		// Generate random number
		rand := intner.Intn(totalWeight)

		// Select a card from the cumulative probability
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

// Cards is a collection of cards of each type. This is used as a simple data container.
type Cards struct {
	BattleCards    []*BattleCard
	StructureCards []*StructureCard
}

// BattleCardPower should ideally be an int, but float64 is more convenient for calculations. It may be displayed to one decimal place.
type BattleCardPower float64

// BattleCardType is a simple string.
type BattleCardType string

// BattleCard is a card played on the Battlefield during combat. This struct is immutable.
type BattleCard struct {
	CardID     CardID
	Experience int
	BasePower  BattleCardPower  // BasePower is the combat power of the card.
	Skill      *BattleCardSkill // Skill is the skill the card possesses.
	Type       BattleCardType   // Type is the card type, such as warrior, mage, or animal. Used to determine the target of a skill's effect.
}

func (c *BattleCard) Level() int {
	return 1 + c.Experience/100
}

func (c *BattleCard) Experiment() {
	c.Experience += (100 / c.Level())
}

func (c *BattleCard) Power() BattleCardPower {
	return c.BasePower * (1 + 0.1*BattleCardPower(c.Level()-1))
}

// StructureCard is a card placed in a Territory. This struct is immutable.
type StructureCard struct {
	CardID              CardID
	DescriptionKey      string
	YieldModifier       YieldModifier       // YieldModifier is a skill that modifies the Yield of a Territory.
	BattlefieldModifier BattlefieldModifier // BattlefieldModifier is a skill that modifies the state of the Battlefield.
}

// CardGenerator is a struct for generating cards.
type CardGenerator struct {
	BattleCards    map[CardID]*BattleCard
	StructureCards map[CardID]*StructureCard
}

// Generate generates cards corresponding to the array of CardIDs given as an argument.
// Returns false if even one corresponding card does not exist.
// If the data is created correctly, Generate will always return true.
func (g *CardGenerator) Generate(cardIDs []CardID) (*Cards, bool) {
	cards := &Cards{
		BattleCards:    make([]*BattleCard, 0),
		StructureCards: make([]*StructureCard, 0),
	}

	for _, cardID := range cardIDs {
		// Check if it exists as a BattleCard
		if battleCard, exists := g.BattleCards[cardID]; exists {
			newBattleCard := *battleCard
			cards.BattleCards = append(cards.BattleCards, &newBattleCard)
			continue
		}

		// Check if it exists as a StructureCard
		if structureCard, exists := g.StructureCards[cardID]; exists {
			newStructureCard := *structureCard
			cards.StructureCards = append(cards.StructureCards, &newStructureCard)
			continue
		}

		// Return false if it does not exist in any type
		return nil, false
	}

	return cards, true
}

// CardDeck is the player's card deck.
type CardDeck struct {
	Cards // Embedded struct
}

// Add adds the given Cards to the CardDeck.
func (cd *CardDeck) Add(cards *Cards) {
	if cards == nil {
		return
	}

	for _, card := range cards.BattleCards {
		found := false
		for _, cardInDeck := range cd.BattleCards {
			if cardInDeck.CardID == card.CardID {
				cardInDeck.Experiment()
				found = true
				break
			}
		}
		if !found {
			cd.BattleCards = append(cd.BattleCards, card)
		}
	}

	cd.StructureCards = append(cd.StructureCards, cards.StructureCards...)
}

// BattleCardSkillID is the identifier for a battle card skill.
type BattleCardSkillID string

// BattleCardSkill is a skill for a battle card.
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

// YieldModifier is a skill that modifies the Yield of a Territory.
type YieldModifier interface {
	Modify(quantity ResourceQuantity) ResourceQuantity // Modifies the argument quantity.
}

// AddYieldModifier adds to the resource quantity.
type AddYieldModifier struct {
	ResourceQuantity ResourceQuantity
}

func (m *AddYieldModifier) Modify(quantity ResourceQuantity) ResourceQuantity {
	return quantity.Add(m.ResourceQuantity)
}

// MultiplyYieldModifier multiplies the resource quantity.
type MultiplyYieldModifier struct {
	FoodMultiply  float64
	MoneyMultiply float64
	WoodMultiply  float64
	IronMultiply  float64
	ManaMultiply  float64
}

func (m *MultiplyYieldModifier) Modify(quantity ResourceQuantity) ResourceQuantity {
	return ResourceQuantity{
		Food:  int(float64(quantity.Food) * (1.0 + m.FoodMultiply)),
		Money: int(float64(quantity.Money) * (1.0 + m.MoneyMultiply)),
		Wood:  int(float64(quantity.Wood) * (1.0 + m.WoodMultiply)),
		Iron:  int(float64(quantity.Iron) * (1.0 + m.IronMultiply)),
		Mana:  int(float64(quantity.Mana) * (1.0 + m.ManaMultiply)),
	}
}

// BattlefieldModifier is a skill that modifies the state of the Battlefield.
type BattlefieldModifier interface {
	Modify(battlefield *Battlefield) *Battlefield // Modifies the argument battlefield.
}

// CardSlotBattlefieldModifier modifies the number of card slots.
type CardSlotBattlefieldModifier struct {
	Value int
}

func (m *CardSlotBattlefieldModifier) Modify(battlefield *Battlefield) *Battlefield {
	battlefield.CardSlot += m.Value
	return battlefield
}

// SupportPowerBattlefieldModifier modifies the support power.
type SupportPowerBattlefieldModifier struct {
	Value float64
}

func (m *SupportPowerBattlefieldModifier) Modify(battlefield *Battlefield) *Battlefield {
	battlefield.BaseSupportPower += m.Value
	return battlefield
}
