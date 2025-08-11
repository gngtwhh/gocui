package main

import (
	"time"

	"github.com/gngtwhh/gocui/font"
	"github.com/gngtwhh/gocui/pb"
)

type Spinnerbar struct {
	cur int
}

func (sb *Spinnerbar) ToString(ctx *pb.Context) string {
	res := []string{"⠧⠤⠴", "⠯⠥⠄", "⠯⠍⠁", "⠏⠉⠙", "⠉⠉⠽", "⠀⠭⠽", "⠤⠤⠽"}[sb.cur]
	res = font.Decorate(res, font.Red+sb.cur)
	sb.cur = (sb.cur + 1) % 7
	return res
}

func main() {
	// test plugin mode
	pb.RegisterToken("%spinbar", &Spinnerbar{})
	pluginBar, _ := pb.NewProgressBar("%spinbar[%bar]%percent|%elapsed", pb.WithPos(1, 0),
		pb.WithStyle(pb.Style{
			Complete:        "#",
			Incomplete:      "-",
			CompleteColor:   font.Green,
			IncompleteColor: font.LightBlack,
		}))
	pluginBar.Iter(100, func() {
		time.Sleep(time.Millisecond * 50) // Simulate some time-consuming task
	})
}
