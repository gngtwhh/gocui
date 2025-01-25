package window

import (
	"fmt"
	"strings"

	"github.com/gngtwhh/gocui/cursor"
)

// ClearArea clears a rectangular area of the screen.
func ClearArea(x, y, width, height int) {
	for i := 0; i < height; i++ {
		cursor.GotoXY(x+i, y)
		fmt.Printf("%s", strings.Repeat(" ", width))
	}
}

// ClearLine clears the line at the given row, or the current line if row is negative.
func ClearLine(row int) {
	if row < 0 {
		fmt.Print("\033[2K")
	} else {
		fmt.Print("\033[s")
		cursor.GotoXY(row, 0)
		fmt.Print("\033[2K\033[u")
	}
}

// ClearLineAfterCursor clear the content from the cursor position to the end of the line
func ClearLineAfterCursor() {
	fmt.Print("\033[K")
}

// ClearLineBeforCursor Clear the cursor position to the beginning of the line
func ClearLineBeforCursor() {
	fmt.Print("\033[1K")
}

func ClearScreen() {
	fmt.Printf("%s", "\033[H\033[J")
}
