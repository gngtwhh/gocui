package pb

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
var DefaultProperty Property

func init() {
	var err error

	DefaultProperty = Property{
		Style: Style{
			Complete:        " ",
			Incomplete:      " ",
			CompleteColor:   font.WhiteBg,
			IncompleteColor: font.RESET,
		},
		Total: 100,
	}

	DefaultBar, err = NewProgressBar("%percent|%bar|%current/%total %elapsed %rate", WithDefault())
	if err != nil {
		panic(err)
	}
}

// Property is the default Property of the progress bar.
// A running instance will be initialized with the progress bar's current Property.
type Property struct {
	Format        string        // Format: The render format of the progress bar(with tokens).
	PosX, PosY    int           // Pos: The position of the progress bar on the screen
	Width         int           // Width: The width of the token:"%bar"
	Uncertain     bool          // Type: Whether the progress bar is Uncertain, default: false
	rate, elapsed time.Duration // Rate: rate is the rate of progress, elapsed is the elapsed time since call to Run()
	// NOTE: rate and elapsed will soon to be deprecated
	Style
	// for UnCertain:
	Total     int // Total: Only available when Uncertain is false
	direction int // Dct: 1(default) for increasing, -1 for decreasing, only available when UnCertain is true
}

// Style is the style struct in the Property struct, used to decorate the token "bar".
type Style struct {
	Complete, Incomplete, UnCertain                string // The style of the progress bar
	CompleteColor, IncompleteColor, UnCertainColor int    // The color of the progress bar
}

// ProgressBar is a simple progress bar implementation.
type ProgressBar struct {
	style    []token      // Parsed style tokens, will not be updated
	running  int          // The number of running instances
	rw       sync.RWMutex // RWMutex to synchronize access to the progress bar
	Property Property
}

// Context is a context created when the progress bar is running.
// Each progress bar instance can create several contexts for reuse.
type Context struct {
	property  Property      // Copy of static progress bar property
	style     []token       // Copy of static progress bar style
	current   int           // current progress
	interrupt chan struct{} // interrupt channel to stop running
	Done      chan struct{} // Done channel to signal completion
}

// NewProgressBar creates a new progress bar that with the given style and total.
func NewProgressBar(style string, mfs ...ModFunc) (pb *ProgressBar, err error) {
	if style == "" {
		return nil, fmt.Errorf("style cannot be empty")
	}

	property := Property{}
	for _, mf := range mfs {
		if mf == nil {
			return nil, fmt.Errorf("modify func cannot be nil")
		}
		mf(&property)
	}

	// revise the properties
	if property.Total == 0 {
		property.Total = 100 // default total is 100
	}
	//if property.Uncertain || (property.Current < 0 || property.Current > property.Total) {
	//	property.Current = 0
	//}
	if property.Uncertain {
		property.direction = 1
	}
	if property.Width <= 0 {
		property.Width = 20 // default width is 20
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
		style:    styleTokens,
		Property: property,
		//direction: 1,
		//interrupt: make(chan struct{}),
		//Done:      make(chan struct{}),
		rw: sync.RWMutex{},
	}
	//close(pb.Done) // Close p.Done initially to indicate that p is not running.
	return
}

