package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/geom"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// MapGridView マップグリッド表示Widget
// 位置: (0,20,520,280) - MainView内で描画
// 5x5のPoint配置（固定）、520x280を5x5に分割（104x56セル）
type MapGridView struct {
	GameState     *core.GameState
	TopLeft       geom.PointF
	CellSize      geom.PointF
	CellLocations []geom.PointF

	// View切り替えのコールバック
	OnPointClicked func(point core.Point)
}

// NewMapGridView MapGridViewを作成する
func NewMapGridView(gameState *core.GameState, onPointClicked func(point core.Point)) *MapGridView {
	cellSize := geom.PointF{X: 520.0 / 5.0, Y: 280.0 / 5.0}
	cellLocations := make([]geom.PointF, 25)
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			// (0,0)が左下になるように、描画Y座標を計算
			// yは論理座標(0が一番下), 4-yで描画座標(0が一番上)に変換
			cellLocations[y*5+x] = geom.PointF{X: float64(x) * cellSize.X, Y: float64(4-y) * cellSize.Y}
		}
	}
	return &MapGridView{
		GameState:      gameState,
		TopLeft:        geom.PointF{X: 0, Y: 20},
		CellSize:       cellSize,
		CellLocations:  cellLocations,
		OnPointClicked: onPointClicked,
	}
}

// HandleInput 入力処理
func (m *MapGridView) HandleInput(input *Input) error {
	justReleased := input.Mouse.IsJustReleased(ebiten.MouseButtonLeft)
	if !justReleased {
		return nil
	}

	cursorX, cursorY := input.Mouse.CursorPosition()
	relativeX := float64(cursorX) - m.TopLeft.X
	relativeY := float64(cursorY) - m.TopLeft.Y

	viewWidth := m.CellSize.X * 5
	viewHeight := m.CellSize.Y * 5

	if relativeX < 0 || relativeX >= viewWidth || relativeY < 0 || relativeY >= viewHeight {
		return nil
	}

	drawGridX := int(relativeX / m.CellSize.X)
	drawGridY := int(relativeY / m.CellSize.Y)

	// (0,0)が左下なのでY座標を反転
	gridY := 4 - drawGridY

	point := m.GameState.MapGrid.GetPoint(drawGridX, gridY)
	if point == nil {
		return nil
	}

	// Point画像の描画領域を計算
	cellTopLeft := m.CellLocations[gridY*5+drawGridX]
	imageX := cellTopLeft.X + (m.CellSize.X-24)/2
	imageY := cellTopLeft.Y + (m.CellSize.Y-24)/2 - 10 // Drawメソッドでのオフセットを考慮
	imageWidth := 24.0
	imageHeight := 24.0

	// クリック位置がPoint画像の範囲内かチェック
	if relativeX >= imageX && relativeX < imageX+imageWidth &&
		relativeY >= imageY && relativeY < imageY+imageHeight {

		// 到達可能なPointかチェック
		if m.GameState.MapGrid.CanInteract(drawGridX, gridY) {
			if m.OnPointClicked != nil {
				m.OnPointClicked(point)
			}
		}
	}

	return nil
}

