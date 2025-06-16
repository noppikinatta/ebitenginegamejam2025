package scene

import (
	"github.com/noppikinatta/bamenn"
)

type Game struct {
	*bamenn.Sequence
	title   *Title
	inGame  *InGame
	current string
}

func CreateSequence() *Game {
	title := &Title{}
	inGame := &InGame{}

	game := &Game{
		title:   title,
		inGame:  inGame,
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
