package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
)

// MarketView Market表示Widget
// 位置: MainView内で描画
type MarketView struct {
	Nation    interface{}      // MyNationまたはOtherNation
	Treasury  *core.Treasury   // 購入可能性判定用
	GameState *core.GameState  // ゲーム状態
	
	// View切り替えのコールバック
	OnBackClicked func() // MapGridViewに戻る
	OnPurchase    func(cardPack *core.CardPack) // CardPack購入時
}

// NewMarketView MarketViewを作成する
func NewMarketView() *MarketView {
	return &MarketView{}
}

// SetNation 表示する国家を設定
func (mv *MarketView) SetNation(nation interface{}) {
	mv.Nation = nation
}

// SetTreasury 国庫を設定
func (mv *MarketView) SetTreasury(treasury *core.Treasury) {
	mv.Treasury = treasury
}

// SetGameState ゲーム状態を設定
func (mv *MarketView) SetGameState(gameState *core.GameState) {
	mv.GameState = gameState
}

// HandleInput 入力処理
func (mv *MarketView) HandleInput(input *Input) error {
	// TODO: マウスクリック処理
	// - 戻るボタンのクリック判定
	// - CardPackのクリック判定と購入処理
	return nil
}

// Draw 描画処理
func (mv *MarketView) Draw(screen *ebiten.Image) {
	if mv.Nation == nil {
		return
	}

	// ヘッダ描画 (0,20,480,40)
	mv.drawHeader(screen)
	
	// 戻るボタン描画 (480,20,40,40)
	mv.drawBackButton(screen)
	
	// CardPack一覧描画
	mv.drawCardPacks(screen)
}

