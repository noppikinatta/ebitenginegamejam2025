package scene

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/bamenn"
	"github.com/noppikinatta/bamenn/bamennutil"
	"github.com/noppikinatta/ebitenginegamejam2025/ui"
)

func CreateSequence(input *ui.Input) ebiten.Game {
	title := NewTitle(input)
	inGame := NewInGame(input)
	result := NewResult(input)
	seq := bamenn.NewSequence(title)
	tran := bamenn.NewLinearTransition(5, 10, bamennutil.LinearFillFadingDrawer{Color: color.Black})

	title.Init(inGame, seq, tran)
	inGame.Init(result, seq, tran)
	result.Init(title, seq, tran)

	return &wrapperGame{
		langSwitcher: &langSwitcher{},
		game:         seq,
	}
}

type wrapperGame struct {
	langSwitcher *langSwitcher
	game         ebiten.Game
}

func (w *wrapperGame) Update() error {
	w.langSwitcher.Update()
	return w.game.Update()
}

func (w *wrapperGame) Draw(screen *ebiten.Image) {
	w.game.Draw(screen)
	w.langSwitcher.Draw(screen)
}

func (w *wrapperGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return w.game.Layout(outsideWidth, outsideHeight)
}
