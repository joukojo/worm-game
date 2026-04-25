# AGENTS.md

This repository is for a terminal-based classical worm/snake game implemented in Go.

## Goal

Build a playable terminal game where:

- The player controls the worm using `W`, `A`, `S`, `D`
- The worm moves continuously on a grid
- The worm grows when it eats food
- The game ends on collision with walls or the worm's own body
- The terminal UI is simple, responsive, and easy to run locally

## Technical Direction

- Language: Go
- Primary interface: terminal / TUI
- Input: single-key controls with `W`, `A`, `S`, `D`
- Keep dependencies minimal unless a library clearly improves terminal input/rendering

Preferred options:

1. Standard library first, if raw terminal input and rendering stay manageable
2. If needed, use a small, well-known terminal library such as `tcell`

## Implementation Guidelines

- Organize the code so game logic is separate from rendering and input handling
- Keep the game loop deterministic and easy to test
- Model the board as a grid with explicit width and height
- Represent the worm as an ordered list of positions
- Prevent instant reversal into the worm's own neck
- Render using plain terminal characters
- Clear and redraw the board each frame unless an incremental approach is clearly better
- Keep timing simple and configurable

## Recommended Project Structure

Use a structure close to this unless there is a strong reason to differ:

```text
.
├── cmd/worm-game/main.go
├── internal/game/
│   ├── game.go
│   ├── model.go
│   ├── input.go
│   └── render.go
├── go.mod
└── README.md
```

Notes:

- `game.go`: main loop and state transitions
- `model.go`: positions, direction, worm state, food, board dimensions
- `input.go`: terminal key mapping for `W`, `A`, `S`, `D`, quit key, and direction validation
- `render.go`: board drawing and HUD text such as score and game-over state

## Gameplay Requirements

- Start with a short worm near the center of the board
- Spawn food on an empty cell only
- Increase score by 1 per food eaten
- Support quitting with `Q` or `Ctrl+C`
- Show a clear game-over message and final score
- Allow restarting after game over if implementation stays simple

## Controls

- `W`: up
- `A`: left
- `S`: down
- `D`: right
- `Q`: quit

Controls should be case-insensitive if practical.

## Quality Bar

- Keep functions small and readable
- Prefer explicit state transitions over implicit side effects
- Add tests for core game logic where practical:
  - movement
  - growth
  - wall collision
  - self collision
  - food spawning rules
  - reverse-direction prevention

## Run Expectations

The project should be runnable with a standard Go workflow, ideally:

```bash
go run ./cmd/worm-game
```

## AI Agent Instructions

When making changes in this repository:

- Preserve a clean separation between game rules and terminal-specific code
- Do not introduce heavy frameworks for a small terminal game
- Prefer straightforward code over clever abstractions
- Update `README.md` when setup or controls change
- Add or update tests when modifying core game logic
- Keep the game playable in a normal terminal without requiring an IDE

If the repository is still empty, the first implementation should:

1. Initialize a Go module
2. Create a runnable terminal game skeleton
3. Implement the game loop, input handling, rendering, and collision logic
4. Add a short README with run instructions and controls
