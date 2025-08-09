package box

import (
	"fmt"
	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/font"
	"github.com/gngtwhh/gocui/utils"
	"github.com/gngtwhh/gocui/window"
	"strings"
)

/**
 * box characters
 * 			0	1	2	3	4	5	6	7	8	9	A	B	C	D	E	F
 * U+250x	─	━	│	┃	┄	┅	┆	┇	┈	┉	┊	┋	┌	┍	┎	┏
 * U+251x	┐	┑	┒	┓	└	┕	┖	┗	┘	┙	┚	┛	├	┝	┞	┟
 * U+252x	┠	┡	┢	┣	┤	┥	┦	┧	┨	┩	┪	┫	┬	┭	┮	┯
 * U+253x	┰	┱	┲	┳	┴	┵	┶	┷	┸	┹	┺	┻	┼	┽	┾	┿
 * U+254x	╀	╁	╂	╃	╄	╅	╆	╇	╈	╉	╊	╋	╌	╍	╎	╏
 * U+255x	═	║	╒	╓	╔	╕	╖	╗	╘	╙	╚	╛	╜	╝	╞	╟
 * U+256x	╠	╡	╢	╣	╤	╥	╦	╧	╨	╩	╪	╫	╬	╭	╮	╯
 * U+257x	╰	╱	╲	╳	╴	╵	╶	╷	╸	╹	╺	╻	╼	╽	╾	╿
 */

// Box types
const (
	FINE = iota
	BOLD
	DOUBLE
	ROUNDED
)

// Title positions
const (
	TopLeft = iota
	TopMid
	TopRight
	BottomLeft
	BottomMid
	BottomRight
	InsideLeft
	InsideMid
	InsideRight
)

// Content align
const (
	Center = iota
	Left
	Right
)

// Types of box
var (
	fine    = []rune{'─', '│', '┌', '┐', '└', '┘', '┬', '┴', '├', '┤', '┼'}
	bold    = []rune{'━', '┃', '┏', '┓', '┗', '┛', '┳', '┻', '┣', '┫', '╋'}
	double  = []rune{'═', '║', '╔', '╗', '╚', '╝', '╦', '╩', '╠', '╣', '╬'}
	rounded = []rune{'─', '│', '╭', '╮', '╰', '╯', '┬', '┴', '├', '┤', '┼'}
)

// default settings
var (
	DefaultBox      *Box
	defaultType     []rune
	DefaultProperty Property
)

// Char contains characters of the box frame
type Char struct {
	TopLeft, TopRight, BottomLeft, BottomRight rune // Characters for box corners
	Top, Bottom, Left, Right                   rune // Characters for box sides
}

// Color contains colors of the box frame
type Color struct {
	TopLeftColor, TopRightColor, BottomLeftColor, BottomRightColor int // color of box corners
	TopColor, BottomColor, LeftColor, RightColor                   int // color of box sides
	TitleColor, InnerColor                                         int // color of the title and the text inside the box
}
type Style struct {
	Char
	Color
}

type Property struct {
	Style          // decorate characters and color
	PadX, PadY int // padding
	Align      int // Align of the content (Center/Left/Right), default Center
	TitlePos   int // title pos (Top/Bottom/Inside x Left/Middle/Right), default TopLeft
	PosX, PosY int // default pos to be print

	BindPos bool // Whether bind the absolute pos, PosX and PosY are valid only when BindPos is true
}

type Box struct {
	Property
}

func init() {
	defaultType = fine
	DefaultProperty = Property{
		PadX: 1, PadY: 1,
		Align:    Center,
		TitlePos: TopLeft,
		Style: Style{
			Char: Char{
				TopLeft: defaultType[2], TopRight: defaultType[3],
				BottomLeft: defaultType[4], BottomRight: defaultType[5],
				Top: defaultType[0], Bottom: defaultType[0],
				Left: defaultType[1], Right: defaultType[1],
			},
			Color: Color{
				TopLeftColor: font.White, TopRightColor: font.White, BottomLeftColor: font.White, BottomRightColor: font.White,
				TopColor: font.White, BottomColor: font.White, LeftColor: font.White, RightColor: font.White,
				TitleColor: font.White, InnerColor: font.White,
			},
		},
	}
	DefaultBox = &Box{DefaultProperty}
}

