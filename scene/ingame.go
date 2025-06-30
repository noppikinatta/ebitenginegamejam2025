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
	// Create a dummy GameState (simplified implementation for Game Jam)
	gameState := load.LoadGameState()

	// Initialize GameUI
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

	if g.gameState.IsVictory() {
		g.sequence.SwitchWithTransition(g.nextScene, g.transition)
		return nil
	}

	// Set mouse position in GameUI
	mouseX, mouseY := ebiten.CursorPosition()
	g.gameUI.SetMousePosition(mouseX, mouseY)

	// Handle input for GameUI
	if err := g.gameUI.HandleInput(g.input); err != nil {
		return err
	}

	// Update GameUI
	if err := g.gameUI.Update(); err != nil {
		return err
	}

	// For debugging: add a test event with the space key
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.gameUI.AddHistoryEvent("Test event triggered")
	}

	return nil
}

func (g *InGame) Draw(screen *ebiten.Image) {
	// All drawing is done in GameUI
	g.gameUI.Draw(screen)
}

func (g *InGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}
