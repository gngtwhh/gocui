package pb

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gngtwhh/gocui/cursor"
	"github.com/gngtwhh/gocui/font"
	"github.com/gngtwhh/gocui/utils"
	"github.com/gngtwhh/gocui/window"
)

// DefaultBar is a pre-created default progress bar tokens
var DefaultBar *ProgressBar
var DefaultProperty Property
var DefaultBarFormat string

func init() {
	var err error

	DefaultBarFormat = "%percent|%bar|%current/%total %elapsed %rate"

	DefaultProperty = Property{
		Style: Style{
			Complete:        " ",
			Incomplete:      " ",
			CompleteColor:   font.WhiteBg,
			IncompleteColor: font.RESET,
		},
		Total: 100,
	}

	DefaultBar, err = NewProgressBar(DefaultBarFormat, WithDefault())
	if err != nil {
		panic(err)
	}
}

// Property is the default Property of the progress bar.
// A running instance will be initialized with the progress bar's current Property.
type Property struct {
	Style             // Style of the "bar" token, including characters and color
	Format     string // Format: The render format of the progress bar(with tokens).
	PosX, PosY int    // Pos: The position of the progress bar on the screen
	Width      int    // Width: The width of the token:"%bar"

	Total int64 // Total: Only available when Uncertain is false

	Uncertain bool // Type: Whether the progress bar is Uncertain, default: false
	Bytes     bool // Type: Whether the progress bar is used for bytes writer, default: false
	BindPos   bool // Whether bind the absolute pos, PosX and PosY are valid only when BindPos is true

	formatChanged bool // Indicates the change in format when updating property
}

// Style is the tokens struct in the Property struct, used to decorate the token "bar".
type Style struct {
	Complete, Incomplete, UnCertain                string // The tokens of the progress bar
	CompleteColor, IncompleteColor, UnCertainColor int    // The color of the progress bar
}

// ProgressBar is a simple progress bar implementation.
type ProgressBar struct {
	property Property
	tokens   []token      // Parsed tokens tokens, will not be updated
	running  int          // The number of running instances
	rw       sync.RWMutex // RWMutex to synchronize access to the progress bar
}

// Context is a context created when the progress bar is running.
// Each progress bar instance can create several contexts for reuse.
type Context struct {
	property  Property      // Copy of static progress bar property
	tokens    []token       // Copy of static progress bar tokens
	current   int64         // current progress,
	startTime time.Time     // start time
	interrupt chan struct{} // interrupt channel to stop running
	// direction: for UnCertain bar to update, 1(default) for increasing, -1 for decreasing, only available when UnCertain is true
	direction int
}

// BytesWriter implements io.Writer interface,
// used to receive data written to the bar and trigger render updates
type BytesWriter struct {
	// bytesChan is the context pointer to which the BytesWriter belongs
	bytesChan chan int
	// closeCh indicate that the BytesWriter has been closed
	closeCh chan struct{}
}

func NewBytesWriter() *BytesWriter {
	return &BytesWriter{
		bytesChan: make(chan int, 1),
		closeCh:   make(chan struct{}),
	}
}

// Write implement io.Writer interface
func (bw *BytesWriter) Write(b []byte) (n int, err error) {
	n = len(b)
	return n, bw.update(b)
}

// update send update current value to the channel to render next status
func (bw *BytesWriter) update(b []byte) error {
	select {
	case <-bw.closeCh:
		return errors.New("channel closed")
	default:
	}
	bw.bytesChan <- len(b)
	return nil
}

// close stop the BytesWriter producer
func (bw *BytesWriter) close() error {
	select {
	case <-bw.closeCh:
		return errors.New("the BytesWriter has been closed")
	default:
	}
	close(bw.closeCh)
	close(bw.bytesChan)
	return nil
}