// Draw 描画処理
func (m *MapGridView) Draw(screen *ebiten.Image) {
	if m.GameState == nil || m.GameState.MapGrid == nil {
		return
	}

	mapGrid := m.GameState.MapGrid

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			point := mapGrid.GetPoint(x, y)
			if point == nil {
				continue
			}

			// セルの左上座標を取得
			cellTopLeft := m.CellLocations[y*5+x]
			screenX := cellTopLeft.X + m.TopLeft.X
			screenY := cellTopLeft.Y + m.TopLeft.Y

			// 到達可能性の線を描画
			if m.GameState.CanInteract(x, y) {
				cellCenterX := screenX + m.CellSize.X/2
				cellCenterY := screenY + m.CellSize.Y/2
				m.drawConnectionLines(screen, x, y, cellCenterX, cellCenterY)
			}
		}
	}

	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			point := mapGrid.GetPoint(x, y)
			if point == nil {
				continue
			}

			// セルの左上座標を取得
			cellTopLeft := m.CellLocations[y*5+x]
			screenX := cellTopLeft.X + m.TopLeft.X
			screenY := cellTopLeft.Y + m.TopLeft.Y

			// Point画像を描画（24x24、セル中央）
			imageX := screenX + (m.CellSize.X-24)/2
			imageY := screenY + (m.CellSize.Y-24)/2 - 10 // 文字のスペースを考慮
			interactive := m.GameState.CanInteract(x, y)
			m.drawPointImage(screen, imageX, imageY, point, interactive)

			// Point名を描画（Point画像の下）
			textX := screenX + 8
			textY := imageY + 24 + 2
			pointName := m.getPointName(x, y, point)

			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(textX, textY)
			drawing.DrawText(screen, pointName, 12, opt)

			// もしコントロールされていない場合は、敵のパワーを描画
			if p, ok := point.(*core.WildernessPoint); ok && !p.Controlled && interactive {
				power := p.Enemy.Power
				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(imageX, imageY+8)
				powerIcon := drawing.Image("ui-power")
				screen.DrawImage(powerIcon, opt)
				opt.GeoM.Translate(16, 0)
				drawing.DrawText(screen, fmt.Sprintf("%.1f", power), 12, opt)
			}

			if p, ok := point.(*core.BossPoint); ok && interactive {
				power := p.Boss.Power
				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Translate(imageX, imageY+8)
				powerIcon := drawing.Image("ui-power")
				screen.DrawImage(powerIcon, opt)
				opt.GeoM.Translate(16, 0)
				drawing.DrawText(screen, fmt.Sprintf("%.1f", power), 12, opt)
			}
		}
	}
}

// drawPointImage Pointの画像を描画
func (m *MapGridView) drawPointImage(screen *ebiten.Image, x, y float64, point core.Point, interactive bool) {
	switch typedPoint := point.(type) {
	case *core.MyNationPoint:
		pointImg := drawing.Image("point-mynation")
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(x, y)
		screen.DrawImage(pointImg, opt)
	case *core.OtherNationPoint:
		pointImg := drawing.Image("point-othernation")
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(x, y)
		if !interactive {
			opt.ColorScale.Scale(0.5, 0.5, 0.5, 1)
		}
		screen.DrawImage(pointImg, opt)
	case *core.WildernessPoint:
		pointImg := drawing.Image(typedPoint.TerrainType)
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(x, y)
		if !interactive {
			opt.ColorScale.Scale(0.5, 0.5, 0.5, 1)
		}
		screen.DrawImage(pointImg, opt)
	case *core.BossPoint:
		pointImg := drawing.Image("point-boss")
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(x, y)
		if !interactive {
			opt.ColorScale.Scale(0.5, 0.5, 0.5, 1)
		}
		screen.DrawImage(pointImg, opt)
	default:
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Scale(24, 24)
		opt.GeoM.Translate(x, y)
		screen.DrawImage(drawing.WhitePixel, opt)
	}
}

// getPointName Point名を取得
func (m *MapGridView) getPointName(x, y int, point core.Point) string {
	switch p := point.(type) {
	case *core.MyNationPoint:
		return lang.Text("nation-mynation")
	case *core.OtherNationPoint:
		return lang.Text(string(p.OtherNation.NationID))
	case *core.WildernessPoint:
		return lang.Text(p.TerrainType)
	case *core.BossPoint:
		return lang.Text("point-boss")
	default:
		return lang.ExecuteTemplate("point-area", map[string]any{"x": x, "y": y})
	}
}

// drawConnectionLines 到達可能なPointへの線を描画
func (m *MapGridView) drawConnectionLines(screen *ebiten.Image, x, y int, centerX, centerY float64) {
	// 隣接する4方向をチェック
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for _, dir := range directions {
		nextX, nextY := x+dir[0], y+dir[1]

		// 範囲内チェック
		if nextX < 0 || nextX >= 5 || nextY < 0 || nextY >= 5 {
			continue
		}

		// 隣接Pointが到達可能かチェック
		if m.GameState.CanInteract(nextX, nextY) {
			// 線を描画
			nextCellTopLeft := m.CellLocations[nextY*5+nextX]
			nextScreenX := nextCellTopLeft.X + m.TopLeft.X
			nextScreenY := nextCellTopLeft.Y + m.TopLeft.Y
			nextCenterX := nextScreenX + m.CellSize.X/2
			nextCenterY := nextScreenY + m.CellSize.Y/2

			m.drawLine(screen, centerX, centerY, nextCenterX, nextCenterY)
		}
	}
}

// drawLine 2点間に線を描画
func (m *MapGridView) drawLine(screen *ebiten.Image, x1, y1, x2, y2 float64) {
	vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 2, color.White, true)
}