//// NewProgressBar creates a new progress bar that with the given style and total.
//func NewProgressBar(style string, property Property) (pb *ProgressBar, err error) {
//	if style == "" {
//		return nil, fmt.Errorf("style cannot be empty")
//	}
//
//	// modify the properties
//	if property.Total == 0 {
//		property.Total = 100
//	}
//	if property.Uncertain || (property.Current < 0 || property.Current > property.Total) {
//		property.Current = 0
//	}
//	if property.Width <= 0 {
//		property.Width = 20
//	}
//	if property.Style.Complete == "" {
//		property.Style.Complete = "#"
//	}
//	if property.Style.Incomplete == "" {
//		property.Style.Incomplete = "-"
//	}
//	if property.Style.UnCertain == "" {
//		property.Style.UnCertain = "<->"
//	}
//	if property.Style.CompleteColor == font.RESET {
//		property.Style.CompleteColor = font.White
//	}
//	if property.Style.IncompleteColor == font.RESET {
//		property.Style.IncompleteColor = font.LightBlack
//	}
//	if property.Style.UnCertainColor == font.RESET {
//		property.Style.UnCertainColor = font.White
//	}
//
//	// generate style tokens
//	styleTokens := unmarshalToken(style)
//	// create progress bar
//	pb = &ProgressBar{
//		style:     styleTokens,
//		Property:  property,
//		direction: 1,
//		interrupt: make(chan struct{}),
//		Done:      make(chan struct{}),
//		rw:        sync.RWMutex{},
//	}
//	close(pb.Done) // Close p.Done initially to indicate that p is not running.
//	return
//}

//// UpdateTo updates the current progress of the progress bar and prints it.
//// In addition, if the bar is running, it will stop.
//func (p *ProgressBar) UpdateTo(ctx Context,current int) {
//	if ctx.running {
//		p.interrupt <- struct{}{}
//		//p.running = false // Needless: p.Run() will set p.running to false after stopping.
//	}
//	<-p.Done // Wait for the previous run to finish
//	p.rw.Lock()
//	if p.Property.Uncertain {
//		p.Property.Current = min(current, p.Property.Width-len(p.Property.Style.UnCertain)) // UnCertain bar use width to update the current progress
//		if p.Property.Current == p.Property.Width-len(p.Property.Style.UnCertain) {
//			p.direction = -1
//		}
//	} else {
//		p.Property.Current = min(current, p.Property.Total) // common bar use total to update the current progress
//	}
//	p.rw.Unlock()
//	p.Print()
//}

// updateCurrent increase the current progress without printing the progress bar.
func (p *ProgressBar) updateCurrent(ctx *Context) {
	//p.rw.Lock()
	if ctx.property.Uncertain {
		ctx.current = min(ctx.current+ctx.property.direction, ctx.property.Width-len(ctx.property.Style.UnCertain)) // UnCertain bar use width to update the current progress
		ctx.current = max(ctx.current, 0)
		if ctx.current == ctx.property.Width-len(ctx.property.Style.UnCertain) || ctx.current == 0 {
			ctx.property.direction = -ctx.property.direction
		}
	} else {
		ctx.current = min(ctx.current+1, ctx.property.Total) // common bar use total to update the current progress
	}
	//p.rw.Unlock()
}

//// Update increase the current progress of the progress bar and prints it.
//// In addition, if the bar is running, it will stop.
//func (p *ProgressBar) Update() {
//	if p.running {
//		p.interrupt <- struct{}{}
//		//p.running = false // Needless: p.Run() will set p.running to false after stopping.
//	}
//	<-p.Done // Wait for the previous run to finish
//	p.updateCurrent()
//	p.Print()
//}

// Print prints the current progress of the progress bar.
func (p *ProgressBar) Print(ctx *Context) {
	//@change: needless lock
	//p.rw.RLock() // Read lock to avoid race condition
	//defer p.rw.RUnlock()

	payloadBuilder := strings.Builder{}
	for _, t := range ctx.style {
		payloadBuilder.WriteString(t.toString(ctx))
	}

	utils.ConsoleMutex.Lock() // Lock the cursor
	defer utils.ConsoleMutex.Unlock()
	{
		window.ClearArea(ctx.property.PosX, ctx.property.PosY, ctx.property.Width, 1)
		cursor.GotoXY(ctx.property.PosX, ctx.property.PosY)
		fmt.Printf("%s", payloadBuilder.String())
	}
}

