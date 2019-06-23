package markdown

import "github.com/MichaelMure/color"

var (
	Bold       = color.New(color.Bold).SprintFunc()
	Italic     = color.New(color.Italic).SprintFunc()
	CrossedOut = color.New(color.CrossedOut).SprintFunc()

	Green        = color.New(color.FgGreen).SprintFunc()
	HiGreen      = color.New(color.FgHiGreen).SprintFunc()
	GreenBold    = color.New(color.FgGreen, color.Bold).SprintFunc()
	Blue         = color.New(color.FgBlue).SprintFunc()
	BlueBgItalic = color.New(color.BgBlue, color.Italic).SprintFunc()
)
