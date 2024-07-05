package progress_bar

import (
	"fmt"
	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/utils"
	"github.com/gngtwhh/gocui/window"
	"strings"
	"sync"
	"time"
)

// bar is a common interface for ProgressBar and UncertainProgressBar.
type bar interface {
	SetPos(posX, posY, width, height int)
	Run(period time.Duration)
	Print()
	Stop()
}

// ProgressBar is a simple progress bar implementation.
type ProgressBar struct {
	Total, Current int
	posX, posY     int
	width, height  int
	rate, elapsed  time.Duration
	tokens         []Token
	running        bool          // Whether the progress bar is running
	interrupt      chan struct{} // interrupt channel to stop Run()
	Done           chan struct{} // Done channel to signal completion
	rw             sync.RWMutex  // RWMutex to synchronize access to the progress bar
}

// // NewProgressBar creates a new progress bar with the given total.
//
//	func NewProgressBar(total int) *ProgressBar {
//		return &ProgressBar{
//			Total:     total,
//			interrupt: make(chan struct{}),
//			Done:      make(chan struct{}),
//			rw:        sync.RWMutex{},
//		}
//	}

// NewProgressBar creates a new progress bar that with the given style and total.
func NewProgressBar(style, config string, total int) (pb *ProgressBar, err error) {
	//styleTokens,configTokens := unmarshalToken(style), unmarshalToken(config)
	styleTokens, _ := unmarshalToken(style), unmarshalToken(config)
	pb = &ProgressBar{
		Total:     total,
		tokens:    styleTokens,
		interrupt: make(chan struct{}),
		Done:      make(chan struct{}),
		rw:        sync.RWMutex{},
	}
	return pb, err
}

// SetPos sets the position of the progress bar.
func (p *ProgressBar) SetPos(posX, posY, width, height int) {
	p.rw.Lock()
	p.posX = posX
	p.posY = posY
	p.width = width
	p.height = height
	p.rw.Unlock()
}

// Update updates the current progress of the progress bar.
func (p *ProgressBar) Update(current int) {
	p.rw.Lock()
	if p.running {
		p.interrupt <- struct{}{}
	}
	p.Current = min(current, p.Total)
	p.rw.Unlock()
	p.Print()
}

// Print prints the current progress of the progress bar.
func (p *ProgressBar) Print() {
	p.rw.RLock()              // Read lock to avoid race condition
	utils.ConsoleMutex.Lock() // Lock the cursor to avoid concurrent access

	//fmt.Printf("%s", cursor.HideCursor())
	//window.ClearArea(p.posX, p.posY, p.width, p.height)
	//for i := 0; i < p.height; i++ {
	//	fmt.Printf("%s", cursor.GotoXY(p.posX+i, p.posY))
	//	block := int(float64(p.Current) / float64(p.Total) * float64(p.width))
	//	fmt.Printf("%s", strings.Repeat("█", block)+strings.Repeat(" ", p.width-block))
	//}
	//fmt.Printf("%s", cursor.ShowCursor())

	var getParamAndToString = func(token Token) (payload string) {
		switch token.(type) {
		case *TokenBar:
			block := int(float64(p.Current) / float64(p.Total) * float64(p.width))
			left, right := 0, p.width-block
			return token.toString(left, right, strings.Repeat("█", block))
		case *TokenCurrent:
			return token.toString(p.Current)
		case *TokenTotal:
			return token.toString(p.Total)
		case *TokenPercent:
			return token.toString(p.Current, p.Total)
		case *TokenElapsed:
			return token.toString(p.elapsed)
		case *TokenRate:
			return token.toString(p.rate)
		case *TokenString:
			return token.toString()
		default:
			return
		}
	}

	cursor.HideCursor()
	window.ClearArea(p.posX, p.posY, p.width, p.height)
	for i := 0; i < p.height; i++ {
		cursor.GotoXY(p.posX+i, p.posY)
		payloadBuilder := strings.Builder{}
		for _, token := range p.tokens {
			payloadBuilder.WriteString(getParamAndToString(token))
		}
		fmt.Printf("%s", payloadBuilder.String())
	}
	//fmt.Printf("%s", cursor.ShowCursor())

	utils.ConsoleMutex.Unlock() // Unlock the cursor
	p.rw.RUnlock()              // Read unlock
}

// Run starts the progress bar.
func (p *ProgressBar) Run(period time.Duration) {
	p.rw.Lock()
	p.rate = period
	p.rw.Unlock()
	if p.running {
		return
	}
	ticker := time.NewTicker(period)
	p.running = true
	p.Print() // Print the initial state
	go func() {
		defer func() {
			ticker.Stop()
			p.running = false
		}()
		for {
			select {
			case <-ticker.C:
				p.Print()
				if p.Current == p.Total {
					p.Done <- struct{}{}
					return
				}
				p.rw.Lock()
				p.Current++
				p.elapsed += period
				p.rw.Unlock()
			case <-p.interrupt:
				return
			}
		}
	}()
}

// Stop stops the progress bar.
func (p *ProgressBar) Stop() {
	p.interrupt <- struct{}{}
	p.running = false
}

// UncertainProgressBar is a progress bar that shows an uncertain progress.
type UncertainProgressBar struct {
	posX, posY         int
	width, height      int
	tokens             []Token
	Current, blockSize int
	direction          int // 1 for increasing, -1 for decreasing
	running            bool
	interrupt          chan struct{}
}

// NewUncertainProgressBar creates a new uncertain progress bar.
func NewUncertainProgressBar() *UncertainProgressBar {
	return &UncertainProgressBar{
		direction: 1,
		interrupt: make(chan struct{}),
	}
}

func (p *UncertainProgressBar) SetPos(posX, posY, width, height int) {
	getBlockSize := func(width int) int {
		if width < 4 {
			return 1
		}
		return width / 5
	}
	p.posX = posX
	p.posY = posY
	p.width = width
	p.height = height
	p.blockSize = getBlockSize(width - 2)
}

func (p *UncertainProgressBar) Print() {
	utils.ConsoleMutex.Lock() // Lock the cursor to avoid concurrent access
	cursor.HideCursor()
	window.ClearArea(p.posX, p.posY, p.width, p.height)
	cursor.GotoXY(p.posX, p.posY)

	leftSpace := p.Current
	rightSpace := p.width - 2 - leftSpace - p.blockSize
	bar := "[" + strings.Repeat(" ", leftSpace) + strings.Repeat("#", p.blockSize) + strings.Repeat(" ", rightSpace) + "]"
	for i := 0; i < p.height; i++ {
		cursor.GotoXY(p.posX+i, p.posY)
		fmt.Printf("%s", bar)
	}
	utils.ConsoleMutex.Unlock() // Unlock the cursor
}

func (p *UncertainProgressBar) update() {
	p.Current = (p.Current + p.direction) % (p.width - 2)
	if p.Current == p.width-2-p.blockSize {
		p.direction = -1
	}
	if p.Current == 0 {
		p.direction = 1
	}
}

func (p *UncertainProgressBar) Run(period time.Duration) {
	if p.running {
		return
	}
	ticker := time.NewTicker(period)
	p.running = true
	go func() {
		defer func() {
			ticker.Stop()
			p.running = false
		}()
		for {
			select {
			case <-ticker.C:
				p.Print()
				p.update()
			case <-p.interrupt:
				return
			}
		}
	}()
}

func (p *UncertainProgressBar) Stop() {
	p.interrupt <- struct{}{}
}
