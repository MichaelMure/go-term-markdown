package markdown

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/MichaelMure/go-term-text"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/fatih/color"
	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/kyokomi/emoji"
)

/*

Here are the possible cases for the AST. You can render it using PlantUML.

@startuml

(*) --> Document
BlockQuote --> BlockQuote
BlockQuote --> CodeBlock
BlockQuote --> List
BlockQuote --> Paragraph
Del --> Emph
Del --> Strong
Del --> Text
Document --> BlockQuote
Document --> CodeBlock
Document --> Heading
Document --> HorizontalRule
Document --> HTMLBlock
Document --> List
Document --> Paragraph
Document --> Table
Emph --> Text
Heading --> Code
Heading --> Del
Heading --> Emph
Heading --> HTMLSpan
Heading --> Image
Heading --> Link
Heading --> Strong
Heading --> Text
Image --> Text
Link --> Image
Link --> Text
ListItem --> List
ListItem --> Paragraph
List --> ListItem
Paragraph --> Code
Paragraph --> Del
Paragraph --> Emph
Paragraph --> Hardbreak
Paragraph --> HTMLSpan
Paragraph --> Image
Paragraph --> Link
Paragraph --> Strong
Paragraph --> Text
Strong --> Emph
Strong --> Text
TableBody --> TableRow
TableCell --> Code
TableCell --> Del
TableCell --> Emph
TableCell --> HTMLSpan
TableCell --> Image
TableCell --> Link
TableCell --> Strong
TableCell --> Text
TableHeader --> TableRow
TableRow --> TableCell
Table --> TableBody
Table --> TableHeader

@enduml

*/

var _ md.Renderer = &renderer{}

