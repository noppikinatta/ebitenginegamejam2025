package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/core"
	"github.com/noppikinatta/ebitenginegamejam2025/drawing"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

type Widget interface {
	HandleInput(input *Input) error
	Draw(screen *ebiten.Image)
}

// Common drawing functions

// DrawResource draws resource quantity (60x20 area)
// x, y: drawing position
// resourceType: resource type (for display)
// value: value to display
// increment: increment (displayed as +2)
func DrawResource(screen *ebiten.Image, x, y float64, resourceType string, value int, increment int) {
	// Left 40x40 for resource icon (currently using text as substitute)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(2.0, 2.0)
	opt.GeoM.Translate(x, y)
	resourceIcon := drawing.Image(resourceType)
	screen.DrawImage(resourceIcon, opt)

	// Right 80x40 for numerical display
	var text string
	if increment != 0 {
		text = fmt.Sprintf("%d(+%d)", value, increment)
	} else {
		text = fmt.Sprintf("%d", value)
	}

	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+40, y+4)
	drawing.DrawText(screen, text, 20, opt)
}

func DrawCardBackground(screen *ebiten.Image, x, y float64, alpha float32) {
	var r, g, b float32
	r = 1.9
	g = 0.8
	b = 0.7
	// Draw frame (rectangle)
	drawing.DrawRect(screen, x+2, y+2, 76, 116, r*alpha, g*alpha, b*alpha, alpha)
	r = 0.9
	g = 0.7
	b = 0.5
	drawing.DrawRect(screen, x+6, y+6, 68, 108, r*alpha, g*alpha, b*alpha, alpha)
}

// DrawCard draws a card (80x120 area)
func DrawCard(screen *ebiten.Image, x, y float64, cardID string) {
	DrawCardBackground(screen, x, y, 1)

	// Draw card name
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(2.0, 2.0)
	opt.GeoM.Translate(x+8, y+8)
	opt.ColorScale.Scale(1, 1, 1, 0.8)
	cardImage := drawing.Image(cardID)
	screen.DrawImage(cardImage, opt)
	cardName := lang.Text(cardID)
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+2, y+96)
	drawing.DrawText(screen, cardName, 18, opt)
}

func DrawBattleCard(screen *ebiten.Image, x, y float64, battleCard *core.BattleCard) {
	DrawCard(screen, x, y, string(battleCard.CardID))
	// Draw card type
	var rectR, rectG, rectB, rectA float32
	rectR, rectG, rectB, rectA = 1, 1, 1, 1
	optText := &ebiten.DrawImageOptions{}
	optText.GeoM.Translate(x+6, y+2)
	cardTypeText := ""
	switch battleCard.Type {
	case "cardtype-str":
		cardTypeText = "S"
		rectR, rectG, rectB, rectA = 1, 0.2, 0.2, 1
	case "cardtype-agi":
		cardTypeText = "A"
		rectR, rectG, rectB, rectA = 0.2, 1, 0.2, 1
	case "cardtype-mag":
		cardTypeText = "M"
		rectR, rectG, rectB, rectA = 0.2, 0.2, 1, 1
	}
	drawing.DrawRect(screen, x+6, y+6, 16, 24, rectR, rectG, rectB, rectA)
	drawing.DrawText(screen, cardTypeText, 20, optText)

	// Draw power
	drawing.DrawRect(screen, x+6, y+68, 68, 28, 0.2, 0.2, 0.2, 0.75)
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Scale(2.0, 2.0)
	opt.GeoM.Translate(x+2, y+64)
	powerIcon := drawing.Image("ui-power")
	screen.DrawImage(powerIcon, opt)
	opt = &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(x+2+32, y+64)
	drawing.DrawText(screen, fmt.Sprintf("%.1f", battleCard.Power()), 24, opt)
}

// DrawButton draws a button (click detection to be implemented later)
func DrawButton(screen *ebiten.Image, x, y, width, height float64, imageKey string) {
	// Draw button frame
	vertices := []ebiten.Vertex{
		{DstX: float32(x), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.8, ColorB: 0.8, ColorA: 1},
		{DstX: float32(x + width), DstY: float32(y), SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.8, ColorB: 0.8, ColorA: 1},
		{DstX: float32(x + width), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.8, ColorB: 0.8, ColorA: 1},
		{DstX: float32(x), DstY: float32(y + height), SrcX: 0, SrcY: 0, ColorR: 0.8, ColorG: 0.8, ColorB: 0.8, ColorA: 1},
	}
	indices := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(vertices, indices, drawing.WhitePixel, &ebiten.DrawTrianglesOptions{})

	// Draw button image
	image := drawing.Image(imageKey)
	imageBounds := image.Bounds()
	imageWidth := float64(imageBounds.Dx())
	imageHeight := float64(imageBounds.Dy())
	imageX := x + width/2 - imageWidth/2
	imageY := y + height/2 - imageHeight/2
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(imageX, imageY)
	screen.DrawImage(image, opt)
}

func DrawCardDescriptionTooltip(screen *ebiten.Image, card interface{}, mouseX, mouseY int) {
	left, top := float64(mouseX-40), float64(mouseY-160)

	switch typedCard := card.(type) {
	case *core.BattleCard:
		drawing.DrawRect(screen, left, top, 480, 120, 0, 0, 0, 0.5)
		typeName := lang.Text(string(typedCard.Type))
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(left, top)
		drawing.DrawText(screen, typeName, 24, opt)

		if typedCard.Skill == nil {
			return
		}
		skillName := lang.Text(string(typedCard.Skill.BattleCardSkillID))
		skillDescription := lang.Text(string(typedCard.Skill.DescriptionKey))
		opt.GeoM.Translate(0, 32)
		drawing.DrawText(screen, skillName, 24, opt)
		opt.GeoM.Translate(0, 32)
		drawing.DrawText(screen, skillDescription, 18, opt)
	case *core.StructureCard:
		drawing.DrawRect(screen, left, top, 800, 120, 0, 0, 0, 0.5)
		description := lang.Text(string(typedCard.DescriptionKey))
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Translate(left, top)
		drawing.DrawText(screen, description, 18, opt)
	}
}
