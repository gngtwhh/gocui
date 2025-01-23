package main

import (
	"fmt"
	"time"

	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/window"
)

func main() {
	window.ClearScreen()
	cursor.GotoXY(0, 0)
	fmt.Print("\n1234567")
	time.Sleep(time.Second)
	window.ClearLine(1)
	fmt.Println("abcdefg")
	fmt.Print("123456789")
}
