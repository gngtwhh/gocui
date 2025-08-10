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
var (
	DefaultBar       *ProgressBar
	DefaultBarFormat string
	DefaultProperty  Property
)

var (
	DefaultUncertainBar         *ProgressBar
	DefaultUncertainBarFormat   string
	DefaultUncertainBarProperty Property
)

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
		total: 100,
	}
	DefaultBar, err = NewProgressBar(DefaultBarFormat, WithDefault())
	if err != nil {
		panic(err)
	}

	DefaultUncertainBarFormat = "[%bar]"
	DefaultUncertainBarProperty = Property{
		Style: Style{
			Incomplete:     " ",
			UnCertain:      "   ",
			UnCertainColor: font.WhiteBg,
		},
	}
	DefaultUncertainBar, err = NewProgressBar(DefaultUncertainBarFormat, WithProperty(DefaultUncertainBarProperty), WithUncertain())
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
	BarWidth   int    // BarWidth: The render width of the token:"%bar"
	Width      int    // Width: The maximum width of the progress bar displayed on the terminal.

	total int64 // Total: Only available when Uncertain is false or Bytes is true

	Uncertain bool // Uncertain: Whether the progress bar is Uncertain, default: false
	Bytes     bool // Type: Whether the progress bar is used for bytes writer, default: false
	BindPos   bool // Whether bind the absolute pos, PosX and PosY are valid only when BindPos is true

	formatChanged bool // Indicates the change in format when updating property
}

// Style is the tokens struct in the Property struct, used to decorate the token "bar".
type Style struct {
	Complete, CompleteHead, Incomplete, UnCertain                     string // The "bar" token style
	CompleteColor, CompleteHeadColor, IncompleteColor, UnCertainColor int    // The color of the bar
}

// ProgressBar is a simple progress bar implementation.
type ProgressBar struct {
	property Property
	tokens   []token      // Parsed tokens tokens, will not be updated
	barPos   []int        // index of "%bar" in tokens
	running  int          // The number of running instances
	rw       sync.RWMutex // RWMutex to synchronize access to the progress bar
}

