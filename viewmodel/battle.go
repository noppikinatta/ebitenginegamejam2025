package viewmodel

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// BattleViewModel provides display information for battle UI
type BattleViewModel struct {
	gameState   *core.GameState
	battlefield *core.Battlefield
	point       core.BattlePoint
}

// NewBattleViewModel creates a new BattleViewModel
func NewBattleViewModel(gameState *core.GameState, battlefield *core.Battlefield, point core.BattlePoint) *BattleViewModel {
	return &BattleViewModel{
		gameState:   gameState,
		battlefield: battlefield,
		point:       point,
	}
}

// Title returns the battle title
func (vm *BattleViewModel) Title() string {
	// Get localized battle title
	return lang.Text("battle_title")
}

// EnemyImage returns the enemy image
func (vm *BattleViewModel) EnemyImage() *ebiten.Image {
	if vm.point == nil || vm.point.Enemy() == nil {
		return nil
	}

	enemy := vm.point.Enemy()
	// Use drawing package to get enemy image by ID
	return drawing.Image("enemy-" + string(enemy.ID()))
}

// EnemyType returns the enemy type name
func (vm *BattleViewModel) EnemyType() string {
	if vm.point == nil || vm.point.Enemy() == nil {
		return ""
	}

	enemy := vm.point.Enemy()
	// Get localized enemy type name
	return lang.Text("enemy_type_" + string(enemy.Type()))
}

// EnemyPower returns the enemy power
func (vm *BattleViewModel) EnemyPower() float64 {
	if vm.point == nil || vm.point.Enemy() == nil {
		return 0.0
	}

	return vm.point.Enemy().Power()
}

// EnemyTalk returns the enemy dialogue
func (vm *BattleViewModel) EnemyTalk() string {
	if vm.point == nil || vm.point.Enemy() == nil {
		return ""
	}

	enemy := vm.point.Enemy()
	// Get localized enemy dialogue
	return lang.Text("enemy_talk_" + string(enemy.ID()))
}

// EnemySkillNames returns the enemy skill names
func (vm *BattleViewModel) EnemySkillNames() []string {
	if vm.point == nil || vm.point.Enemy() == nil {
		return []string{}
	}

	enemy := vm.point.Enemy()
	skills := enemy.Skills()
	names := make([]string, len(skills))

	for i, skill := range skills {
		names[i] = lang.Text("enemy_skill_name_" + string(skill.ID()))
	}

	return names
}

// EnemySkillDescriptions returns the enemy skill descriptions
func (vm *BattleViewModel) EnemySkillDescriptions() []string {
	if vm.point == nil || vm.point.Enemy() == nil {
		return []string{}
	}

	enemy := vm.point.Enemy()
	skills := enemy.Skills()
	descriptions := make([]string, len(skills))

	for i, skill := range skills {
		descriptions[i] = lang.Text("enemy_skill_desc_" + string(skill.ID()))
	}

	return descriptions
}

// CardSlot returns the battle card slot count
func (vm *BattleViewModel) CardSlot() int {
	if vm.battlefield == nil {
		return 0
	}
	return vm.battlefield.CardSlot
}

// CanBeat returns whether the enemy can be defeated
func (vm *BattleViewModel) CanBeat() bool {
	if vm.battlefield == nil {
		return false
	}
	return vm.battlefield.CanBeat()
}

// TotalPower returns the total power of placed battle cards
func (vm *BattleViewModel) TotalPower() float64 {
	if vm.battlefield == nil {
		return 0.0
	}
	return vm.battlefield.CalculateTotalPower()
}

// NumCards returns the number of placed battle cards
func (vm *BattleViewModel) NumCards() int {
	if vm.battlefield == nil {
		return 0
	}
	return len(vm.battlefield.BattleCards)
}

// Card returns battle card view model at the specified index
func (vm *BattleViewModel) Card(idx int) *BattleCardViewModel {
	if vm.battlefield == nil || idx < 0 || idx >= len(vm.battlefield.BattleCards) {
		return nil
	}

	card := vm.battlefield.BattleCards[idx]
	// Calculate the power for this specific card position
	// This might involve skill calculations specific to battlefield position
	calculatedPower := float64(card.Power()) // Simplified - may need battle calculations

	return NewBattleCardViewModel(vm.gameState, card, calculatedPower)
}
