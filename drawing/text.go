package drawing

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/noppikinatta/ebitenginegamejam2025/asset"
	"github.com/noppikinatta/ebitenginegamejam2025/geom"
	"github.com/noppikinatta/ebitenginegamejam2025/lang"
)

func DrawTextTemplate(screen *ebiten.Image, key string, data map[string]any, fontSize float64, opt *ebiten.DrawImageOptions) {
	txt := lang.ExecuteTemplate(key, data)
	DrawText(screen, txt, fontSize, opt)
}

func DrawTextByKey(screen *ebiten.Image, key string, fontSize float64, opt *ebiten.DrawImageOptions) {
	txt := lang.Text(key)
	DrawText(screen, txt, fontSize, opt)
}

func MeasureText(txt string, fontSize float64) geom.PointF {
	txtImg := textImage(txt, fontSize)
	pointI := txtImg.Bounds().Size()

	return geom.PointFFromPoint(pointI)
}

func DrawText(screen *ebiten.Image, txt string, fontSize float64, opt *ebiten.DrawImageOptions) {
	txtImg := textImage(txt, fontSize)
	screen.DrawImage(txtImg, opt)
}

func textImage(txt string, fontSize float64) *ebiten.Image {
	k := textKey{txt, fontSize}
	img, ok := textCache[k]
	if ok {
		return img
	}
	face := asset.FontFace(fontSize)
	lineSpacing := face.Metrics().VLineGap * 1.1
	w, h := text.Measure(txt, face, lineSpacing)
	shadeOffset := fontSize / 10

	img = ebiten.NewImage(int(w+shadeOffset+1), int(h+shadeOffset+1))
	opt := text.DrawOptions{}
	opt.LineSpacing = lineSpacing
	opt.ColorScale.ScaleWithColor(color.Black)
	opt.GeoM.Translate(shadeOffset, shadeOffset)
	text.Draw(img, txt, face, &opt)

	opt = text.DrawOptions{}
	opt.LineSpacing = lineSpacing
	text.Draw(img, txt, face, &opt)
	if textCache == nil {
		textCache = make(map[textKey]*ebiten.Image)
	}
	textCache[k] = img
	return img
}

type textKey struct {
	Text     string
	FontSize float64
}

var textCache map[textKey]*ebiten.Image

func init() {
	if textCache == nil {
		textCache = make(map[textKey]*ebiten.Image)
	}
}
