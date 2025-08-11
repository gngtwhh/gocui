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
 * %bytes: Progress of writing data
 * &percent: Percentage of progress
 **************************************************/

// var legalTokens []string
var registeredTokens = make(map[string]token)

func InitBarToken() {
	// legalTokens = []string{
	// 	"%bar", "%current", "%total", "%percent", "%elapsed", "%rate", "%spinner", "%bytes",
	// }
	registeredTokens["%bar"] = &TokenBar{}
	registeredTokens["%current"] = &TokenCurrent{}
	registeredTokens["%total"] = &TokenTotal{}
	registeredTokens["%percent"] = &TokenPercent{}
	registeredTokens["%elapsed"] = &TokenElapsed{}
	registeredTokens["%rate"] = &TokenRate{minDelay: time.Millisecond * 100}
	registeredTokens["%spinner"] = &TokenSpinner{}
	registeredTokens["%bytes"] = &TokenBytes{}
}

// token is the interface that all tokens must implement.
type token interface {
	ToString(ctx *Context) string
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
type TokenBytes struct{}

// ToString implements the interface
func (b *TokenBar) ToString(ctx *Context) string {
	var repeatStr = func(s string, length int) string {
		if len(s) == 0 {
			return ""
		}
		return strings.Repeat(s, length/len(s)) + s[:length%len(s)]
	}

	p := &ctx.Property
	barWidth := p.BarWidth
	if barWidth == 0 {
		barWidth = ctx.WindowWidth - ctx.WidthWithoutBar
	}
	if p.Uncertain {
		// leftSpace := int(ctx.current)
		// rightSpace := barWidth - leftSpace - len(p.Style.UnCertain)
		// if leftSpace == 0 && ctx.direction == -1 {
		// 	ctx.direction = 1
		// }
		// if rightSpace == 0 && ctx.direction == 1 {
		// 	ctx.direction = -1
		// }
		leftSpace := int(ctx.Current) % barWidth
		uncertainWidth := min(barWidth-leftSpace, len(p.Style.UnCertain))
		rightSpace := max(0, barWidth-leftSpace-len(p.Style.UnCertain))
		return font.Decorate(repeatStr(p.Style.Incomplete, leftSpace), p.Style.IncompleteColor) +
			font.Decorate(repeatStr(p.Style.UnCertain, uncertainWidth), p.Style.UnCertainColor) +
			font.Decorate(repeatStr(p.Style.Incomplete, rightSpace), p.Style.IncompleteColor)
	} else {
		completeLength := int(float64(ctx.Current)/float64(ctx.Total)*float64(barWidth)) - len(p.Style.CompleteHead)
		if completeLength < 0 {
			completeLength = 0
		}
		return font.Decorate(repeatStr(p.Style.Complete, completeLength), p.Style.CompleteColor) +
			font.Decorate(p.Style.CompleteHead, p.Style.CompleteHeadColor) +
			font.Decorate(repeatStr(p.Style.Incomplete, barWidth-completeLength-len(p.Style.CompleteHead)), p.Style.IncompleteColor)
	}
}

func (c *TokenCurrent) ToString(ctx *Context) string {
	return strconv.FormatInt(ctx.Current, 10)
}

func (t *TokenTotal) ToString(ctx *Context) string {
	return strconv.FormatInt(ctx.Total, 10)
}

func (t *TokenPercent) ToString(ctx *Context) string {
	var percent int
	if ctx.Current == 0 {
		percent = 0
	} else {
		percent = int(float64(ctx.Current) / float64(ctx.Total) * 100)
	}
	// 保留2位小数
	return fmt.Sprintf("%3d%%", percent)
}

func (t *TokenElapsed) ToString(ctx *Context) string {
	// return fmt.Sprintf("%5.2fs", ctx.property.elapsed.Seconds())
	return fmt.Sprintf("%.1fs", time.Since(ctx.StartTime).Seconds())
}

func (t *TokenRate) ToString(ctx *Context) string {
	inc := ctx.Current - t.lastProcess
	if t.lastTime.IsZero() {
		t.lastTime = time.Now()
		t.lastRate = "0 it/s"
		return "0 it/s"
	}
	dur := time.Since(t.lastTime)
	if dur < t.minDelay {
		return t.lastRate
	}
	// return dur.String()
	rate := (inc) * int64(time.Second/dur)
	t.lastRate = fmt.Sprintf("%d it/s", rate)
	t.lastProcess = ctx.Current
	t.lastTime = time.Now()
	return t.lastRate
}

func (s *TokenString) ToString(ctx *Context) string {
	return s.payload
}

func (s *TokenSpinner) ToString(ctx *Context) string {
	res := "\\|/-"[s.cur : s.cur+1]
	s.cur = (s.cur + 1) % 4
	return res
}

func (s *TokenBytes) ToString(ctx *Context) string {
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
	if ctx.Property.Uncertain {
		return calStr(ctx.Current)
	}
	return calStr(ctx.Current) + "/" + calStr(ctx.Total)
}

// unmarshalToken converts the token string to a slice of tokens.
func unmarshalToken(format string) (ts []token, barPos []int) {
	if len(format) == 0 {
		return
	}

	ok := false // Whether a valid token is matched
	for len(format) > 0 {
		ok = false
		if format[0] != '%' {
			goto commonString
		}
		for legalToken, legalTokenInstance := range registeredTokens {
			if strings.HasPrefix(format, legalToken) {
				format = format[len(legalToken):]
				var newToken = legalTokenInstance
				ts = append(ts, newToken)
				ok = true
				break
			}
		}
		if ok && len(format) == 0 {
			break
		}
	commonString:
		if format[0] != '%' || !ok {
			if idx := strings.IndexAny(format[1:], "%"); idx == -1 {
				ts = append(ts, &TokenString{payload: format})
				break
			} else {
				ts = append(ts, &TokenString{payload: format[:idx+1]})
				format = format[idx+1:]
			}
		}
	}
	return
}

// RegisterToken allows user to register a new token to achieve the unique effect they desire.
// This means allowing users to use their own custom tokens with specific behaviors in the format,
// as long as they are registered before use.
// WARNING: Registering a token with a name that already exists will overwrite the existing token.
// WARNING: Ensure no existing token becomes a prefix of any newly registered token.
func RegisterToken(name string, token token) {
	registeredTokens[name] = token
}
