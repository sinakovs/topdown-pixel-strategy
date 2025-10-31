package chunkfs

import (
	"encoding/csv"
	"os"
	"strconv"
)

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
