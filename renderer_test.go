package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveLineBreak(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			"hello\nhello",
			"hello hello",
		},
		{
			"hello \nhello",
			"hello hello",
		},
		{
			"hello\n hello",
			"hello hello",
		},
		{
			"    hello    hello   \n   hello  hello   ",
			"    hello    hello hello  hello   ",
		},
		{
			"    hello   ",
			"    hello   ",
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, removeLineBreak(tt.input))
	}
}
