package assetsLoader

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var _ TileProvider = (*tileProvider)(nil)

var ErrTileNotFound = fmt.Errorf("tile not found")

type tileProvider struct {
	tileSet map[TileType]*ebiten.Image
}

// GetTile implements TileProvider.
func (t *tileProvider) GetTile(tileType TileType) (*ebiten.Image, error) {
	img, ok := t.tileSet[tileType]
	if !ok {
		return nil, fmt.Errorf("find tile image of type %s: %w", tileType, ErrTileNotFound)
	}

	return img, nil
}

func NewTileProvider(assetFolderPath string) (*tileProvider, error) {
	provider := &tileProvider{
		tileSet: make(map[TileType]*ebiten.Image),
	}

	load := func(filename string) (*ebiten.Image, error) {
		path := assetFolderPath + "/" + filename
		img, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to load tile image %q: %w", filename, err)
		}
		return img, nil
	}

	img, err := load("grass_flow.png")
	if err != nil {
		return nil, err
	}
	provider.tileSet[TileGrassFlow] = img

	img, err = load("grass.png")
	if err != nil {
		return nil, err
	}
	provider.tileSet[TileGrass] = img

	img, err = load("dirt.png")
	if err != nil {
		return nil, err
	}
	provider.tileSet[TileDirt] = img

	img, err = load("bush.png")
	if err != nil {
		return nil, err
	}
	provider.tileSet[TileBush] = img

	img, err = load("water.png")
	if err != nil {
		return nil, err
	}
	provider.tileSet[TileWater] = img

	img, err = load("rock.png")
	if err != nil {
		return nil, err
	}
	provider.tileSet[TileRock] = img

	return provider, nil
}
