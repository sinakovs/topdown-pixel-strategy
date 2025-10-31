package chunkrepo

import (
	"context"
	"errors"
)

type LayerType uint8

const (
	LayerGround    = 0 // Base floor tiles
	LayerObjects   = 1 // Props, trees, walls
	LayerCollision = 2 // Not drawn; for logic
)

var ErrChunkNotFound = errors.New("chunk not found")

type ChunkKey struct {
	X uint32
	Y uint32
}

type Chunk struct {
	// Add chunk data fields here
}

type RenderChunk struct {
	Key   ChunkKey
	Chunk *Chunk
}

type ChunkRepository interface {
	Load(ctx context.Context, key ChunkKey) (*Chunk, error)
	Save(ctx context.Context, key ChunkKey, chunk *Chunk) error
	Close() error
}
