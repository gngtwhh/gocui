package main

import (
	"time"

	"github.com/gngtwhh/gocui/pb"
)

func main() {
	// test Default Bar
	p := pb.DefaultBar
	p.Iter(100, func() {
		time.Sleep(time.Millisecond * 50) // Simulate some time-consuming task
	})
}
