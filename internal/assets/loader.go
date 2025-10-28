package assets

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	BackgroundImg *ebiten.Image
	UnitImg       *ebiten.Image
)

func init() {
	var err error
	BackgroundImg, _, err = ebitenutil.NewImageFromFile("assets/background/grass.png")
	if err != nil {
		slog.Error("Failed to load background image", "error", err)
	}

	UnitImg, _, err = ebitenutil.NewImageFromFile("assets/units/unit1.png")
	if err != nil {
		slog.Error("Failed to load unit image", "error", err)
	}
}
