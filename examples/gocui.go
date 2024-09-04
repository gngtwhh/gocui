package main

import (
	"bufio"
	"fmt"
	"github.com/gngtwhh/gocui/box"
	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/graph"
	"github.com/gngtwhh/gocui/progress_bar"
	"github.com/gngtwhh/gocui/window"
	"math"
	"os"
	"time"
)

func boxTest() {
	payload := []string{
		"          Books Management System",
		"",
		" 1.Store new books    2.New user registration",
		" 3.Borrow books       4.Return books",
		" 5.All books          6.All user",
		" 7.Delete database    8.Log out",
		"",
		"          Select operation number:",
	}
	window.ClearScreen()
	aBox, _ := box.GetBox(len(payload)+2, 50+2, "bold", payload)
	box.SetBoxAt(aBox, 0, 0)
}

func barTest() {
	window.ClearScreen()
	// test progress bar
	p, _ := progress_bar.NewProgressBar("[%bar] %current/%total-%percent %rate", func(p *progress_bar.Property) {
		p.Style.Complete = "#"
		p.Style.Incomplete = "-"
	})
	p.Run(time.Millisecond * 30)
	// wait
	<-p.Done

	time.Sleep(time.Second * 2)
	window.ClearScreen()

	// test uncertain progress bar
	up, _ := progress_bar.NewProgressBar("[%bar] testing ubar...%spinner", func(p *progress_bar.Property) {
		p.Uncertain = true
		p.Style.Incomplete = " "
		p.Style.UnCertain = "👈🤣👉"
	})
	up.Run(time.Millisecond * 100)

	// wait 3s
	time.Sleep(time.Second * 3)
	up.Stop()

	cursor.GotoXY(1, 0)
	fmt.Println("time out. exit...")
}

func lineTest() {
	drawCoord := func(x, y, length int) {
		graph.Line(x, y, length, '|', 0)
		cursor.GotoXY(x+length, y)
		fmt.Printf("v")
		length = 60
		graph.Line(x, y, length, '-', 1)
		cursor.GotoXY(x, y+length)
		fmt.Printf(">")
		cursor.GotoXY(x, y)
		fmt.Printf("+")
	}

	var x, y, length int

	window.ClearScreen()

	// curve f(x)=x^2
	x, y, length = 1, 2, 10
	drawCoord(x, y, length)
	f := func(x int) int {
		return int(math.Pow(float64(x), 2))
	}
	x, y, length = 1, 2, 10
	graph.Curve(x, y, length, 1, '*', f)

	//time.Sleep(time.Second * 2)
	window.ClearScreen()

	// curve f(x)=sin(x)
	x, y, length = 2, 30, 10
	drawCoord(x, y, length)
	f = func(x int) int {
		return int(math.Sin(float64(x)) * 5)
	}
	x, y, length = 2, 30, 10
	graph.Curve(x, y, length, 1, '*', f)

	//time.Sleep(time.Second * 2)
	window.ClearScreen()

	// curve f(x)=sqrt(3)*x
	x, y, length = 10, 30, 20
	drawCoord(x, y, length)
	f = func(x int) int {
		return int(math.Sqrt(3) * float64(x))
	}
	x, y, length = 10, 30, 15
	graph.Curve(x, y, length, 1, '*', f)
	graph.Curve(x, y, length, -1, '*', f)
}
func windowSizeTest() {
	w, h := window.GetConsoleSize()
	fmt.Printf("weight: %d, height: %d", w, h)
}

func main() {
	//c := '0'
	runList := []string{
		//"barTest",
		"boxTest",
		//"lineTest",
		"windowSizeTest",
	}
	funcs := map[string]func(){
		"barTest":        barTest,
		"lineTest":       lineTest,
		"boxTest":        boxTest,
		"windowSizeTest": windowSizeTest,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for _, s := range runList {
		if f, ok := funcs[s]; ok {
			f()
			scanner.Scan()
		}
	}
}
