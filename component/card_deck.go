package component

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/noppikinatta/ebitenginegamejam2025/entity"
	"github.com/noppikinatta/ebitenginegamejam2025/system"
)

type CardLayout struct {
	cardPositions [][2]int // x, y positions for each card
	cardWidth     int
	cardHeight    int
}

type CardDeck struct {
	cards       []*entity.Card
	cardManager *system.CardManager
	x, y        int
	width       int
	height      int
	layout      *CardLayout
}

func NewCardDeck(cm *system.CardManager, x, y, width, height int) *CardDeck {
	// 初期カードをテンプレートから生成
	initialCards := []*entity.Card{
		cm.CreateCard("Warrior", "Unit", map[string]int{"Gold": 50}),
		cm.CreateCard("Shield Guard", "Unit", map[string]int{"Gold": 40}),
	}

	return &CardDeck{
		cards:       initialCards,
		cardManager: cm,
		x:           x,
		y:           y,
		width:       width,
		height:      height,
		layout: &CardLayout{
			cardWidth:  80,
			cardHeight: 50,
		},
	}
}

func (cd *CardDeck) GetCards() []*entity.Card {
	return cd.cards
}

// Backward-compatible AddCard accepting name only (tests use this)
func (cd *CardDeck) AddCard(cardName string) {
	// デフォルトを Unit として作成
	card := cd.cardManager.CreateCard(cardName, "Unit", map[string]int{"Gold": 0})
	cd.cards = append(cd.cards, card)
	cd.updateLayout()
}

func (cd *CardDeck) AddCardEntity(card *entity.Card) {
	cd.cards = append(cd.cards, card)
	cd.updateLayout()
}

func (cd *CardDeck) GetHorizontalLayout() *CardLayout {
	return cd.layout
}

func (cd *CardDeck) updateLayout() {
	cd.layout.cardPositions = make([][2]int, len(cd.cards))

	cardSpacing := 85 // 5px spacing between cards
	maxCardsPerRow := cd.width / cardSpacing

	for i := range cd.cards {
		row := i / maxCardsPerRow
		col := i % maxCardsPerRow

		cardX := cd.x + col*cardSpacing
		cardY := cd.y + row*55 // 55px per row to accommodate card height + spacing

		cd.layout.cardPositions[i] = [2]int{cardX, cardY}
	}
}

func (cd *CardDeck) Draw(screen *ebiten.Image) {
	// Draw background
	deckBg := ebiten.NewImage(cd.width, cd.height)
	deckBg.Fill(color.RGBA{30, 30, 50, 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(cd.x), float64(cd.y))
	screen.DrawImage(deckBg, op)

	// Draw title
	ebitenutil.DebugPrintAt(screen, "Card Deck", cd.x+5, cd.y+5)

	// Update layout before drawing
	cd.updateLayout()

	// Draw cards
	for i, card := range cd.cards {
		if i < len(cd.layout.cardPositions) {
			pos := cd.layout.cardPositions[i]

			// Draw card background
			cardBg := ebiten.NewImage(cd.layout.cardWidth, cd.layout.cardHeight)
			cardBg.Fill(color.RGBA{80, 80, 120, 255})

			cardOp := &ebiten.DrawImageOptions{}
			cardOp.GeoM.Translate(float64(pos[0]), float64(pos[1]))
			screen.DrawImage(cardBg, cardOp)

			// Draw card name (truncated if too long)
			displayName := card.Name
			if len(displayName) > 10 {
				displayName = displayName[:7] + "..."
			}
			ebitenutil.DebugPrintAt(screen, displayName, pos[0]+2, pos[1]+2)
		}
	}
}
