package misc

import (
	"io/ioutil"
	"os"

	markdown "github.com/MichaelMure/go-term-markdown"
)

func main() {
	path := "Readme.md"
	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = markdown.Render(os.Stdout, source, 80, 6)
	if err != nil {
		panic(err)
	}
}
