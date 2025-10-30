package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	// 640/32 =20 tiles
	// 480/32 =15 tiles
	screenWidth  = 1024
	screenHeight = 918
	tileWidth    = 32
	tileHeight   = 32
)

type Game struct {
	tileSet       map[int]*ebiten.Image
	layers        map[int][][]int
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

	// This fixes stuff
	// tileX -= 0.5
	// tileY += 0.5

	// Clamp or floor to integer tile coordinates
	g.selectedTileX = int(math.Floor(tileX))
	g.selectedTileY = int(math.Floor(tileY))

	return nil
}

// Isometric projection formula:
// https://www.youtube.com/watch?v=04oQ2jOUjkU
func (g *Game) Draw(screen *ebiten.Image) {
	// ✅ FIXED: Fill screen with a real color (gray background)
	screen.Fill(color.RGBA{20, 20, 20, 255})

	centerX := float64(screenWidth) / 2
	centerY := float64(screenHeight) / 4
	anchorX := float64(tileWidth) / 2
	anchorY := 0.0

	widthInTiles := len(g.layers[LayerGround][0])
	heightInTiles := len(g.layers[LayerGround])

	for tileY := 0; tileY < heightInTiles; tileY++ {
		for tileX := 0; tileX < widthInTiles; tileX++ {

			// --- Ground layer ---
			id := g.layers[LayerGround][tileY][tileX]
			if img, ok := g.tileSet[id]; ok {
				opts := &ebiten.DrawImageOptions{}
				screenX := (float64(tileX)*0.5*tileWidth + float64(tileY)*-0.5*tileHeight) + centerX
				screenY := (float64(tileX)*0.25*tileWidth + float64(tileY)*0.25*tileHeight) + centerY
				opts.GeoM.Translate(screenX-anchorX, screenY-anchorY)

				if tileX == g.selectedTileX && tileY == g.selectedTileY {
					opts.ColorScale.Scale(1.5, 1.5, 0.5, 1.0)
				}
				screen.DrawImage(img, opts)
			}

			// --- Object layer ---
			id = g.layers[LayerObjects][tileY][tileX]
			if img, ok := g.tileSet[id]; ok {
				opts := &ebiten.DrawImageOptions{}
				screenX := (float64(tileX)*0.5*tileWidth + float64(tileY)*-0.5*tileHeight) + centerX - 32  // Offset for object width
				screenY := (float64(tileX)*0.25*tileWidth + float64(tileY)*0.25*tileHeight) + centerY - 32 // Offset for object height
				opts.GeoM.Translate(screenX-anchorX, screenY-anchorY)
				screen.DrawImage(img, opts)
			}
		}
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	grassFlowImg, _, err := ebitenutil.NewImageFromFile("cmd/isometric/tiles/grass_flow.png")
	if err != nil {
		log.Fatal(err)
	}
	grassImg, _, err := ebitenutil.NewImageFromFile("cmd/isometric/tiles/grass.png")
	if err != nil {
		log.Fatal(err)
	}
	dirtImg, _, err := ebitenutil.NewImageFromFile("cmd/isometric/tiles/dirt.png")
	if err != nil {
		log.Fatal(err)
	}
	bushImg, _, err := ebitenutil.NewImageFromFile("cmd/isometric/tiles/bush.png")
	if err != nil {
		log.Fatal(err)
	}
	waterImg, _, err := ebitenutil.NewImageFromFile("cmd/isometric/tiles/water.png")
	if err != nil {
		log.Fatal(err)
	}
	rockImg, _, err := ebitenutil.NewImageFromFile("cmd/isometric/tiles/rock.png")
	if err != nil {
		log.Fatal(err)
	}

	groundLayer, err := LoadCSVLayer("cmd/isometric/layers/layer_ground.csv")
	if err != nil {
		log.Fatal(err)
	}
	objectsLayer, err := LoadCSVLayer("cmd/isometric/layers/layer_objects.csv")
	if err != nil {
		log.Fatal(err)
	}
	collisionLayer, err := LoadCSVLayer("cmd/isometric/layers/layer_collision.csv")
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetFullscreen(true)
	game := &Game{
		tileSet: map[int]*ebiten.Image{
			1: grassImg,
			2: waterImg,
			3: bushImg,
			4: dirtImg,
			5: rockImg,
			6: grassFlowImg,
		},
		layers: map[int][][]int{
			LayerGround:    groundLayer,
			LayerObjects:   objectsLayer,
			LayerCollision: collisionLayer,
		},
	}

	err = ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err)
	}
}
