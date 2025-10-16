package main

import (
	"context"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Connect to database
	conn, err := dbconnect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background()) // Close when main() exits

	// Create app
	a := app.New()
	w := a.NewWindow("VGC")
	// Contents
	accueilContent := widget.NewLabel("Accueil")
	jeuxContent := widget.NewLabel("Jeux")
	consolesContent := widget.NewLabel("Consoles")
	accessoiresContent := widget.NewLabel("Accessoires")

	// Sidebar
	sidebar := container.NewAppTabs(
		container.NewTabItem("Accueil", accueilContent),
		container.NewTabItem("Jeux", jeuxContent),
		container.NewTabItem("Consoles", consolesContent),
		container.NewTabItem("Accessoires", accessoiresContent),
	)
	sidebar.SetTabLocation(container.TabLocationLeading)

	// Run app
	w.SetContent(sidebar)
	w.Resize(fyne.NewSize(1600, 900))
	w.ShowAndRun()
}
