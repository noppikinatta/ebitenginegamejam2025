package component

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type History struct {
	entries []string
	x, y    int
	width   int
	height  int
}

func NewHistory(x, y, width, height int) *History {
	return &History{
		entries: []string{
			"Kingdom Year 1000, Month 4: Game Started",
		},
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}

func (h *History) AddEntry(entry string) {
	h.entries = append(h.entries, entry)

	// Keep only the latest 15 entries to fit in the display area
	if len(h.entries) > 15 {
		h.entries = h.entries[len(h.entries)-15:]
	}
}

func (h *History) GetEntries() []string {
	return h.entries
}

func (h *History) Draw(screen *ebiten.Image) {
	// Draw background
	historyBg := ebiten.NewImage(h.width, h.height)
	historyBg.Fill(color.RGBA{50, 50, 50, 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(h.x), float64(h.y))
	screen.DrawImage(historyBg, op)

	// Draw title
	ebitenutil.DebugPrintAt(screen, "History", h.x+5, h.y+5)

	// Draw history entries
	for i, entry := range h.entries {
		if i < 15 { // Limit to 15 visible entries
			// Truncate long entries to fit the width
			displayEntry := entry
			if len(displayEntry) > 18 {
				displayEntry = displayEntry[:15] + "..."
			}
			ebitenutil.DebugPrintAt(screen, displayEntry, h.x+5, h.y+20+i*15)
		}
	}
}
