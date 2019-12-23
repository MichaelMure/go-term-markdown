package debug

import (
	"fmt"
	"io"
	"os"
	"strings"

	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
)

func MarkdownAstSet2PlantUML(node ast.Node) {
	f, err := os.Create("/tmp/markdown_ast_set.puml")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	r := &astSetRenderer{
		f:   f,
		set: make(map[string]struct{}),
	}

	_, _ = fmt.Fprintln(f, "@startuml")
	_, _ = fmt.Fprintln(f, "(*) --> Document")

	md.Render(node, r)

	_, _ = fmt.Fprintln(f, "@enduml")
}

var _ md.Renderer = &astSetRenderer{}

type astSetRenderer struct {
	set map[string]struct{}
	f   *os.File
}

func (a *astSetRenderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	if entering {
		for _, child := range node.GetChildren() {
			str := fmt.Sprintf("%T --> %T\n", node, child)
			if _, has := a.set[str]; !has {
				a.set[str] = struct{}{}
				_, _ = fmt.Fprintf(a.f, strings.Replace(str, "*ast.", "", -1))
			}
		}
	}

	return ast.GoToNext
}

func (a *astSetRenderer) RenderHeader(w io.Writer, ast ast.Node) {}

func (a *astSetRenderer) RenderFooter(w io.Writer, ast ast.Node) {}
