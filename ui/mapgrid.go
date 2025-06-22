package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
)

// MapGridView マップグリッド表示Widget
// 位置: (0,20,520,280) - MainView内で描画
// 5x5のPoint配置（固定）、520x280を5x5に分割（104x56セル）
type MapGridView struct {
	GameState *core.GameState
	
	// View切り替えのコールバック
	OnPointClicked func(point core.Point, viewType ViewType)
}

// NewMapGridView MapGridViewを作成する
func NewMapGridView(gameState *core.GameState) *MapGridView {
	return &MapGridView{
		GameState: gameState,
	}
}

// HandleInput 入力処理
func (mgv *MapGridView) HandleInput(input *Input) error {
	// TODO: マウスクリック処理の実装
	// 現在はnyuuryoku.Mouseの使い方が不明なので後で実装
	return nil
}

// Draw 描画処理
func (mgv *MapGridView) Draw(screen *ebiten.Image) {
	if mgv.GameState == nil || mgv.GameState.MapGrid == nil {
		return
	}

	mapGrid := mgv.GameState.MapGrid
	
	// 5x5のセルに分割して描画
	cellWidth := 520.0 / 5.0   // 104
	cellHeight := 280.0 / 5.0  // 56
	
	for y := 0; y < 5; y++ {
		for x := 0; x < 5; x++ {
			// セルの位置を計算（MainView内での相対位置）
			cellX := float64(x) * cellWidth
			cellY := float64(y) * cellHeight + 20 // MainViewのY座標オフセット
			
			// Pointを取得
			point := mapGrid.GetPoint(x, y)
			if point == nil {
				continue
			}
			
			// Point画像を描画（24x24、セル中央）
			imageX := cellX + (cellWidth-24)/2
			imageY := cellY + (cellHeight-24)/2 - 10 // 文字のスペースを考慮
			mgv.drawPointImage(screen, imageX, imageY, point)
			
			// Point名を描画（Point画像の下）
			textX := cellX + cellWidth/2 - 20 // 中央寄せ（概算）
			textY := imageY + 24 + 5
			pointName := mgv.getPointName(x, y, point)
			
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(textX, textY)
			drawing.DrawText(screen, pointName, 12, opt)
			
			// 到達可能性の線を描画
			if mgv.GameState.CanInteract(x, y) {
				mgv.drawConnectionLines(screen, x, y, cellX+cellWidth/2, cellY+cellHeight/2)
			}
		}
	}
}

// drawPointImage Pointの画像を描画
func (mgv *MapGridView) drawPointImage(screen *ebiten.Image, x, y float64, point core.Point) {
	// 24x24の矩形を描画（後でイラストに差し替え）
	var color [4]float32
	
	switch point.(type) {
	case *core.MyNationPoint:
		color = [4]float32{0.2, 0.8, 0.2, 1} // 緑
	case *core.OtherNationPoint:
		color = [4]float32{0.2, 0.2, 0.8, 1} // 青
	case *core.WildernessPoint:
		wilderness := point.(*core.WildernessPoint)
		if wilderness.Controlled {
			color = [4]float32{0.8, 0.8, 0.2, 1} // 黄（制圧済み）
		} else {
			color = [4]float32{0.8, 0.2, 0.2, 1} // 赤（未制圧）
		}
	case *core.BossPoint:
		color = [4]float32{0.8, 0.2, 0.8, 1} // 紫
	default:
		color = [4]float32{0.5, 0.5, 0.5, 1} // 灰
	}
	
	vertices := []ebiten.Vertex{
		{DstX: float32(x), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
		{DstX: float32(x + 24), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
		{DstX: float32(x + 24), DstY: float32(y + 24), SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
		{DstX: float32(x), DstY: float32(y + 24), SrcX: 0, SrcY: 0, ColorR: color[0], ColorG: color[1], ColorB: color[2], ColorA: color[3]},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
}

// getPointName Point名を取得
func (mgv *MapGridView) getPointName(x, y int, point core.Point) string {
	switch p := point.(type) {
	case *core.MyNationPoint:
		return "My Nation"
	case *core.OtherNationPoint:
		return fmt.Sprintf("Nation %s", p.OtherNation.NationID)
	case *core.WildernessPoint:
		if p.Controlled {
			return fmt.Sprintf("Area %d,%d", x, y)
		} else {
			return fmt.Sprintf("Wild %d,%d", x, y)
		}
	case *core.BossPoint:
		return "Boss"
	default:
		return fmt.Sprintf("Point %d,%d", x, y)
	}
}

// drawConnectionLines 到達可能なPointへの線を描画
func (mgv *MapGridView) drawConnectionLines(screen *ebiten.Image, x, y int, centerX, centerY float64) {
	// 隣接する4方向をチェック
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	
	cellWidth := 520.0 / 5.0
	cellHeight := 280.0 / 5.0
	
	for _, dir := range directions {
		nextX, nextY := x+dir[0], y+dir[1]
		
		// 範囲内チェック
		if nextX < 0 || nextX >= 5 || nextY < 0 || nextY >= 5 {
			continue
		}
		
		// 隣接Pointが到達可能かチェック
		if mgv.GameState.CanInteract(nextX, nextY) {
			// 線を描画
			nextCenterX := float64(nextX)*cellWidth + cellWidth/2
			nextCenterY := float64(nextY)*cellHeight + cellHeight/2 + 20
			
			mgv.drawLine(screen, centerX, centerY, nextCenterX, nextCenterY)
		}
	}
}

// drawLine 2点間に線を描画
func (mgv *MapGridView) drawLine(screen *ebiten.Image, x1, y1, x2, y2 float64) {
	// 簡単な線描画（太さ2ピクセル）
	vertices := []ebiten.Vertex{
		{DstX: float32(x1), DstY: float32(y1), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 0.8},
		{DstX: float32(x2), DstY: float32(y2), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 0.8},
		{DstX: float32(x2), DstY: float32(y2 + 2), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 0.8},
		{DstX: float32(x1), DstY: float32(y1 + 2), SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 0.8},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
}
