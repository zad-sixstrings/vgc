package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create app
	a := app.New()
	w := a.NewWindow("VGC")

	// Sidebar
	sidebar := container.NewVBox(
		widget.NewButton("Collection", func() {}),
		widget.NewButton("Jeux", func() {}),
		widget.NewButton("Consoles", func() {}),
		widget.NewButton("Accessoires", func() {}),
	)

	// Main content view
	mainContent := container.NewStack()

	// Split layout
	split := container.NewHSplit(sidebar, mainContent)
	split.SetOffset((0.1))

	// Run app
	w.SetContent(split)
	w.Resize(fyne.NewSize(1600, 900))
	w.ShowAndRun()
}
