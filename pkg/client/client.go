package clientconn

import (
	"encoding/binary"
	"encoding/gob"
	"log/slog"
	"net"

	chunkrepo "github.com/sinakovs/topdown-pixel-strategy/internal/chunkRepo"
	"github.com/sinakovs/topdown-pixel-strategy/internal/server/world"
)

const magic = "RPCA"
const version = byte(1)

type client struct {
	chunksChan chan []chunkrepo.RenderChunk
}

func New(conn net.Conn) (*client, error) {
	// conn, err := net.Dial("tcp", "localhost:8080")
	// if err != nil {
	// 	return err
	// }

	err := handshake(conn, "a1")
	if err != nil {
		return nil, err
	}

	client := client{
		chunksChan: make(chan []chunkrepo.RenderChunk),
	}

	go func() {
		for {
			var decoded []chunkrepo.RenderChunk
			dec := gob.NewDecoder(conn)
			err = dec.Decode(&decoded)
			if err != nil {
				slog.Error(
					"Failed to decode chunks",
					"error", err.Error(),
				)
			}
			client.chunksChan <- decoded
		}
	}()

	return &client, nil
}

func (c *client) ChunksChan() <-chan []chunkrepo.RenderChunk {
	return c.chunksChan
}

func handshake(conn net.Conn, playerID world.PlayerID) error {

	// send prelude
	_, err := conn.Write([]byte(magic))
	if err != nil {
		conn.Close()
		return err
	}
	_, err = conn.Write([]byte{version})
	if err != nil {
		conn.Close()
		return err
	}

	var l [2]byte

	binary.BigEndian.PutUint16(l[:], uint16(len(playerID)))
	_, err = conn.Write(l[:])
	if err != nil {
		conn.Close()
		return err
	}
	_, err = conn.Write([]byte(playerID))
	if err != nil {
		conn.Close()
		return err
	}

	return nil
}
