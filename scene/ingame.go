package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/bamenn"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/load"
	"github.com/noppikinatta/ebitenginegamejam2025/ui"
)

type InGame struct {
	gameState  *core.GameState
	gameUI     *ui.GameUI
	input      *ui.Input
	canInput   bool
	nextScene  ebiten.Game
	sequence   *bamenn.Sequence
	transition bamenn.Transition
}

func NewInGame(input *ui.Input) *InGame {
	// ダミーGameStateを作成（Game Jam向け簡易実装）
	gameState := load.LoadGameState()

	// GameUIを初期化
	gameUI := ui.NewGameUI(gameState)

	return &InGame{
		gameState: gameState,
		gameUI:    gameUI,
		input:     input,
	}
}

func (g *InGame) Init(nextScene ebiten.Game, sequence *bamenn.Sequence, transition bamenn.Transition) {
	g.nextScene = nextScene
	g.sequence = sequence
	g.transition = transition
}

func (g *InGame) OnArrival() {
	g.canInput = true
}

func (g *InGame) Update() error {
	if !g.canInput {
		return nil
	}

	// マウス位置をGameUIに設定
	mouseX, mouseY := ebiten.CursorPosition()
	g.gameUI.SetMousePosition(mouseX, mouseY)

	// GameUIの入力処理
	if err := g.gameUI.HandleInput(g.input); err != nil {
		return err
	}

	// GameUIの更新処理
	if err := g.gameUI.Update(); err != nil {
		return err
	}

	// デバッグ用：スペースキーでテストイベント追加
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.gameUI.AddHistoryEvent("Test event triggered")
	}

	return nil
}

func (g *InGame) Draw(screen *ebiten.Image) {
	// GameUIで全ての描画を行う
	g.gameUI.Draw(screen)
}

func (g *InGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}
