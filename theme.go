package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// compactTheme
type compactTheme struct{}

func (t *compactTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t *compactTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *compactTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *compactTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	// Text sizes - reduced by 1-2 pixels for compact look
	case theme.SizeNameText:
		return 13 // Default: 14
	case theme.SizeNameHeadingText:
		return 20 // Default: 22
	case theme.SizeNameSubHeadingText:
		return 16 // Default: 18
	case theme.SizeNameCaptionText:
		return 10 // Default: 11

	// Padding & spacing - slightly reduced
	case theme.SizeNamePadding:
		return 3 // Default: 4
	case theme.SizeNameInnerPadding:
		return 6 // Default: 8
	case theme.SizeNameScrollBarSmall:
		return 2 // Default: 3

	// Input widgets
	case theme.SizeNameInputBorder:
		return 1 // Default: 1
	case theme.SizeNameInputRadius:
		return 4 // Default: 5

	// Icons - keep default or slightly smaller
	case theme.SizeNameInlineIcon:
		return 18 // Default: 20

	// Everything else: use default
	default:
		return theme.DefaultTheme().Size(name)
	}
}
