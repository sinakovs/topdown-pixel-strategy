package main

import (
	"image/color"
	_ "image/png"
	"log/slog"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	backgroundImg *ebiten.Image
	unitImg       *ebiten.Image
)

const (
	screenWidth   = 640
	screenHeight  = 480
	windowWidth   = 1024
	windowHeight  = 1024
	gridSize      = 20
	unitGridSize  = gridSize * 2
	gameFPS       = 6
	gameSpeed     = time.Second / gameFPS
	gridLineColor = 0x32323264 // RGBA in hex
)

type Point struct{ X, Y int }

type Game struct {
	unit       Point
	target     Point
	unitActive bool
	lastUpdate time.Time
}

func init() {
	var err error
	backgroundImg, _, err = ebitenutil.NewImageFromFile("assets/background/grass.png")
	if err != nil {
		slog.Error("Failed to load background image", "error", err)
	}

	unitImg, _, err = ebitenutil.NewImageFromFile("assets/units/unit1.png")
	if err != nil {
		slog.Error("Failed to load unit image", "error", err)
	}
}

func (g *Game) Update() error {
	g.handleMouseInput()
	if time.Since(g.lastUpdate) < gameSpeed {
		return nil
	}
	
	g.lastUpdate = time.Now()

	if g.unitActive {
		g.moveUnitTowardsTarget()
	}

	return nil
}

func (g *Game) handleMouseInput() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	x, y := ebiten.CursorPosition()
	g.target = Point{X: x / unitGridSize, Y: y / unitGridSize}

	if g.unit == g.target {
		g.unitActive = !g.unitActive
	}
}

func (g *Game) moveUnitTowardsTarget() {
	if g.unit.X < g.target.X {
		g.unit.X++
	} else if g.unit.X > g.target.X {
		g.unit.X--
	}

	if g.unit.Y < g.target.Y {
		g.unit.Y++
	} else if g.unit.Y > g.target.Y {
		g.unit.Y--
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(backgroundImg, nil)
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
	x, y := float64(g.unit.X*unitGridSize), float64(g.unit.Y*unitGridSize)

	// Draw selection rectangle
	vector.StrokeRect(screen, float32(x), float32(y), unitGridSize, unitGridSize, 1, color.White, true)
	if g.unitActive {
		vector.StrokeRect(screen, float32(x), float32(y), unitGridSize, unitGridSize, 1, color.RGBA{255, 255, 0, 200}, true)
	}

	// Draw the unit image scaled to grid size
	opts := &ebiten.DrawImageOptions{}
	scaleX := float64(unitGridSize) / float64(unitImg.Bounds().Dx())
	scaleY := float64(unitGridSize) / float64(unitImg.Bounds().Dy())
	opts.GeoM.Scale(scaleX, scaleY)
	opts.GeoM.Translate(x, y)
	screen.DrawImage(unitImg, opts)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	initial := screenWidth / unitGridSize / 2
	game := &Game{
		unit:   Point{X: initial, Y: initial},
		target: Point{X: initial, Y: initial},
	}

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Grid-based Medieval Soldier")

	if err := ebiten.RunGame(game); err != nil {
		slog.Error("Game exited with error", "error", err)
	}
}
