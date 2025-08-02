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
	CardID    CardID
	BasePower BattleCardPower  // BasePower is the combat power of the card.
	Skill     *BattleCardSkill // Skill is the skill the card possesses.
	Type      BattleCardType   // Type is the card type, such as warrior, mage, or animal. Used to determine the target of a skill's effect.
}

// NewBattleCard creates a new BattleCard instance.
func NewBattleCard(cardID CardID, basePower BattleCardPower, skill *BattleCardSkill, cardType BattleCardType) *BattleCard {
	return &BattleCard{
		CardID:    cardID,
		BasePower: basePower,
		Skill:     skill,
		Type:      cardType,
	}
}

// ID returns the card ID.
func (c *BattleCard) ID() CardID {
	return c.CardID
}

// Power returns the combat power of the card.
func (c *BattleCard) Power() BattleCardPower {
	return c.BasePower
}

// StructureCard is a card placed in a Territory. This struct is immutable.
type StructureCard struct {
	cardID             CardID
	yieldAdditiveValue ResourceQuantity // Direct additive yield bonus
	yieldModifier      ResourceModifier // Multiplicative yield modifier
	supportPower       float64          // Support power provided to battlefield
	supportCardSlot    int              // Additional card slots provided to battlefield
}

// NewStructureCard creates a new StructureCard instance.
func NewStructureCard(cardID CardID, yieldAdditiveValue ResourceQuantity, yieldModifier ResourceModifier, supportPower float64, supportCardSlot int) *StructureCard {
	return &StructureCard{
		cardID:             cardID,
		yieldAdditiveValue: yieldAdditiveValue,
		yieldModifier:      yieldModifier,
		supportPower:       supportPower,
		supportCardSlot:    supportCardSlot,
	}
}

// ID returns the card ID.
func (c *StructureCard) ID() CardID {
	return c.cardID
}

// YieldAdditiveValue returns the direct additive yield bonus.
func (c *StructureCard) YieldAdditiveValue() ResourceQuantity {
	return c.yieldAdditiveValue
}

// YieldModifier returns the multiplicative yield modifier.
func (c *StructureCard) YieldModifier() ResourceModifier {
	return c.yieldModifier
}

// SupportPower returns the support power provided to battlefield.
func (c *StructureCard) SupportPower() float64 {
	return c.supportPower
}

// SupportCardSlot returns the additional card slots provided to battlefield.
func (c *StructureCard) SupportCardSlot() int {
	return c.supportCardSlot
}

// CardDictionary is a struct for generating cards.
type CardDictionary struct {
	battleCards    map[CardID]*BattleCard
	structureCards map[CardID]*StructureCard
}

func NewCardDictionary(battleCards []*BattleCard, structureCards []*StructureCard) *CardDictionary {
	dict := &CardDictionary{
		battleCards:    make(map[CardID]*BattleCard),
		structureCards: make(map[CardID]*StructureCard),
	}

	for _, card := range battleCards {
		dict.battleCards[card.ID()] = card
	}

	for _, card := range structureCards {
		dict.structureCards[card.ID()] = card
	}

	return dict
}

func (d *CardDictionary) BattleCard(cardID CardID) (*BattleCard, bool) {
	card, exists := d.battleCards[cardID]
	return card, exists
}

func (d *CardDictionary) StructureCard(cardID CardID) (*StructureCard, bool) {
	card, exists := d.structureCards[cardID]
	return card, exists
}

// CardDeck is the player's card deck.
type CardDeck struct {
	hand map[CardID]int
}

// NewCardDeck creates a new CardDeck instance.
func NewCardDeck() *CardDeck {
	return &CardDeck{
		hand: make(map[CardID]int),
	}
}

// Add adds a card to the deck by CardID.
func (cd *CardDeck) Add(cardID CardID) {
	cd.hand[cardID]++
}

// Remove removes a card from the deck by CardID.
// Returns true if the card was successfully removed, false if the card doesn't exist.
func (cd *CardDeck) Remove(cardID CardID) bool {
	if cd.hand[cardID] <= 0 {
		return false
	}
	cd.hand[cardID]--
	if cd.hand[cardID] == 0 {
		delete(cd.hand, cardID)
	}
	return true
}

// Count returns the number of cards with the given CardID in the deck.
func (cd *CardDeck) Count(cardID CardID) int {
	return cd.hand[cardID]
}

func (cd *CardDeck) CountTypesInHand() int {
	return len(cd.hand)
}

// GetAllCardIDs returns all CardIDs in the deck, with duplicates for multiple copies.
func (cd *CardDeck) GetAllCardIDs() []CardID {
	var cardIDs []CardID
	for cardID, count := range cd.hand {
		for i := 0; i < count; i++ {
			cardIDs = append(cardIDs, cardID)
		}
	}
	return cardIDs
}

// GetAllCardCounts returns a copy of the internal card count map.
func (cd *CardDeck) GetAllCardCounts() map[CardID]int {
	result := make(map[CardID]int)
	for cardID, count := range cd.hand {
		result[cardID] = count
	}
	return result
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

type BattleCardSkillCalculatorCondition struct {
	Condition  func(options *BattleCardSkillCalculationOptions) bool
	Calculator BattleCardSkillCalculator
}

func (c *BattleCardSkillCalculatorCondition) Calculate(options *BattleCardSkillCalculationOptions) {
	if !c.Condition(options) {
		return
	}
	c.Calculator.Calculate(options)
}

type BattleCardSkillCalculatorEffectSelf struct {
	Effect *BattleCardSkillEffect
}

func (c *BattleCardSkillCalculatorEffectSelf) Calculate(options *BattleCardSkillCalculationOptions) {
	c.Effect.Apply(options.BattleCardPowerModifiers[options.BattleCardIndex])
}

type BattleCardSkillCalculatorEffectIdxs struct {
	IdxDeltas []int
	Effect    *BattleCardSkillEffect
}

func (c *BattleCardSkillCalculatorEffectIdxs) Calculate(options *BattleCardSkillCalculationOptions) {
	for _, idxDelta := range c.IdxDeltas {
		modifierIdx := options.BattleCardIndex + idxDelta
		if modifierIdx < 0 || modifierIdx >= len(options.BattleCardPowerModifiers) {
			continue
		}
		c.Effect.Apply(options.BattleCardPowerModifiers[modifierIdx])
	}
}

type BattleCardSkillCalculatorEffectAll struct {
	Effect *BattleCardSkillEffect
}

func (c *BattleCardSkillCalculatorEffectAll) Calculate(options *BattleCardSkillCalculationOptions) {
	for _, m := range options.BattleCardPowerModifiers {
		c.Effect.Apply(m)
	}
}

type BattleCardSkillCalculatorEffectAllCondition struct {
	Condition func(idx int, options *BattleCardSkillCalculationOptions) bool
	Effect    *BattleCardSkillEffect
}

func (c *BattleCardSkillCalculatorEffectAllCondition) Calculate(options *BattleCardSkillCalculationOptions) {
	for i, m := range options.BattleCardPowerModifiers {
		if !c.Condition(i, options) {
			continue
		}
		c.Effect.Apply(m)
	}
}

type BattleCardSkillEffect struct {
	Modifier *BattleCardPowerModifier
}

func (e *BattleCardSkillEffect) Apply(modifier *BattleCardPowerModifier) {
	modifier.Union(e.Modifier)
}
