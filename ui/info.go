package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
	"github.com/noppikinatta/ebitenginegamejam2025/viewmodel"
)

// InfoViewMode is the display mode of InfoView.
type InfoViewMode int

const (
	InfoModeHistory InfoViewMode = iota
	InfoModeCardInfo
	InfoModeNationPoint
	InfoModeWildernessPoint
	InfoModeEnemySkill
)

// InfoView is a widget for displaying information.
// Position: (520,20,120,280).
// Changes the content of the information displayed according to the situation.
type InfoView struct {
	CurrentMode InfoViewMode
	viewModel   *viewmodel.HistoryViewModel
}

// NewInfoView creates an InfoView.
func NewInfoView(viewModel *viewmodel.HistoryViewModel) *InfoView {
	return &InfoView{
		CurrentMode: InfoModeHistory, // The default is HistoryView.
		viewModel:   viewModel,
	}
}

// HandleInput handles input.
func (iv *InfoView) HandleInput(input *Input) error {
	// InfoView basically does not accept input (display only).
	return nil
}

// Draw handles drawing.
func (iv *InfoView) Draw(screen *ebiten.Image) {
	// Draw background.
	iv.drawBackground(screen)

	// Draw the content according to the current mode.
	switch iv.CurrentMode {
	case InfoModeHistory:
		iv.drawHistoryView(screen)
	case InfoModeCardInfo:
		iv.drawCardInfoView(screen)
	case InfoModeNationPoint:
		iv.drawNationPointView(screen)
	case InfoModeWildernessPoint:
		iv.drawWildernessPointView(screen)
	case InfoModeEnemySkill:
		iv.drawEnemySkillView(screen)
	}
}

// drawBackground draws the background.
func (iv *InfoView) drawBackground(screen *ebiten.Image) {
	// InfoView background (1040,40,240,680).
	vertices := []ebiten.Vertex{
		{DstX: 1040, DstY: 40, SrcX: 0, SrcY: 0, ColorR: 0.15, ColorG: 0.15, ColorB: 0.2, ColorA: 1},
		{DstX: 1280, DstY: 40, SrcX: 0, SrcY: 0, ColorR: 0.15, ColorG: 0.15, ColorB: 0.2, ColorA: 1},
		{DstX: 1280, DstY: 720, SrcX: 0, SrcY: 0, ColorR: 0.15, ColorG: 0.15, ColorB: 0.2, ColorA: 1},
		{DstX: 1040, DstY: 720, SrcX: 0, SrcY: 0, ColorR: 0.15, ColorG: 0.15, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
}

// drawHistoryView draws the HistoryView.
func (iv *InfoView) drawHistoryView(screen *ebiten.Image) {
	// Title.
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(1050, 40)
	drawing.DrawText(screen, lang.Text("ui-history"), 24, opt)

	historyLen := iv.viewModel.HistoryLen()

	// Display when there is no history.
	if historyLen == 0 {
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(1050, 80)
		drawing.DrawText(screen, lang.Text("ui-no-events"), 18, opt)
		return
	}

	// Display history events.
	historyDelta := 0
	if historyLen > 7 {
		historyDelta = historyLen - 7
		historyLen = 7
	}

	for i := range historyLen {
		historyIdx := i + historyDelta
		dateText := iv.viewModel.HistoryDateText(historyIdx)
		eventText := iv.viewModel.HistoryEventText(historyIdx)

		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(1050, 80.0+float64(i)*70)
		drawing.DrawText(screen, dateText, 18, opt)

		y := 80.0 + float64(i)*70
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(1050, y)
		drawing.DrawText(screen, dateText, 18, opt)

		y += 24
		opt = &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(1050, y)
		drawing.DrawText(screen, eventText, 16, opt)
	}
}

// drawCardInfoView draws the CardInfoView.
func (iv *InfoView) drawCardInfoView(screen *ebiten.Image) {
	// TODO: implement
}

// drawBattleCardInfo draws the detailed information of a BattleCard.
func (iv *InfoView) drawBattleCardInfo(screen *ebiten.Image, card *core.BattleCard, y float64) {
	// TODO: implement
}

// drawStructureCardInfo draws the detailed information of a StructureCard.
func (iv *InfoView) drawStructureCardInfo(screen *ebiten.Image, card *core.StructureCard, y float64) {
	// TODO: implement
}

// drawNationPointView draws the NationPointView.
func (iv *InfoView) drawNationPointView(screen *ebiten.Image) {
	// TODO: implement
}

// drawWildernessPointView draws the WildernessPointView.
func (iv *InfoView) drawWildernessPointView(screen *ebiten.Image) {
	// TODO: implement
}

// drawEnemySkillView draws the EnemySkillView.
func (iv *InfoView) drawEnemySkillView(screen *ebiten.Image) {
	// TODO: implement
}
