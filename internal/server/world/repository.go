package world

import (
	"context"

	chunkrepo "github.com/sinakovs/topdown-pixel-strategy/internal/chunkRepo"
)

type World interface {
	Start()
	Stop()
	Subscribe(playerId PlayerID) <-chan struct{}
	Unsubscribe(playerId PlayerID)
	ChunksForPlayer(ctx context.Context, playerID PlayerID) ([]chunkrepo.RenderChunk, error)
}
