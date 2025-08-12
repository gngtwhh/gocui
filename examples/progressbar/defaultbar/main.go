package main

import (
	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/pb"
	"time"
)

func main() {
	// test Default Bar
	cursor.HideCursor()
	p := pb.DefaultBar
	p.Iter(
		100, func() {
			// time.Sleep(time.Second * 5)
			time.Sleep(time.Millisecond * 50) // Simulate some time-consuming task
		},
	)
}
