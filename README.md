# worm-game

A terminal-based worm/snake game written in Go. It uses `tcell` for responsive single-key input and terminal rendering while keeping the game rules isolated in `internal/game`.

## Run

```bash
go run ./cmd/worm-game
```

## Controls

- `W` or Up Arrow: move up
- `A` or Left Arrow: move left
- `S` or Down Arrow: move down
- `D` or Right Arrow: move right
- `Q` or `Ctrl+C`: quit
- `R`: restart after game over

Controls are case-insensitive.

## Gameplay

- The worm starts near the center of the board and moves continuously.
- Eating `*` grows the worm and increases the score by `1`.
- The game ends if the worm hits a wall or its own body.
- Reverse moves into the worm's neck are ignored.
- Food only spawns on empty cells.

## Development

```bash
go test ./...
```

## Documentation

- [Game structure and design notes](docs/architecture.md)
