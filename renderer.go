package markdown

import (
	"bytes"
	"fmt"
	stdcolor "image/color"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/MichaelMure/go-term-text"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	mdrender "github.com/yuin/goldmark/renderer"
	mdtext "github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"golang.org/x/net/html"

	htmlWalker "github.com/MichaelMure/go-term-markdown/html"
)

/*

Here are the possible cases for the AST. You can render it using PlantUML.

@startuml

(*) --> Document
Blockquote --> Blockquote
Blockquote --> CodeBlock
Blockquote --> List
Blockquote --> Paragraph
CodeSpan --> Text
Document --> Blockquote
Document --> CodeBlock
Document --> FencedCodeBlock
Document --> Heading
Document --> HTMLBlock
Document --> List
Document --> Paragraph
Document --> Table
Document --> TextBlock
Document --> ThematicBreak
Emphasis --> Emphasis
Emphasis --> Text
Heading --> CodeSpan
Heading --> Emphasis
Heading --> Image
Heading --> Link
Heading --> RawHTML
Heading --> Strikethrough
Heading --> Text
Image --> Text
Link --> Image
Link --> Text
ListItem --> List
ListItem --> Paragraph
ListItem --> TextBlock
List --> ListItem
Paragraph --> AutoLink
Paragraph --> CodeSpan
Paragraph --> Emphasis
Paragraph --> Image
Paragraph --> Link
Paragraph --> RawHTML
Paragraph --> Strikethrough
Paragraph --> Text
Strikethrough --> Emphasis
Strikethrough --> Text
TableCell --> CodeSpan
TableCell --> Emphasis
TableCell --> Image
TableCell --> Link
TableCell --> RawHTML
TableCell --> Strikethrough
TableCell --> Text
TableHeader --> TableCell
TableRow --> TableCell
Table --> TableHeader
Table --> TableRow
TextBlock --> AutoLink
TextBlock --> CodeSpan
TextBlock --> Emphasis
TextBlock --> Link
TextBlock --> TaskCheckBox
TextBlock --> Text

@enduml

*/

var _ mdrender.NodeRenderer = &renderer{}

type renderer struct {
	// maximum line width allowed
	lineWidth int
	// constant left padding to apply
	leftPad int

	// Dithering mode for ansimage
	// Default is fine directly through a terminal
	// DitheringWithBlocks is recommended if a terminal UI library is used
	imageDithering ansimage.DitheringMode

	// all the custom left paddings, without the fixed space from leftPad
	padAccumulator []string

	// one-shot indent for the first line of the inline content
	indent string

	// for Heading, Paragraph, HTMLBlock and TableCell, accumulate the content of
	// the child nodes (Link, Text, Image, formatting ...). The result
	// is then rendered appropriately when exiting the node.
	inlineAccumulator strings.Builder

	// record and render the heading numbering
	headingNumbering headingNumbering
	headingShade     levelShadeFmt

	blockQuoteLevel int
	blockQuoteShade levelShadeFmt

	linksRefs *linkRefs

	table *tableRenderer
}

func (r *renderer) RegisterFuncs(reg mdrender.NodeRendererFuncRegisterer) {
	// blocks
	reg.Register(ast.KindDocument, r.render)
	reg.Register(ast.KindHeading, r.render)
	reg.Register(ast.KindBlockquote, r.render)
	reg.Register(ast.KindCodeBlock, r.render)
	reg.Register(ast.KindFencedCodeBlock, r.render)
	reg.Register(ast.KindHTMLBlock, r.render)
	reg.Register(ast.KindList, r.render)
	reg.Register(ast.KindListItem, r.render)
	reg.Register(ast.KindParagraph, r.render)
	reg.Register(ast.KindTextBlock, r.render)
	reg.Register(ast.KindThematicBreak, r.render)

	// inlines
	reg.Register(ast.KindAutoLink, r.render)
	reg.Register(ast.KindCodeSpan, r.render)
	reg.Register(ast.KindEmphasis, r.render)
	reg.Register(ast.KindImage, r.render)
	reg.Register(ast.KindLink, r.render)
	reg.Register(ast.KindRawHTML, r.render)
	reg.Register(ast.KindText, r.render)
	reg.Register(ast.KindString, r.render)

	// extensions
	reg.Register(extast.KindTaskCheckBox, r.render)
	reg.Register(extast.KindStrikethrough, r.render)
	reg.Register(extast.KindTable, r.render)
	reg.Register(extast.KindTableHeader, r.render)
	reg.Register(extast.KindTableRow, r.render)
	reg.Register(extast.KindTableCell, r.render)
	reg.Register(extast.KindTaskCheckBox, r.render)
}

