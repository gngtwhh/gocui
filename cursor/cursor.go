package cursor

import (
	"fmt"
)

//var (
//	csi = "\033["
//)

// GotoXY returns the escape sequence to move the cursor to the given position.
func GotoXY(x, y int) {
	fmt.Printf("\033[%d;%dH", x+1, y+1)
}

// Up returns the escape sequence to move the cursor up by n lines.
func Up(n int) {
	fmt.Printf("\033[%dA", n)
}

// Down returns the escape sequence to move the cursor down by n lines.
func Down(n int) {
	fmt.Printf("\033[%dB", n)
}

// Left returns the escape sequence to move the cursor left by n columns.
func Left(n int) {
	fmt.Printf("\033[%dD", n)
}

// Right returns the escape sequence to move the cursor right by n columns.
func Right(n int) {
	fmt.Printf("\033[%dC", n)
}

// HideCursor returns the escape sequence to hide the cursor.
func HideCursor() {
	fmt.Printf("\033[?25l")
}

// ShowCursor returns the escape sequence to show the cursor.
func ShowCursor() {
	fmt.Printf("\033[?25h")
}