// context is a context created when the progress bar is running.
// Each progress bar instance can create several contexts for reuse.
type context struct {
	property Property // Copy of static progress bar property
	tokens   []token  // Copy of static progress bar tokens
	// barPos          []int         // Copy of static progress bar barPos
	current         int64         // current progress
	windowWidth     int           // window width, set by window.GetConsoleSize()
	WidthWithoutBar int           // accumulated render width without bar
	startTime       time.Time     // start time
	interrupt       chan struct{} // interrupt channel to stop running
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

func NewContext(p *ProgressBar) context {
	style := make([]token, len(p.tokens))
	copy(style, p.tokens)

	ctx := context{
		property: p.property,
		tokens:   style,
		// barPos:          DefaultBarPos,
		current:   0,
		startTime: time.Now(),
		interrupt: make(chan struct{}),
		direction: 1,
	}
	ctx.windowWidth, _ = window.GetConsoleSize()
	return ctx
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
	if property.total == 0 {
		property.total = 100 // default total is 100
	}
	if property.BarWidth < 0 {
		property.BarWidth = 0 // default 0 means full width
	}
	x, _ := window.GetConsoleSize()
	if property.Width <= 0 && property.Width > x {
		property.Width = x
	}
	if property.Style.Complete == "" {
		property.Style.Complete = "="
	}
	// if property.Style.CompleteHead == "" {
	// 	property.Style.CompleteHead = ">"
	// }
	if property.Style.Incomplete == "" {
		property.Style.Incomplete = "-"
	}
	if property.Style.UnCertain == "" {
		property.Style.UnCertain = "<->"
	}
	if property.Style.CompleteColor == font.RESET {
		property.Style.CompleteColor = font.White
	}
	if property.Style.CompleteHeadColor == font.RESET {
		property.Style.CompleteHeadColor = font.White
	}
	if property.Style.IncompleteColor == font.RESET {
		property.Style.IncompleteColor = font.LightBlack
	}
	if property.Style.UnCertainColor == font.RESET {
		property.Style.UnCertainColor = font.White
	}
	// generate tokens tokens
	styleTokens, barPos := unmarshalToken(style)
	// create progress bar
	pb = &ProgressBar{
		tokens:   styleTokens,
		property: property,
		barPos:   barPos,
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
	if property.total == 0 {
		property.total = 100 // default total is 100
	}
	if property.BarWidth < 0 {
		property.BarWidth = 0 // default 0 means full width
	}
	x, _ := window.GetConsoleSize()
	if property.Width <= 0 && property.Width > x {
		property.Width = x
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
		p.tokens, p.barPos = unmarshalToken(property.Format)
	}
	return p, nil
}

// updateCurrent increase the current progress without printing the progress bar.
func (p *ProgressBar) updateCurrent(ctx *context) {
	if ctx.property.Uncertain {
		// ctx.current = min(
		// 	ctx.current+int64(ctx.direction), int64(ctx.property.BarWidth-len(ctx.property.Style.UnCertain)),
		// ) // UnCertain bar use width to update the current progress
		// ctx.current = max(ctx.current, 0)
		// if ctx.current == int64(ctx.property.BarWidth-len(ctx.property.Style.UnCertain)) || ctx.current == 0 {
		// 	ctx.direction = -ctx.direction
		// }
		ctx.current = ctx.current + int64(ctx.direction)
	} else {
		// common bar use total to update the current progress
		ctx.current = min(ctx.current+1, ctx.property.total)
	}
}

// updateCurrentWithAdd increase the current progress by add
func (p *ProgressBar) updateCurrentWithAdd(ctx *context, add int) {
	if ctx.property.Uncertain {
		ctx.current = min(
			ctx.current+int64(ctx.direction*add), int64(ctx.property.BarWidth-len(ctx.property.Style.UnCertain)),
		) // UnCertain bar use width to update the current progress
		ctx.current = max(ctx.current, 0)
		if ctx.current == int64(ctx.property.BarWidth-len(ctx.property.Style.UnCertain)) || ctx.current == 0 {
			ctx.direction = -ctx.direction
		}
	} else {
		ctx.current = min(
			ctx.current+int64(add), ctx.property.total,
		) // common bar use total to update the current progress
	}
}

// Print prints the current progress of the progress bar.
func (p *ProgressBar) Print(ctx *context) {
	var payloadBuilder0 strings.Builder
	payloadBuilder1 := strings.Builder{}
	var barToken *TokenBar
	for _, t := range ctx.tokens {
		if _, ok := t.(*TokenBar); ok && barToken == nil {
			barToken = t.(*TokenBar)
			payloadBuilder0 = payloadBuilder1
			payloadBuilder1 = strings.Builder{}
		} else {
			payloadBuilder1.WriteString(t.toString(ctx))
		}
	}
	ctx.WidthWithoutBar = payloadBuilder0.Len() + payloadBuilder1.Len()
	barStr := barToken.toString(ctx)

	utils.ConsoleMutex.Lock() // Lock the cursor
	defer utils.ConsoleMutex.Unlock()
	{
		// window.ClearLine(-1)
		if ctx.property.BindPos {
			cursor.GotoXY(ctx.property.PosX, ctx.property.PosY)
		} else {
			fmt.Print("\r")
		}
		fmt.Print(payloadBuilder0.String())
		fmt.Print(barStr)
		fmt.Print(payloadBuilder1.String())
		window.ClearLineAfterCursor()
	}
}

// Stop stops the progress bar.
func (p *ProgressBar) stop(ctx *context) {
	// ctx.interrupt <- struct{}{}
	p.rw.Lock()
	p.running--
	p.rw.Unlock()
}

// iter starts a progress bar iteration, and returns two channels:
// iter <-chan int: to be used to iterate over the progress bar;
// stop chan<- struct{}: to stop the progress bar.
// This method should not be used if the progress bar is uncertain,
// otherwise, the returned iter channel will be closed, and the stop channel will be nil.
func (p *ProgressBar) iter(n int) (iter <-chan int64, stop chan<- struct{}) {
	ch := make(chan int64)
	var ctx context
	p.rw.Lock()
	{
		p.property.total = int64(n)
		if p.property.Uncertain {
			close(ch) // Uncertain bar cannot be iterated
			p.rw.Unlock()
			return ch, nil
		}
		ctx = NewContext(p)
		p.running++
	}
	p.rw.Unlock()

	ctx.startTime = time.Now()
	go func() {
		defer p.stop(&ctx)
		defer close(ch)
		for i := ctx.current; i <= ctx.property.total; i++ {
			p.Print(&ctx)
			select {
			case ch <- i:
			case <-ctx.interrupt:
				return
			}
			p.updateCurrent(&ctx)
		}
	}()
	return ch, ctx.interrupt
}

// Iter WILL BLOCK, start an default progress bar over the param function and render the bar.
// The bar will iterate n times and call f for each iteration.
func (p *ProgressBar) Iter(n int, f func()) {
	if n <= 0 || f == nil {
		return
	}
	it, stop := p.iter(n)
	for range it {
		f()
	}
	close(stop)
}

// Run starts an uncertain progress bar, and returns a channel to stop the progress bar.
// period is the time interval between each update, pass 0 to use the default period(100ms).
// If the progress bar is not uncertain, the returned channel will be nil.
func (p *ProgressBar) Run(period time.Duration) (stop chan<- struct{}) {
	var ctx context
	p.rw.Lock()
	{
		if !p.property.Uncertain {
			p.rw.Unlock()
			return nil
		}

		ctx = NewContext(p)
		if period == 0 {
			period = time.Millisecond * 100 // default period is 100ms
		}
		p.running++
	}
	p.rw.Unlock()

	ctx.startTime = time.Now()
	ticker := time.NewTicker(period)

	go func() {
		defer p.stop(&ctx)
		for {
			select {
			case <-ticker.C:
				p.Print(&ctx)
				p.updateCurrent(&ctx)
			case <-ctx.interrupt: // interrupt got a signal or closed
				ticker.Stop()
				return
			}
		}
	}()
	return ctx.interrupt
}

// RunWithWriter automatically start a progress bar with writing bytes data.
// param n: the total bytes to write.
// It returns a writer for user to write data and a stop channel indicate exit.
// This method should not be called if the progress bar is not with writer or is uncertain.
func (p *ProgressBar) RunWithWriter(n int64) (writer *BytesWriter, stop chan<- struct{}) {
	if p.property.Uncertain {
		return
	}
	var ctx context
	p.rw.Lock()
	{
		// check if the bar is with writer
		if !p.property.Bytes {
			p.rw.Unlock()
			return
		}
		p.property.total = n // n bytes to receive
		ctx = NewContext(p)
		p.running++
	}
	p.rw.Unlock()

	ctx.startTime = time.Now()

	bw := NewBytesWriter()

	go func() {
		defer p.stop(&ctx)
		p.Print(&ctx) // print 0
		for {
			select {
			case add := <-bw.bytesChan:
				p.updateCurrentWithAdd(&ctx, add)
				p.Print(&ctx)
				if ctx.current == ctx.property.total {
					bw.close()
					return
				}
			case <-ctx.interrupt: // interrupt got a signal or closed
				bw.close()
				return
			}
		}
	}()
	return bw, ctx.interrupt
}

// Go WILL BLOCK, start the uncertain bar over the param function and render the bar until f finish.
// This method will panic if the progress bar is not uncertain.
func (p *ProgressBar) Go(f func()) {
	stop := p.Run(0)
	if stop == nil {
		panic("progress bar is not uncertain")
	}
	f()
	close(stop)
}

// Go WILL BLOCK, start an default uncertain bar over the param function and render the bar until f finish.
// This method will panic if the progress bar is not uncertain.
func Go(f func()) {
	pb := DefaultUncertainBar
	stop := pb.Run(0)
	if stop == nil {
		panic("progress bar is not uncertain")
	}
	f()
	close(stop)
}

// Iter WILL BLOCK, start an default progress bar over the param function and render the bar.
// The bar will iterate n times and call f for each iteration.
func Iter(n int, f func()) {
	if n <= 0 || f == nil {
		return
	}
	p := DefaultBar
	p.Iter(n, f)
}
