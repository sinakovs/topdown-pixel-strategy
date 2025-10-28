package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

func (g *Game) handleMouseInput() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	x, y := ebiten.CursorPosition()
	g.target = Point{x: x / unitGridSize, y: y / unitGridSize}

	if g.unit == g.target {
		g.unitActive = !g.unitActive
	}
}
