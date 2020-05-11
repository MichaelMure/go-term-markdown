package markdown

import "io"

// [text]: http://example.com/  "title"
//  |      |                     |
//  text   destination           title
type link struct {
	text        string
	destination string
	title       string
}

type linkRefs struct {
	links []link
}

func NewLinkRefs() *linkRefs {
	return &linkRefs{}
}

func (lr *linkRefs) Reset() {
	lr.links = nil
}

func (lr *linkRefs) Add(text, destination, title string) {
	lr.links = append(lr.links, link{
		text:        text,
		destination: destination,
		title:       title,
	})
}

func (lr *linkRefs) Render(w io.Writer, leftPad int, lineWidth int) {

}
