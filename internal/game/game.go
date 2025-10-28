package game

import "time"

type Point struct{ x, y int }

type Game struct {
	unit       Point
	target     Point
	unitActive bool
	lastUpdate time.Time
}
