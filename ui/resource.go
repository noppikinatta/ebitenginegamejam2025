package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
)

// ResourceView 資源表示Widget
// 位置: (0,0,300,20)
// 5種類の資源を60x20ずつで表示
type ResourceView struct {
	Treasury  *core.Treasury           // 国庫情報
	Increment core.ResourceQuantity    // 増分情報
}

// NewResourceView ResourceViewを作成する
func NewResourceView(treasury *core.Treasury) *ResourceView {
	return &ResourceView{
		Treasury:  treasury,
		Increment: core.ResourceQuantity{}, // デフォルトは増分なし
	}
}

// SetIncrement 増分情報を設定する
func (rv *ResourceView) SetIncrement(increment core.ResourceQuantity) {
	rv.Increment = increment
}

// HandleInput 入力処理（ResourceViewは入力を受け付けない）
func (rv *ResourceView) HandleInput(input *Input) error {
	return nil
}

// Draw 描画処理
func (rv *ResourceView) Draw(screen *ebiten.Image) {
	if rv.Treasury == nil {
		return
	}

	resources := rv.Treasury.Resources
	
	// 5種類の資源を60x20ずつで表示
	// Money (0, 0, 60, 20)
	DrawResource(screen, 0, 0, GetResourceIcon("Money"), resources.Money, rv.Increment.Money)
	
	// Food (60, 0, 60, 20)  
	DrawResource(screen, 60, 0, GetResourceIcon("Food"), resources.Food, rv.Increment.Food)
	
	// Wood (120, 0, 60, 20)
	DrawResource(screen, 120, 0, GetResourceIcon("Wood"), resources.Wood, rv.Increment.Wood)
	
	// Iron (180, 0, 60, 20)
	DrawResource(screen, 180, 0, GetResourceIcon("Iron"), resources.Iron, rv.Increment.Iron)
	
	// Mana (240, 0, 60, 20)
	DrawResource(screen, 240, 0, GetResourceIcon("Mana"), resources.Mana, rv.Increment.Mana)
}
