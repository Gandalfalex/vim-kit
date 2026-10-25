package visual

import "testing"

func TestSelection(t *testing.T) {
	lines := []string{"alpha", "bravo", "charlie", "delta"}
	cases := []struct {
		name   string
		mode   Mode
		anchor Pos
		cursor Pos
		want   string
	}{
		{"line down", Line, Pos{1, 3}, Pos{2, 0}, "bravo\ncharlie"},
		{"line up", Line, Pos{2, 0}, Pos{0, 4}, "alpha\nbravo\ncharlie"},
		{"char same line", Char, Pos{2, 1}, Pos{2, 3}, "har"},
		{"char reversed", Char, Pos{2, 3}, Pos{2, 1}, "har"},
		{"char multi", Char, Pos{0, 3}, Pos{2, 1}, "ha\nbravo\nch"},
		{"char past end", Char, Pos{3, 2}, Pos{3, 99}, "lta"},
		{"clamped", Line, Pos{2, 0}, Pos{9, 0}, "charlie\ndelta"},
	}
	for _, c := range cases {
		s := Selection{Mode: c.mode, Anchor: c.anchor}
		if got := s.Text(c.cursor, lines); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}

	var s Selection
	if s.Text(Pos{}, lines) != "" || s.Active() {
		t.Error("none should select nothing")
	}
	s.Start(Line, Pos{1, 2})
	s.Start(Char, Pos{3, 0}) // switching keeps the anchor
	if s.Mode != Char || s.Anchor != (Pos{1, 2}) {
		t.Errorf("switch: %+v", s)
	}
	if !s.Contains(Pos{2, 1}, lines, 1, 2) || s.Contains(Pos{2, 1}, lines, 1, 1) || !s.Contains(Pos{2, 1}, lines, 2, 0) {
		t.Error("Contains char")
	}
	s.Toggle(Char, Pos{})
	if s.Active() {
		t.Error("toggle should clear")
	}
}
