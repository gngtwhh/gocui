package main

import (
	"fmt"
	"github.com/gngtwhh/gocui/box"
	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/graph"
	"github.com/gngtwhh/gocui/progress_bar"
	"github.com/gngtwhh/gocui/window"
	"math"
	"time"
)

func boxTest() {
	/*payload := []string{
		"                       图书管理系统        ",
		"",
		"                1.采编入库     2.添加用户   ",
		"                3.借阅图书     4.归还图书   ",
		"                5.所有图书     6.所有用户   ",
		"                7.删除文件     8.退出系统   ",
		"",
		"                    请输入编号:",
	}*/
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
	p, _ := progress_bar.NewProgressBar("[%bar] %current/%total-%percent %rate", "", 100)
	p.SetPos(0, 0, 52, 1)
	p.Run(time.Millisecond * 100)

	// test uncertain progress bar
	//up := progress_bar.NewUncertainProgressBar()
	//up.SetPos(11, 0, 52, 1)
	//up.Run(time.Millisecond * 100)

	// wait
	<-p.Done
	//up.Stop()
	fmt.Printf("%s", cursor.GotoXY(1, 0))
	fmt.Println("time out. exit...")
}

func lineTest() {
	drawCoord := func(x, y, length int) {
		graph.Line(x, y, length, '|', 0)
		fmt.Printf("%s", cursor.GotoXY(x+length, y)+"v")
		length = 60
		graph.Line(x, y, length, '-', 1)
		fmt.Printf("%s", cursor.GotoXY(x, y+length)+">")
		fmt.Printf("%s", cursor.GotoXY(x, y)+"+")
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
func main() {
	//boxTest()
	barTest()
	//window.ClearScreen()
	//lineTest()
	c := '0'
	_, _ = fmt.Scanf("%c", &c)
}
