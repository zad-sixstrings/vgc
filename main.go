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
	defer conn.Close(context.Background())

	// Fetch data
	games, err := getGames(conn)
	if err != nil {
		log.Fatal(err)
	}

	consoles, err := getConsoles(conn)
	if err != nil {
		log.Fatal(err)
	}

	// Create app
	a := app.New()
	w := a.NewWindow("VGC")

	// Build tab contents
	accueilContent := widget.NewLabel("Dashboard") // TODO
	jeuxContent := buildJeuxTab(games)
	consolesContent := buildConsolesTab(consoles)
	accessoiresContent := widget.NewLabel("Accessoires") // TODO

	// Create sidebar with tabs
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
