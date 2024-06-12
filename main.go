package main

import (
	"fmt"
	"github.com/gngtwhh/gocui/box"
	"github.com/gngtwhh/gocui/progress_bar"
	"github.com/gngtwhh/gocui/window"
	"time"
)

func main() {
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

	p := progress_bar.NewProgressBar(100)
	p.SetPos(10, 0, 52, 1)
	/*for i := 0; i < 100; i++ {
		p.Update(i)
		time.Sleep(time.Millisecond * 100)
	}*/
	p.Run(time.Millisecond * 100)
	<-p.Done
	fmt.Println("\ntime out. exit...")
}
