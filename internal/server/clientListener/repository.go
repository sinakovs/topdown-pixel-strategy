package clientlistener

import "github.com/sinakovs/topdown-pixel-strategy/internal/server/client"

type ClientListener interface {
	Accept() (client.Client, error)
}