type renderer struct {
	// maximum line width allowed
	lineWidth int
	// constant left padding to apply
	leftPad int

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

	blockQuoteLevel int

	table *tableRenderer
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

func (r *renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	// TODO: remove
	// fmt.Println(node.Type, string(node.Literal), entering)

	switch node := node.(type) {
	case *ast.Document:
		// Nothing to do

	case *ast.BlockQuote:
		// set and remove a colored bar on the left
		if entering {
			r.blockQuoteLevel++
			r.addPad(quoteShade(r.blockQuoteLevel)("┃ "))
		} else {
			r.blockQuoteLevel--
			r.popPad()
		}

	case *ast.List:
		if next := ast.GetNextNode(node); !entering && next != nil {
			_, parentIsListItem := node.GetParent().(*ast.ListItem)
			_, nextIsList := next.(*ast.List)
			if !nextIsList && !parentIsListItem {
				_, _ = fmt.Fprintln(w)
			}
		}

	case *ast.ListItem:
		// write the prefix, add a padding if needed, and let Paragraph handle the rest
		if entering {
			switch {
			// numbered list
			case node.ListFlags&ast.ListTypeOrdered != 0:
				itemNumber := 1
				siblings := node.GetParent().GetChildren()
				for _, sibling := range siblings {
					if sibling == node {
						break
					}
					itemNumber++
				}
				prefix := fmt.Sprintf("%d. ", itemNumber)
				r.indent = r.pad() + Green(prefix)
				r.addPad(strings.Repeat(" ", text.WordLen(prefix)))

			// header of a definition
			case node.ListFlags&ast.ListTypeTerm != 0:
				r.inlineAccumulator.WriteString(greenOn)

			// content of a definition
			case node.ListFlags&ast.ListTypeDefinition != 0:
				r.addPad("  ")

			// no flags means it's the normal bullet point list
			default:
				r.indent = r.pad() + Green("• ")
				r.addPad("  ")
			}
		} else {
			switch {
			// numbered list
			case node.ListFlags&ast.ListTypeOrdered != 0:
				r.popPad()

			// header of a definition
			case node.ListFlags&ast.ListTypeTerm != 0:
				r.inlineAccumulator.WriteString(colorOff)

			// content of a definition
			case node.ListFlags&ast.ListTypeDefinition != 0:
				r.popPad()
				_, _ = fmt.Fprintln(w)

			// no flags means it's the normal bullet point list
			default:
				r.popPad()
			}
		}

	case *ast.Paragraph:
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

			// extra line break in some cases
			if next := ast.GetNextNode(node); next != nil {
				switch next.(type) {
				case *ast.Paragraph, *ast.Heading, *ast.HorizontalRule,
					*ast.CodeBlock, *ast.HTMLBlock:
					_, _ = fmt.Fprintln(w)
				}
			}
		}

	case *ast.Heading:
		if !entering {
			content := r.inlineAccumulator.String()
			r.inlineAccumulator.Reset()

			// render the full line with the headingNumbering
			r.headingNumbering.Observe(node.Level)
			rendered := fmt.Sprintf("%s%s %s", r.pad(), r.headingNumbering.Render(), content)

			// wrap if needed
			wrapped, _ := text.Wrap(rendered, r.lineWidth)
			colored := headingShade(node.Level)(wrapped)
			_, _ = fmt.Fprintln(w, colored)

			// render the underline, if any
			if node.Level == 1 {
				_, _ = fmt.Fprintf(w, "%s%s\n", r.pad(), strings.Repeat("─", r.lineWidth-r.leftPad))
			}

			_, _ = fmt.Fprintln(w)
		}

	case *ast.HorizontalRule:
		_, _ = fmt.Fprintf(w, "%s%s\n\n", r.pad(), strings.Repeat("─", r.lineWidth-r.leftPad))

	case *ast.Emph:
		if entering {
			r.inlineAccumulator.WriteString(italicOn)
		} else {
			r.inlineAccumulator.WriteString(italicOff)
		}

	case *ast.Strong:
		if entering {
			r.inlineAccumulator.WriteString(boldOn)
		} else {
			r.inlineAccumulator.WriteString(boldOff)
		}

	case *ast.Del:
		if entering {
			r.inlineAccumulator.WriteString(crossedOutOn)
		} else {
			r.inlineAccumulator.WriteString(crossedOutOff)
		}

	case *ast.Link:
		if entering {
			r.inlineAccumulator.WriteString("[")
			r.inlineAccumulator.WriteString(string(ast.GetFirstChild(node).AsLeaf().Literal))
			r.inlineAccumulator.WriteString("](")
			r.inlineAccumulator.WriteString(Blue(string(node.Destination)))
			if len(node.Title) > 0 {
				r.inlineAccumulator.WriteString(" ")
				r.inlineAccumulator.WriteString(string(node.Title))
			}
			r.inlineAccumulator.WriteString(")")
			return ast.SkipChildren
		}

	case *ast.Image:

	case *ast.Text:
		content := string(node.Literal)
		if shouldCleanText(node) {
			content = removeLineBreak(content)
		}
		// emoji support !
		emojed := emoji.Sprint(content)
		r.inlineAccumulator.WriteString(emojed)

	case *ast.HTMLBlock:
		content := Red(string(node.Literal))
		out, _ := text.WrapWithPad(content, r.lineWidth, r.pad())
		_, _ = fmt.Fprint(w, out, "\n\n")

	case *ast.CodeBlock:
		r.renderCodeBlock(w, node)

	case *ast.Softbreak:
		// not actually implemented in blackfriday/gomarkdown
		r.inlineAccumulator.WriteString("\n")

	case *ast.Hardbreak:
		r.inlineAccumulator.WriteString("\n")

	case *ast.Code:
		r.inlineAccumulator.WriteString(BlueBgItalic(string(node.Literal)))

	case *ast.HTMLSpan:
		fmt.Println("SPAN:", string(node.Literal))
		r.inlineAccumulator.WriteString(Red(string(node.Literal)))

	case *ast.Table:
		if entering {
			r.table = newTableRenderer()
		} else {
			r.table.Render(w, r.leftPad, r.lineWidth)
			r.table = nil
		}

	case *ast.TableCell:
		if !entering {
			content := r.inlineAccumulator.String()
			r.inlineAccumulator.Reset()

			if node.IsHeader {
				r.table.AddHeaderCell(content, node.Align)
			} else {
				r.table.AddBodyCell(content)
			}
		}

	case *ast.TableHeader:
		// nothing to do

	case *ast.TableBody:
		// nothing to do

	case *ast.TableRow:
		if _, ok := node.Parent.(*ast.TableBody); ok && entering {
			r.table.NextBodyRow()
		}

	default:
		panic(fmt.Sprintf("Unknown node type %T", node))
	}

	return ast.GoToNext
}

func (*renderer) RenderHeader(w io.Writer, node ast.Node) {}

func (*renderer) RenderFooter(w io.Writer, node ast.Node) {}

func (r *renderer) renderCodeBlock(w io.Writer, node *ast.CodeBlock) {
	code := string(node.Literal)
	var lexer chroma.Lexer
	// try to get the lexer from the language tag if any
	if len(node.Info) > 0 {
		lexer = lexers.Get(string(node.Info))
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
		case *ast.BlockQuote:
			return false

		case *ast.Heading, *ast.Image, *ast.Link,
			*ast.TableCell, *ast.Document, *ast.ListItem:
			return true
		}

		node = node.GetParent()
	}

	panic("bad markdown document or missing case")
}
