// Package visual tracks vim visual-mode selections (V and v) over lines of
// text and extracts the selected text.
package visual

import "strings"

// Mode is the kind of selection.
type Mode int

const (
	None Mode = iota
	Line      // V: whole lines
	Char      // v: from one character to another
)

// Pos is a position in text; Col counts runes.
type Pos struct{ Line, Col int }

func (p Pos) before(q Pos) bool { return p.Line < q.Line || (p.Line == q.Line && p.Col < q.Col) }

// Selection is an active selection anchored where it started; the other end
// is the caller's cursor.
type Selection struct {
	Mode   Mode
	Anchor Pos
}

// Start begins (or switches to) mode at the cursor. Switching between Line
// and Char keeps the anchor line, as in vim.
func (s *Selection) Start(mode Mode, cursor Pos) {
	if s.Mode == None {
		s.Anchor = cursor
	}
	s.Mode = mode
}

// Toggle starts mode, or ends the selection when mode is already active.
func (s *Selection) Toggle(mode Mode, cursor Pos) {
	if s.Mode == mode {
		s.Clear()
		return
	}
	s.Start(mode, cursor)
}

// Clear ends the selection.
func (s *Selection) Clear() { *s = Selection{} }

// Active reports whether a selection is in progress.
func (s Selection) Active() bool { return s.Mode != None }

// Bounds returns the ordered ends of the selection, clamped to lines.
// For Line mode the columns are meaningless.
func (s Selection) Bounds(cursor Pos, lines []string) (from, to Pos) {
	from, to = s.Anchor, cursor
	if s.Mode == Line {
		from.Col, to.Col = 0, 0
	}
	if to.before(from) {
		from, to = to, from
	}
	from.Line = max(from.Line, 0)
	to.Line = min(to.Line, len(lines)-1)
	return from, to
}

// Contains reports whether line l (and, in Char mode, column c) is selected.
func (s Selection) Contains(cursor Pos, lines []string, l, c int) bool {
	if s.Mode == None {
		return false
	}
	from, to := s.Bounds(cursor, lines)
	if l < from.Line || l > to.Line {
		return false
	}
	if s.Mode == Line {
		return true
	}
	p := Pos{l, c}
	return !p.before(from) && !to.before(p)
}

// Text returns the selected text, lines joined with "\n".
func (s Selection) Text(cursor Pos, lines []string) string {
	if s.Mode == None || len(lines) == 0 {
		return ""
	}
	from, to := s.Bounds(cursor, lines)
	if from.Line > to.Line {
		return ""
	}
	if s.Mode == Line {
		return strings.Join(lines[from.Line:to.Line+1], "\n")
	}
	if from.Line == to.Line {
		return slice(lines[from.Line], from.Col, to.Col+1)
	}
	parts := []string{slice(lines[from.Line], from.Col, -1)}
	parts = append(parts, lines[from.Line+1:to.Line]...)
	parts = append(parts, slice(lines[to.Line], 0, to.Col+1))
	return strings.Join(parts, "\n")
}

// slice returns runes [i, j) of s, clamped; j < 0 means to the end.
func slice(s string, i, j int) string {
	r := []rune(s)
	if j < 0 || j > len(r) {
		j = len(r)
	}
	i = max(min(i, j), 0)
	return string(r[i:j])
}
