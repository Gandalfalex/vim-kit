package cmdline

// Line is the state of an open command line besides its text: the completion
// menu with the selected entry, and the position in the history.
type Line struct {
	History []string // executed lines, oldest first
	Items   []Item   // completion candidates for the current text
	Index   int      // selected candidate, -1 for none

	hist int // history position; len(History) is the line being typed
}

// Open starts a new line: no selection, history at the newest end.
func (l *Line) Open(items []Item) {
	l.hist = len(l.History)
	l.SetItems(items)
}

// Close drops the completion menu; the history stays.
func (l *Line) Close() { l.SetItems(nil) }

// SetItems replaces the candidates, e.g. after the text changed.
func (l *Line) SetItems(items []Item) { l.Items, l.Index = items, -1 }

// Cycling reports whether a candidate is selected, i.e. the text on the line
// was put there by Cycle and the menu still applies to it.
func (l *Line) Cycling() bool { return len(l.Items) > 0 && l.Index >= 0 }

// Cycle selects the next (dir 1) or previous (dir -1) candidate, wrapping
// around, and returns the text to put on the line.
func (l *Line) Cycle(dir int) (string, bool) {
	n := len(l.Items)
	if n == 0 {
		return "", false
	}
	if l.Index < 0 && dir < 0 {
		l.Index = 0 // going back from "none" lands on the last candidate
	}
	l.Index = ((l.Index+dir)%n + n) % n
	return l.Items[l.Index].Value, true
}

// Recall moves through the history, older with dir -1 and newer with dir 1,
// and returns the text to put on the line. Past the newest entry the line
// is empty again.
func (l *Line) Recall(dir int) (string, bool) {
	if len(l.History) == 0 {
		return "", false
	}
	l.hist = min(max(l.hist+dir, 0), len(l.History))
	if l.hist == len(l.History) {
		return "", true
	}
	return l.History[l.hist], true
}

// Commit adds an executed line to the history, skipping empty lines and
// immediate repeats.
func (l *Line) Commit(line string) {
	if n := len(l.History); line != "" && (n == 0 || l.History[n-1] != line) {
		l.History = append(l.History, line)
	}
}
