package client

import (
	chunkrepo "github.com/sinakovs/topdown-pixel-strategy/internal/chunkRepo"
	"github.com/sinakovs/topdown-pixel-strategy/internal/server/world"
)

type Client interface {
	SendSnapshot(snapshot []chunkrepo.RenderChunk) error
	PlayerID() world.PlayerID
	Close()
}
