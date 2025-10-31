package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {


	// groundLayer, err := LoadCSVLayer("cmd/isometric/layers/layer_ground.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// objectsLayer, err := LoadCSVLayer("cmd/isometric/layers/layer_objects.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// collisionLayer, err := LoadCSVLayer("cmd/isometric/layers/layer_collision.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ebiten.SetFullscreen(true)
	game := &Game{
		tileSet:
		layers: map[int][][]int{
			LayerGround:    groundLayer,
			LayerObjects:   objectsLayer,
			LayerCollision: collisionLayer,
		},
	}

	err = ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err)
	}
}
