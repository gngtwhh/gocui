package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/font"
	"github.com/gngtwhh/gocui/graph"
	"github.com/gngtwhh/gocui/pb"
	"github.com/gngtwhh/gocui/window"
)

func boxTest() {
	return
	// payload := []string{
	// 	"",
	// 	" 1.Store new books    2.New user registration ",
	// 	" 3.Borrow books       4.Return books",
	// 	" 5.All books          6.All user",
	// 	" 7.Delete database    8.Log out",
	// 	"",
	// 	"          Select operation number:",
	// }
	// window.ClearScreen()

	// aBox, _ := box.NewBox(box.Property{
	// 	BoxType: "rounded",
	// 	Title:   "Books Management System",
	// 	Payload: payload,
	// })
	// box.SetBoxAt(aBox, 0, 0)
}

//func barTestOld() {
//	window.ClearScreen()
//	cursor.HideCursor()
//
//	// test progress bar
//	p, _ := pb.NewProgressBar("[%bar] %current/%total-%percent %rate", pb.property{
//		Style: pb.Style{
//			Complete:        "#",
//			Incomplete:      "-",
//			CompleteColor:   font.Green,
//			IncompleteColor: font.LightBlack,
//		},
//	})
//	for _ = range p.Iter() {
//		time.Sleep(time.Millisecond * 30)
//	}
//	//p.Run(time.Millisecond * 30)
//	//// wait
//	//<-p.Done
//
//	time.Sleep(time.Second * 2)
//	window.ClearScreen()
//
//	// test uncertain progress bar
//	up, _ := pb.NewProgressBar("[%bar] testing ubar...%spinner", pb.property{
//		Uncertain: true,
//		Style: pb.Style{
//			Incomplete: " ",
//			UnCertain:  "👈🤣👉",
//		},
//	})
//	up.Run(time.Millisecond * 100)
//
//	// wait 3s
//	time.Sleep(time.Second * 3)
//	up.Stop()
//
//	// test Default Bar
//	p = pb.DefaultBar
//	for i := range p.Iter() {
//		//fmt.Printf("i=%d\n", i)
//		fmt.Printf("%d", i)
//		time.Sleep(time.Millisecond * 100)
//	}
//	//p.Run(time.Millisecond * 100)
//	//// wait
//	//<-p.Done
//
//	cursor.GotoXY(1, 0)
//	fmt.Println("time out. exit...")
//}

func barTest() {
	cursor.HideCursor()
	window.ClearScreen()

	// test Default Bar
	p := pb.DefaultBar
	it, _ := p.Iter()
	for range it {
		//fmt.Printf("i=%d\n", i)
		time.Sleep(time.Millisecond * 50) // Simulate some time-consuming task
		// time.Sleep(time.Second * 10) // Simulate some time-consuming task
	}
	// test progress bar
	p, _ = pb.NewProgressBar("[%bar] %current/%total-%percent %rate", pb.WithPos(1, 0),
		pb.WithStyle(pb.Style{
			Complete:        "#",
			Incomplete:      "-",
			CompleteColor:   font.Green,
			IncompleteColor: font.LightBlack,
		}))
	it, _ = p.Iter()
	for range it {
		time.Sleep(time.Millisecond * 30) // Simulate some time-consuming task
	}

	//test uncertain progress bar
	up, _ := pb.NewProgressBar("[%bar] testing ubar...%spinner", pb.WithUncertain(), pb.WithPos(2, 0),
		pb.WithStyle(pb.Style{
			Incomplete: " ",
			UnCertain:  "👈🤣👉",
		}))
	stop := up.Run(time.Millisecond * 100)
	_, _ = up.UpdateProperty(pb.WithPos(5, 0))
	stop2 := up.Run(time.Millisecond * 200)
	time.Sleep(time.Second * 3) // Simulate a 3-second time-consuming task
	// close(stop2)
	stop2 <- struct{}{}
	time.Sleep(time.Second * 3)
	close(stop)

	cursor.GotoXY(3, 0)
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
	time.Sleep(time.Second * 2)

	window.ClearScreen()

	// curve f(x)=sin(x)
	x, y, length = 2, 30, 10
	drawCoord(x, y, length)
	f = func(x int) int {
		return int(math.Sin(float64(x)) * 5)
	}
	x, y, length = 2, 30, 10
	graph.Curve(x, y, length, 1, '*', f)
	time.Sleep(time.Second * 2)

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
	fmt.Printf("Command info: weight: %d, height: %d\n", w, h)
}

