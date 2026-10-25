// Package keys maps key presses to named actions, with user-configurable
// bindings, vim count prefixes (5j) and multi-key sequences (gg).
package keys

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

type Action string

type Bindings map[Action][]string

func (b Bindings) Clone() Bindings {
	c := make(Bindings, len(b))
	for a, ks := range b {
		c[a] = append([]string(nil), ks...)
	}
	return c
}

type Keymap struct {
	actions  map[string]Action
	prefixes map[string]bool
}

func NewKeymap(layers ...Bindings) (*Keymap, error) {
	merged := Bindings{}
	for _, l := range layers {
		for a, ks := range l {
			merged[a] = ks
		}
	}
	km := &Keymap{actions: map[string]Action{}, prefixes: map[string]bool{}}
	var errs []string
	for _, a := range sortedActions(merged) {
		for _, k := range merged[a] {
			if k == "" {
				continue
			}
			if prev, dup := km.actions[k]; dup && prev != a {
				errs = append(errs, fmt.Sprintf("key %q bound to both %q and %q", k, prev, a))
				continue
			}
			km.actions[k] = a
			if seq := Split(k); len(seq) > 1 {
				for i := 1; i < len(seq); i++ {
					km.prefixes[strings.Join(seq[:i], "")] = true
				}
			}
		}
	}
	for p := range km.prefixes {
		if a, ok := km.actions[p]; ok {
			errs = append(errs, fmt.Sprintf("key %q (%s) also starts a sequence", p, a))
		}
	}
	if len(errs) > 0 {
		sort.Strings(errs)
		return nil, fmt.Errorf("keymap: %s", strings.Join(errs, "; "))
	}
	return km, nil
}

// MustKeymap is NewKeymap for built-in bindings that are known to be valid.
func MustKeymap(layers ...Bindings) *Keymap {
	km, err := NewKeymap(layers...)
	if err != nil {
		panic(err)
	}
	return km
}

// Lookup returns the action bound to a complete key or sequence.
func (km *Keymap) Lookup(seq string) (Action, bool) {
	a, ok := km.actions[seq]
	return a, ok
}

// IsPrefix reports whether seq starts a longer bound sequence.
func (km *Keymap) IsPrefix(seq string) bool { return km.prefixes[seq] }

// Keys returns the keys bound to a, sorted, e.g. for help screens.
func (km *Keymap) Keys(a Action) []string {
	var ks []string
	for k, x := range km.actions {
		if x == a {
			ks = append(ks, k)
		}
	}
	sort.Strings(ks)
	return ks
}

// named are multi-character key names that are not sequences.
var named = map[string]bool{
	"up": true, "down": true, "left": true, "right": true, "enter": true, "esc": true,
	"tab": true, "space": true, "backspace": true, "delete": true, "insert": true,
	"home": true, "end": true, "pgup": true, "pgdown": true,
}

// Split breaks a binding into its keys: "gg" → [g g], "ctrl+d" → [ctrl+d].
func Split(k string) []string {
	if strings.Contains(k, "+") || named[k] || isFKey(k) || utf8.RuneCountInString(k) < 2 {
		return []string{k}
	}
	seq := make([]string, 0, len(k))
	for _, r := range k {
		seq = append(seq, string(r))
	}
	return seq
}

func isFKey(k string) bool {
	if len(k) < 2 || k[0] != 'f' {
		return false
	}
	for _, c := range k[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func sortedActions(b Bindings) []Action {
	as := make([]Action, 0, len(b))
	for a := range b {
		as = append(as, a)
	}
	sort.Slice(as, func(i, j int) bool { return as[i] < as[j] })
	return as
}