func newRenderer(lineWidth int, leftPad int, opts ...Options) *renderer {
	r := &renderer{
		lineWidth:       lineWidth,
		leftPad:         leftPad,
		padAccumulator:  make([]string, 0, 10),
		headingShade:    shade(defaultHeadingShades),
		blockQuoteShade: shade(defaultQuoteShades),
		linksRefs:       NewLinkRefs(),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *renderer) pad() string {
	return strings.Repeat(" ", r.leftPad) + strings.Join(r.padAccumulator, "")
}

func (r *renderer) addPad(pad string) {
	r.padAccumulator = append(r.padAccumulator, pad)
}

func (r *renderer) popPad() {
	r.padAccumulator = r.padAccumulator[:len(r.padAccumulator)-1]
}

func (r *renderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {

	switch node := node.(type) {
	case *ast.Document:
		// Nothing to do

	case *ast.Blockquote:
		// set and remove a colored bar on the left
		if entering {
			r.blockQuoteLevel++
			r.addPad(r.blockQuoteShade(r.blockQuoteLevel)("┃ "))
		} else {
			r.blockQuoteLevel--
			r.popPad()
			// extra line break in some cases to separate elements in the terminal
			// rendering.
			if node.Parent().Kind() == ast.KindDocument && node.NextSibling() != nil {
				_, _ = fmt.Fprintln(w)
			}
		}

	case *ast.List:
		if !entering {
			// extra new line at the end of a list *if* there is something next and
			// we are not in a nested list
			_, parentIsListItem := node.Parent().(*ast.ListItem)
			if !parentIsListItem && node.NextSibling() != nil {
				_, _ = fmt.Fprintln(w)
			}
		}

	case *ast.ListItem:
		list := node.Parent().(*ast.List)

		// write the prefix, add a padding if needed, and let Paragraph handle the rest
		if entering {
			switch {
			// numbered list
			case list.IsOrdered():
				itemNumber := list.Start
				for sibling := list.FirstChild(); sibling != nil; sibling = sibling.NextSibling() {
					if sibling == node {
						break
					}
					itemNumber++
				}
				prefix := fmt.Sprintf("%d. ", itemNumber)
				r.indent = r.pad() + Green(prefix)
				r.addPad(strings.Repeat(" ", text.Len(prefix)))

			// normal bullet point list
			default:
				r.indent = r.pad() + Green("• ")
				r.addPad("  ")
			}
		} else {
			r.popPad()
		}

	case *extast.TaskCheckBox:
		if entering {
			if node.IsChecked {
				r.inlineAccumulator.WriteString("[X] ")
			} else {
				r.inlineAccumulator.WriteString("[ ] ")
			}
		}

	// TextBlock: block of text within something else (ListItem)
	case *ast.TextBlock:
		if isLinkRefDefinition(node) {
			if !entering {
				r.linksRefs.Render(w, r.leftPad, r.lineWidth)
				r.linksRefs.Reset()

				// extra line break in some cases to separate elements in the terminal
				// rendering.
				if node.Parent().Kind() == ast.KindDocument && node.NextSibling() != nil {
					_, _ = fmt.Fprintln(w)
				}
			}
			break
		}

		// on exiting, collect and format the accumulated content
		if !entering {
			content := r.inlineAccumulator.String()
			r.inlineAccumulator.Reset()

			var out string
			if r.indent != "" {
				out, _ = text.WrapWithPadIndent(content, r.lineWidth, r.indent, r.pad())
				r.indent = ""
			} else {
				out, _ = text.WrapWithPad(content, r.lineWidth, r.pad())
			}
			_, _ = fmt.Fprint(w, out, "\n")
		}

	// Paragraph: standalone block of text
	case *ast.Paragraph:
		// on exiting, collect and format the accumulated content
		if !entering {
			content := r.inlineAccumulator.String()
			r.inlineAccumulator.Reset()

			// emoji support !
			content = emoji.Sprint(content)

			var out string
			if r.indent != "" {
				out, _ = text.WrapWithPadIndent(content, r.lineWidth, r.indent, r.pad())
				r.indent = ""
			} else {
				out, _ = text.WrapWithPad(content, r.lineWidth, r.pad())
			}
			_, _ = fmt.Fprint(w, out, "\n")

			// extra line break in some cases to separate elements in the terminal
			// rendering.
			if node.Parent().Kind() == ast.KindDocument && node.NextSibling() != nil {
				_, _ = fmt.Fprintln(w)
			}
		}

	case *ast.Heading:
		if !entering {
			r.renderHeading(w, node.Level)
		}

	case *ast.ThematicBreak:
		if entering {
			r.renderHorizontalRule(w)
		}

	case *ast.Emphasis:
		switch node.Level {
		case 1: // italic
			if entering {
				r.inlineAccumulator.WriteString(italicOn)
			} else {
				r.inlineAccumulator.WriteString(italicOff)
			}
		case 2: // strong/bold
			if entering {
				r.inlineAccumulator.WriteString(boldOn)
			} else {
				// This is super silly but some terminals, instead of having
				// the ANSI code SGR 21 do "bold off" like the logic would guide,
				// do "double underline" instead. This is madness.

				// To resolve that problem, we take a snapshot of the escape state,
				// remove the bold, then output "reset all" + snapshot
				es := text.EscapeState{}
				es.Witness(r.inlineAccumulator.String())
				es.Bold = false
				r.inlineAccumulator.WriteString(resetAll)
				if !es.IsZero() {
					r.inlineAccumulator.WriteString(es.String())
				}
			}
		}

	case *extast.Strikethrough:
		if entering {
			r.inlineAccumulator.WriteString(crossedOutOn)
		} else {
			r.inlineAccumulator.WriteString(crossedOutOff)
		}

	case *ast.Link:
		// The spec say that the text

		// TODO: resolve this madness

		if entering {
			r.inlineAccumulator.WriteString("[")
		} else {
			r.inlineAccumulator.WriteString("](")
			r.inlineAccumulator.WriteString(Blue(string(node.Destination)))
			if len(node.Title) > 0 {
				r.inlineAccumulator.WriteString(" ")
				r.inlineAccumulator.WriteString(string(node.Title))
			}
			r.inlineAccumulator.WriteString(")")
		}

	case *ast.AutoLink:
		// TODO

	case *ast.Image:
	// TODO: only render if child of Document
	// if entering {
	// 	// the alt text/title is weirdly parsed and is actually
	// 	// a child text of this node
	// 	var title string
	// 	if len(node.Children) == 1 {
	// 		if t, ok := node.Children[0].(*ast.Text); ok {
	// 			title = string(t.Literal)
	// 		}
	// 	}
	//
	// 	str, rendered := r.renderImage(
	// 		string(node.Destination), title,
	// 		r.lineWidth-r.leftPad,
	// 	)
	//
	// 	if rendered {
	// 		r.inlineAccumulator.WriteString("\n")
	// 		r.inlineAccumulator.WriteString(str)
	// 		r.inlineAccumulator.WriteString("\n\n")
	// 	} else {
	// 		r.inlineAccumulator.WriteString(str)
	// 		r.inlineAccumulator.WriteString("\n")
	// 	}
	//
	// 	return ast.SkipChildren
	// }

	// a single inline piece of text
	case *ast.Text:
		val := node.Segment.Value(source)
		if entering {
			if node.IsRaw() {
				RawWrite(&r.inlineAccumulator, val)
			} else {
				Write(&r.inlineAccumulator, val)

				// the parser remove all new lines so we need to put them back on
				// softbreak.
				if node.HardLineBreak() || node.SoftLineBreak() {
					r.inlineAccumulator.WriteByte('\n')
				}
			}
		}

	case *ast.FencedCodeBlock:
		if entering {
			r.renderCodeBlock(w, source, node.Lines(), string(node.Language(source)))
		} else {
			// extra line break in some cases to separate elements in the terminal
			// rendering.
			if node.Parent().Kind() == ast.KindDocument && node.NextSibling() != nil {
				_, _ = fmt.Fprintln(w)
			}
		}

	case *ast.CodeBlock:
		if entering {
			r.renderCodeBlock(w, source, node.Lines(), "")
		} else {
			// extra line break in some cases to separate elements in the terminal
			// rendering.
			if node.Parent().Kind() == ast.KindDocument && node.NextSibling() != nil {
				_, _ = fmt.Fprintln(w)
			}
		}

	case *ast.CodeSpan:
		if entering {
			for c := node.FirstChild(); c != nil; c = c.NextSibling() {
				segment := c.(*ast.Text).Segment
				value := segment.Value(source)
				if bytes.HasSuffix(value, []byte("\n")) {
					r.inlineAccumulator.WriteString(BlueBgItalic(string(value[:len(value)-1])))
					if c != node.LastChild() {
						r.inlineAccumulator.Write([]byte(" "))
					}
				} else {
					r.inlineAccumulator.WriteString(BlueBgItalic(string(value)))
				}
			}
			return ast.WalkSkipChildren, nil
		}

	case *ast.HTMLBlock:
		if entering {
			r.renderHTMLBlock(w, source, node)
		}

	// something that looks like a HTML tag
	// Sadly even when concatenated they are parsed into independent pieces in the AST,
	// which sort of make sense because where does a "piece of HTML" end within a block
	// of text ? What happen when the closing tag don't match ?
	case *ast.RawHTML:
		if entering {
			content := accumulateSegments(node.Segments, source)
			r.inlineAccumulator.WriteString(Red(content))
		}

	case *extast.Table:
		if entering {
			r.table = newTableRenderer()
		} else {
			r.table.Render(w, r.leftPad, r.lineWidth)
			r.table = nil
		}

	case *extast.TableCell:
		if !entering {
			content := r.inlineAccumulator.String()
			r.inlineAccumulator.Reset()

			align := CellAlignLeft
			switch node.Alignment {
			case extast.AlignRight:
				align = CellAlignRight
			case extast.AlignCenter:
				align = CellAlignCenter
			}

			if node.Parent().Kind() == extast.KindTableHeader {
				r.table.AddHeaderCell(content, align)
			} else {
				r.table.AddBodyCell(content, CellAlignCopyHeader)
			}
		}

	case *extast.TableHeader:
		// nothing to do

	case *extast.TableRow:
		if entering {
			r.table.NextBodyRow()
		}

	case *ast.String:
		panic("this type is for the typographer extension which is not enabled ")

	default:
		panic(fmt.Sprintf("Unknown node type %T", node))
	}

	return ast.WalkContinue, nil
}

func (r *renderer) renderHorizontalRule(w io.Writer) {
	_, _ = fmt.Fprintf(w, "%s%s\n\n", r.pad(), strings.Repeat("─", r.lineWidth-r.leftPad))
}

// level = 1 = h1 = title
// level = 2 = h2
// ...
func (r *renderer) renderHeading(w io.Writer, level int) {
	content := r.inlineAccumulator.String()
	r.inlineAccumulator.Reset()

	// render the full line with the headingNumbering
	r.headingNumbering.Observe(level)
	content = fmt.Sprintf("%s %s", r.headingNumbering.Render(), content)

	if level == 1 { // = title
		// wrap if needed
		wrapped, lines := text.WrapWithPadAlign(content, r.lineWidth-r.leftPad-2, "", text.AlignCenter)
		split := strings.Split(wrapped, "\n")

		_, _ = fmt.Fprintf(w, "%s╔%s╗\n", r.pad(), strings.Repeat("═", r.lineWidth-r.leftPad-2))
		for i := 0; i < lines; i++ {
			_, _ = fmt.Fprintf(w, "%s║%s", r.pad(), split[i])
			_, _ = fmt.Fprintf(w, "%s║\n", strings.Repeat(" ", r.lineWidth-r.leftPad-2-text.Len(split[i])))
		}
		_, _ = fmt.Fprintf(w, "%s╚%s╝\n", r.pad(), strings.Repeat("═", r.lineWidth-r.leftPad-2))
	} else {
		content = r.headingShade(level)(content)
		// wrap if needed
		wrapped, _ := text.WrapWithPad(content, r.lineWidth, r.pad())
		_, _ = fmt.Fprintln(w, wrapped)
	}

	_, _ = fmt.Fprintln(w)
}

func (r *renderer) renderCodeBlock(w io.Writer, source []byte, lines *mdtext.Segments, language string) {
	code := accumulateSegments(lines, source)
	// code = strings.TrimSuffix(code, "\n")

	var lexer chroma.Lexer
	// try to get the lexer from the language tag if any
	if len(language) > 0 {
		lexer = lexers.Get(language)
	}
	// fallback on detection
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}
	// all failed :-(
	if lexer == nil {
		lexer = lexers.Fallback
	}
	// simplify the lexer output
	lexer = chroma.Coalesce(lexer)

	var formatter chroma.Formatter
	if color.NoColor {
		formatter = formatters.Fallback
	} else {
		formatter = formatters.TTY8
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		// Something failed, falling back to no highlight render
		r.renderFormattedCodeBlock(w, code)
		return
	}

	buf := &bytes.Buffer{}

	err = formatter.Format(buf, styles.Pygments, iterator)
	if err != nil {
		// Something failed, falling back to no highlight render
		r.renderFormattedCodeBlock(w, code)
		return
	}

	r.renderFormattedCodeBlock(w, buf.String())
}

func (r *renderer) renderFormattedCodeBlock(w io.Writer, code string) {
	// remove the trailing line break
	code = strings.TrimRight(code, "\n")

	r.addPad(GreenBold("┃ "))
	output, _ := text.WrapWithPad(code, r.lineWidth, r.pad())
	r.popPad()

	_, _ = fmt.Fprintln(w, output)
}

func (r *renderer) renderHTMLBlock(w io.Writer, source []byte, node *ast.HTMLBlock) {
	var buf bytes.Buffer

	flushInline := func() {
		if r.inlineAccumulator.Len() <= 0 {
			return
		}
		content := r.inlineAccumulator.String()
		r.inlineAccumulator.Reset()
		out, _ := text.WrapWithPad(content, r.lineWidth, r.pad())
		_, _ = fmt.Fprint(&buf, out, "\n\n")
	}

	data := accumulateSegments(node.Lines(), source)
	doc, err := html.Parse(strings.NewReader(data))
	if err != nil {
		// if there is a parsing error, fallback to a simple render
		r.inlineAccumulator.Reset()
		content := Red(data)
		out, _ := text.WrapWithPad(content, r.lineWidth, r.pad())
		_, _ = fmt.Fprint(w, out, "\n\n")
		return
	}

	htmlWalker.WalkFunc(doc, func(node *html.Node, entering bool) htmlWalker.WalkStatus {
		// if node.Type != html.TextNode {
		// 	fmt.Println(node.Type, "(", node.Data, ")", entering)
		// }

		switch node.Type {
		case html.CommentNode, html.DoctypeNode:
			// Not rendered

		case html.DocumentNode:

		case html.ElementNode:
			switch node.Data {
			case "html", "body":
				return htmlWalker.GoToNext

			case "head":
				return htmlWalker.SkipChildren

			case "div", "p":
				if entering {
					flushInline()
				} else {
					content := r.inlineAccumulator.String()
					r.inlineAccumulator.Reset()
					if len(content) == 0 {
						return htmlWalker.GoToNext
					}
					// remove all line breaks, those are fully managed in HTML
					content = strings.Replace(content, "\n", "", -1)
					align := getDivHTMLAttr(node.Attr)
					content, _ = text.WrapWithPadAlign(content, r.lineWidth, r.pad(), align)
					_, _ = fmt.Fprint(&buf, content, "\n\n")
				}

			case "h1":
				if !entering {
					r.renderHeading(&buf, 1)
				}
			case "h2":
				if !entering {
					r.renderHeading(&buf, 2)
				}
			case "h3":
				if !entering {
					r.renderHeading(&buf, 3)
				}
			case "h4":
				if !entering {
					r.renderHeading(&buf, 4)
				}
			case "h5":
				if !entering {
					r.renderHeading(&buf, 5)
				}
			case "h6":
				if !entering {
					r.renderHeading(&buf, 6)
				}

			case "img":
				flushInline()
				src, title := getImgHTMLAttr(node.Attr)
				str, _ := r.renderImage(src, title, r.lineWidth-len(r.pad()))
				r.inlineAccumulator.WriteString(str)

			case "hr":
				flushInline()
				r.renderHorizontalRule(&buf)

			case "ul", "ol":
				if !entering {
					if node.NextSibling == nil {
						_, _ = fmt.Fprint(&buf, "\n")
						return htmlWalker.GoToNext
					}
					switch node.NextSibling.Data {
					case "ul", "ol":
					default:
						_, _ = fmt.Fprint(&buf, "\n")
					}
				}

			case "li":
				if entering {
					switch node.Parent.Data {
					case "ul":
						r.indent = r.pad() + Green("• ")
						r.addPad("  ")

					case "ol":
						itemNumber := 1
						previous := node.PrevSibling
						for previous != nil {
							itemNumber++
							previous = previous.PrevSibling
						}
						prefix := fmt.Sprintf("%d. ", itemNumber)
						r.indent = r.pad() + Green(prefix)
						r.addPad(strings.Repeat(" ", text.Len(prefix)))

					default:
						r.inlineAccumulator.WriteString(Red(renderRawHtml(node)))
						return htmlWalker.SkipChildren
					}
				} else {
					switch node.Parent.Data {
					case "ul", "ol":
						content := r.inlineAccumulator.String()
						r.inlineAccumulator.Reset()
						out, _ := text.WrapWithPadIndent(content, r.lineWidth, r.indent, r.pad())
						r.indent = ""
						_, _ = fmt.Fprint(&buf, out, "\n")
						r.popPad()
					}
				}

			case "a":
				if entering {
					r.inlineAccumulator.WriteString("[")
				} else {
					href, alt := getAHTMLAttr(node.Attr)
					r.inlineAccumulator.WriteString("](")
					r.inlineAccumulator.WriteString(Blue(href))
					if len(alt) > 0 {
						r.inlineAccumulator.WriteString(" ")
						r.inlineAccumulator.WriteString(alt)
					}
					r.inlineAccumulator.WriteString(")")
				}

			case "br":
				if entering {
					r.inlineAccumulator.WriteString("\n")
				}

			case "table":
				if entering {
					flushInline()
					r.table = newTableRenderer()
				} else {
					r.table.Render(&buf, r.leftPad, r.lineWidth)
					r.table = nil
				}

			case "thead", "tbody":
				// nothing to do

			case "tr":
				if entering && node.Parent.Data != "thead" {
					r.table.NextBodyRow()
				}

			case "th":
				if !entering {
					content := r.inlineAccumulator.String()
					r.inlineAccumulator.Reset()

					align := getTdHTMLAttr(node.Attr)
					r.table.AddHeaderCell(content, align)
				}

			case "td":
				if !entering {
					content := r.inlineAccumulator.String()
					r.inlineAccumulator.Reset()

					align := getTdHTMLAttr(node.Attr)
					r.table.AddBodyCell(content, align)
				}

			case "strong", "b":
				if entering {
					r.inlineAccumulator.WriteString(boldOn)
				} else {
					// This is super silly but some terminals, instead of having
					// the ANSI code SGR 21 do "bold off" like the logic would guide,
					// do "double underline" instead. This is madness.

					// To resolve that problem, we take a snapshot of the escape state,
					// remove the bold, then output "reset all" + snapshot
					es := text.EscapeState{}
					es.Witness(r.inlineAccumulator.String())
					es.Bold = false
					r.inlineAccumulator.WriteString(resetAll)
					r.inlineAccumulator.WriteString(es.String())
				}

			case "i", "em":
				if entering {
					r.inlineAccumulator.WriteString(italicOn)
				} else {
					r.inlineAccumulator.WriteString(italicOff)
				}

			case "s":
				if entering {
					r.inlineAccumulator.WriteString(crossedOutOn)
				} else {
					r.inlineAccumulator.WriteString(crossedOutOff)
				}

			default:
				r.inlineAccumulator.WriteString(Red(renderRawHtml(node)))
			}

		case html.TextNode:
			t := strings.TrimSpace(node.Data)
			t = strings.ReplaceAll(t, "\n", "")
			r.inlineAccumulator.WriteString(t)

		default:
			panic("unhandled case")
		}

		return htmlWalker.GoToNext
	})

	flushInline()
	_, _ = fmt.Fprint(w, buf.String())
	r.inlineAccumulator.Reset()

	// 		// dl + (dt+dd)
	//
	// 		// details
	// 		// summary
	//
}

func getDivHTMLAttr(attrs []html.Attribute) text.Alignment {
	for _, attr := range attrs {
		switch attr.Key {
		case "align":
			switch attr.Val {
			case "left":
				return text.AlignLeft
			case "center":
				return text.AlignCenter
			case "right":
				return text.AlignRight
			}
		}
	}
	return text.AlignLeft
}

func getImgHTMLAttr(attrs []html.Attribute) (src, title string) {
	for _, attr := range attrs {
		switch attr.Key {
		case "src":
			src = attr.Val
		case "alt":
			title = attr.Val
		}
	}
	return
}

func getAHTMLAttr(attrs []html.Attribute) (href, alt string) {
	for _, attr := range attrs {
		switch attr.Key {
		case "href":
			href = attr.Val
		case "alt":
			alt = attr.Val
		}
	}
	return
}

func getTdHTMLAttr(attrs []html.Attribute) CellAlign {
	for _, attr := range attrs {
		switch attr.Key {
		case "align":
			switch attr.Val {
			case "right":
				return CellAlignRight
			case "left":
				return CellAlignLeft
			case "center":
				return CellAlignCenter
			}

		case "style":
			for _, pair := range strings.Split(attr.Val, " ") {
				split := strings.Split(pair, ":")
				if split[0] != "text-align" || len(split) != 2 {
					continue
				}
				switch split[1] {
				case "right":
					return CellAlignRight
				case "left":
					return CellAlignLeft
				case "center":
					return CellAlignCenter
				}
			}
		}
	}
	return CellAlignLeft
}

func renderRawHtml(node *html.Node) string {
	var result strings.Builder
	openContent := make([]string, 0, 8)

	openContent = append(openContent, node.Data)
	for _, attr := range node.Attr {
		openContent = append(openContent, fmt.Sprintf("%s=\"%s\"", attr.Key, attr.Val))
	}

	result.WriteString("<")
	result.WriteString(strings.Join(openContent, " "))

	if node.FirstChild == nil {
		result.WriteString("/>")
		return result.String()
	}

	result.WriteString(">")

	child := node.FirstChild
	for child != nil {
		if child.Type == html.TextNode {
			t := strings.TrimSpace(child.Data)
			result.WriteString(t)
			child = child.NextSibling
			continue
		}

		switch node.Data {
		case "ul", "p":
			result.WriteString("\n  ")
		}

		result.WriteString(renderRawHtml(child))
		child = child.NextSibling
	}

	switch node.Data {
	case "ul", "p":
		result.WriteString("\n")
	}

	result.WriteString("</")
	result.WriteString(node.Data)
	result.WriteString(">")

	return result.String()
}

func (r *renderer) renderImage(dest string, title string, lineWidth int) (result string, rendered bool) {
	title = strings.ReplaceAll(title, "\n", "")
	title = strings.TrimSpace(title)
	dest = strings.ReplaceAll(dest, "\n", "")
	dest = strings.TrimSpace(dest)

	fallback := func() (string, bool) {
		return fmt.Sprintf("![%s](%s)", title, Blue(dest)), false
	}

	reader, err := imageFromDestination(dest)
	if err != nil {
		return fallback()
	}

	x := lineWidth

	if r.imageDithering == ansimage.DitheringWithChars || r.imageDithering == ansimage.DitheringWithBlocks {
		// not sure why this is needed by ansimage
		// x *= 4
	}

	img, err := ansimage.NewScaledFromReader(reader, math.MaxInt32, x,
		stdcolor.Black, ansimage.ScaleModeFit, r.imageDithering)

	if err != nil {
		return fallback()
	}

	if title != "" {
		return fmt.Sprintf("%s%s: %s", img.Render(), title, Blue(dest)), true
	}
	return fmt.Sprintf("%s%s", img.Render(), Blue(dest)), true
}

func imageFromDestination(dest string) (io.ReadCloser, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	if strings.HasPrefix(dest, "http://") || strings.HasPrefix(dest, "https://") {
		res, err := client.Get(dest)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("http: %v", http.StatusText(res.StatusCode))
		}

		return res.Body, nil
	}

	return os.Open(dest)
}

