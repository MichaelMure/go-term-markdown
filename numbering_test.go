package markdown

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_numbering(t *testing.T) {
	var n headingNumbering

	assert.Equal(t, "", n.Render())
	n.Observe(1)
	assert.Equal(t, "1", n.Render())
	n.Observe(1)
	assert.Equal(t, "2", n.Render())
	n.Observe(1)
	assert.Equal(t, "3", n.Render())
	n.Observe(2)
	assert.Equal(t, "3.1", n.Render())
	n.Observe(2)
	assert.Equal(t, "3.2", n.Render())
	n.Observe(2)
	assert.Equal(t, "3.3", n.Render())
	n.Observe(4)
	assert.Equal(t, "3.3.0.1", n.Render())
	n.Observe(4)
	assert.Equal(t, "3.3.0.2", n.Render())
	n.Observe(3)
	assert.Equal(t, "3.3.1", n.Render())
	n.Observe(3)
	assert.Equal(t, "3.3.2", n.Render())
	n.Observe(1)
	assert.Equal(t, "4", n.Render())
}
