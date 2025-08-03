package viewmodel

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// CardDeckViewModel provides display information for card deck UI
type CardDeckViewModel struct {
	gameState          *core.GameState
	cardViewModelCache *CardViewModel
}

// NewCardDeckViewModel creates a new CardDeckViewModel
func NewCardDeckViewModel(gameState *core.GameState) *CardDeckViewModel {
	return &CardDeckViewModel{
		gameState: gameState,
	}
}

func (vm *CardDeckViewModel) CardID(idx int) (core.CardID, bool) {
	if vm.cardViewModelCache == nil {
		vm.cardViewModelCache = &CardViewModel{}
	}

	if idx < 0 || idx >= len(vm.gameState.CardDisplayOrder) {
		return "", false
	}

	return vm.gameState.CardDisplayOrder[idx], true
}

// BattleCard returns battle card view model at the specified index
func (vm *CardDeckViewModel) Card(idx int) (*CardViewModel, bool) {
	cardID, ok := vm.CardID(idx)
	if !ok {
		return nil, false
	}

	vm.cardViewModelCache.Duplicates = vm.gameState.CardDeck.Count(cardID)

	battleCard, ok := vm.gameState.CardDictionary.BattleCard(cardID)
	if ok {
		vm.cardViewModelCache.FromBattleCard(battleCard)
		return vm.cardViewModelCache, true
	}

	structureCard, ok := vm.gameState.CardDictionary.StructureCard(cardID)
	if ok {
		vm.cardViewModelCache.FromStructureCard(structureCard)
		return vm.cardViewModelCache, true
	}

	return nil, false
}

func (vm *CardDeckViewModel) CountTypesInHand() int {
	return vm.gameState.CardDeck.CountTypesInHand()
}

func (vm *CardDeckViewModel) IsBattleCard(cardID core.CardID) bool {
	_, ok := vm.gameState.CardDictionary.BattleCard(cardID)
	return ok
}

func (vm *CardDeckViewModel) IsStructureCard(cardID core.CardID) bool {
	_, ok := vm.gameState.CardDictionary.StructureCard(cardID)
	return ok
}

type CardViewModel struct {
	Image            *ebiten.Image
	Name             string
	Duplicates       int
	HasCardType      bool
	CardTypeColor    drawing.ColorF32
	CardTypeName     string
	HasPower         bool
	Power            float64
	HasSkill         bool
	SkillName        string
	SkillDescription string
}

func (c *CardViewModel) reset() {
	c.Image = nil
	c.Name = ""
	c.HasCardType = false
	c.CardTypeColor = drawing.ColorF32{}
	c.CardTypeName = ""
	c.HasPower = false
	c.HasSkill = false
	c.SkillName = ""
	c.SkillDescription = ""
}

func (c *CardViewModel) FromBattleCard(battleCard *core.BattleCard) {
	c.reset()
	c.Image = drawing.Image(string(battleCard.CardID))
	c.Name = lang.Text(string(battleCard.CardID))
	c.HasCardType = true
	switch battleCard.Type {
	case "cardtype-str":
		c.CardTypeColor = drawing.NewColorF32(1, 0.2, 0.2, 1)
	case "cardtype-agi":
		c.CardTypeColor = drawing.NewColorF32(0.2, 1, 0.2, 1)
	case "cardtype-mag":
		c.CardTypeColor = drawing.NewColorF32(0.2, 0.2, 1, 1)
	}
	c.CardTypeName = lang.Text(string(battleCard.Type))
	c.HasPower = true
	c.Power = float64(battleCard.Power())
	if battleCard.Skill != nil {
		c.HasSkill = true
		c.SkillName = lang.Text(string(battleCard.Skill.BattleCardSkillID))
		c.SkillDescription = lang.Text(battleCard.Skill.DescriptionKey) // TODO: description should be moved from core package
	}
}

func (c *CardViewModel) FromStructureCard(structureCard *core.StructureCard) {
	c.reset()
	c.Image = drawing.Image(string(structureCard.ID()))
	c.Name = lang.Text(string(structureCard.ID()))
}
