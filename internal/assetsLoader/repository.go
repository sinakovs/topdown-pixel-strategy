package assetsLoader

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type TileType uint8

const (
	TileGrass TileType = iota
	TileWater
	TileBush
	TileDirt
	TileRock
	TileGrassFlow
)

type TileProvider interface {
	GetTile(tileType TileType) (*ebiten.Image, error)
}

func (t TileType) String() string {
	switch t {
	case TileGrass:
		return "Grass"
	case TileWater:
		return "Water"
	case TileBush:
		return "Bush"
	case TileDirt:
		return "Dirt"
	case TileRock:
		return "Rock"
	case TileGrassFlow:
		return "GrassFlow"
	default:
		return fmt.Sprintf("UnknownTile(%d)", t)
	}
}
