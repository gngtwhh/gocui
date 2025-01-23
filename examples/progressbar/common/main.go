package main

import (
	"time"

	"github.com/gngtwhh/gocui/font"
	"github.com/gngtwhh/gocui/pb"
)

func main() {
	// test progress bar
	p, _ := pb.NewProgressBar("%spinner[%bar] %percent %rate [%elapsed]",
		pb.WithStyle(pb.Style{
			Complete:        ">",
			Incomplete:      "-",
			CompleteColor:   font.Green,
			IncompleteColor: font.LightBlack,
		}))
	it, _ := p.Iter()
	for range it {
		time.Sleep(time.Millisecond * 50) // Simulate some time-consuming task
	}
}
