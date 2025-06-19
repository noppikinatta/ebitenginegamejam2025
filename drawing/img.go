package drawing

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/asset"
)

var (
	dummyImageBase = ebiten.NewImage(3, 3)

	// WhitePixel is useful to draw fill shape with DrawTriangles.
	WhitePixel = dummyImageBase.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

	fallbackImage = ebiten.NewImage(32, 32)
)

func init() {
	dummyImageBase.Fill(color.White)
	fallbackImage.Fill(color.RGBA{G: 255, A: 255})
	DrawText(fallbackImage, "IMAGE\n NOT\n  FOUND", 9, &ebiten.DrawImageOptions{})
}

func Image(key string) *ebiten.Image {
	img, ok := asset.Images()[key]
	if !ok {
		return fallbackImage
	}
	return img
}