// NewProgressBar creates a new progress bar that with the given tokens and total.
func NewProgressBar(style string, mfs ...ModFunc) (pb *ProgressBar, err error) {
	if style == "" {
		return nil, fmt.Errorf("tokens cannot be empty")
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
	// generate tokens tokens
	styleTokens := unmarshalToken(style)
	// create progress bar
	pb = &ProgressBar{
		tokens:   styleTokens,
		property: property,
		rw:       sync.RWMutex{},
	}
	return
}

// UpdateProperty updates Property of the progress bar.
// If the bar has instances running, it will return an error and do nothing.
func (p *ProgressBar) UpdateProperty(mfs ...ModFunc) (NewBar *ProgressBar, err error) {
	p.rw.Lock()
	defer p.rw.Unlock()

	property := &p.property
	for _, mf := range mfs {
		if mf == nil {
			return p, fmt.Errorf("modify func cannot be nil")
		}
		mf(property)
	}

	// revise the properties
	if property.Total == 0 {
		property.Total = 100 // default total is 100
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
	// generate tokens tokens
	if property.formatChanged || property.Format != "" {
		property.formatChanged = false
		p.tokens = unmarshalToken(property.Format)
	}
	return p, nil
}

// updateCurrent increase the current progress without printing the progress bar.
func (p *ProgressBar) updateCurrent(ctx *Context) {
	if ctx.property.Uncertain {
		ctx.current = min(ctx.current+int64(ctx.direction), int64(ctx.property.Width-len(ctx.property.Style.UnCertain))) // UnCertain bar use width to update the current progress
		ctx.current = max(ctx.current, 0)
		if ctx.current == int64(ctx.property.Width-len(ctx.property.Style.UnCertain)) || ctx.current == 0 {
			ctx.direction = -ctx.direction
		}
	} else {
		ctx.current = min(ctx.current+1, ctx.property.Total) // common bar use total to update the current progress
	}
}

// updateCurrentWithAdd increase the current progress by add
func (p *ProgressBar) updateCurrentWithAdd(ctx *Context, add int) {
	if ctx.property.Uncertain {
		ctx.current = min(ctx.current+int64(ctx.direction*add), int64(ctx.property.Width-len(ctx.property.Style.UnCertain))) // UnCertain bar use width to update the current progress
		ctx.current = max(ctx.current, 0)
		if ctx.current == int64(ctx.property.Width-len(ctx.property.Style.UnCertain)) || ctx.current == 0 {
			ctx.direction = -ctx.direction
		}
	} else {
		ctx.current = min(ctx.current+int64(add), ctx.property.Total) // common bar use total to update the current progress
	}
}

// Print prints the current progress of the progress bar.
func (p *ProgressBar) Print(ctx *Context) {
	payloadBuilder := strings.Builder{}
	for _, t := range ctx.tokens {
		payloadBuilder.WriteString(t.toString(ctx))
	}

	utils.ConsoleMutex.Lock() // Lock the cursor
	defer utils.ConsoleMutex.Unlock()
	{
		window.ClearLine(-1)
		if ctx.property.BindPos {
			cursor.GotoXY(ctx.property.PosX, ctx.property.PosY)
		} else {
			fmt.Print("\r")
		}
		fmt.Printf("%s", payloadBuilder.String())
	}
}

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
func (p *ProgressBar) Iter() (iter <-chan int64, stop chan<- struct{}) {
	ch := make(chan int64)
	var ctx Context
	p.rw.Lock()
	{
		if p.property.Uncertain {
			close(ch) // Uncertain bar cannot be iterated
			p.rw.Unlock()
			return ch, nil
		}

		style := make([]token, len(p.tokens))
		copy(style, p.tokens)
		ctx = Context{
			property:  p.property,
			tokens:    style,
			current:   0,
			interrupt: make(chan struct{}),
			direction: 1,
		}
		p.running++
	}
	p.rw.Unlock()

	ctx.startTime = time.Now()
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
		if !p.property.Uncertain {
			close(ch) // certain bar cannot be run automatically
			p.rw.Unlock()
			return nil
		}

		style := make([]token, len(p.tokens))
		copy(style, p.tokens)
		ctx = Context{
			property:  p.property,
			tokens:    style,
			current:   0,
			interrupt: make(chan struct{}),
			direction: 1,
		}
		if period == 0 {
			period = time.Millisecond * 100 // default period is 100ms
		}
		p.running++
	}
	p.rw.Unlock()

	ctx.startTime = time.Now()
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

// RunWithWriter automatically start a progress bar with writing bytes data
// returns a writer for user to write data and a stop channel indicate exit
func (p *ProgressBar) RunWithWriter() (writer *BytesWriter, stop chan<- struct{}) {
	ch := make(chan int)
	var ctx Context
	p.rw.Lock()
	{
		// check if the bar is with writer
		if !p.property.Bytes {
			close(ch)
			p.rw.Unlock()
			return
		}

		style := make([]token, len(p.tokens))
		copy(style, p.tokens)
		ctx = Context{
			property:  p.property,
			tokens:    style,
			current:   0,
			interrupt: make(chan struct{}),
			direction: 1,
		}
		p.running++
	}
	p.rw.Unlock()

	ctx.startTime = time.Now()

	bw := NewBytesWriter()

	go func() {
		for {
			select {
			case add := <-bw.bytesChan:
				p.Print(&ctx)
				p.updateCurrentWithAdd(&ctx, add)
			case _, ok := <-ctx.interrupt: // interrupt got a signal or closed
				if ok {
					close(ch)
				}
				bw.close()
				return
			}
		}
	}()
	return bw, ctx.interrupt
}
