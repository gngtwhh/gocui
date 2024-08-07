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

// Property is the Property of the progress bar.
type Property struct {
	Total, Current int           // Only available when Uncertain is false
	PosX, PosY     int           // The position of the progress bar on the screen
	Width          int           // The Width of the token:"%bar"
	Uncertain      bool          // Whether the progress bar is Uncertain
	rate, elapsed  time.Duration // rate is the rate of progress, elapsed is the elapsed time since call to Run()
	Style          struct {
		BarComplete, BarIncomplete, UnCertain string // The style of the progress bar
	}
}

// ModFunc is a function that modifies the Property of the progress bar.
type ModFunc func(p *Property)

// ProgressBar is a simple progress bar implementation.
type ProgressBar struct {
	// private fields
	style     []token
	running   bool          // Whether the progress bar is running
	direction int           // 1(default) for increasing, -1 for decreasing, only available when Uncertain is true
	interrupt chan struct{} // interrupt channel to stop Run()
	rw        sync.RWMutex  // RWMutex to synchronize access to the progress bar
	// public fields
	Property Property
	Done     chan struct{} // Done channel to signal completion
}

// NewProgressBar creates a new progress bar that with the given style and total.
func NewProgressBar(style string, mod ModFunc) (pb *ProgressBar, err error) {
	if style == "" {
		return nil, fmt.Errorf("style cannot be empty")
	}

	// modify the properties
	// default style
	property := Property{
		Total: 100,
		PosX:  0, PosY: 0, Width: 20,
		Style: struct {
			BarComplete, BarIncomplete, UnCertain string
		}{"#", "-", "<->"},
	}
	if mod != nil {
		mod(&property)
	}

	styleTokens := unmarshalToken(style)
	pb = &ProgressBar{
		style:     styleTokens,
		Property:  property,
		direction: 1,
		interrupt: make(chan struct{}),
		Done:      make(chan struct{}),
		rw:        sync.RWMutex{},
	}
	return
}

// Update updates the current progress of the progress bar.
// //If call by uncertain bar, it will just stop but do nothing.
func (p *ProgressBar) Update(current int) {
	if p.running {
		p.interrupt <- struct{}{}
	}
	p.rw.Lock()
	if !p.Property.Uncertain {
		p.Property.Current = min(current, p.Property.Width) // Uncertain bar use width to update the current progress
	} else {
		p.Property.Current = min(current, p.Property.Total) // common bar use total to update the current progress
	}
	p.rw.Unlock()
	p.Print()
}

// Print prints the current progress of the progress bar.
func (p *ProgressBar) Print() {
	p.rw.RLock()              // Read lock to avoid race condition
	utils.ConsoleMutex.Lock() // Lock the cursor to avoid concurrent access
	{
		cursor.HideCursor()
		window.ClearArea(p.Property.PosX, p.Property.PosY, p.Property.Width, 1)
		cursor.GotoXY(p.Property.PosX, p.Property.PosY)
		payloadBuilder := strings.Builder{}
		for _, t := range p.style {
			payloadBuilder.WriteString(t.toString(&p.Property))
		}
		fmt.Printf("%s", payloadBuilder.String())
	}
	utils.ConsoleMutex.Unlock() // Unlock the cursor
	p.rw.RUnlock()              // Read unlock
}

// Run starts the progress bar.
func (p *ProgressBar) Run(period time.Duration) {
	if p.running {
		return
	}

	var update = func() {
		p.rw.Lock()
		defer p.rw.Unlock()

		if p.Property.Uncertain {
			p.Property.Current = (p.Property.Current + p.direction) % (p.Property.Width)
			if p.Property.Current == p.Property.Width-len(p.Property.Style.UnCertain) {
				p.direction = -1
			}
			if p.Property.Current == 0 {
				p.direction = 1
			}
		} else {
			p.Property.Current++
		}
		p.Property.elapsed += period
	}

	p.rw.Lock()
	p.Property.rate = period
	p.rw.Unlock()

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
				if !p.Property.Uncertain && p.Property.Current == p.Property.Total ||
					(p.Property.Uncertain && p.Property.Current == p.Property.Width) {
					p.Done <- struct{}{}
					return
				}
				update()
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
