package progress_bar

import (
	"fmt"
	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/window"
	"strings"
	"sync"
	"time"
)

// ProgressBar is a simple progress bar implementation.
type ProgressBar struct {
	Total, Current int
	posX, posY     int
	width, height  int
	running        bool          // Whether the progress bar is running
	interrupt      chan struct{} // interrupt channel to stop Run()
	Done           chan struct{} // Done channel to signal completion
	rw             sync.RWMutex  // RWMutex to synchronize access to the progress bar
}

// NewProgressBar creates a new progress bar with the given total.
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		Total:     total,
		interrupt: make(chan struct{}),
		Done:      make(chan struct{}),
		rw:        sync.RWMutex{},
	}
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
	cursor.CursorMutex.Lock() // Lock the cursor to avoid concurrent access

	fmt.Printf("%s", cursor.HideCursor())
	window.ClearArea(p.posX, p.posY, p.width, p.height)
	for i := 0; i < p.height; i++ {
		fmt.Printf("%s", cursor.GotoXY(p.posX+i, p.posY))
		block := int(float64(p.Current) / float64(p.Total) * float64(p.width))
		fmt.Printf("%s", strings.Repeat("█", block)+strings.Repeat(" ", p.width-block))
	}
	cursor.CursorMutex.Unlock() // Unlock the cursor
	p.rw.RUnlock()              // Read unlock
}

// Run starts the progress bar.
func (p *ProgressBar) Run(period time.Duration) {
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
				if p.Current == p.Total {
					p.Done <- struct{}{}
					return
				}
				p.rw.Lock()
				p.Current++
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
	cursor.CursorMutex.Lock() // Lock the cursor to avoid concurrent access
	fmt.Printf("%s", cursor.HideCursor())
	window.ClearArea(p.posX, p.posY, p.width, p.height)
	fmt.Printf("%s", cursor.GotoXY(p.posX, p.posY))

	leftSpace := p.Current
	rightSpace := p.width - 2 - leftSpace - p.blockSize
	bar := "[" + strings.Repeat(" ", leftSpace) + strings.Repeat("█", p.blockSize) + strings.Repeat(" ", rightSpace) + "]"
	for i := 0; i < p.height; i++ {
		fmt.Printf("%s", cursor.GotoXY(p.posX+i, p.posY))
		fmt.Printf("%s", bar)
	}
	cursor.CursorMutex.Unlock() // Unlock the cursor
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
