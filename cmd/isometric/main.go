package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	// 640/32 =20 tiles
	// 480/32 =15 tiles
	screenWidth  = 640
	screenHeight = 480
	tileWidth    = 32
	tileHeight   = 32
)

type Game struct {
	groundTile    *ebiten.Image
	selectedTileX int
	selectedTileY int
}

type Drawable struct {
	img  *ebiten.Image
	opts *ebiten.DrawImageOptions
}

func (g *Game) Update() error {
	// Get current mouse position
	mouseX, mouseY := ebiten.CursorPosition()

	// Center offsets (must match Draw)
	centerX := float64(screenWidth) / 2
	centerY := float64(screenHeight) / 4

	// Convert screen coords to local map coords
	sx := float64(mouseX) - centerX
	sy := float64(mouseY) - centerY

	// Apply inverse isometric projection
	tileX := (sx/(tileWidth/2) + sy/(tileHeight/4)) / 2
	tileY := (sy/(tileHeight/4) - sx/(tileWidth/2)) / 2

	// Clamp or floor to integer tile coordinates
	g.selectedTileX = int(math.Floor(tileX))
	g.selectedTileY = int(math.Floor(tileY))

	return nil
}

// Isometric projection formula:
// https://www.youtube.com/watch?v=04oQ2jOUjkU
func (g *Game) Draw(screen *ebiten.Image) {
	// Calculate map size in tiles
	widthInTiles := screenWidth / tileWidth
	heightInTiles := screenHeight / tileHeight

	// Offset values to center the map on screen
	centerX := float64(screenWidth) / 2
	centerY := float64(screenHeight) / 4

	// Loop through tiles and compute their screen positions
	for tileY := 0; tileY < widthInTiles; tileY++ {
		for tileX := 0; tileX < heightInTiles; tileX++ {

			// --- Convert tile coordinates to screen coordinates (isometric projection)
			screenX := (float64(tileX)*0.5*tileWidth + float64(tileY)*-0.5*tileHeight) + centerX
			screenY := (float64(tileX)*0.25*tileWidth + float64(tileY)*0.25*tileHeight) + centerY

			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(screenX, screenY)

			// --- Highlight the selected tile by applying a color scale
			if tileX == g.selectedTileX && tileY == g.selectedTileY {
				opts.ColorScale.Scale(1.5, 1.5, 0.5, 1.0) // brighten yellowish
			}

			// --- Draw the base tile
			screen.DrawImage(g.groundTile, opts)
		}
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	img, _, err := ebitenutil.NewImageFromFile("cmd/isometric/tile_040.png")
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetFullscreen(true)
	game := &Game{
		groundTile: img,
	}

	err = ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err)
	}
}