func removeLineBreak(text string) string {
	lines := strings.Split(text, "\n")

	if len(lines) <= 1 {
		return text
	}

	for i, l := range lines {
		switch i {
		case 0:
			lines[i] = strings.TrimRightFunc(l, unicode.IsSpace)
		case len(lines) - 1:
			lines[i] = strings.TrimLeftFunc(l, unicode.IsSpace)
		default:
			lines[i] = strings.TrimFunc(l, unicode.IsSpace)
		}
	}
	return strings.Join(lines, " ")
}

func shouldCleanText(node ast.Node) bool {
	for node != nil {
		switch node.(type) {
		case *ast.Blockquote:
			return false

		case *ast.Heading, *ast.Image, *ast.Link,
			*extast.TableCell, *ast.Document, *ast.ListItem:
			return true
		}

		node = node.Parent()
	}

	panic("bad markdown document or missing case")
}

func accumulateSegments(lines *mdtext.Segments, source []byte) string {
	var builder strings.Builder
	l := lines.Len()
	for i := 0; i < l; i++ {
		line := lines.At(i)
		builder.Write(line.Value(source))
	}
	return builder.String()
}

func isLinkRefDefinition(node ast.Node) bool {
	// Don't ask me why but somehow the parser return an empty
	// TextBlock when a Link reference definition block happen.

	// Yes, this makes me sad as well

	tb, ok := node.(*ast.TextBlock)

	if ok && tb.Lines().Len() == 0 &&
		node.Parent() != nil &&
		node.Parent().Kind() == ast.KindDocument {
		return true
	}
	return false
}
