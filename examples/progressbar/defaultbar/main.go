package main

import (
	"time"

	"github.com/gngtwhh/gocui/pb"
)

func main() {
	// test Default Bar
	p := pb.DefaultBar
	it, _ := p.Iter()
	for range it {
		time.Sleep(time.Millisecond * 50) // Simulate some time-consuming task
	}
}
