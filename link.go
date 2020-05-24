package markdown

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	text "github.com/MichaelMure/go-term-text"
)

// link:
//
// foo [text][label] bar
//      |     |
//      text  label
//
//
// reference list:
//
// [label]: http://example.com/  "title"
//  |      |                     |
//  label   destination           title
type link struct {
	label       string
	destination string
	title       string
}

type linkRefs struct {
	labelCounter int
	links        []link
}

func NewLinkRefs() *linkRefs {
	return &linkRefs{
		labelCounter: 1,
	}
}

func (lr *linkRefs) Add(destination, title string) string {
	label := strconv.Itoa(lr.labelCounter)
	lr.labelCounter++

	lr.links = append(lr.links, link{
		label:       label,
		destination: destination,
		title:       title,
	})

	return label
}

func (lr *linkRefs) Count() int {
	return len(lr.links)
}

func (lr *linkRefs) Reset() {
	lr.links = nil
}

func (lr *linkRefs) Render(w io.Writer, leftPad int, lineWidth int) {
	if len(lr.links) == 0 {
		return
	}

	var content strings.Builder
	for _, l := range lr.links {
		content.WriteString("[")
		content.WriteString(l.label)
		content.WriteString("]: ")
		content.WriteString(l.destination)
		if len(l.title) > 0 {
			content.WriteString(" \"")
			content.WriteString(l.title)
			content.WriteString("\"")
		}
		content.WriteString("\n")
	}

	pad := strings.Repeat(" ", leftPad)
	out, _ := text.WrapWithPad(content.String(), lineWidth, pad)
	_, _ = fmt.Fprint(w, out, "\n")
}
