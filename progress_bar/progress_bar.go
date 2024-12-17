package progress_bar

import (
	"fmt"
	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/font"
	"github.com/gngtwhh/gocui/utils"
	"github.com/gngtwhh/gocui/window"
	"strings"
	"sync"
	"time"
)

// DefaultBar is a pre-created default progress bar style
var DefaultBar *ProgressBar

func init() {
	var err error
	DefaultBar, err = NewProgressBar("%percent|%bar|%current/%total %elapsed %rate", Property{
		Style: Style{
			Complete:        " ",
			Incomplete:      " ",
			CompleteColor:   font.WhiteBg,
			IncompleteColor: font.RESET,
		},
		Total: 100,
	})
	if err != nil {
		panic(err)
	}
}

// Property is the Property of the progress bar.
type Property struct {
	Total, Current int           // Only available when Uncertain is false
	PosX, PosY     int           // The position of the progress bar on the screen
	Width          int           // The Width of the token:"%bar"
	Uncertain      bool          // Whether the progress bar is Uncertain, default: false
	rate, elapsed  time.Duration // rate is the rate of progress, elapsed is the elapsed time since call to Run()
	Style
}

// Style is the style struct in the Property struct, used to decorate the token "bar".
type Style struct {
	Complete, Incomplete, UnCertain                string // The style of the progress bar
	CompleteColor, IncompleteColor, UnCertainColor int    // The color of the progress bar
}

// ModFunc is a function that modifies the Property of the progress bar.
type ModFunc func(p *Property)

// ProgressBar is a simple progress bar implementation.
type ProgressBar struct {
	// private fields
	style     []token
	running   bool          // Whether the progress bar is running
	direction int           // 1(default) for increasing, -1 for decreasing, only available when UnCertain is true
	interrupt chan struct{} // interrupt channel to stop Run()
	rw        sync.RWMutex  // RWMutex to synchronize access to the progress bar
	// public fields
	Property Property
	Done     chan struct{} // Done channel to signal completion
}

// NewProgressBar creates a new progress bar that with the given style and total.
func NewProgressBar(style string, property Property) (pb *ProgressBar, err error) {
	if style == "" {
		return nil, fmt.Errorf("style cannot be empty")
	}

	// modify the properties
	if property.Total == 0 {
		property.Total = 100
	}
	if property.Uncertain || (property.Current < 0 || property.Current > property.Total) {
		property.Current = 0
	}
	if property.Width <= 0 {
		property.Width = 20
	}
	if property.Style.Complete == "" {
		property.Style.Complete = "#"
	}
	if property.Style.Incomplete == "" {
		property.Style.Incomplete = "-"
	}
	if property.Style.UnCertain == "" {
		property.Style.UnCertain = "<->"
	}
	if property.Style.CompleteColor == font.RESET {
		property.Style.CompleteColor = font.White
	}
	if property.Style.IncompleteColor == font.RESET {
		property.Style.IncompleteColor = font.LightBlack
	}
	if property.Style.UnCertainColor == font.RESET {
		property.Style.UnCertainColor = font.White
	}

	// generate style tokens
	styleTokens := unmarshalToken(style)
	// create progress bar
	pb = &ProgressBar{
		style:     styleTokens,
		Property:  property,
		direction: 1,
		interrupt: make(chan struct{}),
		Done:      make(chan struct{}),
		rw:        sync.RWMutex{},
	}
	close(pb.Done) // Close p.Done initially to indicate that p is not running.
	return
}

// Update updates the current progress of the progress bar and prints it.
// In addition, if the bar is running, it will stop.
func (p *ProgressBar) Update(current int) {
	if p.running {
		p.interrupt <- struct{}{}
		//p.running = false // Needless: p.Run() will set p.running to false after stopping.
	}
	<-p.Done // Wait for the previous run to finish
	p.rw.Lock()
	if p.Property.Uncertain {
		p.Property.Current = min(current, p.Property.Width-len(p.Property.Style.UnCertain)) // UnCertain bar use width to update the current progress
		if p.Property.Current == p.Property.Width-len(p.Property.Style.UnCertain) {
			p.direction = -1
		}
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
			p.Property.Current = p.Property.Current + p.direction
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
	p.Done = make(chan struct{})
	p.Print() // Print the initial state
	go func() {
		defer func() {
			ticker.Stop()
			p.running = false
			close(p.Done) // close p.Done to broadcast completion signal
		}()
		for {
			select {
			case <-ticker.C:
				p.Print()
				if !p.Property.Uncertain && p.Property.Current == p.Property.Total ||
					(p.Property.Uncertain && p.Property.Current == p.Property.Width) {
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
