// Package motion implements vim cursor motions over lists, pagers and lines
// of text. Motions are keys.Actions, so apps bind and rebind them like any
// other action.
package motion

import "vimkit/keys"

// Vertical motions, for lists (Window) and pagers (Pager).
const (
	Down         keys.Action = "down"
	Up           keys.Action = "up"
	HalfPageDown keys.Action = "half_page_down"
	HalfPageUp   keys.Action = "half_page_up"
	PageDown     keys.Action = "page_down"
	PageUp       keys.Action = "page_up"
	ScrollDown   keys.Action = "scroll_down" // view moves, cursor stays (ctrl+e)
	ScrollUp     keys.Action = "scroll_up"
	Top          keys.Action = "top"    // or line N with a count
	Bottom       keys.Action = "bottom" // or line N with a count
	ScreenTop    keys.Action = "screen_top"
	ScreenMiddle keys.Action = "screen_middle"
	ScreenBottom keys.Action = "screen_bottom"
	Center       keys.Action = "center"        // zz
	CursorTop    keys.Action = "cursor_top"    // zt
	CursorBottom keys.Action = "cursor_bottom" // zb
)

// Horizontal motions within one line of text.
const (
	Left         keys.Action = "left"
	Right        keys.Action = "right"
	WordForward  keys.Action = "word_forward"
	WordBackward keys.Action = "word_backward"
	WordEnd      keys.Action = "word_end"
	LineStart    keys.Action = "line_start"
	LineEnd      keys.Action = "line_end"
)

// VerticalBindings returns the default vim keys for vertical motions.
func VerticalBindings() keys.Bindings {
	return keys.Bindings{
		Down:         {"j", "down", "ctrl+n"},
		Up:           {"k", "up", "ctrl+p"},
		HalfPageDown: {"ctrl+d"},
		HalfPageUp:   {"ctrl+u"},
		PageDown:     {"ctrl+f", "pgdown"},
		PageUp:       {"ctrl+b", "pgup"},
		ScrollDown:   {"ctrl+e"},
		ScrollUp:     {"ctrl+y"},
		Top:          {"gg", "home"},
		Bottom:       {"G", "end"},
		ScreenTop:    {"H"},
		ScreenMiddle: {"M"},
		ScreenBottom: {"L"},
		Center:       {"zz"},
		CursorTop:    {"zt"},
		CursorBottom: {"zb"},
	}
}

// HorizontalBindings returns the default vim keys for in-line motions.
func HorizontalBindings() keys.Bindings {
	return keys.Bindings{
		Left:         {"h", "left"},
		Right:        {"l", "right"},
		WordForward:  {"w"},
		WordBackward: {"b"},
		WordEnd:      {"e"},
		LineStart:    {"0"},
		LineEnd:      {"$"},
	}
}
