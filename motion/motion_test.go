package motion

import (
	"testing"

	"vimkit/keys"
)

func ev(a keys.Action, count ...int) keys.Event {
	e := keys.Event{Action: a, Count: 1}
	if len(count) > 0 {
		e.Count, e.HasCount = count[0], true
	}
	return e
}

func TestDefaultBindingsValid(t *testing.T) {
	keys.MustKeymap(VerticalBindings(), HorizontalBindings())
}

func TestWindow(t *testing.T) {
	const n, h = 100, 20
	cases := []struct {
		start Window
		ev    keys.Event
		want  int
	}{
		{Window{}, ev(Down, 3), 3},
		{Window{}, ev(Up), 0},
		{Window{Cursor: 50}, ev(HalfPageDown), 60},
		{Window{Cursor: 50}, ev(PageUp), 30},
		{Window{}, ev(Bottom), 99},
		{Window{}, ev(Bottom, 42), 41},
		{Window{Cursor: 50}, ev(Top), 0},
		{Window{Cursor: 50}, ev(Top, 7), 6},
		{Window{Cursor: 50, Top: 40}, ev(ScreenTop), 40},
		{Window{Cursor: 50, Top: 40}, ev(ScreenMiddle), 50},
		{Window{Cursor: 50, Top: 40}, ev(ScreenBottom), 59},
		{Window{Cursor: 95}, ev(Down, 10), 99},
		{Window{Cursor: 0, Top: 0}, ev(ScrollDown, 5), 5},
	}
	for _, c := range cases {
		w := c.start
		if !w.Apply(c.ev, n, h) {
			t.Fatalf("%s not a motion", c.ev.Action)
		}
		if w.Cursor != c.want {
			t.Errorf("%+v %s: cursor %d, want %d", c.start, c.ev.Action, w.Cursor, c.want)
		}
	}
	if (&Window{}).Apply(ev("yank"), n, h) {
		t.Error("yank is not a motion")
	}

	w := Window{Cursor: 50}
	w.Apply(ev(Center), n, h)
	if w.Top != 40 {
		t.Errorf("center: top %d, want 40", w.Top)
	}
	w.Cursor = 90
	w.Scroll(n, h)
	if w.Top != 73 {
		t.Errorf("scrolloff: top %d, want 73", w.Top)
	}
}

type fakePager struct{ y, lines int }

func (p *fakePager) ScrollDown(n int) { p.y = min(p.y+n, p.lines-1) }
func (p *fakePager) ScrollUp(n int)   { p.y = max(p.y-n, 0) }
func (p *fakePager) HalfPageDown()    { p.ScrollDown(5) }
func (p *fakePager) HalfPageUp()      { p.ScrollUp(5) }
func (p *fakePager) PageDown()        { p.ScrollDown(10) }
func (p *fakePager) PageUp()          { p.ScrollUp(10) }
func (p *fakePager) SetYOffset(n int) { p.y = n }
func (p *fakePager) Top()             { p.y = 0 }
func (p *fakePager) Bottom()          { p.y = p.lines - 1 }

func TestPager(t *testing.T) {
	cases := []struct {
		ev   keys.Event
		want int
	}{
		{ev(Down, 3), 23},
		{ev(ScrollUp), 19},
		{ev(HalfPageDown, 2), 30},
		{ev(PageUp), 10},
		{ev(Top), 0},
		{ev(Bottom), 99},
		{ev(Bottom, 5), 4},
	}
	for _, c := range cases {
		p := &fakePager{y: 20, lines: 100}
		if !ApplyPager(p, c.ev) || p.y != c.want {
			t.Errorf("%s: y %d, want %d", c.ev.Action, p.y, c.want)
		}
	}
	if ApplyPager(&fakePager{}, ev(ScreenTop)) {
		t.Error("pager has no screen motions")
	}
}

func TestText(t *testing.T) {
	const line = "foo.bar  baz_qux ä!"
	//            0123456789012345678
	cases := []struct {
		ev   keys.Event
		col  int
		want int
	}{
		{ev(WordForward), 0, 3},
		{ev(WordForward), 3, 4},
		{ev(WordForward), 4, 9},
		{ev(WordForward, 2), 0, 4},
		{ev(WordForward), 17, 18},
		{ev(WordBackward), 9, 4},
		{ev(WordBackward), 4, 3},
		{ev(WordBackward), 0, 0},
		{ev(WordEnd), 0, 2},
		{ev(WordEnd), 4, 6},
		{ev(WordEnd), 6, 15},
		{ev(LineStart), 10, 0},
		{ev(LineEnd), 0, 18},
		{ev(Left), 0, 0},
		{ev(Right, 50), 0, 18},
	}
	for _, c := range cases {
		got, ok := ApplyText(c.ev, line, c.col)
		if !ok || got != c.want {
			t.Errorf("%s from %d: got %d, want %d", c.ev.Action, c.col, got, c.want)
		}
	}
	if got, _ := ApplyText(ev(WordForward), "", 0); got != 0 {
		t.Errorf("empty line: %d", got)
	}
	if ClampCol("ab", 9) != 1 || ClampCol("", 3) != 0 {
		t.Error("ClampCol")
	}
}
