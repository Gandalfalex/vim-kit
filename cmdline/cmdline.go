// Package cmdline implements the vim ":" command line without the text
// widget: a set of ex-commands with aliases and unique-prefix lookup,
// completion candidates, and the history and completion-cycling state of an
// open line. The app owns the input field and decides what a command returns.
package cmdline

import (
	"fmt"
	"strconv"
	"strings"
)

// Item is one completion candidate.
type Item struct {
	Value string // text the command line becomes when the item is chosen
	Label string // text shown in a completion menu
	Desc  string // optional explanation next to the label
}

// Command is one ex-command. C is what handlers work on (usually the app's
// model) and R what running a command yields (e.g. a tea.Cmd).
type Command[C, R any] struct {
	Name     string
	Aliases  []string
	Args     string // usage of the argument, e.g. "<dir>"
	Desc     string
	Complete func(ctx C, arg string) []Item // candidates for the argument
	Run      func(ctx C, arg string) R
}

// Usage returns the command as shown in menus and help: "open (o) <name>".
func (c Command[C, R]) Usage() string {
	s := c.Name
	if len(c.Aliases) > 0 {
		s += " (" + strings.Join(c.Aliases, ",") + ")"
	}
	if c.Args != "" {
		s += " " + c.Args
	}
	return s
}

// Set is the list of commands an app understands, in display order.
type Set[C, R any] []Command[C, R]

// Error is returned by Resolve; its text follows vim's messages.
type Error struct {
	Name       string
	Candidates []string // set when Name is an ambiguous prefix
}

func (e *Error) Error() string {
	if len(e.Candidates) > 0 {
		return fmt.Sprintf("E464: Ambiguous use of command: %s (%s)", e.Name, strings.Join(e.Candidates, ", "))
	}
	return "E492: Not an editor command: " + e.Name
}

// Resolve finds a command by name or alias, else by unique name prefix
// (":buc" for ":bucket").
func (s Set[C, R]) Resolve(name string) (*Command[C, R], error) {
	for i := range s {
		if s[i].Name == name {
			return &s[i], nil
		}
		for _, a := range s[i].Aliases {
			if a == name {
				return &s[i], nil
			}
		}
	}
	var hit *Command[C, R]
	var names []string
	for i := range s {
		if name != "" && strings.HasPrefix(s[i].Name, name) {
			hit = &s[i]
			names = append(names, s[i].Name)
		}
	}
	switch len(names) {
	case 1:
		return hit, nil
	case 0:
		return nil, &Error{Name: name}
	}
	return nil, &Error{Name: name, Candidates: names}
}

// Split breaks a command line into the command name and its argument.
func Split(line string) (name, arg string) {
	name, arg, _ = strings.Cut(strings.TrimSpace(line), " ")
	return name, strings.TrimSpace(arg)
}

// LineNumber reports whether name is a line address, as in ":42".
func LineNumber(name string) (int, bool) {
	n, err := strconv.Atoi(name)
	return n, err == nil && n >= 0
}

// Complete returns the candidates for the text typed so far: command names
// while the name is being typed, then the command's own argument candidates.
// Every Value is a whole command line.
func (s Set[C, R]) Complete(ctx C, input string) []Item {
	name, arg, hasArg := strings.Cut(strings.TrimLeft(input, " "), " ")
	if !hasArg {
		if _, ok := LineNumber(name); ok {
			return []Item{{Value: name, Label: name, Desc: "jump to line " + name}}
		}
		var items, exact []Item
		for _, c := range s {
			it := Item{Value: c.Name, Label: c.Usage(), Desc: c.Desc}
			isAlias := false
			for _, a := range c.Aliases {
				isAlias = isAlias || a == name
			}
			if isAlias { // exact alias typed: surface its command first
				exact = append(exact, it)
			} else {
				items = append(items, it)
			}
		}
		return append(exact, Match(items, name)...)
	}
	c, err := s.Resolve(name)
	if err != nil || c.Complete == nil {
		return nil
	}
	items := c.Complete(ctx, arg)
	out := make([]Item, len(items))
	for i, it := range items {
		it.Value = name + " " + it.Value
		out[i] = it
	}
	return out
}

// Match keeps the items whose value starts with prefix (ignoring case),
// followed by those that merely contain it.
func Match(items []Item, prefix string) []Item {
	p := strings.ToLower(prefix)
	var starts, contains []Item
	for _, it := range items {
		v := strings.ToLower(it.Value)
		switch {
		case strings.HasPrefix(v, p):
			starts = append(starts, it)
		case p != "" && strings.Contains(v, p):
			contains = append(contains, it)
		}
	}
	return append(starts, contains...)
}

// MatchStrings is Match for plain values without descriptions.
func MatchStrings(vals []string, prefix string) []Item {
	items := make([]Item, len(vals))
	for i, v := range vals {
		items[i] = Item{Value: v, Label: v}
	}
	return Match(items, prefix)
}
