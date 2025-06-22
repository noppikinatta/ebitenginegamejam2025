package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
)

// CalendarView カレンダー表示Widget
// 位置: (520,0,120,20)
// 現在のTurnを年月表示する
type CalendarView struct {
	CurrentTurn int // 現在のターン
}

// NewCalendarView CalendarViewを作成する
func NewCalendarView() *CalendarView {
	return &CalendarView{
		CurrentTurn: 0, // デフォルトは1年1月
	}
}

// SetCurrentTurn 現在のターンを設定する
func (cv *CalendarView) SetCurrentTurn(turn int) {
	cv.CurrentTurn = turn
}

// HandleInput 入力処理（CalendarViewは入力を受け付けない）
func (cv *CalendarView) HandleInput(input *Input) error {
	return nil
}

// Draw 描画処理
func (cv *CalendarView) Draw(screen *ebiten.Image) {
	// Turnから年月を取得
	turn := core.Turn(cv.CurrentTurn)
	year, month := turn.YearMonth()

	// YYYY/MM形式で表示
	text := fmt.Sprintf("%04d/%02d", year, month)

	// 位置: (520,0,120,20)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(520, 0)
	drawing.DrawText(screen, text, 12, opt)
}
