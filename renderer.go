package markdown

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"gopkg.in/russross/blackfriday.v2"

	"github.com/MichaelMure/git-bug/util/text"
)

/*

Here are the possible cases for the AST. You can render it using PlantUML.

@startuml

(*) --> Document
BlockQuote --> BlockQuote
BlockQuote --> CodeBlock
BlockQuote --> List
BlockQuote --> Paragraph
Del --> Text
Document --> BlockQuote
Document --> CodeBlock
Document --> Heading
Document --> HorizontalRule
Document --> HTMLBlock
Document --> List
Document --> Paragraph
Emph --> Text
Heading --> Text
Item --> List
Item --> Paragraph
Link --> Text
List --> Item
Paragraph --> Code
Paragraph --> Del
Paragraph --> Emph
Paragraph --> HTMLSpan
Paragraph --> Link
Paragraph --> Strong
Paragraph --> Text
Strong --> Emph
Strong --> Text

@enduml

*/

var _ blackfriday.Renderer = &renderer{}

type renderer struct {
	// maximum line width allowed
	lineWidth int
	// constant left padding to apply
	leftPad int

	// all the custom left paddings, without the fixed space from leftPad
	padAccumulator []string

	// Count the number of line in the rendered output
	// lines int

	// record and render the heading numbering
	headingNumbering headingNumbering

	paragraph strings.Builder

	blockQuoteLevel int
}

func newRenderer(lineWidth int, leftPad int) *renderer {
	return &renderer{
		lineWidth:      lineWidth,
		leftPad:        leftPad,
		padAccumulator: make([]string, 0, 10),
	}
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

func (r *renderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	fmt.Println(node.Type, string(node.Literal), entering)

	green := color.New(color.FgGreen)

	switch node.Type {
	case blackfriday.Document:
		// Nothing to do

	case blackfriday.BlockQuote:
		// set and remove a colored bar on the left
		if entering {
			r.blockQuoteLevel++
			r.addPad(quoteShade(r.blockQuoteLevel)("┃ "))
		} else {
			r.blockQuoteLevel--
			r.popPad()
		}

	case blackfriday.List:
		if !entering && node.Next != nil {
			if node.Next.Type != blackfriday.List && node.Parent.Type != blackfriday.Item {
				_, _ = fmt.Fprintln(w)
			}
		}

	case blackfriday.Item:
		fmt.Println(node.ListFlags)

		// write the prefix, add a padding if needed, and let Paragraph handle the rest
		if entering {
			switch {
			// numbered list
			case node.ListData.ListFlags&blackfriday.ListTypeOrdered != 0:
				itemNumber := 1
				for prev := node.Prev; prev != nil; prev = prev.Prev {
					itemNumber++
				}
				prefix := fmt.Sprintf("%d. ", itemNumber)
				r.paragraph.WriteString(prefix)
				r.addPad(strings.Repeat(" ", text.WordLen(prefix)))

			// content of a definition
			case node.ListData.ListFlags&blackfriday.ListTypeTerm != 0:
				r.paragraph.WriteString("  ")
				r.addPad("  ")

			// header of a definition
			case node.ListData.ListFlags&blackfriday.ListTypeDefinition != 0:
				r.paragraph.WriteString("  ")
				r.paragraph.WriteString(green.Format())
				r.addPad("  ")

			// no flags means it's the normal bullet point list
			default:
				r.paragraph.WriteString("• ")
				r.addPad("  ")
			}
		} else {
			r.popPad()
			if node.ListData.ListFlags&blackfriday.ListTypeDefinition != 0 {

			}
		}

	case blackfriday.Paragraph:
		if !entering {
			out, _ := text.WrapWithPad(r.paragraph.String(), r.lineWidth, r.pad())
			_, _ = fmt.Fprint(w, out, "\n")

			// extra line break in some cases
			if node.Next != nil {
				switch node.Next.Type {
				case blackfriday.Paragraph, blackfriday.Heading, blackfriday.HorizontalRule,
					blackfriday.CodeBlock:
					_, _ = fmt.Fprintln(w)
				}

				if node.Next.Type == blackfriday.List && node.Parent.Type != blackfriday.Item {
					_, _ = fmt.Fprintln(w)
				}
			}
			r.paragraph.Reset()
		}

	case blackfriday.Heading:
		// the child node of a heading is a blackfriday.Text. We render the whole thing
		// in one go and skip the child.

		// render the full line with the headingNumbering
		r.headingNumbering.Observe(node.Level)
		rendered := fmt.Sprintf("%s%s %s", r.pad(), r.headingNumbering.Render(), string(node.FirstChild.Literal))

		// output the text, truncated if needed, no line break
		truncated := runewidth.Truncate(rendered, r.lineWidth, "…")
		colored := headingShade(node.Level)(truncated)
		_, _ = fmt.Fprintln(w, colored)

		// render the underline, if any
		if node.Level == 1 {
			_, _ = fmt.Fprintf(w, "%s%s\n", r.pad(), strings.Repeat("─", r.lineWidth-r.leftPad))
		}

		_, _ = fmt.Fprintln(w)

		return blackfriday.SkipChildren

	case blackfriday.HorizontalRule:
		_, _ = fmt.Fprintf(w, "%s%s\n\n", r.pad(), strings.Repeat("─", r.lineWidth-r.leftPad))

	case blackfriday.Emph:
		r.paragraph.WriteString(Italic(string(node.FirstChild.Literal)))
		return blackfriday.SkipChildren

	case blackfriday.Strong:
		r.paragraph.WriteString(Bold(string(node.FirstChild.Literal)))
		return blackfriday.SkipChildren

	case blackfriday.Del:
		r.paragraph.WriteString(CrossedOut(string(node.FirstChild.Literal)))
		return blackfriday.SkipChildren

	case blackfriday.Link:
		r.paragraph.WriteString("[")
		r.paragraph.WriteString(string(node.FirstChild.Literal))
		r.paragraph.WriteString("](")
		r.paragraph.WriteString(Blue(string(node.LinkData.Destination)))
		if len(node.LinkData.Title) > 0 {
			r.paragraph.WriteString(" ")
			r.paragraph.WriteString(string(node.LinkData.Title))
		}
		r.paragraph.WriteString(")")
		return blackfriday.SkipChildren

	case blackfriday.Image:

	case blackfriday.Text:
		r.paragraph.Write(node.Literal)

	case blackfriday.HTMLBlock:

	case blackfriday.CodeBlock:
		r.renderCodeBlock(w, node)

	case blackfriday.Softbreak:

	case blackfriday.Hardbreak:

	case blackfriday.Code:
		r.paragraph.WriteString(BlueBgItalic(string(node.Literal)))

	case blackfriday.HTMLSpan:

	case blackfriday.Table:

	case blackfriday.TableCell:

	case blackfriday.TableHead:

	case blackfriday.TableBody:

	case blackfriday.TableRow:

	default:
		panic("Unknown node type " + node.Type.String())
	}

	return blackfriday.GoToNext
}

func (*renderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	fmt.Println(ast)
}

func (*renderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	fmt.Println(ast)
}

func (r *renderer) renderCodeBlock(w io.Writer, node *blackfriday.Node) {
	code := string(node.Literal)
	var lexer chroma.Lexer
	// try to get the lexer from the language tag if any
	if len(node.CodeBlockData.Info) > 0 {
		lexer = lexers.Get(string(node.CodeBlockData.Info))
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

	_, _ = fmt.Fprint(w, output)

	_, _ = fmt.Fprintf(w, "\n\n")
}
