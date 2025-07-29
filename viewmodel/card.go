package viewmodel

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// BattleCardViewModel provides display information for battle cards
type BattleCardViewModel struct {
	gameState       *core.GameState
	card            *core.BattleCard
	calculatedPower float64
}

// NewBattleCardViewModel creates a new BattleCardViewModel
func NewBattleCardViewModel(gameState *core.GameState, card *core.BattleCard, calculatedPower float64) *BattleCardViewModel {
	return &BattleCardViewModel{
		gameState:       gameState,
		card:            card,
		calculatedPower: calculatedPower,
	}
}

// Image returns the card image
func (vm *BattleCardViewModel) Image() *ebiten.Image {
	// Use drawing package to get card image
	return drawing.Image("battlecard-" + string(vm.card.CardID))
}

// Name returns the localized card name
func (vm *BattleCardViewModel) Name() string {
	// Get localized card name
	return lang.Text("battlecard_name_" + string(vm.card.CardID))
}

// Duplicates returns the number of this card in the deck
func (vm *BattleCardViewModel) Duplicates() int {
	return vm.gameState.CardDeck.Count(vm.card.CardID)
}

// CardTypeImage returns the card type image
func (vm *BattleCardViewModel) CardTypeImage() *ebiten.Image {
	// Use drawing package to get card type image
	return drawing.Image("cardtype-" + string(vm.card.Type))
}

// CardTypeName returns the localized card type name
func (vm *BattleCardViewModel) CardTypeName() string {
	// Get localized card type name
	return lang.Text("cardtype_name_" + string(vm.card.Type))
}

// Power returns the calculated power of the card
func (vm *BattleCardViewModel) Power() float64 {
	return vm.calculatedPower
}

// SkillName returns the localized skill name
func (vm *BattleCardViewModel) SkillName() string {
	if vm.card.Skill == nil {
		return ""
	}
	
	// Get localized skill name
	return lang.Text("battlecard_skill_name_" + string(vm.card.Skill.BattleCardSkillID))
}

// SkillDescription returns the localized skill description
func (vm *BattleCardViewModel) SkillDescription() string {
	if vm.card.Skill == nil {
		return ""
	}
	
	// Get localized skill description
	return lang.Text("battlecard_skill_desc_" + string(vm.card.Skill.BattleCardSkillID))
}

// StructureCardViewModel provides display information for structure cards
type StructureCardViewModel struct {
	gameState *core.GameState
	card      *core.StructureCard
}

// NewStructureCardViewModel creates a new StructureCardViewModel
func NewStructureCardViewModel(gameState *core.GameState, card *core.StructureCard) *StructureCardViewModel {
	return &StructureCardViewModel{
		gameState: gameState,
		card:      card,
	}
}

// Image returns the card image
func (vm *StructureCardViewModel) Image() *ebiten.Image {
	// Use drawing package to get structure card image
	return drawing.Image("structurecard-" + string(vm.card.ID()))
}

// Name returns the localized card name
func (vm *StructureCardViewModel) Name() string {
	// Get localized card name
	return lang.Text("structurecard_name_" + string(vm.card.ID()))
}

// Duplicates returns the number of this card in the deck
func (vm *StructureCardViewModel) Duplicates() int {
	return vm.gameState.CardDeck.Count(vm.card.ID())
} 