func FontTest() {
	fmt.Println("-----Testing text color-----")
	fmt.Println(font.Decorate("Black", font.Black))
	fmt.Println(font.Decorate("Red", font.Red))
	fmt.Println(font.Decorate("Green", font.Green))
	fmt.Println(font.Decorate("Yellow", font.Yellow))
	fmt.Println(font.Decorate("Blue", font.Blue))
	fmt.Println(font.Decorate("Magenta", font.Magenta))
	fmt.Println(font.Decorate("Cyan", font.Cyan))
	fmt.Println(font.Decorate("White", font.White))
	fmt.Println(font.Decorate("LightBlack", font.LightBlack))
	fmt.Println(font.Decorate("LightRed", font.LightRed))
	fmt.Println(font.Decorate("LightGreen", font.LightGreen))
	fmt.Println(font.Decorate("LightYellow", font.LightYellow))
	fmt.Println(font.Decorate("LightBlue", font.LightBlue))
	fmt.Println(font.Decorate("LightMagenta", font.LightMagenta))
	fmt.Println(font.Decorate("LightCyan", font.LightCyan))
	fmt.Println(font.Decorate("LightWhite", font.LightWhite))

	fmt.Println("-----Testing background style-----")
	fmt.Println(font.Decorate("BlackBg", font.BlackBg))
	fmt.Println(font.Decorate("RedBg", font.RedBg))
	fmt.Println(font.Decorate("GreenBg", font.GreenBg))
	fmt.Println(font.Decorate("YellowBg", font.YellowBg))
	fmt.Println(font.Decorate("BlueBg", font.BlueBg))
	fmt.Println(font.Decorate("MagentaBg", font.MagentaBg))
	fmt.Println(font.Decorate("CyanBg", font.CyanBg))
	fmt.Println(font.Decorate("WhiteBg", font.WhiteBg))
	fmt.Println(font.Decorate("LightBlackBg", font.LightBlackBg))
	fmt.Println(font.Decorate("LightRedBg", font.LightRedBg))
	fmt.Println(font.Decorate("LightGreenBg", font.LightGreenBg))
	fmt.Println(font.Decorate("LightYellowBg", font.LightYellowBg))
	fmt.Println(font.Decorate("LightBlueBg", font.LightBlueBg))
	fmt.Println(font.Decorate("LightMagentaBg", font.LightMagentaBg))
	fmt.Println(font.Decorate("LightCyanBg", font.LightCyanBg))
	fmt.Println(font.Decorate("LightWhiteBg", font.LightWhiteBg))

	fmt.Println("-----Testing font style-----")
	fmt.Println(font.Decorate("Bold", font.Bold))
	fmt.Println(font.Decorate("Dim", font.Dim))
	fmt.Println(font.Decorate("Italic", font.Italic))
	fmt.Println(font.Decorate("Underline", font.Underline))
	fmt.Println(font.Decorate("BlinkSlow", font.BlinkSlow))
	fmt.Println(font.Decorate("BlinkFast", font.BlinkFast))
	fmt.Println(font.Decorate("Reverse", font.Reverse))
	fmt.Println(font.Decorate("Hide", font.Hide))
}

func main() {
	//c := '0'
	runList := []string{
		"barTest",
		// "boxTest",
		// "lineTest",
		// "windowSizeTest",
		// "FontTest",
	}
	funcs := map[string]func(){
		"barTest":        barTest,
		"lineTest":       lineTest,
		"boxTest":        boxTest,
		"windowSizeTest": windowSizeTest,
		"FontTest":       FontTest,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for _, s := range runList {
		if f, ok := funcs[s]; ok {
			f()
			scanner.Scan()
			window.ClearScreen()
		}
	}

	cursor.ShowCursor()
}
