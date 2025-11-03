package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sinakovs/topdown-pixel-strategy/internal/assetsLoader"
)

const (
	// 640/32 =20 tiles
	// 480/32 =15 tiles
	screenWidth  = 1024
	screenHeight = 918
	tileWidth    = 32
	tileHeight   = 32
)

var _ ebiten.Game = (*game)(nil)

type game struct {
	tileProvider  *assetsLoader.TileProvider
	layers        map[int][][]int
	selectedTileX int
	selectedTileY int
}
