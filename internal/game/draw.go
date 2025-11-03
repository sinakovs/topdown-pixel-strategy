package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	chunkrepo "github.com/sinakovs/topdown-pixel-strategy/internal/chunkRepo"
)

type Drawable struct {
	img  *ebiten.Image
	opts *ebiten.DrawImageOptions
}

// Isometric projection formula:
// https://www.youtube.com/watch?v=04oQ2jOUjkU
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 20, 255}) // background

	centerX := float64(screenWidth) / 2
	centerY := float64(screenHeight) / 4
	anchorX := float64(tileWidth) / 2
	widthInTiles := len(g.layers[chunkrepo.LayerGround][0])
	heightInTiles := len(g.layers[chunkrepo.LayerGround])

	for y := 0; y < heightInTiles; y++ {
		for x := 0; x < widthInTiles; x++ {
			g.drawTile(screen, x, y, centerX, centerY, anchorX, false) // ground
			g.drawTile(screen, x, y, centerX, centerY, anchorX, true)  // object
		}
	}
}

// drawTile renders either a ground or object tile at (x, y).
func (g *game) drawTile(
	screen *ebiten.Image,
	x, y int,
	centerX, centerY, anchorX float64,
	isObject bool,
) {
	layer := chunkrepo.LayerGround
	yOffset := 0.0

	if isObject {
		layer = chunkrepo.LayerObjects
		yOffset = -0.5 * float64(tileHeight) // same as “-16” magic number
	}

	id := g.layers[layer][y][x]
	img, ok := g.tileSet[id]
	if !ok || img == nil {
		return
	}

	screenX := (float64(x)*0.5*tileWidth + float64(y)*-0.5*tileHeight) + centerX
	screenY := (float64(x)*0.25*tileWidth + float64(y)*0.25*tileHeight) + centerY + yOffset

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(screenX-anchorX, screenY)

	// Highlight selected tile
	if x == g.selectedTileX && y == g.selectedTileY {
		if isObject {
			opts.ColorScale.Scale(1.5, 0.5, 0.5, 1.0) // reddish highlight
		} else {
			opts.ColorScale.Scale(1.5, 1.5, 0.5, 1.0) // yellowish highlight
		}
	}

	screen.DrawImage(img, opts)
}
