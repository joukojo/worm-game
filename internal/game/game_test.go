package game

import (
	"math/rand"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestStepMovesForward(t *testing.T) {
	g := New(20, 10, rand.New(rand.NewSource(1)))
	head := g.Head()

	g.Step()

	if got := g.Head(); got != (Point{X: head.X + 1, Y: head.Y}) {
		t.Fatalf("expected head to move right to %+v, got %+v", Point{X: head.X + 1, Y: head.Y}, got)
	}
	if len(g.Worm) != 3 {
		t.Fatalf("expected worm length 3, got %d", len(g.Worm))
	}
}

func TestEatingFoodGrowsWorm(t *testing.T) {
	g := New(20, 10, rand.New(rand.NewSource(1)))
	g.Food = Point{X: g.Head().X + 1, Y: g.Head().Y}

	g.Step()

	if g.Score != 1 {
		t.Fatalf("expected score 1, got %d", g.Score)
	}
	if len(g.Worm) != 4 {
		t.Fatalf("expected worm length 4, got %d", len(g.Worm))
	}
}

func TestWallCollisionEndsGame(t *testing.T) {
	g := New(4, 4, rand.New(rand.NewSource(1)))
	g.Worm = []Point{{X: 3, Y: 1}, {X: 2, Y: 1}, {X: 1, Y: 1}}
	g.Dir = Right
	g.NextDir = Right

	g.Step()

	if !g.Over {
		t.Fatal("expected game over after wall collision")
	}
}

func TestSelfCollisionEndsGame(t *testing.T) {
	g := New(8, 8, rand.New(rand.NewSource(1)))
	g.Worm = []Point{
		{X: 3, Y: 3},
		{X: 3, Y: 4},
		{X: 2, Y: 4},
		{X: 2, Y: 3},
		{X: 2, Y: 2},
		{X: 3, Y: 2},
	}
	g.Dir = Up
	g.NextDir = Left

	g.Step()

	if !g.Over {
		t.Fatal("expected game over after self collision")
	}
}

func TestReverseDirectionIsIgnored(t *testing.T) {
	g := New(20, 10, rand.New(rand.NewSource(1)))

	g.SetDirection(Left)
	g.Step()

	if g.Dir != Right {
		t.Fatalf("expected direction to remain right, got %v", g.Dir)
	}
}

func TestQueuedReverseDirectionIsIgnored(t *testing.T) {
	g := New(20, 10, rand.New(rand.NewSource(1)))

	g.SetDirection(Up)
	g.SetDirection(Left)
	g.Step()

	if g.Dir != Up {
		t.Fatalf("expected queued reverse turn to be ignored and direction to be up, got %v", g.Dir)
	}
	if got := g.Head(); got != (Point{X: 10, Y: 4}) {
		t.Fatalf("expected head to move up to %+v, got %+v", Point{X: 10, Y: 4}, got)
	}
}

func TestFoodNeverSpawnsOnWorm(t *testing.T) {
	g := New(3, 2, rand.New(rand.NewSource(1)))
	g.Worm = []Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 2, Y: 0},
		{X: 0, Y: 1},
		{X: 1, Y: 1},
	}

	g.spawnFood()

	if g.Food != (Point{X: 2, Y: 1}) {
		t.Fatalf("expected only free cell to be selected, got %+v", g.Food)
	}
}

func TestEatingLastFreeCellEndsRound(t *testing.T) {
	g := New(2, 2, rand.New(rand.NewSource(1)))
	g.Worm = []Point{
		{X: 1, Y: 0},
		{X: 0, Y: 0},
		{X: 0, Y: 1},
	}
	g.Dir = Down
	g.NextDir = Down
	g.Food = Point{X: 1, Y: 1}

	g.Step()

	if !g.Over {
		t.Fatal("expected game to end when the board is full")
	}
	if g.Score != 1 {
		t.Fatalf("expected score 1, got %d", g.Score)
	}
	if len(g.Worm) != 4 {
		t.Fatalf("expected worm length 4, got %d", len(g.Worm))
	}
}

func TestDrawPlacesBordersAtBoardEdges(t *testing.T) {
	g := New(10, 6, rand.New(rand.NewSource(1)))
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("failed to init simulation screen: %v", err)
	}
	screen.SetSize(g.Width+2, g.Height+4)

	Draw(screen, g)

	assertRuneAt(t, screen, 0, 1, '#')
	assertRuneAt(t, screen, g.Width+1, 1, '#')
	assertRuneAt(t, screen, 0, g.Height+2, '#')
	assertRuneAt(t, screen, g.Width+1, g.Height+2, '#')
}

func assertRuneAt(t *testing.T, screen tcell.Screen, x, y int, want rune) {
	t.Helper()
	got, _, _, _ := screen.GetContent(x, y)
	if got != want {
		t.Fatalf("expected %q at (%d,%d), got %q", want, x, y, got)
	}
}
