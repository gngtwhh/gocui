package pb

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gngtwhh/gocui/font"
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
		"%bar", "%current", "%total", "%percent", "%elapsed", "%rate", "%spinner", "%bytes",
	}
}

// token is the interface that all tokens must implement.
type token interface {
	//toString(p *Property) string
	toString(ctx *Context) string
}

// Here are the tokens that can be used in the format string.
// All tokens use the toString method to convert the tokens to a string for print.

type TokenBar struct{}
type TokenCurrent struct{}
type TokenTotal struct{}
type TokenPercent struct{}
type TokenElapsed struct{}
type TokenRate struct {
	minDelay    time.Duration
	lastRate    string
	lastTime    time.Time
	lastProcess int64
}
type TokenString struct{ payload string }
type TokenSpinner struct{ cur int8 }
type TokenPercentage struct{}
type TokenBytes struct{}

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
		leftSpace := int(ctx.current)
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
	return strconv.FormatInt(ctx.current, 10)
}

func (t *TokenTotal) toString(ctx *Context) string {
	return strconv.FormatInt(ctx.property.Total, 10)
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
	//return fmt.Sprintf("%5.2fs", ctx.property.elapsed.Seconds())
	return fmt.Sprintf("%.1fs", time.Since(ctx.startTime).Seconds())
}

func (t *TokenRate) toString(ctx *Context) string {
	//defer func() { t.lastTime = time.Now() }()
	//if t.lastTime.IsZero() {
	//	return "0 it/s"
	//}
	//inc := 1 / time.Since(t.lastTime).Seconds()
	//t.lastTime = time.Now()
	//// return fmt.Sprintf("%.1f it/s", inc)
	//return fmt.Sprintf("%.0f it/s", inc+0.5)

	//dur := time.Since(ctx.startTime)
	//rate := (ctx.current) * int64(time.Second/dur)
	//return fmt.Sprintf("%d it/s", rate)

	defer func() {
		t.lastProcess = ctx.current
		t.lastTime = time.Now()
	}()
	inc := ctx.current - t.lastProcess
	if t.lastTime.IsZero() {
		t.lastRate = "0 it/s"
		return "0 it/s"
	}
	dur := time.Since(t.lastTime)
	if dur < t.minDelay {
		return t.lastRate
	}
	rate := (inc) * int64(time.Second/dur)
	t.lastRate = fmt.Sprintf("%d it/s", rate)
	return t.lastRate
}

func (s *TokenString) toString(ctx *Context) string {
	return s.payload
}

func (s *TokenSpinner) toString(ctx *Context) string {
	res := "\\|/-"[s.cur : s.cur+1]
	s.cur = (s.cur + 1) % 4
	return res
}

func (s *TokenBytes) toString(ctx *Context) string {
	calStr := func(b int64) string {
		if b == 0 {
			return "0 B"
		}
		sizes := []string{" B", " kB", " MB", " GB", " TB", " PB", " EB"}
		base := 1024.0
		e := math.Floor(math.Log(float64(b)) / math.Log(base))
		unit := sizes[int(e)]
		val := math.Floor(float64(b)/math.Pow(base, e)*10+0.5) / 10
		return fmt.Sprintf("%.1f%s", val, unit)
	}
	if ctx.property.Uncertain {
		return calStr(ctx.current)
	}
	return calStr(ctx.current) + "/" + calStr(ctx.property.Total)
}

// unmarshalToken converts the token string to a slice of tokens.
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
					ts = append(ts, &TokenRate{minDelay: time.Millisecond})
				case "%spinner":
					ts = append(ts, &TokenSpinner{})
				case "%bytes":
					ts = append(ts, &TokenBytes{})
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
