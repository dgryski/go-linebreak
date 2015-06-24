// Package linebreak wraps text at a given width
package linebreak

import "strings"

// This code is a translation of `linear` from http://xxyxyz.org/line-breaking/
// https://en.wikipedia.org/wiki/SMAWK_algorithm
// A. Aggarwal, T. Tokuyama. Consecutive interval query and dynamic programming on intervals. Discrete Applied Mathematics 85, 1998.

// Wrap formats text at the given width in linear time
func Wrap(text string, width int) string {
	words := strings.Fields(text)
	count := len(words)
	offsets := []int{0}
	for i, w := range words {
		offsets = append(offsets, offsets[i]+len(w))
	}

	minima := make([]int, count+1)
	for i := 1; i < len(minima); i++ {
		minima[i] = 1000000000000000000
	}
	breaks := make([]int, count+1)

	// closes over offsets, minima
	cost := func(i, j int) int {
		w := offsets[j] - offsets[i] + j - i - 1
		if w > width {
			c := 10000000000 * (w - width)
			return c
		}
		c := minima[i] + (width-w)*(width-w)
		return c
	}

	var smawk func([]int, []int)
	// smawk closes over cost, minima, breaks
	smawk = func(rows, columns []int) {
		var stack []int
		i := 0
		for i < len(rows) {
			if len(stack) > 0 {
				c := columns[len(stack)-1]
				if cost(peek(stack), c) < cost(rows[i], c) {
					if len(stack) < len(columns) {
						stack = push(stack, rows[i])
					}
					i++
				} else {
					stack = pop(stack)
				}
			} else {
				stack = push(stack, rows[i])
				i++
			}
		}
		rows = stack

		if len(columns) > 1 {
			smawk(rows, step(columns[1:], 2))
		}

		i = 0
		var j int
		for j < len(columns) {
			var end int
			if j+1 < len(columns) {
				end = breaks[columns[j+1]]
			} else {
				end = rows[len(rows)-1]
			}
			c := cost(rows[i], columns[j])
			if c < minima[columns[j]] {
				minima[columns[j]] = c
				breaks[columns[j]] = rows[i]
			}
			if rows[i] < end {
				i++
			} else {
				j += 2
			}
		}
	}

	n := count + 1
	i := 0
	offset := 0
	for {
		r := min(n, 1<<uint(i+1))
		edge := (1 << uint(i)) + offset
		smawk(genrange(0+offset, edge), genrange(edge, r+offset))
		x := minima[r-1+offset]
		// because python code has 'for ... else'
		var terminatedFor bool
		for _, j := range genrange(1<<uint(i), r-1) {
			y := cost(j+offset, r-1+offset)
			if y <= x {
				n -= j
				i = 0
				offset += j
				terminatedFor = true
				break
			}
		}
		if !terminatedFor {
			if r == n {
				break
			}
			i++
		}
	}

	var lines []string
	j := count
	for j > 0 {
		i = breaks[j]
		lines = append(lines, strings.Join(words[i:j], " "))
		j = i
	}

	// reverse lines
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	return strings.Join(lines, "\n")
}

// trivial int stack
func push(s []int, i int) []int { return append(s, i) }
func pop(s []int) []int         { return s[:len(s)-1] }
func peek(s []int) int          { return s[len(s)-1] }

// python list[a::b]
func step(ints []int, step int) []int {
	var r []int

	for i := 0; i < len(ints); i += step {
		r = append(r, ints[i])
	}
	return r
}

// python range(a,b)
func genrange(start, stop int) []int {
	var r []int
	for i := start; i < stop; i++ {
		r = append(r, i)
	}
	return r
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
