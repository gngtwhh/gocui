package font

import (
	"fmt"
	"strings"
)

// Style of font
const (
	Bold       = 1
	Dim        = 2
	Italic     = 3
	Underline  = 4
	BlinkSlow  = 5
	BlinkFast  = 6
	Reverse    = 7
	Hide       = 8
	CrossedOut = 9
)

// 8/16 colors
const (
	RESET = 0

	// foreground colors
	Black        = 30
	Red          = 31
	Green        = 32
	Yellow       = 33
	Blue         = 34
	Magenta      = 35
	Cyan         = 36
	White        = 37
	LightBlack   = 90
	LightRed     = 91
	LightGreen   = 92
	LightYellow  = 93
	LightBlue    = 94
	LightMagenta = 95
	LightCyan    = 96
	LightWhite   = 97

	// background colors
	BlackBg        = 40
	RedBg          = 41
	GreenBg        = 42
	YellowBg       = 43
	BlueBg         = 44
	MagentaBg      = 45
	CyanBg         = 46
	WhiteBg        = 47
	LightBlackBg   = 100
	LightRedBg     = 101
	LightGreenBg   = 102
	LightYellowBg  = 103
	LightBlueBg    = 104
	LightMagentaBg = 105
	LightCyanBg    = 106
	LightWhiteBg   = 107
)

func SetColor(color int) {
	fmt.Printf("\033[%dm", color)
}

func SetColorRgb(r, g, b int, bg bool) {
	if bg {
		fmt.Printf("\033[48;2;%d;%d;%dm", r, g, b)
	} else {
		fmt.Printf("\033[38;2;%d;%d;%dm", r, g, b)
	}
}

func ResetColor() {
	fmt.Printf("\033[%dm", RESET)
}

func SetStyle(style int) string {
	return fmt.Sprintf("\033[%dm", style)
}

func Decorate(text string, style ...int) string {
	buf := strings.Builder{}
	buf.WriteString("\033[")
	l := len(style)
	for i, s := range style {
		buf.WriteString(fmt.Sprintf("%d", s))
		if i < l-1 {
			buf.WriteString(";")
		}
	}
	buf.WriteString("m")
	buf.WriteString(text)
	buf.WriteString("\033[0m")
	return buf.String()
}

func Splice(ks ...interface{}) string {
	buf := strings.Builder{}
	for _, k := range ks {
		switch v := k.(type) {
		case int:
			buf.WriteString(fmt.Sprintf("\033[%dm", v))
		case rune:
			buf.WriteRune(v)
		case string:
			buf.WriteString(v)
		}
	}
	return buf.String()
}
