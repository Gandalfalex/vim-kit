package motion

import "github.com/Gandalfalex/vim-kit/keys"

// Pager is a scrollable text view without a cursor, e.g. a bubbles viewport
// behind a small adapter.
type Pager interface {
	ScrollDown(n int)
	ScrollUp(n int)
	HalfPageDown()
	HalfPageUp()
	PageDown()
	PageUp()
	SetYOffset(n int)
	Top()
	Bottom()
}

// ApplyPager runs ev's motion on p and reports whether ev was one it supports.
// With no cursor, cursor and view motions (j vs ctrl+e) both scroll.
func ApplyPager(p Pager, ev keys.Event) bool {
	c := ev.Count
	switch ev.Action {
	case Down, ScrollDown:
		p.ScrollDown(c)
	case Up, ScrollUp:
		p.ScrollUp(c)
	case HalfPageDown:
		times(c, p.HalfPageDown)
	case HalfPageUp:
		times(c, p.HalfPageUp)
	case PageDown:
		times(c, p.PageDown)
	case PageUp:
		times(c, p.PageUp)
	case Top, Bottom:
		switch {
		case ev.HasCount:
			p.SetYOffset(c - 1)
		case ev.Action == Top:
			p.Top()
		default:
			p.Bottom()
		}
	default:
		return false
	}
	return true
}

func times(n int, f func()) {
	for range n {
		f()
	}
}
