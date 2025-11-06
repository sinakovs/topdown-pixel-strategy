package custom

import (
	"encoding/gob"
	"fmt"
	"net"

	chunkrepo "github.com/sinakovs/topdown-pixel-strategy/internal/chunkRepo"
	"github.com/sinakovs/topdown-pixel-strategy/internal/server/client"
	"github.com/sinakovs/topdown-pixel-strategy/internal/server/world"
)

var _ client.Client = (*CustomClient)(nil)

type CustomClient struct {
	conn     net.Conn
	playerID world.PlayerID
}

func New(conn net.Conn, playerID world.PlayerID) *CustomClient {
	return &CustomClient{
		conn:     conn,
		playerID: playerID,
	}
}

// Close implements client.Client.
func (c *CustomClient) Close() {
	panic("unimplemented")
}

// PlayerID implements client.Client.
func (c *CustomClient) PlayerID() world.PlayerID {
	return c.playerID
}

// SendSnapshot implements client.Client.
func (c *CustomClient) SendSnapshot(snapshot []chunkrepo.RenderChunk) error {
	enc := gob.NewEncoder(c.conn)
	err := enc.Encode(snapshot)
	if err != nil {
		return fmt.Errorf("encode snapshot: %w", err)
	}

	return nil
}
