package debug

import (
	"fmt"
	"os"

	"golang.org/x/net/html"

	htmlWalker "github.com/MichaelMure/go-term-markdown/html"
)

func HtmlAst2PlantUML(node *html.Node) {
	f, err := os.Create("/tmp/html_ast.puml")
	if err != nil {
		panic(err)
	}

	_, _ = fmt.Fprintln(f, "@startuml")
	_, _ = fmt.Fprintln(f, "(*) -> Blah")

	htmlWalker.WalkFunc(node, func(node *html.Node, entering bool) htmlWalker.WalkStatus {
		type2str := func(nodeType html.NodeType) string {
			return [...]string{"ErrorNode", "TextNode", "DocumentNode", "ElementNode", "CommentNode", "DoctypeNode", "scopeMarkerNode"}[nodeType]
		}
		if entering {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				t := type2str(node.Type)
				if node.Type == html.ElementNode {
					t = node.Data
				}
				tc := type2str(child.Type)
				if child.Type == html.ElementNode {
					tc = child.Data
				}

				_, _ = fmt.Fprintf(f, "\"%s-%p\" --> \"%s-%p\"\n", t, node, tc, child)
			}
		}
		return htmlWalker.GoToNext
	})
	_, _ = fmt.Fprintln(f, "@enduml")
}
