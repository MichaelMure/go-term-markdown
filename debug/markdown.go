package debug

import (
	"fmt"
	"io"
	"os"
	"strings"

	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
)

func MarkdownAst2PlantUML(node ast.Node) {
	f, err := os.Create("/tmp/markdown_ast.puml")
	if err != nil {
		panic(err)
	}

	r := &astRenderer{
		f: f,
	}

	md.Render(node, r)
}

var _ md.Renderer = &astRenderer{}

type astRenderer struct {
	f *os.File
}

func (a *astRenderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	if entering {
		for _, child := range node.GetChildren() {
			str := fmt.Sprintf("%T --> %T\n", node, child)
			_, _ = fmt.Fprintf(a.f, strings.Replace(str, "*ast.", "", -1))
		}
	}

	return ast.GoToNext
}

func (a *astRenderer) RenderHeader(w io.Writer, ast ast.Node) {
	_, _ = fmt.Fprintln(a.f, "@startuml")
	_, _ = fmt.Fprintln(a.f, "(*) --> Document")
}

func (a *astRenderer) RenderFooter(w io.Writer, ast ast.Node) {
	_, _ = fmt.Fprintln(a.f, "@enduml")
	_ = a.f.Close()
}
