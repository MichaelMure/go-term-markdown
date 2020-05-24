package markdown

import "strconv"

type headingNumbering struct {
	levels [5]int
}

// Observe register the event of a new level with the given depth and
// adjust the numbering accordingly
func (hn *headingNumbering) Observe(level int) {
	if level <= 0 {
		panic("level start at 1, to mimic h1,h2 ... from html")
	}
	if level > 6 {
		panic("Markdown is limited to 6 levels of heading")
	}

	if level == 1 {
		// only one title with no numbering, nothing to count
		return
	}

	hn.levels[level-2]++
	for i := level - 1; i < 5; i++ {
		hn.levels[i] = 0
	}
}

// Render render the current headings numbering.
func (hn *headingNumbering) Render() string {
	slice := hn.levels[:]

	// pop the last zero levels
	for i := 4; i >= 0; i-- {
		if hn.levels[i] != 0 {
			break
		}
		slice = slice[:len(slice)-1]
	}

	var result string

	for i := range slice {
		if i > 0 {
			result += "."
		}
		result += strconv.Itoa(slice[i])
	}

	return result
}
