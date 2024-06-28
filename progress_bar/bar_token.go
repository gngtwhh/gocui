package progress_bar

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Param is an alias for the empty interface and is used to give the toString method any type of parameter value.
type Param = interface{}

// Token is the interface that all tokens must implement.
type Token interface {
	toString(p ...Param) string
}

// Here are the tokens that can be used in the format string.
// All tokens use the toString method to convert the tokens to a string for print.

type TokenBar struct{}
type TokenCurrent struct{}
type TokenTotal struct{}
type TokenPercent struct{}
type TokenElapsed struct{}
type TokenRate struct{}
type TokenString struct {
	payload string
}

// toString prints the bar.
// p should be left int,right int,block string
// return " "*left + block + " "*right
func (b *TokenBar) toString(p ...Param) string {
	return strings.Repeat(" ", p[0].(int)) + p[2].(string) + strings.Repeat(" ", p[1].(int))
}

// p should be current int
// return "%d" of current
func (c *TokenCurrent) toString(p ...Param) string {
	return strconv.Itoa(p[0].(int))
}

// p should be total int
// return "%d" of total
func (t *TokenTotal) toString(p ...Param) string {
	return strconv.Itoa(p[0].(int))
}

// p should be current int,total int
// return "dd.dd%"
func (t *TokenPercent) toString(p ...Param) string {
	var percent float64
	if p[0].(int) == 0 {
		percent = 0
	} else {
		percent = float64(p[0].(int)) / float64(p[1].(int)) * 100
	}
	// 保留2位小数
	return fmt.Sprintf("%5.2f%%", percent)
}

// p should be time.Duration
// return "dd.dds"
func (t *TokenElapsed) toString(p ...Param) string {
	return fmt.Sprintf("%5.2fs", p[0].(time.Duration).Seconds())
}

// p should be time.Duration
// return "dd.d ops/s"
func (t *TokenRate) toString(p ...Param) string {
	// ops/s
	count := float64(time.Second) / float64(p[0].(time.Duration))
	return fmt.Sprintf("%4.1f ops/s", count)
}

// p should be void
// just return the payload of the string token
func (s *TokenString) toString(p ...Param) string {
	return s.payload
}

// unmarshalToken converts the token string to a slice of tokens.
func unmarshalToken(token string) (ts []Token) {
	if len(token) == 0 {
		return
	}
	legalTokens := []string{
		"%bar", "%current", "%total", "%percent", "%elapsed", "%rate",
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

func formatBarString(thisBar bar) (barString string) {
	switch thisBar.(type) {
	case *ProgressBar:
		for _, token := range thisBar.(*ProgressBar).tokens {
			barString += token.toString(thisBar)
		}
	case *UncertainProgressBar:
		for _, token := range thisBar.(*UncertainProgressBar).tokens {
			barString += token.toString(thisBar)
		}
	}
	return
}
