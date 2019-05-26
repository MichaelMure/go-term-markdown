package markdown

import (
	"github.com/MichaelMure/git-bug/util/colors"
)

var headingShades = []func(a ...interface{}) string{
	colors.GreenBold,
	colors.GreenBold,
	colors.HiGreen,
	colors.Green,
}

// Return the color function corresponding to the level.
// Beware, level start counting from 1.
func headingShade(level int) func(a ...interface{}) string {
	if level < 1 {
		level = 1
	}
	if level > len(headingShades) {
		level = len(headingShades)
	}
	return headingShades[level-1]
}

var quoteShades = []func(a ...interface{}) string{
	colors.GreenBold,
	colors.GreenBold,
	colors.HiGreen,
	colors.Green,
}

// Return the color function corresponding to the level.
func quoteShade(level int) func(a ...interface{}) string {
	if level < 1 {
		level = 1
	}
	if level > len(headingShades) {
		level = len(headingShades)
	}
	return headingShades[level-1]
}
