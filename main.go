package main

import box "github.com/gngtwhh/gocui/box"

func main() {
	aBox, _ := box.GetBox(5, 10, "fine")
	box.SetBoxAt(aBox, 1, 2)
}
