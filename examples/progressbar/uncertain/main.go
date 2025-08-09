package main

import (
	"fmt"
	"time"

	"github.com/gngtwhh/gocui/pb"
)

func main() {
	//test uncertain progress bar
	up, _ := pb.NewProgressBar("[%bar] waiting operation...%spinner", pb.WithUncertain(),
		pb.WithStyle(pb.Style{
			Incomplete: " ",
			UnCertain:  "👈🤣👉",
		}))
	stop := up.Run(time.Millisecond * 100)
	// Simulate a 3-second time-consuming task
	time.Sleep(time.Second * 3)
	close(stop)
	fmt.Printf("\ndone\n")
}
