package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// ResourceView 資源表示Widget
// 位置: (0,0,300,20)
// 5種類の資源を60x20ずつで表示
type ResourceView struct {
	GameState *core.GameState
}

// NewResourceView ResourceViewを作成する
func NewResourceView(gameState *core.GameState) *ResourceView {
	return &ResourceView{
		GameState: gameState,
	}
}

// HandleInput 入力処理（ResourceViewは入力を受け付けない）
func (rv *ResourceView) HandleInput(input *Input) error {
	return nil
}

// Draw 描画処理
func (rv *ResourceView) Draw(screen *ebiten.Image) {
	resources := rv.GameState.Treasury.Resources
	yield := rv.GameState.GetYield()

	// 5種類の資源を80x20ずつで表示
	// Money (0, 0, 80, 20)
	DrawResource(screen, 0, 0, "resource-money", resources.Money, yield.Money)

	// Food (80, 0, 80, 20)
	DrawResource(screen, 80, 0, "resource-food", resources.Food, yield.Food)

	// Wood (160, 0, 80, 20)
	DrawResource(screen, 160, 0, "resource-wood", resources.Wood, yield.Wood)

	// Iron (240, 0, 80, 20)
	DrawResource(screen, 240, 0, "resource-iron", resources.Iron, yield.Iron)

	// Mana (320, 0, 80, 20)
	DrawResource(screen, 320, 0, "resource-mana", resources.Mana, yield.Mana)
}
