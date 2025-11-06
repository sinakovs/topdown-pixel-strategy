package clientconn

import (
	"net"
	"testing"

	chunkrepo "github.com/sinakovs/topdown-pixel-strategy/internal/chunkRepo"
	clientlistener "github.com/sinakovs/topdown-pixel-strategy/internal/server/clientListener"
)

func TestClientConn(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	go func() {
		server, err := clientlistener.HandleConn(serverConn)
		if err != nil {
			t.Error(err)
		}

		chunks := []chunkrepo.RenderChunk{
			chunkrepo.RenderChunk{
				Key: chunkrepo.ChunkKey{
					X: 1,
					Y: 1,
				},
			},
		}

		err = server.SendSnapshot(chunks)
		if err != nil {
			t.Error(err)
		}
	}()

	client, err := New(clientConn)
	if err != nil {
		t.Error(err)
	}

	a := <-client.ChunksChan()

	t.Error(a)

}
