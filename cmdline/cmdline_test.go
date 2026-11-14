package cmdline

import (
	"strings"
	"testing"
)

type app struct{ hosts []string }

var set = Set[*app, string]{
	{Name: "bucket", Aliases: []string{"b"}, Args: "<name>", Desc: "switch bucket",
		Complete: func(a *app, arg string) []Item { return MatchStrings(a.hosts, arg) },
		Run:      func(_ *app, arg string) string { return "bucket:" + arg }},
	{Name: "buckets", Run: func(*app, string) string { return "buckets" }},
	{Name: "quit", Aliases: []string{"q", "q!"}, Run: func(*app, string) string { return "quit" }},
}

func TestResolve(t *testing.T) {
	cases := map[string]string{
		"bucket": "bucket", "b": "bucket", "buckets": "buckets", "q!": "quit", "qu": "quit",
	}
	for in, want := range cases {
		c, err := set.Resolve(in)
		if err != nil || c.Name != want {
			t.Errorf("%q: got %v, %v; want %s", in, c, err, want)
		}
	}
	if _, err := set.Resolve("buc"); err == nil || !strings.Contains(err.Error(), "E464") {
		t.Errorf("buc should be ambiguous, got %v", err)
	}
	for _, in := range []string{"nope", ""} {
		if _, err := set.Resolve(in); err == nil || !strings.Contains(err.Error(), "E492") {
			t.Errorf("%q should be unknown, got %v", in, err)
		}
	}
}

func TestSplitAndLineNumber(t *testing.T) {
	if name, arg := Split("  cd   a b  "); name != "cd" || arg != "a b" {
		t.Errorf("split: %q %q", name, arg)
	}
	if n, ok := LineNumber("42"); !ok || n != 42 {
		t.Errorf("42: %d %v", n, ok)
	}
	if _, ok := LineNumber("q"); ok {
		t.Error("q is not a line number")
	}
}

func values(items []Item) string {
	var vs []string
	for _, it := range items {
		vs = append(vs, it.Value)
	}
	return strings.Join(vs, "|")
}

func TestComplete(t *testing.T) {
	a := &app{hosts: []string{"prod", "staging", "reprod"}}
	cases := map[string]string{
		"":         "bucket|buckets|quit",
		"bu":       "bucket|buckets",
		"q":        "quit", // exact alias
		"b":        "bucket|buckets",
		"b ":       "b prod|b staging|b reprod",
		"bucket p": "bucket prod|bucket reprod", // prefix match before substring
		"quit x":   "",
		"12":       "12",
	}
	for in, want := range cases {
		if got := values(set.Complete(a, in)); got != want {
			t.Errorf("%q: got %q, want %q", in, got, want)
		}
	}
	if got := set.Complete(a, "")[0].Label; got != "bucket (b) <name>" {
		t.Errorf("label %q", got)
	}
}

func TestLine(t *testing.T) {
	var l Line
	if _, ok := l.Cycle(1); ok {
		t.Error("cycle without items")
	}
	if _, ok := l.Recall(-1); ok {
		t.Error("recall without history")
	}

	l.Commit("one")
	l.Commit("one") // repeat
	l.Commit("")
	l.Commit("two")
	if len(l.History) != 2 {
		t.Fatalf("history %v", l.History)
	}

	l.Open([]Item{{Value: "a"}, {Value: "b"}})
	steps := []struct {
		dir  int
		want string
	}{{-1, "b"}, {-1, "a"}, {1, "b"}, {1, "a"}}
	for _, s := range steps {
		if got, _ := l.Cycle(s.dir); got != s.want {
			t.Errorf("cycle %d: %q, want %q", s.dir, got, s.want)
		}
	}
	hist := []struct {
		dir  int
		want string
	}{{-1, "two"}, {-1, "one"}, {-1, "one"}, {1, "two"}, {1, ""}, {1, ""}}
	for _, s := range hist {
		if got, _ := l.Recall(s.dir); got != s.want {
			t.Errorf("recall %d: %q, want %q", s.dir, got, s.want)
		}
	}
	if l.SetItems(nil); l.Index != -1 {
		t.Error("SetItems should clear the selection")
	}
}
