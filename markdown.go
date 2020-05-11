package markdown

import (
	"io"

	mdext "github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	mdrender "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	mdutil "github.com/yuin/goldmark/util"

	"github.com/MichaelMure/go-term-markdown/debug"
)

// block parsers required for core + GFM
var bp = append(
	parser.DefaultBlockParsers(),
	// mdutil.Prioritized(mdext.NewDefinitionListParser(), 101),
	// mdutil.Prioritized(mdext.NewDefinitionDescriptionParser(), 102),
	// mdutil.Prioritized(mdext.NewFootnoteBlockParser(), 999),
)

// inline parsers required for core + GFM
var ip = append(
	parser.DefaultInlineParsers(),
	mdutil.Prioritized(mdext.NewTaskCheckBoxParser(), 0),
	// mdutil.Prioritized(mdext.NewFootnoteParser(), 101),
	mdutil.Prioritized(mdext.NewStrikethroughParser(), 500),
	mdutil.Prioritized(mdext.NewLinkifyParser(), 999),
)

// paragraph transformers required for core + GFM
var pt = append(
	parser.DefaultParagraphTransformers(),
	mdutil.Prioritized(mdext.NewTableParagraphTransformer(), 200),
)

// AST transformers required
var at = []mdutil.PrioritizedValue{
	// mdutil.Prioritized(mdext.NewFootnoteASTTransformer(), 999),
}

var defaultParser = parser.NewParser(
	parser.WithBlockParsers(bp...),
	parser.WithInlineParsers(ip...),
	parser.WithParagraphTransformers(pt...),
	parser.WithASTTransformers(at...),
)

var d = debug.NewMarkdownAstSet2PlantUML()

func Render(w io.Writer, source []byte, lineWidth int, leftPad int, opts ...Options) error {
	reader := text.NewReader(source)
	nodes := defaultParser.Parse(reader)

	renderer := mdrender.NewRenderer(mdrender.WithNodeRenderers(mdutil.Prioritized(
		newRenderer(lineWidth, leftPad, opts...), 1000)),
	)

	d.Render(source, nodes)

	return renderer.Render(w, source, nodes)
}
