package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewGame() *game {
	return &game{}
}

func (g *game) Update() error {
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

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
