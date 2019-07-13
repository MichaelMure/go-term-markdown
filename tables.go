package markdown

import (
	"io"
	"strings"

	"gopkg.in/russross/blackfriday.v2"

	"github.com/MichaelMure/git-bug/util/text"
)

type tableCell struct {
	content   string
	alignment blackfriday.CellAlignFlags
}

type tableRenderer struct {
	header []tableCell
	body   [][]tableCell
}

func newTableRenderer() *tableRenderer {
	return &tableRenderer{}
}

func (tr *tableRenderer) AddHeaderCell(content string, alignment blackfriday.CellAlignFlags) {
	tr.header = append(tr.header, tableCell{
		content:   content,
		alignment: alignment,
	})
}

func (tr *tableRenderer) NextBodyRow() {
	tr.body = append(tr.body, nil)
}

func (tr *tableRenderer) AddBodyCell(content string) {
	row := tr.body[len(tr.body)-1]
	column := len(row)
	row = append(row, tableCell{
		content:   content,
		alignment: tr.header[column].alignment,
	})
	tr.body[len(tr.body)-1] = row
}

func (tr *tableRenderer) Render(w io.Writer, lineWidth int) {
	columnWidths := tr.columnWidths(lineWidth)

	drawTopLine(w, columnWidths)

	drawRow(w, tr.header, columnWidths)

	drawHeaderUnderline(w, columnWidths)

	for i, row := range tr.body {
		drawRow(w, row, columnWidths)
		if i != len(tr.body)-1 {
			drawRowLine(w, columnWidths)
		}
	}

	drawBottomLine(w, columnWidths)
}

func (tr *tableRenderer) columnWidths(lineWidth int) []int {
	maxWidth := make([]int, len(tr.header))

	for i, cell := range tr.header {
		maxWidth[i] = max(maxWidth[i], text.MaxLineLen(cell.content))
	}

	for _, row := range tr.body {
		for i, cell := range row {
			maxWidth[i] = max(maxWidth[i], text.MaxLineLen(cell.content))
		}
	}

	sum := 0
	for _, width := range maxWidth {
		sum += width
	}

	available := lineWidth - len(tr.header) - 1

	// the easy case, content is not large enough to overflow
	if sum <= available {
		return maxWidth
	}

	// We have an overflow. First, we take as is the columns that are thinner
	// than the space equally divided.
	// Integer division, rounded lower.
	fairSpace := available / len(tr.header)

	result := make([]int, len(tr.header))
	remainingColumn := len(tr.header)

	for i, width := range maxWidth {
		if width <= fairSpace {
			result[i] = width
			available -= width
			remainingColumn--
		} else {
			// Mark the column as non-allocated yet
			result[i] = -1
		}
	}

	// Now we allocate evenly the remaining space to the remaining columns
	for i, width := range result {
		if width == -1 {
			width = available / remainingColumn
			result[i] = width
			available -= width
			remainingColumn--
		}
	}

	return result
}

func drawTopLine(w io.Writer, columnWidths []int) {
	_, _ = w.Write([]byte("┌"))
	for i, width := range columnWidths {
		_, _ = w.Write([]byte(strings.Repeat("─", width)))
		if i != len(columnWidths)-1 {
			_, _ = w.Write([]byte("┬"))
		}
	}
	_, _ = w.Write([]byte("┐"))
	_, _ = w.Write([]byte("\n"))
}

func drawHeaderUnderline(w io.Writer, columnWidths []int) {
	_, _ = w.Write([]byte("╞"))
	for i, width := range columnWidths {
		_, _ = w.Write([]byte(strings.Repeat("═", width)))
		if i != len(columnWidths)-1 {
			_, _ = w.Write([]byte("╪"))
		}
	}
	_, _ = w.Write([]byte("╡"))
	_, _ = w.Write([]byte("\n"))
}

func drawBottomLine(w io.Writer, columnWidths []int) {
	_, _ = w.Write([]byte("└"))
	for i, width := range columnWidths {
		_, _ = w.Write([]byte(strings.Repeat("─", width)))
		if i != len(columnWidths)-1 {
			_, _ = w.Write([]byte("┴"))
		}
	}
	_, _ = w.Write([]byte("┘"))
	_, _ = w.Write([]byte("\n"))
}

func drawRowLine(w io.Writer, columnWidths []int) {
	_, _ = w.Write([]byte("├"))
	for i, width := range columnWidths {
		_, _ = w.Write([]byte(strings.Repeat("─", width)))
		if i != len(columnWidths)-1 {
			_, _ = w.Write([]byte("┼"))
		}
	}
	_, _ = w.Write([]byte("┤"))
	_, _ = w.Write([]byte("\n"))
}

func drawRow(w io.Writer, cells []tableCell, columnWidths []int) {
	contents := make([][]string, len(cells))

	// As we draw the row line by line, we need a way to reset and recover
	// the formatting when we alternate between cells. To do that, we accumulate
	// the ongoing series of ANSI escape sequence for each cell and output them
	// again each time we switch to the next cell so we end up in the exact same
	// state. Inefficient but works.
	// TODO: this could use a way to compress the sequences into a few escape
	// sequences instead of growing unbounded.
	formatting := make([]strings.Builder, len(cells))

	accFormatting := func(cellIndex int, items []text.EscapeItem) {
		for _, item := range items {
			formatting[cellIndex].WriteString(item.Item)
		}
	}

	maxLength := 0

	// Wrap each cell content into multiple lines, depending on
	// how wide each cell is.
	for i, cell := range cells {
		wrapped, lines := text.Wrap(cell.content, columnWidths[i])
		contents[i] = strings.Split(wrapped, "\n")
		maxLength = max(maxLength, lines)
	}

	// Draw the row line by line
	for i := 0; i < maxLength; i++ {
		_, _ = w.Write([]byte("│"))
		for j, width := range columnWidths {
			content := ""
			if len(contents[j]) > i {
				content = contents[j][i]
				_, _ = w.Write([]byte(formatting[j].String()))
				_, _ = w.Write([]byte(content))
				// extract and accumulate the formatting
				_, seqs := text.ExtractTermEscapes(content)
				accFormatting(j, seqs)
				_, _ = w.Write([]byte(resetAll))
			}
			padding := strings.Repeat(" ", width-text.WordLen(content))
			_, _ = w.Write([]byte(padding))
			_, _ = w.Write([]byte("│"))
		}
		_, _ = w.Write([]byte("\n"))
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
