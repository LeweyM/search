package graphdraw

import "fmt"

type point struct {
	x, y int
}

type Square struct {
	x, y, size int
}

func (s Square) lines() []line {
	return []line{
		{a: point{x: s.x - padding, y: s.y + padding}, b: point{x: s.x + padding, y: s.y + padding}}, // top line
		{a: point{x: s.x - padding, y: s.y - padding}, b: point{x: s.x + padding, y: s.y - padding}}, // bottom line
		{a: point{x: s.x - padding, y: s.y - padding}, b: point{x: s.x - padding, y: s.y + padding}}, // left line
		{a: point{x: s.x + padding, y: s.y - padding}, b: point{x: s.x + padding, y: s.y + padding}}, // right line
	}
}

type character struct {
	x, y int
	char rune
}

func (s *Square) chars() []character {
	lines := s.lines()
	var res []character
	for _, l := range lines {
		res = append(res, getCharsFromLine(l)...)
	}
	return res
}

func getCharsFromLine(l line) []character {
	var res []character
	var ch rune

	o := l.getOrientation()
	if o == horizontal {
		ch = '─'
	} else {
		ch = '│'
	}

	// make sure that a.x is lower than b.x
	// make sure that a.y is lower than b.y
	a := l.a
	b := l.b
	if a.x > b.x {
		l.a, l.b = b, a
	}
	if a.y > b.y {
		l.a, l.b = b, a
	}

	x := a.x
	for x < b.x {
		res = append(res, character{
			x:    x,
			y:    a.y,
			char: ch,
		})
		x++
	}

	y := a.y
	for y < b.y {
		res = append(res, character{
			x:    a.x,
			y:    y,
			char: ch,
		})
		y++
	}

	return res
}

type line struct {
	a, b point
}

type orientation int

const (
	horizontal orientation = iota
	vertical
)

func (l line) getOrientation() orientation {
	// lots of assumptions here...
	if l.a.y > l.b.y {
		return vertical
	} else {
		return horizontal
	}
}

type Board struct {
	shapes []Square
}

func (b Board) draw() string {
	pixels := map[int]map[int]rune{}

	for _, shape := range b.shapes {
		for _, c := range shape.chars() {
			_, rowExists := pixels[c.y]
			if !rowExists {
				pixels[c.y] = map[int]rune{}
			}

			pixels[c.y][c.x] = c.char
		}
	}

	return fmt.Sprintf("%s", pixelsToString(pixels))
}

func pixelsToString(pixels map[int]map[int]rune) string {
	seen := 0

	return pixels
}
