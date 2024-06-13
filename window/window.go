package window

import (
	"fmt"
	"github.com/gngtwhh/gocui/cursor"
	"strings"
)

// ClearArea clears a rectangular area of the screen.
func ClearArea(x, y, width, height int) {
	for i := 0; i < height; i++ {
		fmt.Printf("%s", cursor.GotoXY(x+i, y)+strings.Repeat(" ", width))
	}
}

// ClearLine clears the line at the given row, or the current line if row is negative.
func ClearLine(row int) {
	if row < 0 {
		fmt.Printf("%s", "\033[s\033[K\033[u")
	} else {
		fmt.Printf("%s", "\033[s"+cursor.GotoXY(row, 0)+"\033[K\033[u")
	}
}

func ClearScreen() {
	fmt.Printf("%s", "\033[H\033[J")
}