//// Run starts the progress bar.
//func (p *ProgressBar) Run(period time.Duration) {
//	if p.running {
//		return
//	}
//
//	var update = func() {
//		p.rw.Lock()
//		defer p.rw.Unlock()
//
//		if p.Property.Uncertain {
//			p.Property.Current = p.Property.Current + p.direction
//			if p.Property.Current == p.Property.Width-len(p.Property.Style.UnCertain) {
//				p.direction = -1
//			}
//			if p.Property.Current == 0 {
//				p.direction = 1
//			}
//		} else {
//			p.Property.Current++
//		}
//		p.Property.elapsed += period
//	}
//
//	p.rw.Lock()
//	p.Property.rate = period
//	p.rw.Unlock()
//
//	ticker := time.NewTicker(period)
//	p.running = true
//	p.Done = make(chan struct{})
//	p.Print() // Print the initial state
//	go func() {
//		defer func() {
//			ticker.Stop()
//			p.running = false
//			close(p.Done) // close p.Done to broadcast completion signal
//		}()
//		for {
//			select {
//			case <-ticker.C:
//				p.Print()
//				if !p.Property.Uncertain && p.Property.Current == p.Property.Total ||
//					(p.Property.Uncertain && p.Property.Current == p.Property.Width) {
//					return
//				}
//				update()
//			case <-p.interrupt:
//				return
//			}
//		}
//	}()
//}

// Stop stops the progress bar.
func (p *ProgressBar) Stop(ctx *Context) {
	ctx.interrupt <- struct{}{}
	p.rw.Lock()
	p.running--
	p.rw.Unlock()
}

// Iter starts a progress bar iteration, and returns two channels:
// iter <-chan int: to be used to iterate over the progress bar;
// stop chan<- struct{}: to stop the progress bar.
// This method should not be used if the progress bar is uncertain,
// otherwise, the returned iter channel will be closed, and the stop channel will be nil.
func (p *ProgressBar) Iter() (iter <-chan int, stop chan<- struct{}) {
	ch := make(chan int)
	var ctx Context
	p.rw.Lock()
	{
		if p.Property.Uncertain {
			close(ch) // Uncertain bar cannot be iterated
			p.rw.Unlock()
			return ch, nil
		}

		style := make([]token, len(p.style))
		copy(style, p.style)
		ctx = Context{
			property:  p.Property,
			style:     style,
			current:   0,
			interrupt: make(chan struct{}),
			Done:      make(chan struct{}),
		}
		p.running++
	}
	p.rw.Unlock()

	go func() {
		defer close(ch)
		for i := ctx.current; i <= ctx.property.Total; i++ {
			p.Print(&ctx)
		BREAK:
			for {
				select {
				case ch <- i:
					break BREAK
				case <-ctx.interrupt:
					return
				}
			}
			p.updateCurrent(&ctx)
		}
	}()
	return ch, ctx.interrupt
}

// Run starts an uncertain progress bar, and returns a channel to stop the progress bar.
// period is the time interval between each update, pass 0 to use the default period(100ms).
// If the progress bar is not uncertain, the returned channel will be nil.
func (p *ProgressBar) Run(period time.Duration) (stop chan<- struct{}) {
	ch := make(chan int)
	var ctx Context
	p.rw.Lock()
	{
		if !p.Property.Uncertain {
			close(ch) // certain bar cannot be run automatically
			p.rw.Unlock()
			return nil
		}

		style := make([]token, len(p.style))
		copy(style, p.style)
		ctx = Context{
			property:  p.Property,
			style:     style,
			current:   0,
			interrupt: make(chan struct{}),
			Done:      make(chan struct{}),
		}
		if period == 0 {
			period = time.Millisecond * 100 // default period is 100ms
		}
		ctx.property.rate = period
		p.running++
	}
	p.rw.Unlock()

	ticker := time.NewTicker(period)
	go func() {
		for {
			select {
			case <-ticker.C:
				p.Print(&ctx)
				p.updateCurrent(&ctx)
			case _, ok := <-ctx.interrupt: // interrupt got a signal or closed
				if ok {
					close(ch)
				}
				ticker.Stop()
				return
			}
		}
	}()
	return ctx.interrupt
}
