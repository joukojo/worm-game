# Game Structure

This project is a small terminal worm/snake game. The code is intentionally split so the rules of the game can be understood and tested without needing a real terminal.

## Package Layout

```text
.
├── cmd/worm-game/main.go     # executable entry point and terminal event loop
├── internal/game/model.go    # core game data types
├── internal/game/game.go     # game state transitions and collision rules
├── internal/game/input.go    # keyboard event to game command mapping
├── internal/game/render.go   # terminal drawing
└── internal/game/game_test.go # tests for game rules and rendering basics
```

## Responsibilities

### `cmd/worm-game/main.go`

`main.go` is the executable shell around the game. It creates and initializes the `tcell` screen, starts a goroutine for terminal events, creates the game from the current terminal size, and runs the main loop.

The loop has three modes:

- If the terminal is too small, it draws a resize message and waits for input or resize events.
- If the game is over, it draws the final state and waits for restart or quit.
- If the game is active, it reacts to key events and advances the game on each ticker tick.

This file deliberately knows about terminal size, timers, event polling, and process exit. It should not contain the detailed rules for movement, growth, food placement, or collision.

### `internal/game/model.go`

`model.go` defines the core state:

- `Point` is an explicit grid coordinate.
- `Direction` is a small enum for movement.
- `Game` stores board dimensions, worm body, current direction, queued direction, food position, score, game-over state, and random number generator.

The worm is represented as an ordered slice of `Point` values, with the head at index `0`. This makes movement easy to reason about: add the new head to the front, then remove the tail unless the worm ate food.

### `internal/game/game.go`

`game.go` owns the rules. `New` builds a starting game, `Step` advances one tick, `SetDirection` queues a direction change, and helper functions check walls, self-collision, and food spawning.

`Step` follows a fixed order:

1. Apply the queued direction.
2. Compute the next head position.
3. Check whether the move eats food.
4. Check wall and self-collision.
5. Add the new head.
6. Either grow and spawn food, or remove the tail.

That order matters. A normal move is allowed to enter the current tail cell because the tail moves away during the same tick. When growing, the tail does not move, so the same position must count as a collision. That is why `hitSelf` receives the `growing` flag.

### `internal/game/input.go`

`input.go` translates `tcell.EventKey` values into game commands:

- direction commands for `W`, `A`, `S`, `D` and arrow keys
- quit for `Q` or `Ctrl+C`
- restart for `R`

The rest of the program handles `Command` values instead of raw key events. This keeps key mapping in one place and makes it easier to change controls later.

### `internal/game/render.go`

`render.go` draws the current state into a `tcell.Screen`. It clears the screen, draws a header, board border, worm, food, and footer. The visual representation is intentionally plain:

- `@` for the head
- `o` for the body
- `*` for food
- `#` for the border

Rendering is separate from rule updates. Drawing the board should not change game state.

### `internal/game/game_test.go`

The tests focus on the behavior most likely to break:

- normal movement
- growth and scoring
- wall collision
- self-collision
- reverse-direction prevention
- food spawning only on free cells
- full-board behavior
- basic border rendering

Most tests exercise `internal/game` directly, which is possible because the game rules are not buried inside terminal event handling.

## Important Design Choices

### Why use `tcell`?

The game needs responsive single-key input without waiting for Enter, terminal resize events, screen clearing, and portable terminal cleanup. Doing all of that with only the standard library would require raw terminal mode handling and platform-specific details.

`tcell` is still a small terminal-focused dependency, so it improves input and rendering without turning the project into a large framework.

### Why separate rules from the executable loop?

The terminal loop is inherently stateful and tied to real time. Game rules are deterministic state transitions. Keeping them separate makes the important behavior testable with simple unit tests.

This also keeps `main.go` from becoming a large file where input handling, rendering, timing, and collision logic are mixed together.

### Why store both `Dir` and `NextDir`?

The worm moves on fixed ticks, but keyboard input can arrive between ticks. `NextDir` records the latest accepted direction, and `Dir` remains the direction used for the current movement.

This keeps movement deterministic: each `Step` applies at most one queued direction and then moves exactly once.

### Why prevent reverse movement in `SetDirection`?

Instantly reversing into the worm's own neck would create an unavoidable self-collision and feels bad in a snake-style game. Rejecting opposite directions at input time keeps invalid movement out of the game state.

### Why spawn food by building a list of free cells?

The board is small, so scanning every cell is simple and reliable. It guarantees food only appears on empty cells and avoids retry loops that can become inefficient or awkward when the board is almost full.

For a terminal game of this size, clarity is more valuable than a more complex data structure.

### Why use a deterministic random source in tests?

`Game` accepts a `*rand.Rand` so tests can provide a seeded generator. That makes food placement predictable when needed. The executable uses a time-based seed so normal gameplay is varied.

### Why resize the board to the terminal?

The game uses the available terminal size instead of a hard-coded board. This makes it feel natural in different terminals while still enforcing a minimum playable size.

When the terminal becomes too small, the game pauses behind a clear message instead of drawing a broken board.

## Extension Notes

Good places for future changes:

- Add difficulty settings by changing `DefaultTick` or making tick duration configurable.
- Add levels or fixed board sizes in `cmd/worm-game/main.go`.
- Add new visual styles in `render.go`.
- Add new commands in `input.go`.
- Add new rule behavior in `game.go`, covered by focused tests.

Keep the same boundary: terminal concerns in `cmd/worm-game` and rendering/input files, game rules in `game.go` and `model.go`.
