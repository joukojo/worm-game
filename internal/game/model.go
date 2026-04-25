package game

import "math/rand"

type Point struct {
	X int
	Y int
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Game struct {
	Width   int
	Height  int
	Worm    []Point
	Dir     Direction
	NextDir Direction
	Food    Point
	Score   int
	Over    bool
	rng     *rand.Rand
}
