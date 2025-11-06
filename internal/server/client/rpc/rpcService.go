package rpcservice

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	chunkrepo "github.com/sinakovs/topdown-pixel-strategy/internal/chunkRepo"
	"github.com/sinakovs/topdown-pixel-strategy/internal/server/world"
)

var _ Client = (*ClientRPC)(nil)

// Args to request chunks for a player
type PlayerArgs struct {
}

// Reply containing chunks
type ChunksReply struct {
	Chunks []chunkrepo.RenderChunk // make sure Chunk is exported!
}

type ClientRPC struct {
	playerID world.PlayerID
	world    world.World
	subsChan <-chan struct{}
}

func (c *ClientRPC) GetSnapshot(_ PlayerArgs, snapshot *ChunksReply) error {
	<-c.subsChan
	ctx := context.Background()
	chunks, err := c.world.ChunksForPlayer(ctx, c.playerID)
	if err != nil {
		return fmt.Errorf("get chunks: %w", err)
	}

	snapshot.Chunks = chunks
	return nil
}

func (c *ClientRPC) PlayerID() world.PlayerID {
	return c.playerID
}

func New(conn net.Conn, worldRepo world.World) (*ClientRPC, error) {
	const magic = "RPCA"
	const version = byte(1)
	// ——— INSECURE handshake: MAGIC | VERSION | idLen(uint16) | idBytes ———

	// Read and check magic
	hdr := make([]byte, 4)

	_, err := io.ReadFull(conn, hdr)
	if err != nil {
		return nil, fmt.Errorf("failed to read magic: %w", err)
	}
	if string(hdr) != magic {
		return nil, fmt.Errorf("invalid magic: got %q, want %q", string(hdr), magic)
	}

	// Read and check version
	var ver [1]byte
	_, err = io.ReadFull(conn, ver[:])
	if err != nil {
		return nil, fmt.Errorf("failed to read version: %w", err)
	}
	if ver[0] != version {
		return nil, fmt.Errorf("unsupported version: %d", ver[0])
	}

	// Read client ID length
	var idLen uint16
	err = binary.Read(conn, binary.BigEndian, &idLen)
	if err != nil {
		return nil, fmt.Errorf("failed to read id length: %w", err)
	}
	if idLen == 0 || idLen > 1024 {
		return nil, fmt.Errorf("invalid id length: %d", idLen)
	}

	// Read client ID
	id := make([]byte, idLen)

	_, err = io.ReadFull(conn, id)
	if err != nil {
		return nil, fmt.Errorf("failed to read client id: %w", err)
	}

	playerID := world.PlayerID(id)

	subsChan := worldRepo.Subscribe(playerID)

	client := ClientRPC{
		playerID: playerID,
		subsChan: subsChan,
	}

	return &client, nil
}
