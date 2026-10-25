package keys

import "strconv"

// Event is a complete key press or sequence with its count.
type Event struct {
	Seq      string // "j", "gg", "ctrl+d"
	Action   Action // "" when Seq is not bound
	Count    int    // always >= 1
	HasCount bool
}

// Input turns single key presses into Events, collecting count prefixes
// (5j, 10G) and multi-key sequences (gg) along the way.
type Input struct {
	NoCount bool // treat digits as plain keys
	count   string
	pending string
}

// Feed consumes one key against km. It returns the completed Event and true,
// or false while a count or sequence is still being typed.
//
// A digit starts a count unless it is bound on its own (like "1" switching
// tabs); "0" without a count is always a plain key, as in vim.
func (in *Input) Feed(km *Keymap, k string) (Event, bool) {
	if in.pending == "" && !in.NoCount && isDigit(k) && (in.count != "" || (k != "0" && !in.bound(km, k))) {
		if len(in.count) < 6 {
			in.count += k
		}
		return Event{}, false
	}
	seq := in.pending + k
	if km.IsPrefix(seq) {
		in.pending = seq
		return Event{}, false
	}
	ev := Event{Seq: seq, Count: 1, HasCount: in.count != ""}
	if ev.HasCount {
		ev.Count, _ = strconv.Atoi(in.count)
		ev.Count = max(ev.Count, 1)
	}
	ev.Action, _ = km.Lookup(seq)
	in.Reset()
	return ev, true
}

// Reset drops any partial count or sequence.
func (in *Input) Reset() { in.count, in.pending = "", "" }

// Partial returns typed-but-incomplete input for display, e.g. "5g".
func (in *Input) Partial() string { return in.count + in.pending }

func (in *Input) bound(km *Keymap, k string) bool {
	_, ok := km.Lookup(k)
	return ok
}

func isDigit(k string) bool { return len(k) == 1 && k[0] >= '0' && k[0] <= '9' }