// func NewBox(p Property) (box []string, err error) {
//	payloadCnt := len(p.Payload)
//	var useChar []rune
//	switch p.BoxType {
//	case "fine":
//		useChar = fine
//	case "bold":
//		useChar = bold
//	case "double":
//		useChar = double
//	case "rounded":
//		useChar = rounded
//	default:
//		useChar = fine
//	}
//	// calculate row and col
//	row, col := p.Row, p.Col
//	if (row < 2 && row != 0) || (col < 2 && col != 0) {
//		box = []string{}
//		err = fmt.Errorf("row and col must be greater than 2")
//		return
//	}
//	if row == 0 {
//		row = len(p.Payload) + 2
//	}
//	if col == 0 {
//		for _, p := range p.Payload {
//			col = max(len(p), col)
//		}
//		col += 2
//	}
//	if len(p.Title) > col-2 {
//		p.Title = p.Title[:col-2] // 标题过长，截断
//	}
//	// generate title line of box
//	{
//		aLen := min(2, col-2-len(p.Title))
//		box = append(box, string(useChar[2])+strings.Repeat(string(useChar[0]), aLen)+p.Title+
//			strings.Repeat(string(useChar[0]), col-2-len(p.Title)-aLen)+string(useChar[3]))
//	}
//	for i := 0; i < row-2; i++ {
//		//暂时无法在包含有非ASCII字符时正确对齐
//
//		if i < payloadCnt {
//			copyLen := min(col-2, len([]byte(p.Payload[i])))
//			a := string([]byte(p.Payload[i])[:copyLen])
//			b := strings.Repeat(" ", col-2-copyLen)
//			box = append(box, string(useChar[1])+a+b+string(useChar[1]))
//		} else {
//			box = append(box, string(useChar[1])+strings.Repeat(" ", col-2)+string(useChar[1]))
//		}
//	}
//	box = append(box, string(useChar[4])+strings.Repeat(string(useChar[0]), col-2)+string(useChar[5]))
//	return box, err
// }

// NewBox create a box template with several modify functions
func NewBox(mfs ...ModFunc) (box *Box, err error) {
	var p Property
	for _, mf := range mfs {
		if mf == nil {
			return nil, fmt.Errorf("modify func cannot be nil")
		}
		mf(&p)
	}
	// revise the characters
	if p.Style.Top == rune(0) {
		p.Style.Top = defaultType[0]
	}
	if p.Style.Bottom == rune(0) {
		p.Style.Bottom = defaultType[0]
	}
	if p.Style.Left == rune(0) {
		p.Style.Left = defaultType[1]
	}
	if p.Style.Right == rune(0) {
		p.Style.Right = defaultType[1]
	}
	if p.Style.TopLeft == rune(0) {
		p.Style.TopLeft = defaultType[2]
	}
	if p.Style.TopRight == rune(0) {
		p.Style.TopRight = defaultType[3]
	}
	if p.Style.BottomLeft == rune(0) {
		p.Style.BottomLeft = defaultType[4]
	}
	if p.Style.BottomRight == rune(0) {
		p.Style.BottomRight = defaultType[5]
	}
	// revise the colors
	if p.Style.TopLeftColor == 0 {
		p.Style.TopLeftColor = font.White
	}
	if p.Style.TopRightColor == 0 {
		p.Style.TopRightColor = font.White
	}
	if p.Style.BottomLeftColor == 0 {
		p.Style.BottomLeftColor = font.White
	}
	if p.Style.BottomRightColor == 0 {
		p.Style.BottomRightColor = font.White
	}
	if p.Style.TopColor == 0 {
		p.Style.TopColor = font.White
	}
	if p.Style.BottomColor == 0 {
		p.Style.BottomColor = font.White
	}
	if p.Style.LeftColor == 0 {
		p.Style.LeftColor = font.White
	}
	if p.Style.RightColor == 0 {
		p.Style.RightColor = font.White
	}
	if p.Style.TitleColor == 0 {
		p.Style.TitleColor = font.White
	}
	if p.Style.InnerColor == 0 {
		p.Style.InnerColor = font.White
	}
	// revise other props
	if p.PadX < 0 {
		p.PadX = 0
	}
	if p.PadY < 0 {
		p.PadY = 0
	}
	if p.Align < 0 || p.Align > 2 {
		p.Align = Center
	}
	if p.TitlePos < 0 || p.TitlePos > 8 {
		p.TitlePos = TopLeft
	}
	maxY, maxX := window.GetConsoleSize()
	if p.PosX < 0 || p.PosX >= maxX {
		p.PosX = 0
	}
	if p.PosY < 0 || p.PosY >= maxY {
		p.PosY = 0
	}
	return &Box{p}, nil
}

func (box *Box) Print(title string, payload []string) {
	p := box.Property
	maxLen := 0
	for _, str := range payload {
		maxLen = max(maxLen, len(str))
	}
	x, y := 0, 0
	if box.Property.BindPos {
		x, y = p.PosX, p.PosY
	}
	var line string
	utils.ConsoleMutex.Lock()
	defer utils.ConsoleMutex.Unlock()
	// first line
	line = font.Splice(
		p.TopLeftColor, p.TopLeft, p.TopColor, strings.Repeat(string(p.Top), maxLen+p.PadX*2), p.TopRightColor,
		p.TopRight,
	)
	cursor.GotoXY(x, y)
	fmt.Print(line)
	y++

}
