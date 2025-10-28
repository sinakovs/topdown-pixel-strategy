package game

import "time"

const (
	screenWidth   = 640
	screenHeight  = 480
	windowWidth   = 1024
	windowHeight  = 1024
	gridSize      = 20
	unitGridSize  = gridSize * 2
	gameFPS       = 60
	gameSpeed     = time.Second / gameFPS
	gridLineColor = 0x32323264 // RGBA in hex
)
