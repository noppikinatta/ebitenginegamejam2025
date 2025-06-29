package ui

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

// MarketView Market表示Widget
// 位置: MainView内で描画
type MarketView struct {
	Nation    core.Nation     // MyNationまたはOtherNation
	Treasury  *core.Treasury  // 購入可能性判定用
	GameState *core.GameState // ゲーム状態

	// View切り替えのコールバック
	OnBackClicked func()                        // MapGridViewに戻る
	OnPurchase    func(cardPack *core.CardPack) // CardPack購入時
}

// NewMarketView MarketViewを作成する
func NewMarketView(onBackClicked func()) *MarketView {
	return &MarketView{
		OnBackClicked: onBackClicked,
	}
}

// SetNation 表示する国家を設定
func (mv *MarketView) SetNation(nation core.Nation) {
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
	if input.Mouse.IsJustReleased(ebiten.MouseButtonLeft) {
		cursorX, cursorY := input.Mouse.CursorPosition()

		// 戻るボタンのクリック判定 (480,20,40,40)
		if cursorX >= 480 && cursorX < 520 && cursorY >= 20 && cursorY < 60 {
			if mv.OnBackClicked != nil {
				mv.OnBackClicked()
				return nil
			}
		}

		// CardPackのクリック判定と購入処理
		mv.handleMarketItemClick(cursorX, cursorY)
	}
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
	mv.drawMarketItems(screen)
}

// drawHeader Nation名ヘッダを描画
func (mv *MarketView) drawHeader(screen *ebiten.Image) {
	nationName := mv.Nation.Name()

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
	DrawButton(screen, 480, 20, 40, 40, "ui-close")
}

// drawMarketItems MarketItem一覧を描画
func (mv *MarketView) drawMarketItems(screen *ebiten.Image) {
	marketItems := mv.getAllMarketItems()

	// CardPack表示領域: 260x80 × 6個
	// 配置: (0,60,260,80), (260,60,260,80), (0,140,260,80), (260,140,260,80), (0,220,260,80), (260,220,260,80)
	positions := [][4]float64{
		{0, 60, 260, 80},    // 左上
		{260, 60, 260, 80},  // 右上
		{0, 140, 260, 80},   // 左中
		{260, 140, 260, 80}, // 右中
		{0, 220, 260, 80},   // 左下
		{260, 220, 260, 80}, // 右下
	}

	for i, item := range marketItems {
		if i >= 6 { // 最大6つまで表示
			break
		}

		pos := positions[i]
		mv.drawMarketItem(screen, item, i, pos[0], pos[1], pos[2], pos[3])
	}
}

