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

// BattleCard returns battle card view model at the specified index
func (vm *CardDeckViewModel) Card(idx int) (*CardViewModel, bool) {
	if vm.cardViewModelCache == nil {
		vm.cardViewModelCache = &CardViewModel{}
	}

	if idx < 0 || idx >= len(vm.gameState.CardDisplayOrder) {
		return nil, false
	}

	cardID := vm.gameState.CardDisplayOrder[idx]
	vm.cardViewModelCache.Duplicates = vm.gameState.CardDeck.Count(cardID)

	battleCard, ok := vm.gameState.CardDictionary.BattleCard(cardID)
	if ok {
		vm.cardViewModelCache.FromBattleCard(battleCard)
		return vm.cardViewModelCache, true
	}

	structureCard, ok := vm.gameState.CardDictionary.StructureCard(cardID)
	if ok {
		vm.setStructureCard(structureCard)
		return vm.cardViewModelCache, true
	}

	return nil, false
}

func (vm *CardDeckViewModel) CountTypesInHand() int {
	return vm.gameState.CardDeck.CountTypesInHand()
}

type CardViewModel struct {
	Image            *ebiten.Image
	Name             string
	Duplicates       int
	HasCardType      bool
	CardTypeImage    *ebiten.Image
	CardTypeName     string
	HasPower         bool
	Power            float64
	HasSkill         bool
	SkillName        string
	SkillDescription string
}

func (c *CardViewModel) FromBattleCard(battleCard *core.BattleCard) {
	c.Image = drawing.Image(string(battleCard.CardID))
	c.Name = lang.Text(string(battleCard.CardID))
	c.HasCardType = true
	c.CardTypeImage = drawing.Image(string(battleCard.Type))
	c.CardTypeName = lang.Text(string(battleCard.Type))
	c.HasPower = true
	c.Power = float64(battleCard.Power())
	c.HasSkill = true
	c.SkillName = lang.Text(string(battleCard.Skill.BattleCardSkillID))

}
