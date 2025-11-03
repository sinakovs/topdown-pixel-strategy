package gameserver

import (
	"context"
	"log/slog"

	"github.com/sinakovs/topdown-pixel-strategy/internal/server/client"
	clientlistener "github.com/sinakovs/topdown-pixel-strategy/internal/server/clientListener"
	"github.com/sinakovs/topdown-pixel-strategy/internal/server/world"
)

type gameserver struct {
	world          world.World
	clientListener clientlistener.ClientListener
}

func New(
	world world.World,
	clientListener clientlistener.ClientListener,
) *gameserver {
	return &gameserver{
		world:          world,
		clientListener: clientListener,
	}
}

func (g *gameserver) ListenAndServe() {
	g.world.Start()
	for {
		client, err := g.clientListener.Accept()
		if err != nil {
			slog.Error(
				"Failed to accept client",
				"error", err.Error(),
			)
		}

		go g.handleClient(client)
	}
}

// TODO: gracefully close clients on shutdown
func (g *gameserver) Shutdown() {
	g.world.Stop()
}

// TODO: spawn unit (inside Subscribe())
func (g *gameserver) handleClient(client client.Client) {
	playerID := client.PlayerID()
	updateChan := g.world.Subscribe(playerID)

	for range updateChan {
		ctx := context.Background()
		chunks, err := g.world.ChunksForPlayer(ctx, playerID)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"Failed to get chunks for player",
				"player_id", playerID,
				"error", err.Error(),
			)
			continue
		}

		err = client.SendSnapshot(chunks)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"Failed to send chunks to player",
				"player_id", playerID,
				"error", err.Error(),
			)
			continue
		}
	}
}
