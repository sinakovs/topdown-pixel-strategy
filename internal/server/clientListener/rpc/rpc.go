package clientlistenerrpc

import (
	"fmt"
	"log/slog"
	"net"
	"net/rpc"

	rpcservice "github.com/sinakovs/topdown-pixel-strategy/internal/server/client/rpc"
	"github.com/sinakovs/topdown-pixel-strategy/internal/server/world"
)

func StartRPC(worldRepo world.World) error {

	listener, err := net.Listen("tcp", ":1235")
	if err != nil {
		return fmt.Errorf("listen tcp : %w", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error(
				"Failed to accept conn",
				"error", err,
			)
			continue
		}

		client, err := rpcservice.New(conn, worldRepo)
		if err != nil {
			slog.Error(
				"Failed to create new client",
				"error", err,
			)
			continue
		}

		server := rpc.NewServer()
		err = server.Register(client)
		if err != nil {
			return fmt.Errorf("register RPC : %w", err)
		}

		go server.ServeConn(conn)
	}
}
