package pb

import "github.com/gngtwhh/gocui/window"

// ModFunc is a function that modifies the Property of the progress bar.
type ModFunc func(p *Property)

// WithDefault sets the default property of the progress bar.
// Should be used before any other ModFunc, otherwise it will overwrite all the other ModFuncs
// that in front of it.
func WithDefault() ModFunc {
	return func(p *Property) {
		*p = DefaultProperty
	}
}

// WithProperty sets the property of the progress bar to given property.
// NOTE: This will cause the current property of the progress bar to be overwritten.
func WithProperty(p Property) ModFunc {
	return func(p2 *Property) {
		*p2 = p
	}
}

// WithPos sets the position of the progress bar on the screen.
// If set, the progress bar will be placed at the specified position,
// otherwise, it will refresh at the current line of the cursor(by default).
// param x, y: the position of the progress bar on the screen, must be within [0, screen width/height),
// if x or y is out of range, it will not be set.
func WithPos(x, y int) ModFunc {
	return func(p *Property) {
		w, h := window.GetConsoleSize()
		if x < 0 || y < 0 || x >= h || y >= w {
			return
		}
		p.BindPos = true
		p.PosX = x
		p.PosY = y
	}
}

// WithWidth sets the width of the progress bar shown on the screen.
// The supplied width argument w must be a positive integer,
// should be larger than or equal to the actual width of the expected render.
func WithWidth(w int) ModFunc {
	return func(p *Property) {
		p.Width = w
	}
}

// WithUncertain sets the type of progress bar to uncertain.
func WithUncertain() ModFunc {
	return func(p *Property) {
		p.Uncertain = true
	}
}

// WithStyle sets the style of the progress bar.
func WithStyle(s Style) ModFunc {
	return func(p *Property) {
		p.Style = s
	}
}

// WithFormat sets the format of the progress bar.
// The format string will be parsed to tokens before rendering.
func WithFormat(f string) ModFunc {
	return func(p *Property) {
		p.Format = f
		p.formatChanged = true
	}
}

// WithWriter set the bar type as IO bar
func WithWriter() ModFunc {
	return func(p *Property) {
		p.Bytes = true
	}
}

// WithBarWidth set the width of the bar width
func WithBarWidth(w int) ModFunc {
	return func(p *Property) {
		p.BarWidth = w
	}
}
