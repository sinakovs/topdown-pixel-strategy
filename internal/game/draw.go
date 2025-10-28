package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sinakovs/topdown-pixel-strategy/internal/assets"
)

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.BackgroundImg, nil)
	g.drawGrid(screen)
	g.drawUnit(screen)
	ebitenutil.DebugPrintAt(screen, "Click to move soldier", 5, 5)
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	lineColor := color.RGBA{50, 50, 50, 100}
	for i := 0; i <= screenWidth/gridSize; i++ {
		vector.StrokeLine(screen, float32(i*gridSize), 0, float32(i*gridSize), float32(screenHeight), 1, lineColor, false)
	}
	for j := 0; j <= screenHeight/gridSize; j++ {
		vector.StrokeLine(screen, 0, float32(j*gridSize), float32(screenWidth), float32(j*gridSize), 1, lineColor, false)
	}
}

func (g *Game) drawUnit(screen *ebiten.Image) {
	x, y := float64(g.unit.x*unitGridSize), float64(g.unit.y*unitGridSize)

	// Draw selection rectangle
	vector.StrokeRect(screen, float32(x), float32(y), unitGridSize, unitGridSize, 1, color.White, true)
	if g.unitActive {
		vector.StrokeRect(screen, float32(x), float32(y), unitGridSize, unitGridSize, 1, color.RGBA{255, 255, 0, 200}, true)
	}

	// Draw the unit image scaled to grid size
	opts := &ebiten.DrawImageOptions{}
	scaleX := float64(unitGridSize) / float64(assets.UnitImg.Bounds().Dx())
	scaleY := float64(unitGridSize) / float64(assets.UnitImg.Bounds().Dy())
	opts.GeoM.Scale(scaleX, scaleY)
	opts.GeoM.Translate(x, y)
	screen.DrawImage(assets.UnitImg, opts)
}
