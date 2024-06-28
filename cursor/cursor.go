package cursor

import (
	"fmt"
)

// GotoXY returns the escape sequence to move the cursor to the given position.
func GotoXY(x, y int) string {
	return fmt.Sprintf("\033[%d;%dH", x+1, y+1)
}

// Up returns the escape sequence to move the cursor up by n lines.
func Up(n int) string {
	return fmt.Sprintf("\033[%dA", n)
}

// Down returns the escape sequence to move the cursor down by n lines.
func Down(n int) string {
	return fmt.Sprintf("\033[%dB", n)
}

// Left returns the escape sequence to move the cursor left by n columns.
func Left(n int) string {
	return fmt.Sprintf("\033[%dD", n)
}

// Right returns the escape sequence to move the cursor right by n columns.
func Right(n int) string {
	return fmt.Sprintf("\033[%dC", n)
}

// HideCursor returns the escape sequence to hide the cursor.
func HideCursor() string {
	return "\033[?25l"
}

// ShowCursor returns the escape sequence to show the cursor.
func ShowCursor() string {
	return "\033[?25h"
}
