package progress_bar

import (
	"fmt"
	"github.com/gngtwhh/gocui/font"
	"strconv"
	"strings"
	"time"
)

// token is the interface that all style must implement.
type token interface {
	toString(p *Property) string
}

// Here are the style that can be used in the format string.
// All style use the toString method to convert the style to a string for print.

type TokenBar struct{}
type TokenCurrent struct{}
type TokenTotal struct{}
type TokenPercent struct{}
type TokenElapsed struct{}
type TokenRate struct{}
type TokenString struct {
	payload string
}
type TokenSpinner struct {
	cur int8
}

// toString prints the bar.
func (b *TokenBar) toString(p *Property) string {
	var repeatStr = func(s string, length int) string {
		if len(s) == 0 {
			return ""
		}
		return strings.Repeat(s, length/len(s)) + s[:length%len(s)]
	}
	if p.Uncertain {
		leftSpace := p.Current
		rightSpace := p.Width - leftSpace - len(p.Style.UnCertain)
		return font.Decorate(repeatStr(p.Style.Incomplete, leftSpace), p.Style.IncompleteColor) +
			font.Decorate(p.Style.UnCertain, p.Style.UnCertainColor) +
			font.Decorate(repeatStr(p.Style.Incomplete, rightSpace), p.Style.IncompleteColor)
	} else {
		completeLength := int(float64(p.Current) / float64(p.Total) * float64(p.Width))
		return font.Decorate(repeatStr(p.Style.Complete, completeLength), p.Style.CompleteColor) +
			font.Decorate(repeatStr(p.Style.Incomplete, p.Width-completeLength), p.Style.IncompleteColor)
	}
}

func (c *TokenCurrent) toString(p *Property) string {
	return strconv.Itoa(p.Current)
}

func (t *TokenTotal) toString(p *Property) string {
	return strconv.Itoa(p.Total)
}

func (t *TokenPercent) toString(p *Property) string {
	var percent float64
	if p.Current == 0 {
		percent = 0
	} else {
		percent = float64(p.Current) / float64(p.Total) * 100
	}
	// 保留2位小数
	return fmt.Sprintf("%5.2f%%", percent)
}

func (t *TokenElapsed) toString(p *Property) string {
	return fmt.Sprintf("%5.2fs", p.elapsed.Seconds())
}

func (t *TokenRate) toString(p *Property) string {
	// ops/s
	count := float64(time.Second) / float64(p.rate)
	return fmt.Sprintf("%4.1f ops/s", count)
}

func (s *TokenString) toString(p *Property) string {
	return s.payload
}

func (s *TokenSpinner) toString(p *Property) string {
	res := "\\|/-"[s.cur : s.cur+1]
	s.cur = (s.cur + 1) % 4
	return res
}

// unmarshalToken converts the token string to a slice of style.
func unmarshalToken(token string) (ts []token) {
	if len(token) == 0 {
		return
	}
	legalTokens := []string{
		"%bar", "%current", "%total", "%percent", "%elapsed", "%rate", "%spinner",
	}

	for len(token) > 0 {
		if token[0] != '%' {
			if idx := strings.IndexAny(token, "%"); idx == -1 {
				ts = append(ts, &TokenString{payload: token})
				break
			} else {
				ts = append(ts, &TokenString{payload: token[:idx]})
				token = token[idx:]
			}
		}
		ok := false // matched a legal token
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
		if !ok {
			ts = append(ts, &TokenString{payload: "%"})
			token = token[1:]
		}
	}
	return
}
