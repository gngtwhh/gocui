package pb

import (
	"fmt"
	"github.com/gngtwhh/gocui/font"
	"strconv"
	"strings"
	"time"
)

/**************************************************
 * Currently supported tokens:
 * %string: Regular string
 * %bar: The body of the progress bar
 * %current: Current progress value
 * %total: Total progress value
 * %elapsed: The elapsed time of the progress bar
 * %rate: Speed of the progress bar
 * %spinner: A rotator character
 **************************************************/

var legalTokens []string

func init() {
	legalTokens = []string{
		"%bar", "%current", "%total", "%percent", "%elapsed", "%rate", "%spinner",
	}
}

// token is the interface that all style must implement.
type token interface {
	//toString(p *Property) string
	toString(ctx *Context) string
}

// Here are the style that can be used in the format string.
// All style use the toString method to convert the style to a string for print.

type TokenBar struct{}
type TokenCurrent struct{}
type TokenTotal struct{}
type TokenPercent struct{}
type TokenElapsed struct{}
type TokenRate struct{}
type TokenString struct{ payload string }
type TokenSpinner struct{ cur int8 }
type TokenPercentage struct{}

// toString implements the interface
func (b *TokenBar) toString(ctx *Context) string {
	var repeatStr = func(s string, length int) string {
		if len(s) == 0 {
			return ""
		}
		return strings.Repeat(s, length/len(s)) + s[:length%len(s)]
	}
	p := &ctx.property
	if p.Uncertain {
		leftSpace := ctx.current
		rightSpace := p.Width - leftSpace - len(p.Style.UnCertain)
		return font.Decorate(repeatStr(p.Style.Incomplete, leftSpace), p.Style.IncompleteColor) +
			font.Decorate(p.Style.UnCertain, p.Style.UnCertainColor) +
			font.Decorate(repeatStr(p.Style.Incomplete, rightSpace), p.Style.IncompleteColor)
	} else {
		completeLength := int(float64(ctx.current) / float64(p.Total) * float64(p.Width))
		return font.Decorate(repeatStr(p.Style.Complete, completeLength), p.Style.CompleteColor) +
			font.Decorate(repeatStr(p.Style.Incomplete, p.Width-completeLength), p.Style.IncompleteColor)
	}
}

func (c *TokenCurrent) toString(ctx *Context) string {
	return strconv.Itoa(ctx.current)
}

func (t *TokenTotal) toString(ctx *Context) string {
	return strconv.Itoa(ctx.property.Total)
}

func (t *TokenPercent) toString(ctx *Context) string {
	var percent int
	if ctx.current == 0 {
		percent = 0
	} else {
		percent = int(float64(ctx.current) / float64(ctx.property.Total) * 100)
	}
	// 保留2位小数
	return fmt.Sprintf("%3d%%", percent)
}

func (t *TokenElapsed) toString(ctx *Context) string {
	return fmt.Sprintf("%5.2fs", ctx.property.elapsed.Seconds())
}

func (t *TokenRate) toString(ctx *Context) string {
	// ops/s
	count := float64(time.Second) / float64(ctx.property.rate)
	return fmt.Sprintf("%4.1f ops/s", count)
}

func (s *TokenString) toString(ctx *Context) string {
	return s.payload
}

func (s *TokenSpinner) toString(ctx *Context) string {
	res := "\\|/-"[s.cur : s.cur+1]
	s.cur = (s.cur + 1) % 4
	return res
}

// unmarshalToken converts the token string to a slice of style.
func unmarshalToken(token string) (ts []token) {
	if len(token) == 0 {
		return
	}

	ok := false // Whether a valid token is matched
	for len(token) > 0 {
		ok = false
		if token[0] != '%' {
			goto commonString
		}
		for _, legalToken := range legalTokens {
			if strings.HasPrefix(token, legalToken) {
				token = token[len(legalToken):]
				switch legalToken {
				case "%bar":
					ts = append(ts, &TokenBar{})
				case "%current":
					ts = append(ts, &TokenCurrent{})
				case "%total":
					ts = append(ts, &TokenTotal{})
				case "%percent":
					ts = append(ts, &TokenPercent{})
				case "%elapsed":
					ts = append(ts, &TokenElapsed{})
				case "%rate":
					ts = append(ts, &TokenRate{})
				case "%spinner":
					ts = append(ts, &TokenSpinner{})
				}
				ok = true
				break
			}
		}
		if ok && len(token) == 0 {
			break
		}
	commonString:
		if token[0] != '%' || !ok {
			if idx := strings.IndexAny(token[1:], "%"); idx == -1 {
				ts = append(ts, &TokenString{payload: token})
				break
			} else {
				ts = append(ts, &TokenString{payload: token[:idx+1]})
				token = token[idx+1:]
			}
		}
	}
	return
}
