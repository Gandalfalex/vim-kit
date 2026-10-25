package motion

import "vimkit/keys"

// Window is a cursor over n rows shown h rows at a time, starting at Top.
type Window struct {
	Cursor int
	Top    int
}

// Clamp keeps the cursor within [0, n).
func (w *Window) Clamp(n int) {
	w.Cursor = min(w.Cursor, n-1)
	w.Cursor = max(w.Cursor, 0)
}

// Move moves the cursor by d rows.
func (w *Window) Move(d, n int) {
	w.Cursor += d
	w.Clamp(n)
}

// Scroll adjusts Top so the cursor stays visible with a small scrolloff.
func (w *Window) Scroll(n, h int) {
	if h <= 0 {
		return
	}
	off := 2
	if h < 7 {
		off = 0
	}
	if w.Cursor < w.Top+off {
		w.Top = w.Cursor - off
	}
	if w.Cursor > w.Top+h-1-off {
		w.Top = w.Cursor - h + 1 + off
	}
	w.Top = min(w.Top, n-h)
	w.Top = max(w.Top, 0)
}

// Center puts the cursor in the middle of the view (zz).
func (w *Window) Center(n, h int) {
	w.Top = w.Cursor - h/2
	w.Scroll(n, h)
}

// Apply runs ev's motion over n rows with view height h and reports whether
// ev was a vertical motion. Call Scroll afterwards to keep the cursor in view.
func (w *Window) Apply(ev keys.Event, n, h int) bool {
	c := ev.Count
	switch ev.Action {
	case Down:
		w.Move(c, n)
	case Up:
		w.Move(-c, n)
	case HalfPageDown:
		w.Move(h/2*c, n)
	case HalfPageUp:
		w.Move(-h/2*c, n)
	case PageDown:
		w.Move(h*c, n)
	case PageUp:
		w.Move(-h*c, n)
	case ScrollDown:
		w.Top += c
		w.Cursor = max(w.Cursor, w.Top)
		w.Clamp(n)
	case ScrollUp:
		w.Top = max(0, w.Top-c)
		w.Cursor = min(w.Cursor, w.Top+h-1)
		w.Clamp(n)
	case Top:
		w.Cursor = lineOr(ev, 0)
		w.Clamp(n)
	case Bottom:
		w.Cursor = lineOr(ev, n-1)
		w.Clamp(n)
	case ScreenTop:
		w.Cursor = w.Top
		w.Clamp(n)
	case ScreenMiddle:
		w.Cursor = w.Top + min(h, n-w.Top)/2
		w.Clamp(n)
	case ScreenBottom:
		w.Cursor = w.Top + min(h, n-w.Top) - 1
		w.Clamp(n)
	case Center:
		w.Center(n, h)
	case CursorTop:
		w.Top = w.Cursor
		w.Scroll(n, h)
	case CursorBottom:
		w.Top = w.Cursor - h + 1
		w.Scroll(n, h)
	default:
		return false
	}
	return true
}

// lineOr returns the 0-based line named by ev's count, or def without one.
func lineOr(ev keys.Event, def int) int {
	if ev.HasCount {
		return ev.Count - 1
	}
	return def
}
