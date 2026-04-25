package game

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func Draw(screen tcell.Screen, g *Game) {
	screen.Clear()
	style := tcell.StyleDefault
	borderStyle := style.Foreground(tcell.ColorWhite)
	headStyle := style.Foreground(tcell.ColorGreen)
	bodyStyle := style.Foreground(tcell.ColorOlive)
	foodStyle := style.Foreground(tcell.ColorRed)
	textStyle := style.Foreground(tcell.ColorWhite)

	header := fmt.Sprintf("Worm Game  Score:%d  Move:WASD  Quit:Q", g.Score)
	if g.Over {
		header += "  Restart:R"
	}
	drawText(screen, 0, 0, fitLine(header, g.Width+2), textStyle)
	drawHorizontalBorder(screen, 0, 1, g.Width+2, borderStyle)

	for y := 0; y < g.Height; y++ {
		screen.SetContent(0, y+2, '#', nil, borderStyle)
		for x := 0; x < g.Width; x++ {
			p := Point{X: x, Y: y}
			switch {
			case p == g.Head():
				screen.SetContent(x+1, y+2, '@', nil, headStyle)
			case p == g.Food:
				screen.SetContent(x+1, y+2, '*', nil, foodStyle)
			case contains(g.Worm[1:], p):
				screen.SetContent(x+1, y+2, 'o', nil, bodyStyle)
			default:
				screen.SetContent(x+1, y+2, ' ', nil, style)
			}
		}
		screen.SetContent(g.Width+1, y+2, '#', nil, borderStyle)
	}

	drawHorizontalBorder(screen, 0, g.Height+2, g.Width+2, borderStyle)

	footer := "Eat food, avoid walls, and do not bite yourself."
	if g.Over {
		footer = fmt.Sprintf("Game over. Final score: %d. Press R to restart or Q to quit.", g.Score)
	}
	drawText(screen, 0, g.Height+3, fitLine(footer, g.Width+2), textStyle)
}

func fitLine(text string, width int) string {
	if width <= 0 {
		return ""
	}
	if len(text) <= width {
		return text
	}
	if width <= 3 {
		return text[:width]
	}
	return text[:width-3] + "..."
}

func DrawTooSmall(screen tcell.Screen, message string, minWidth, minHeight int) {
	screen.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	drawText(screen, 0, 0, "Terminal too small for Worm Game.", style)
	drawText(screen, 0, 1, fitLine(message, max(1, screenWidth(screen))), style)
	drawText(screen, 0, 2, fmt.Sprintf("Resize to at least %dx%d.", minWidth, minHeight), style)
	drawText(screen, 0, 3, "Press Q to quit.", style)
}

func drawHorizontalBorder(screen tcell.Screen, x, y, width int, style tcell.Style) {
	for i := 0; i < width; i++ {
		screen.SetContent(x+i, y, '#', nil, style)
	}
}

func drawText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, r := range text {
		screen.SetContent(x+i, y, r, nil, style)
	}
}

func screenWidth(screen tcell.Screen) int {
	width, _ := screen.Size()
	return width
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
