package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("VGC")

	w.SetContent(widget.NewLabel("Ca fonctionne genre?"))
	w.Resize(fyne.NewSize(1600, 900))
	w.ShowAndRun()
}
