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
	p.rw.RLock() // Read lock to avoid race condition

	fmt.Printf("%s", cursor.HideCursor())
	window.ClearArea(p.posX, p.posY, p.width, p.height)
	fmt.Printf("%s", cursor.GotoXY(p.posX, p.posY))
	for i := 0; i < p.height; i++ {
		fmt.Printf("%s", cursor.GotoXY(p.posX+i, p.posY))
		block := int(float64(p.Current) / float64(p.Total) * float64(p.width))
		fmt.Printf("%s", strings.Repeat("█", block)+strings.Repeat(" ", p.width-block))
	}
	p.rw.RUnlock() // Read unlock
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
	posX, posY    int
	width, height int
	Interrupt     chan struct{}
	Done          chan struct{}
}
