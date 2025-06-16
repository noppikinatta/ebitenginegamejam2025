package scene

import (
	"github.com/noppikinatta/bamenn"
)

type Game struct {
	*bamenn.Sequence
	title   *Title
	inGame  *InGame
	result  *Result
	current string
}

func CreateSequence() *Game {
	title := NewTitle()
	inGame := NewInGame()
	result := NewResult()

	game := &Game{
		title:   title,
		inGame:  inGame,
		result:  result,
		current: "title",
	}

	// 最初はタイトル画面から開始
	seq := bamenn.NewSequence(title)
	game.Sequence = seq

	return game
}

func (g *Game) Update() error {
	// 現在のシーンに応じて遷移をチェック
	if g.current == "title" {
		nextScene := g.title.GetNextScene()
		if nextScene == "ingame" {
			g.current = "ingame"
			g.Sequence.Switch(g.inGame)
		}
	}

	return g.Sequence.Update()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}

func (g *Game) GetCurrentScene() string {
	return g.current
}

func (g *Game) SetCurrentScene(sceneName string) {
	g.current = sceneName
	switch sceneName {
	case "title":
		g.Sequence.Switch(g.title)
	case "ingame":
		g.Sequence.Switch(g.inGame)
	case "result":
		g.Sequence.Switch(g.result)
	}
}

func (g *Game) GetTitleScene() *Title {
	return g.title
}

func (g *Game) GetInGameScene() *InGame {
	return g.inGame
}

func (g *Game) GetResultScene() *Result {
	return g.result
}

func (g *Game) TransitionTo(sceneName string) {
	g.SetCurrentScene(sceneName)
}
