package box

type ModFunc func(p *Property)

// WithDefault sets the default property of the box.
// Should be used before any other ModFunc, otherwise it will overwrite all the other ModFuncs
// that in front of it.
func WithDefault() ModFunc {
	return func(p *Property) {
		*p = DefaultProperty
	}
}

// WithProperty sets the property of the box to given property.
// NOTE: This will cause the current property of the box to be overwritten.
func WithProperty(p Property) ModFunc {
	return func(p2 *Property) {
		*p2 = p
	}
}

// WithPos sets the position of the box on the screen.
// param x, y: the position of the TopLeft corner of the box on the screen, must be within [0, screen width/height),
// if x or y is out of range, it will not be set.
func WithPos(x, y int) ModFunc {
	return func(p *Property) {
		p.PosX = x
		p.PosY = y
		p.BindPos = true
	}
}

// WithStyle sets the style of the box.
func WithStyle(s Style) ModFunc {
	return func(p *Property) {
		p.Style = s
	}
}
