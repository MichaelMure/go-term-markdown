package fuzzing

import (
	"io/ioutil"

	markdown "github.com/MichaelMure/go-term-markdown"
)

func Fuzz(data []byte) int {
	err := markdown.Render(ioutil.Discard, data, 50, 4)
	if err != nil {
		panic(err)
	}
	return 1
}
