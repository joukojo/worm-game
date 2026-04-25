package game

import "github.com/gdamore/tcell/v2"

type CommandType int

const (
	CommandDirection CommandType = iota
	CommandQuit
	CommandRestart
)

type Command struct {
	Type      CommandType
	Direction Direction
}

func CommandFromKeyEvent(ev *tcell.EventKey) (*Command, bool) {
	switch ev.Key() {
	case tcell.KeyCtrlC:
		return &Command{Type: CommandQuit}, true
	case tcell.KeyUp:
		return &Command{Type: CommandDirection, Direction: Up}, true
	case tcell.KeyDown:
		return &Command{Type: CommandDirection, Direction: Down}, true
	case tcell.KeyLeft:
		return &Command{Type: CommandDirection, Direction: Left}, true
	case tcell.KeyRight:
		return &Command{Type: CommandDirection, Direction: Right}, true
	}

	switch ev.Rune() {
	case 'w', 'W':
		return &Command{Type: CommandDirection, Direction: Up}, true
	case 's', 'S':
		return &Command{Type: CommandDirection, Direction: Down}, true
	case 'a', 'A':
		return &Command{Type: CommandDirection, Direction: Left}, true
	case 'd', 'D':
		return &Command{Type: CommandDirection, Direction: Right}, true
	case 'q', 'Q':
		return &Command{Type: CommandQuit}, true
	case 'r', 'R':
		return &Command{Type: CommandRestart}, true
	default:
		return nil, false
	}
}
