package markdown

import "github.com/eliukblau/pixterm/pkg/ansimage"

type Options func(cfg *config)

// DitheringMode type is used for image scale dithering mode constants.
type DitheringMode uint8

const (
	NoDithering = DitheringMode(iota)
	DitheringWithBlocks
	DitheringWithChars
)

// Dithering mode for ansimage
// Default is fine directly through a terminal
// DitheringWithBlocks is recommended if a terminal UI library is used
func WithImageDithering(mode DitheringMode) Options {
	return func(cfg *config) {
		cfg.imageDithering = ansimage.DitheringMode(mode)
	}
}

// Use a custom collection of ANSI colors for the headings
func WithHeadingShades(shades []shadeFmt) Options {
	return func(cfg *config) {
		cfg.headingShade = shade(shades)
	}
}

// Use a custom collection of ANSI colors for the blockquotes
func WithBlockquoteShades(shades []shadeFmt) Options {
	return func(cfg *config) {
		cfg.blockQuoteShade = shade(shades)
	}
}
