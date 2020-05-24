package debug

import (
	"fmt"
	"os"
	"strings"

	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
	mdrender "github.com/yuin/goldmark/renderer"
	mdutil "github.com/yuin/goldmark/util"
)

func NewMarkdownAstSet2PlantUML() *mdAst2Puml {
	f, err := os.Create("/tmp/markdown_ast_set.puml")
	if err != nil {
		panic(err)
	}

	_, _ = fmt.Fprintln(f, "@startuml")
	_, _ = fmt.Fprintln(f, "(*) --> Document")

	r := mdrender.NewRenderer(mdrender.WithNodeRenderers(
		mdutil.Prioritized(&astSetRenderer{
			set: make(map[string]struct{}),
		}, 1000)),
	)

	return &mdAst2Puml{r: r, f: f}
}

var _ mdrender.NodeRenderer = &astSetRenderer{}

type astSetRenderer struct {
	set map[string]struct{}
}

func (r *astSetRenderer) RegisterFuncs(reg mdrender.NodeRendererFuncRegisterer) {
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

func (r *astSetRenderer) render(w mdutil.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			str := fmt.Sprintf("%T --> %T\n", node, child)
			if _, has := r.set[str]; !has {
				r.set[str] = struct{}{}
				_, _ = fmt.Fprintf(w, strings.Replace(str, "*ast.", "", -1))
			}
		}
	}

	return ast.WalkContinue, nil
}
