package motion

import (
	"unicode"
	"unicode/utf8"

	"github.com/Gandalfalex/vim-kit/keys"
)

// ApplyText runs ev's in-line motion from column col (in runes) of line and
// returns the new column, or false when ev is not a horizontal motion.
func ApplyText(ev keys.Event, line string, col int) (int, bool) {
	n := utf8.RuneCountInString(line)
	switch ev.Action {
	case Left:
		col -= ev.Count
	case Right:
		col += ev.Count
	case WordForward:
		col = repeat(ev.Count, col, func(c int) int { return WordStartForward(line, c) })
	case WordBackward:
		col = repeat(ev.Count, col, func(c int) int { return WordStartBackward(line, c) })
	case WordEnd:
		col = repeat(ev.Count, col, func(c int) int { return WordEndForward(line, c) })
	case LineStart:
		col = 0
	case LineEnd:
		col = n - 1
	default:
		return col, false
	}
	return ClampCol(line, col), true
}

// ClampCol keeps col on a character of line (0 for an empty line), e.g.
// after moving the cursor to a shorter line.
func ClampCol(line string, col int) int {
	n := utf8.RuneCountInString(line)
	return max(min(col, n-1), 0)
}

func repeat(n, col int, step func(int) int) int {
	for range n {
		col = step(col)
	}
	return col
}

// charClass returns the vim word class: 0 space, 1 word [letters digits _],
// 2 punctuation.
func charClass(r rune) int {
	switch {
	case unicode.IsSpace(r):
		return 0
	case r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r):
		return 1
	}
	return 2
}

// WordStartForward returns the column of the next word start (vim w).
func WordStartForward(line string, col int) int {
	runes := []rune(line)
	n := len(runes)
	if col >= n-1 {
		return max(n-1, 0)
	}
	col = max(col, 0)
	i := col
	for cls := charClass(runes[col]); i < n && charClass(runes[i]) == cls; i++ {
	}
	for i < n && unicode.IsSpace(runes[i]) {
		i++
	}
	return min(i, n-1)
}

// WordStartBackward returns the column of the previous word start (vim b).
func WordStartBackward(line string, col int) int {
	runes := []rune(line)
	n := len(runes)
	if col <= 0 || n == 0 {
		return 0
	}
	i := min(col, n) - 1
	for i > 0 && unicode.IsSpace(runes[i]) {
		i--
	}
	for cls := charClass(runes[i]); i > 0 && charClass(runes[i-1]) == cls; i-- {
	}
	return i
}

// WordEndForward returns the column of the end of the current or next word
// (vim e).
func WordEndForward(line string, col int) int {
	runes := []rune(line)
	n := len(runes)
	if col >= n-1 {
		return max(n-1, 0)
	}
	i := max(col, 0) + 1
	for i < n && unicode.IsSpace(runes[i]) {
		i++
	}
	if i >= n {
		return n - 1
	}
	for cls := charClass(runes[i]); i < n-1 && charClass(runes[i+1]) == cls; i++ {
	}
	return i
}
