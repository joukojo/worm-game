package game

import (
	"math/rand"
	"time"
)

const DefaultTick = 120 * time.Millisecond

func New(width, height int, rng *rand.Rand) *Game {
	if rng == nil {
		rng = rand.New(rand.NewSource(1))
	}

	centerX := width / 2
	centerY := height / 2

	g := &Game{
		Width:   width,
		Height:  height,
		Worm:    []Point{{X: centerX, Y: centerY}, {X: centerX - 1, Y: centerY}, {X: centerX - 2, Y: centerY}},
		Dir:     Right,
		NextDir: Right,
		rng:     rng,
	}
	g.spawnFood()
	return g
}

func (g *Game) Restart() *Game {
	return New(g.Width, g.Height, g.rng)
}

func (g *Game) Head() Point {
	return g.Worm[0]
}

func (g *Game) SetDirection(dir Direction) {
	if g.Over {
		return
	}
	if isOpposite(g.Dir, dir) {
		return
	}
	g.NextDir = dir
}

func (g *Game) Step() {
	if g.Over {
		return
	}

	g.Dir = g.NextDir
	head := g.Head()
	next := nextPoint(head, g.Dir)
	growing := next == g.Food

	if g.hitWall(next) || g.hitSelf(next, growing) {
		g.Over = true
		return
	}

	g.Worm = append([]Point{next}, g.Worm...)
	if growing {
		g.Score++
		g.spawnFood()
		return
	}

	g.Worm = g.Worm[:len(g.Worm)-1]
}

func isOpposite(a, b Direction) bool {
	return (a == Up && b == Down) ||
		(a == Down && b == Up) ||
		(a == Left && b == Right) ||
		(a == Right && b == Left)
}

func nextPoint(p Point, dir Direction) Point {
	switch dir {
	case Up:
		return Point{X: p.X, Y: p.Y - 1}
	case Down:
		return Point{X: p.X, Y: p.Y + 1}
	case Left:
		return Point{X: p.X - 1, Y: p.Y}
	default:
		return Point{X: p.X + 1, Y: p.Y}
	}
}

func (g *Game) hitWall(p Point) bool {
	return p.X < 0 || p.X >= g.Width || p.Y < 0 || p.Y >= g.Height
}

func (g *Game) hitSelf(next Point, growing bool) bool {
	limit := len(g.Worm)
	if !growing {
		limit--
	}
	for i := 0; i < limit; i++ {
		if g.Worm[i] == next {
			return true
		}
	}
	return false
}

func (g *Game) spawnFood() {
	free := make([]Point, 0, g.Width*g.Height-len(g.Worm))
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			p := Point{X: x, Y: y}
			if !contains(g.Worm, p) {
				free = append(free, p)
			}
		}
	}

	if len(free) == 0 {
		g.Over = true
		return
	}

	g.Food = free[g.rng.Intn(len(free))]
}

func contains(points []Point, target Point) bool {
	for _, p := range points {
		if p == target {
			return true
		}
	}
	return false
}
