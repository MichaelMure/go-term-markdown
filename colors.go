package markdown

import "github.com/MichaelMure/color"

var (
	Bold       = color.New(color.Bold)
	Italic     = color.New(color.Italic)
	CrossedOut = color.New(color.CrossedOut)

	Green        = color.New(color.FgGreen).SprintFunc()
	HiGreen      = color.New(color.FgHiGreen).SprintFunc()
	GreenBold    = color.New(color.FgGreen, color.Bold).SprintFunc()
	Blue         = color.New(color.FgBlue).SprintFunc()
	BlueBgItalic = color.New(color.BgBlue, color.Italic).SprintFunc()
)
