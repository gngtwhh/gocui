package pb

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
// otherwise, it will refresh line by line(by default).
// param x, y: the position of the progress bar on the screen, must be within [0, screen width/height),
// if x or y is out of range, it will not be set.
func WithPos(x, y int) ModFunc {
	return func(p *Property) {
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

// WithCount sets the total count of the progress bar to be iterated.
// The supplied count argument c must be a positive integer, default is 100.
func WithCount(c int) ModFunc {
	return func(p *Property) {
		p.Total = c
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
// The format string will be parsed before rendering.
func WithFormat(f string) ModFunc {
	return func(p *Property) {
		p.Format = f
	}
}
