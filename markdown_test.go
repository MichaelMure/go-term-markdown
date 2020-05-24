package markdown

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	color.NoColor = false

	sourcepath := "testdata_source/"
	resultpath := "testdata_result/"

	err := filepath.Walk(sourcepath, func(fullpath string, info os.FileInfo, err error) error {
		require.NoError(t, err)

		if info.IsDir() {
			return nil
		}

		_, file := filepath.Split(fullpath)
		name := strings.TrimRight(file, ".md")

		t.Run(name, func(t *testing.T) {
			source, err := ioutil.ReadFile(path.Join(sourcepath, name+".md"))
			require.NoError(t, err)

			expected, err := ioutil.ReadFile(path.Join(resultpath, name+".txt"))
			require.NoError(t, err)

			var output bytes.Buffer
			err = Render(&output, source, 40, 4)
			require.NoError(t, err)

			assert.Equal(t, string(expected), output.String())
		})

		return nil
	})

	require.NoError(t, err)
}

func Test__DoRender(t *testing.T) {
	// This is not a real test, it's here to create the output testdata.
	// uncomment to generate a new test case
	// t.SkipNow()

	color.NoColor = false

	sourcepath := "testdata_source/"
	resultpath := "testdata_result/"

	err := filepath.Walk(sourcepath, func(fullpath string, info os.FileInfo, err error) error {
		require.NoError(t, err)

		if info.IsDir() {
			return nil
		}

		_, file := filepath.Split(fullpath)
		name := strings.TrimRight(file, ".md")

		// if name != "image" {
		// 	// if name != "README" {
		// 	// if name != "Amps and angle encoding" {
		// 	// if name != "Links, shortcut references" {
		// 	// if name != "Links, inline style" {
		// 	// if name != "Table" {
		// 	return nil
		// }

		source, err := ioutil.ReadFile(path.Join(sourcepath, name+".md"))
		require.NoError(t, err)

		var output bytes.Buffer
		err = Render(&output, source, 40, 4)

		err = ioutil.WriteFile(path.Join(resultpath, name+".txt"), output.Bytes(), 0666)
		require.NoError(t, err)

		return nil
	})

	require.NoError(t, err)
}
