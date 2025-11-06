package clientlistener

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/sinakovs/topdown-pixel-strategy/internal/server/client"
	"github.com/sinakovs/topdown-pixel-strategy/internal/server/client/custom"
	"github.com/sinakovs/topdown-pixel-strategy/internal/server/world"
)

var _ ClientListener = (*customListener)(nil)

type customListener struct {
	listener net.Listener
}

func New(listener net.Listener) *customListener {
	return &customListener{
		listener: listener,
	}
}

func (l *customListener) Accept() (client.Client, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, fmt.Errorf("accepted connection: %w", err)
	}

	return HandleConn(conn)
}

func HandleConn(conn net.Conn) (client.Client, error) {
	playerID, err := handshake(conn)
	if err != nil {
		return nil, fmt.Errorf("handshake with client: %w", err)
	}

	client := custom.New(conn, playerID)

	return client, nil
}

func handshake(conn net.Conn) (world.PlayerID, error) {
	const magic = "RPCA"
	const version = byte(1)
	// ——— INSECURE handshake: MAGIC | VERSION | idLen(uint16) | idBytes ———

	// Read and check magic
	hdr := make([]byte, 4)

	_, err := io.ReadFull(conn, hdr)
	if err != nil {
		return "", fmt.Errorf("failed to read magic: %w", err)
	}
	if string(hdr) != magic {
		return "", fmt.Errorf("invalid magic: got %q, want %q", string(hdr), magic)
	}

	// Read and check version
	var ver [1]byte
	_, err = io.ReadFull(conn, ver[:])
	if err != nil {
		return "", fmt.Errorf("failed to read version: %w", err)
	}
	if ver[0] != version {
		return "", fmt.Errorf("unsupported version: %d", ver[0])
	}

	// Read client ID length
	var idLen uint16
	err = binary.Read(conn, binary.BigEndian, &idLen)
	if err != nil {
		return "", fmt.Errorf("failed to read id length: %w", err)
	}
	if idLen == 0 || idLen > 1024 {
		return "", fmt.Errorf("invalid id length: %d", idLen)
	}

	// Read client ID
	id := make([]byte, idLen)

	_, err = io.ReadFull(conn, id)
	if err != nil {
		return "", fmt.Errorf("failed to read client id: %w", err)
	}

	return world.PlayerID(id), nil
}
