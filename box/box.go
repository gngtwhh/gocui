package box

import (
	"fmt"
	"strings"

	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/utils"
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

const (
	FINE = iota
	BOLD
	DOUBLE
	ROUNDED
)

// Types of box
var fine = []rune{'─', '│', '┌', '┐', '└', '┘', '┬', '┴', '├', '┤', '┼'}
var bold = []rune{'━', '┃', '┏', '┓', '┗', '┛', '┳', '┻', '┣', '┫', '╋'}
var double = []rune{'═', '║', '╔', '╗', '╚', '╝', '╦', '╩', '╠', '╣', '╬'}
var rounded = []rune{'─', '│', '╭', '╮', '╰', '╯', '┬', '┴', '├', '┤', '┼'}

type Property struct {
	Row, Col int    // min size of the box, 0 for auto
	BoxType  string // box type
	Title    string
	Payload  []string
}

func GetBox(p Property) (box []string, err error) {
	payloadCnt := len(p.Payload)
	var useChar []rune
	switch p.BoxType {
	case "fine":
		useChar = fine
	case "bold":
		useChar = bold
	case "double":
		useChar = double
	case "rounded":
		useChar = rounded
	default:
		useChar = fine
	}
	// calculate row and col
	row, col := p.Row, p.Col
	if (row < 2 && row != 0) || (col < 2 && col != 0) {
		box = []string{}
		err = fmt.Errorf("row and col must be greater than 2")
		return
	}
	if row == 0 {
		row = len(p.Payload) + 2
	}
	if col == 0 {
		for _, p := range p.Payload {
			col = max(len(p), col)
		}
		col += 2
	}
	if len(p.Title) > col-2 {
		p.Title = p.Title[:col-2] // 标题过长，截断
	}
	// generate title line of box
	{
		aLen := min(2, col-2-len(p.Title))
		box = append(box, string(useChar[2])+strings.Repeat(string(useChar[0]), aLen)+p.Title+
			strings.Repeat(string(useChar[0]), col-2-len(p.Title)-aLen)+string(useChar[3]))
	}
	for i := 0; i < row-2; i++ {
		//暂时无法在包含有非ASCII字符时正确对齐
		if i < payloadCnt {
			copyLen := min(col-2, len([]byte(p.Payload[i])))
			a := string([]byte(p.Payload[i])[:copyLen])
			b := strings.Repeat(" ", col-2-copyLen)
			box = append(box, string(useChar[1])+a+b+string(useChar[1]))
		} else {
			box = append(box, string(useChar[1])+strings.Repeat(" ", col-2)+string(useChar[1]))
		}
	}
	box = append(box, string(useChar[4])+strings.Repeat(string(useChar[0]), col-2)+string(useChar[5]))
	return box, err
}

func SetBoxAt(box []string, x, y int) {
	row := len(box)
	utils.ConsoleMutex.Lock()
	defer utils.ConsoleMutex.Unlock()
	for i := 0; i < row; i++ {
		cursor.GotoXY(x+i, y)
		fmt.Print(box[i])
	}
}
