package box

import (
	"fmt"
	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/utils"
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

func GetBox(row, col int, boxType string, payload []string) (box []string, err error) {
	var useChar []rune
	payloadCnt := len(payload)
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
		//暂时无法在包含有非ASCII字符时正确对齐
		if i < payloadCnt {
			copyLen := min(col-2, len([]byte(payload[i])))
			a := string([]byte(payload[i])[:copyLen])
			//a := fmt.Sprintf("%*s", copyLen, payload[i]) //无法截断
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
