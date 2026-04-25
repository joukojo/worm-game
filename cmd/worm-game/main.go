package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"

	"worm-game/internal/game"
)

const (
	minBoardWidth  = 10
	minBoardHeight = 6
	statusLines    = 4
	borderColumns  = 2
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create screen: %v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	g, boardErr := newGameForScreen(screen, rng)

	events := make(chan tcell.Event, 16)
	go pollScreenEvents(screen, events)

	ticker := time.NewTicker(game.DefaultTick)
	defer ticker.Stop()

	for {
		minScreenWidth, minScreenHeight := requiredScreenSize(g)
		if boardErr != nil {
			game.DrawTooSmall(screen, boardErr.Error(), minScreenWidth, minScreenHeight)
		} else {
			game.Draw(screen, g)
		}
		screen.Show()

		if boardErr != nil {
			ev, ok := <-events
			if !ok {
				return
			}
			if handleEvent(screen, ev, &g, rng, &boardErr) {
				return
			}
			continue
		}

		if g.Over {
			ev, ok := <-events
			if !ok {
				return
			}
			if handleEvent(screen, ev, &g, rng, &boardErr) {
				return
			}
			continue
		}

		select {
		case ev, ok := <-events:
			if !ok {
				return
			}
			if handleEvent(screen, ev, &g, rng, &boardErr) {
				return
			}
		case <-ticker.C:
			g.Step()
		}
	}
}

func pollScreenEvents(screen tcell.Screen, out chan<- tcell.Event) {
	defer close(out)
	for {
		ev := screen.PollEvent()
		if ev == nil {
			return
		}
		out <- ev
	}
}

func handleEvent(screen tcell.Screen, ev tcell.Event, g **game.Game, rng *rand.Rand, boardErr *error) bool {
	switch typed := ev.(type) {
	case *tcell.EventResize:
		screen.Sync()
		if *g == nil {
			*g, *boardErr = newGameForScreen(screen, rng)
			return false
		}
		*boardErr = gameFitsScreen(screen, *g)
	case *tcell.EventKey:
		cmd, ok := game.CommandFromKeyEvent(typed)
		if !ok {
			return false
		}
		if cmd.Type == game.CommandQuit {
			return true
		}
		if *boardErr == nil {
			handleCommand(cmd, g, screen, rng, boardErr)
		}
	}
	return false
}

func handleCommand(cmd *game.Command, g **game.Game, screen tcell.Screen, rng *rand.Rand, boardErr *error) {
	switch cmd.Type {
	case game.CommandRestart:
		if *g != nil && *boardErr == nil {
			*g = (*g).Restart()
			return
		}
		restarted, err := newGameForScreen(screen, rng)
		*boardErr = err
		if err == nil {
			*g = restarted
		}
	case game.CommandDirection:
		(*g).SetDirection(cmd.Direction)
	}
}

func newGameForScreen(screen tcell.Screen, rng *rand.Rand) (*game.Game, error) {
	width, height := screen.Size()
	boardWidth := width - borderColumns
	boardHeight := height - statusLines
	if boardWidth < minBoardWidth || boardHeight < minBoardHeight {
		return nil, fmt.Errorf(
			"terminal too small: need at least %dx%d, have %dx%d",
			minBoardWidth+borderColumns,
			minBoardHeight+statusLines,
			width,
			height,
		)
	}

	return game.New(boardWidth, boardHeight, rng), nil
}

func gameFitsScreen(screen tcell.Screen, g *game.Game) error {
	width, height := screen.Size()
	minWidth, minHeight := requiredScreenSize(g)
	if width < minWidth || height < minHeight {
		return fmt.Errorf(
			"terminal too small: need at least %dx%d, have %dx%d",
			minWidth,
			minHeight,
			width,
			height,
		)
	}
	return nil
}

func requiredScreenSize(g *game.Game) (int, int) {
	if g == nil {
		return minBoardWidth + borderColumns, minBoardHeight + statusLines
	}
	return g.Width + borderColumns, g.Height + statusLines
}
