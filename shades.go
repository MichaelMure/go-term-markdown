package markdown

var defaultHeadingShades = []ShadeFmt{
	GreenBold,
	GreenBold,
	HiGreen,
	Green,
}

var defaultQuoteShades = []ShadeFmt{
	GreenBold,
	GreenBold,
	HiGreen,
	Green,
}

type ShadeFmt func(a ...interface{}) string

type levelShadeFmt func(level int) ShadeFmt

// Return a function giving the color function corresponding to the level.
// Beware, level start counting from 1.
func shade(shades []ShadeFmt) levelShadeFmt {
	return func(level int) ShadeFmt {
		if level < 1 {
			level = 1
		}
		if level > len(shades) {
			level = len(shades)
		}
		return shades[level-1]
	}
}
