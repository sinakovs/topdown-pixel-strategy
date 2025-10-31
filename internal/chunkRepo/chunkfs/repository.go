package chunkfs

import "github.com/sinakovs/topdown-pixel-strategy/internal/assetsLoader"

type LayerType int

type MapChunk struct {
	Layers map[LayerType]*Layer
}

type Layer struct {
	Type LayerType
	Grid [][]assetsLoader.TileType
}
