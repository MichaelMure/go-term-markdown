package fuzzing

import markdown "github.com/MichaelMure/go-term-markdown"

func Fuzz(data []byte) int {
	markdown.Render(string(data), 50, 4)
	return 1
}