// drawHeader Nation名ヘッダを描画
func (mv *MarketView) drawHeader(screen *ebiten.Image) {
	var nationName string
	
	switch n := mv.Nation.(type) {
	case *core.MyNation:
		nationName = "My Nation"
	case *core.OtherNation:
		nationName = fmt.Sprintf("Nation %s", n.NationID)
	default:
		nationName = "Unknown Nation"
	}
	
	// ヘッダ背景
	vertices := []ebiten.Vertex{
		{DstX: 0, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.3, ColorA: 1},
		{DstX: 480, DstY: 20, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.3, ColorA: 1},
		{DstX: 480, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.3, ColorA: 1},
		{DstX: 0, DstY: 60, SrcX: 0, SrcY: 0, ColorR: 0.3, ColorG: 0.3, ColorB: 0.3, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
	
	// Nation名テキスト
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(10, 30)
	drawing.DrawText(screen, nationName, 16, opt)
}

// drawBackButton 戻るボタンを描画
func (mv *MarketView) drawBackButton(screen *ebiten.Image) {
	DrawButton(screen, 480, 20, 40, 40, "X")
}

// drawCardPacks CardPack一覧を描画
func (mv *MarketView) drawCardPacks(screen *ebiten.Image) {
	cardPacks := mv.getVisibleCardPacks()
	
	// CardPack表示領域: 260x80 × 6個
	// 配置: (0,60,260,80), (260,60,260,80), (0,140,260,80), (260,140,260,80), (0,220,260,80), (260,220,260,80)
	positions := [][4]float64{
		{0, 60, 260, 80},     // 左上
		{260, 60, 260, 80},   // 右上
		{0, 140, 260, 80},    // 左中
		{260, 140, 260, 80},  // 右中
		{0, 220, 260, 80},    // 左下
		{260, 220, 260, 80},  // 右下
	}
	
	for i, cardPack := range cardPacks {
		if i >= 6 { // 最大6つまで表示
			break
		}
		
		pos := positions[i]
		mv.drawCardPack(screen, cardPack, i, pos[0], pos[1], pos[2], pos[3])
	}
}

// drawCardPack 個別のCardPackを描画
func (mv *MarketView) drawCardPack(screen *ebiten.Image, cardPack *core.CardPack, index int, x, y, width, height float64) {
	// CardPack枠を描画
	vertices := []ebiten.Vertex{
		{DstX: float32(x), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 0.9, ColorG: 0.9, ColorB: 0.9, ColorA: 1},
		{DstX: float32(x + width), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 0.9, ColorG: 0.9, ColorB: 0.9, ColorA: 1},
		{DstX: float32(x + width), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: 0.9, ColorG: 0.9, ColorB: 0.9, ColorA: 1},
		{DstX: float32(x), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: 0.9, ColorG: 0.9, ColorB: 0.9, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
	
	// CardPack画像 (0,60,40,40) -> 相対位置(0,0,40,40)
	mv.drawCardPackImage(screen, x, y, 40, 40)
	
	// CardPack名 (40,60,220,20) -> 相対位置(40,0,220,20)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+40, y)
	cardPackName := fmt.Sprintf("Pack %s", cardPack.CardPackID)
	drawing.DrawText(screen, cardPackName, 14, opt)
	
	// CardPack説明 (40,80,220,40) -> 相対位置(40,20,220,40)
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+40, y+20)
	description := "Card pack with various cards" // ダミーテキスト
	drawing.DrawText(screen, description, 10, opt)
	
	// CardPackの値段 (0,120,260,20) -> 相対位置(0,60,260,20)
	mv.drawCardPackPrice(screen, index, x, y+60, 260, 20)
}

// drawCardPackImage CardPack画像を描画
func (mv *MarketView) drawCardPackImage(screen *ebiten.Image, x, y, width, height float64) {
	// 24x32の画像（ダミーとして矩形）
	imageX := x + (width-24)/2
	imageY := y + (height-32)/2
	
	vertices := []ebiten.Vertex{
		{DstX: float32(imageX), DstY: float32(imageY), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: float32(imageX + 24), DstY: float32(imageY), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: float32(imageX + 24), DstY: float32(imageY + 32), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
		{DstX: float32(imageX), DstY: float32(imageY + 32), SrcX: 0, SrcY: 0, ColorR: 0.6, ColorG: 0.4, ColorB: 0.2, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})
}

// drawCardPackPrice CardPackの値段を描画
func (mv *MarketView) drawCardPackPrice(screen *ebiten.Image, index int, x, y, width, height float64) {
	// 値段情報を取得
	price, canPurchase := mv.getCardPackPrice(index)
	if price == nil {
		return
	}
	
	// Resource1種類につき60x20で表示
	resourceTypes := []struct {
		name  string
		value int
	}{
		{"Money", price.Money},
		{"Food", price.Food},
		{"Wood", price.Wood},
		{"Iron", price.Iron},
		{"Mana", price.Mana},
	}
	
	currentX := x
	for _, resource := range resourceTypes {
		if resource.value > 0 && currentX < x+width-60 {
			// Resource画像(20x20)とPrice数字(40x20)
			icon := GetResourceIcon(resource.name)
			
			// Resourceアイコン
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(currentX, y)
			drawing.DrawText(screen, icon, 12, opt)
			
			// Price数字（購入不可能な場合は赤文字）
			opt = &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(currentX+20, y)
			if !canPurchase {
				// TODO: 赤文字の実装（現在は通常色）
			}
			priceText := fmt.Sprintf("%d", resource.value)
			drawing.DrawText(screen, priceText, 12, opt)
			
			currentX += 60
		}
	}
}

// getVisibleCardPacks 表示可能なCardPack一覧を取得
func (mv *MarketView) getVisibleCardPacks() []*core.CardPack {
	switch n := mv.Nation.(type) {
	case *core.MyNation:
		return n.VisibleCardPacks()
	case *core.OtherNation:
		return n.VisibleCardPacks()
	default:
		return []*core.CardPack{}
	}
}

// getCardPackPrice CardPackの価格と購入可能性を取得
func (mv *MarketView) getCardPackPrice(index int) (*core.ResourceQuantity, bool) {
	if mv.Treasury == nil {
		return nil, false
	}
	
	switch n := mv.Nation.(type) {
	case *core.MyNation:
		if n.Market != nil && index < len(n.Market.Items) {
			item := n.Market.Items[index]
			canPurchase := n.CanPurchase(index, mv.Treasury)
			return &item.Price, canPurchase
		}
	case *core.OtherNation:
		if n.Market != nil && index < len(n.Market.Items) {
			item := n.Market.Items[index]
			canPurchase := n.CanPurchase(index, mv.Treasury)
			return &item.Price, canPurchase
		}
	}
	
	return nil, false
}
