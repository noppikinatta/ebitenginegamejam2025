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

	// 5種類の資源を60x20ずつで表示
	// Money (0, 0, 60, 20)
	DrawResource(screen, 0, 0, GetResourceIcon("Money"), resources.Money, yield.Money)

	// Food (60, 0, 60, 20)
	DrawResource(screen, 60, 0, GetResourceIcon("Food"), resources.Food, yield.Food)

	// Wood (120, 0, 60, 20)
	DrawResource(screen, 120, 0, GetResourceIcon("Wood"), resources.Wood, yield.Wood)

	// Iron (180, 0, 60, 20)
	DrawResource(screen, 180, 0, GetResourceIcon("Iron"), resources.Iron, yield.Iron)

	// Mana (240, 0, 60, 20)
	DrawResource(screen, 240, 0, GetResourceIcon("Mana"), resources.Mana, yield.Mana)
}
