package markdown

import "gopkg.in/russross/blackfriday.v2"

func Render(source string, lineWidth int, leftPad int) []byte {
	renderer := newRenderer(lineWidth, leftPad)

	return blackfriday.Run([]byte(source), blackfriday.WithRenderer(renderer))
}
