package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

type LayerType int

const (
	LayerGround    = 0 // Base floor tiles
	LayerObjects   = 1 // Props, trees, walls
	LayerCollision = 2 // Not drawn; for logic
)

type MapChunk struct {
	Width, Height int
	Layers        map[LayerType][][]int // Each layer is a 2D grid of tile IDs
}

type Map struct {
	width, height int
	tiles         [][]int
}

func LoadCSVLayer(path string) ([][]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	layer := make([][]int, len(records))
	for y, row := range records {
		layer[y] = make([]int, len(row))
		for x, val := range row {
			num, _ := strconv.Atoi(val)
			layer[y][x] = num
		}
	}

	return layer, nil
}
