package box

import (
	"fmt"
	"github.com/gngtwhh/gocui/cursor"
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

const (
	FINE = iota
	BOLD
	DOUBLE
)

var fine = []rune{'─', '│', '┌', '┐', '└', '┘', '┬', '┴', '├', '┤', '┼'}
var bold = []rune{'━', '┃', '┏', '┓', '┗', '┛', '┳', '┻', '┣', '┫', '╋'}
var double = []rune{'═', '║', '╔', '╗', '╚', '╝', '╦', '╩', '╠', '╣', '╬'}

func GetBox(row, col int, boxType string) (box []string, err error) {
	var useChar []rune
	switch boxType {
	case "fine":
		useChar = fine
	case "bold":
		useChar = bold
	case "double":
		useChar = double
	default:
		useChar = fine
	}
	if row < 3 || col < 3 {
		err = fmt.Errorf("row and col must be greater than 3")
	}
	box = append(box, string(useChar[2])+strings.Repeat(string(useChar[0]), col-2)+string(useChar[3]))
	for i := 0; i < row-2; i++ {
		box = append(box, string(useChar[1])+strings.Repeat(" ", col-2)+string(useChar[1]))
	}
	box = append(box, string(useChar[4])+strings.Repeat(string(useChar[0]), col-2)+string(useChar[5]))
	return box, err
}

func SetBoxAt(box []string, x, y int) {
	row := len(box)
	for i := 0; i < row; i++ {
		fmt.Printf(cursor.GotoXY(x+i, y) + box[i])
	}
}
