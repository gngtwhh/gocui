package graph

import (
	"fmt"
	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/utils"
)

// Line Draws a line
// x, y - starting point
// length - length of the line
// ch - character to draw, 0 - horizontal, 1 - vertical
func Line(x, y, length int, ch rune, lineType uint8) {
	//utils.ConsoleMutex.Lock()

	utils.ConsoleMutex.Lock()

	if lineType == 0 {
		for i := 0; i < length; i++ {
			fmt.Printf("%s", cursor.GotoXY(x+i, y)+string(ch))
		}
	} else {
		for i := 0; i < length; i++ {
			fmt.Printf("%s", cursor.GotoXY(x, y+i)+string(ch))
		}
	}

	utils.ConsoleMutex.Unlock()
}

// Curve Draws a curve
// x, y - starting point
// length - length of the curve
// ch - character to draw
// f - function that returns the y coordinate by the x coordinate
func Curve(x, y, length, sign int, ch rune, f func(int) int) {
	utils.ConsoleMutex.Lock()

	//for i := 0; i < length; i++ {
	//	fmt.Printf("%s", cursor.GotoXY(x+i, y+f(i))+string(ch))
	//}

	if sign < 0 {
		for i := -length; i < 0; i++ {
			if x+i < 0 {
				continue
			}
			fmt.Printf("%s", cursor.GotoXY(x+i, y+f(i))+string(ch))
		}
	} else {
		for i := 0; i < length; i++ {
			fmt.Printf("%s", cursor.GotoXY(x+i, y+f(i))+string(ch))
		}
	}

	utils.ConsoleMutex.Unlock()
}