// drawMarketItem 個別のMarketItemを描画
func (mv *MarketView) drawMarketItem(screen *ebiten.Image, item *core.MarketItem, index int, x, y, width, height float64) {
	isAvailable := mv.isMarketItemAvailable(item)

	// CardPack枠を描画（レベル不足の場合は暗くする）
	var colorR, colorG, colorB float32 = 0.9, 0.9, 0.9
	if !isAvailable {
		colorR, colorG, colorB = 0.5, 0.5, 0.5 // 暗くする
	}

	vertices := []ebiten.Vertex{
		{DstX: float32(x), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: float32(x + width), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: float32(x + width), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
		{DstX: float32(x), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: colorR, ColorG: colorG, ColorB: colorB, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// CardPack画像 (0,60,40,40) -> 相対位置(0,0,40,40)
	mv.drawCardPackImage(screen, x, y, 40, 40)

	// CardPack名 (40,60,220,20) -> 相対位置(40,0,220,20)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+40, y)
	cardPackName := lang.Text("cardpack-" + string(item.CardPack.CardPackID))
	drawing.DrawText(screen, cardPackName, 14, opt)

	// CardPack説明 (40,80,220,40) -> 相対位置(40,20,220,40)
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+40, y+20)
	var description string
	if !isAvailable {
		description = lang.ExecuteTemplate("market-required-level", map[string]any{"level": item.RequiredLevel})
	} else {
		description = lang.Text("market-card-pack-desc")
	}
	drawing.DrawText(screen, description, 10, opt)

	// CardPackの値段 (0,120,260,20) -> 相対位置(0,60,260,20)
	mv.drawCardPackPrice(screen, item, index, x, y+60, 260, 20)
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
func (mv *MarketView) drawCardPackPrice(screen *ebiten.Image, item *core.MarketItem, index int, x, y, width, height float64) {
	// 値段情報を取得
	_, canPurchase := mv.getCardPackPrice(index)
	price := item.Price

	// Resource1種類につき60x20で表示
	resourceTypes := []struct {
		name  string
		value int
	}{
		{"resource-money", price.Money},
		{"resource-food", price.Food},
		{"resource-wood", price.Wood},
		{"resource-iron", price.Iron},
		{"resource-mana", price.Mana},
	}

	currentX := x
	for _, resource := range resourceTypes {
		if resource.value > 0 && currentX < x+width-60 {
			// Resource画像(20x20)とPrice数字(40x20)
			icon := drawing.Image(resource.name)

			// Resourceアイコン
			opt := &ebiten.DrawImageOptions{}
			opt.GeoM.Translate(currentX, y)
			screen.DrawImage(icon, opt)

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

// getAllMarketItems 全てのMarketItem一覧を取得（レベル不足含む）
func (mv *MarketView) getAllMarketItems() []*core.MarketItem {
	if mv.Nation == nil {
		return []*core.MarketItem{}
	}
	market := mv.Nation.GetMarket()
	if market == nil {
		return []*core.MarketItem{}
	}
	return market.Items
}

// getCardPackPrice CardPackの価格と購入可能性を取得
func (mv *MarketView) getCardPackPrice(index int) (*core.ResourceQuantity, bool) {
	if mv.Treasury == nil || mv.Nation == nil {
		return nil, false
	}

	market := mv.Nation.GetMarket()
	if market != nil && index < len(market.Items) {
		item := market.Items[index]
		canPurchase := mv.Nation.CanPurchase(index, mv.Treasury)
		return &item.Price, canPurchase
	}

	return nil, false
}

// handleMarketItemClick MarketItemのクリック処理
func (mv *MarketView) handleMarketItemClick(cursorX, cursorY int) {
	positions := [][4]int{
		{0, 60, 260, 80},    // 左上
		{260, 60, 260, 80},  // 右上
		{0, 140, 260, 80},   // 左中
		{260, 140, 260, 80}, // 右中
		{0, 220, 260, 80},   // 左下
		{260, 220, 260, 80}, // 右下
	}

	marketItems := mv.getAllMarketItems()

	for i, pos := range positions {
		if i >= len(marketItems) {
			break
		}

		if cursorX >= pos[0] && cursorX < pos[0]+pos[2] &&
			cursorY >= pos[1] && cursorY < pos[1]+pos[3] {
			// MarketItemがクリックされた
			item := marketItems[i]

			// レベル不足の場合は購入できない
			if !mv.isMarketItemAvailable(item) {
				return // 何もしない
			}

			if err := mv.PurchaseCardPack(item); err == nil {
				// 購入成功時、MapGridViewに戻る
				if mv.OnBackClicked != nil {
					mv.OnBackClicked()
				}
			}
			break
		}
	}
}

// simpleRand は Intner インターフェースを実装する簡単な乱数生成器
type simpleRand struct {
	*rand.Rand
}

func newSimpleRand() *simpleRand {
	return &simpleRand{rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (sr *simpleRand) Intn(n int) int {
	return sr.Rand.Intn(n)
}

// PurchaseCardPack カードパック購入処理
func (mv *MarketView) PurchaseCardPack(item *core.MarketItem) error {
	if mv.GameState == nil || mv.Nation == nil {
		return fmt.Errorf("GameState or Nation is nil")
	}

	market := mv.Nation.GetMarket()
	if market == nil {
		return fmt.Errorf("Market is nil")
	}

	// アイテムのインデックスを見つける
	itemIndex := -1
	for i, marketItem := range market.Items {
		if marketItem == item {
			itemIndex = i
			break
		}
	}

	if itemIndex == -1 {
		return fmt.Errorf("Item not found in market")
	}

	// 購入処理
	cardPack, ok := mv.Nation.Purchase(itemIndex, mv.GameState.Treasury)
	if !ok {
		return fmt.Errorf("Purchase failed")
	}

	// CardPackを開いてCardsを取得
	rng := newSimpleRand()
	cardIDs := cardPack.Open(rng)

	cards, ok := mv.GameState.CardGenerator.Generate(cardIDs)
	if !ok {
		return fmt.Errorf("Card generation failed")
	}

	// GameState.CardDeckに追加
	mv.GameState.CardDeck.Add(cards)

	mv.GameState.NextTurn()

	return nil
}

// isMarketItemAvailable MarketItemが利用可能かどうかを判定
func (mv *MarketView) isMarketItemAvailable(item *core.MarketItem) bool {
	if mv.Nation == nil {
		return false
	}
	market := mv.Nation.GetMarket()
	if market == nil {
		return false
	}
	return market.Level >= item.RequiredLevel
}
