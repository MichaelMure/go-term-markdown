package markdown

import (
	"strings"
	"testing"

	"github.com/russross/blackfriday"
	"github.com/stretchr/testify/assert"
)

func TestColumnWidths(t *testing.T) {
	const lineWidth = 40

	cases := []struct {
		cellWidths []int
		expected   []int
		truncated  bool
	}{
		{
			[]int{0},
			[]int{0},
			false,
		},
		{
			[]int{0, 0, 0, 0, 0},
			[]int{0, 0, 0, 0, 0},
			false,
		},
		{
			[]int{1, 2, 3, 4, 5},
			[]int{1, 2, 3, 4, 5},
			false,
		},
		{
			// overflow, one column
			[]int{60},
			[]int{38},
			false,
		},
		{
			// overflow, multiple columns
			[]int{30, 30, 30},
			[]int{12, 12, 12}, // (40-4)/3
			false,
		},
		{
			// overflow, different columns
			[]int{30, 60, 30},
			[]int{12, 12, 12},
			false,
		},
		{
			// overflow, different columns with one small enough
			[]int{10, 30, 30},
			[]int{10, 13, 13},
			false,
		},
		{
			// overflow, different columns with one small enough
			[]int{10, 60, 30},
			[]int{10, 13, 13},
			false,
		},
		{
			// too much columns
			[]int{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10},
			[]int{5, 5, 5, 5, 5, 5},
			true,
		},
	}

	for _, tc := range cases {
		tr := newTableRenderer()

		for _, w := range tc.cellWidths {
			tr.AddHeaderCell(strings.Repeat("a", w), blackfriday.TableAlignmentLeft)
		}
		tr.NextBodyRow()
		for _, w := range tc.cellWidths {
			tr.AddBodyCell(strings.Repeat("a", w))
		}

		result, truncated := tr.columnWidths(lineWidth)

		assert.Equal(t, tc.expected, result)
		assert.Equal(t, tc.truncated, truncated)
	}
}
