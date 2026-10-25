package keys

import (
	"strings"
	"testing"
)

var testBindings = Bindings{
	"down":   {"j", "down"},
	"top":    {"gg"},
	"center": {"zz"},
	"tab1":   {"1"},
	"start":  {"0"},
	"quit":   {"ctrl+c"},
}

func feed(in *Input, km *Keymap, ks ...string) (Event, bool) {
	var ev Event
	var ok bool
	for _, k := range ks {
		ev, ok = in.Feed(km, k)
	}
	return ev, ok
}

func TestInput(t *testing.T) {
	km := MustKeymap(testBindings)
	cases := []struct {
		keys []string
		want Event
	}{
		{[]string{"j"}, Event{"j", "down", 1, false}},
		{[]string{"5", "j"}, Event{"j", "down", 5, true}},
		{[]string{"2", "0", "j"}, Event{"j", "down", 20, true}},
		{[]string{"0"}, Event{"0", "start", 1, false}},          // 0 alone is a key
		{[]string{"1"}, Event{"1", "tab1", 1, false}},           // bound digit wins
		{[]string{"5", "1", "j"}, Event{"j", "down", 51, true}}, // ...unless a count started
		{[]string{"g", "g"}, Event{"gg", "top", 1, false}},
		{[]string{"5", "g", "g"}, Event{"gg", "top", 5, true}},
		{[]string{"g", "x"}, Event{"gx", "", 1, false}}, // unbound sequence
		{[]string{"x"}, Event{"x", "", 1, false}},
		{[]string{"ctrl+c"}, Event{"ctrl+c", "quit", 1, false}},
	}
	for _, c := range cases {
		var in Input
		got, ok := feed(&in, km, c.keys...)
		if !ok || got != c.want {
			t.Errorf("%v: got %+v %v, want %+v", c.keys, got, ok, c.want)
		}
		if in.Partial() != "" {
			t.Errorf("%v: leftover %q", c.keys, in.Partial())
		}
	}

	var in Input
	if _, ok := feed(&in, km, "3", "z"); ok || in.Partial() != "3z" {
		t.Errorf("partial: ok=%v %q", ok, in.Partial())
	}
	in.Reset()
	if in.Partial() != "" {
		t.Error("reset left partial")
	}

	in = Input{NoCount: true}
	if ev, _ := feed(&in, km, "5"); ev.Seq != "5" || ev.HasCount {
		t.Errorf("NoCount: got %+v", ev)
	}
}

func TestKeymap(t *testing.T) {
	cases := []struct {
		name   string
		layers []Bindings
		err    string
	}{
		{"ok", []Bindings{testBindings}, ""},
		{"dup", []Bindings{{"a": {"x"}, "b": {"x"}}}, `key "x" bound to both "a" and "b"`},
		{"prefix", []Bindings{{"a": {"g"}, "b": {"gg"}}}, `key "g" (a) also starts a sequence`},
		{"override fixes dup", []Bindings{{"a": {"x"}, "b": {"y"}}, {"b": {"x"}}, {"a": {"z"}}}, ""},
		{"same action twice", []Bindings{{"a": {"x", "x"}}}, ""},
	}
	for _, c := range cases {
		_, err := NewKeymap(c.layers...)
		switch {
		case c.err == "" && err != nil:
			t.Errorf("%s: unexpected %v", c.name, err)
		case c.err != "" && (err == nil || !strings.Contains(err.Error(), c.err)):
			t.Errorf("%s: got %v, want %q", c.name, err, c.err)
		}
	}

	km := MustKeymap(testBindings, Bindings{"down": {}}) // unbind
	if _, ok := km.Lookup("j"); ok {
		t.Error("empty layer should unbind")
	}
	if got := MustKeymap(testBindings).Keys("down"); strings.Join(got, ",") != "down,j" {
		t.Errorf("Keys: %v", got)
	}
}

func TestSplit(t *testing.T) {
	cases := map[string]string{
		"j": "j", "gg": "g g", "zz": "z z", "ctrl+d": "ctrl+d", "up": "up",
		"pgdown": "pgdown", "f12": "f12", "shift+tab": "shift+tab", " ": " ", "ä": "ä",
	}
	for in, want := range cases {
		if got := strings.Join(Split(in), " "); got != want {
			t.Errorf("Split(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFromStruct(t *testing.T) {
	type nav struct {
		Up   []string `yaml:"up"`
		Skip []string `yaml:"-"`
	}
	cfg := struct {
		Nav    nav
		Search []string `yaml:"search,omitempty"`
		Plain  []string
		Other  int
	}{nav{Up: []string{"k"}}, []string{"?"}, []string{"p"}, 3}
	b := FromStruct(&cfg, "yaml")
	if len(b) != 3 || b["up"][0] != "k" || b["search"][0] != "?" || b["Plain"][0] != "p" {
		t.Errorf("FromStruct: %v", b)
	}
	if o := b.Only("up", "missing"); len(o) != 1 || o["up"][0] != "k" {
		t.Errorf("Only: %v", o)
	}
}
