package rpcservice

import (
	"github.com/sinakovs/topdown-pixel-strategy/internal/server/world"
)

type Client interface {
	GetSnapshot(_ PlayerArgs, snapshot *ChunksReply) error
	PlayerID() world.PlayerID
	// Close()
}
