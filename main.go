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

	// Create app
	a := app.New()
	w := a.NewWindow("VGC")

	// Create sidebar with tabs
	sidebar := container.NewAppTabs(
		container.NewTabItem("Accueil", widget.NewLabel("Dashboard")),
		container.NewTabItem("Jeux", widget.NewLabel("Loading...")),
		container.NewTabItem("Consoles", widget.NewLabel("Loading...")),
		container.NewTabItem("Accessoires", widget.NewLabel("Loading...")),
	)
	sidebar.SetTabLocation(container.TabLocationLeading)

	// Declare refresh functions as variables first
	var refreshGamesTab func()
	var refreshConsolesTab func()
	var refreshAccessoriesTab func()

	// Now define them
	refreshGamesTab = func() {
		games, err := getGames(conn)
		if err != nil {
			log.Println("Error fetching games:", err)
			return
		}
		sidebar.Items[1].Content = buildJeuxTab(w, conn, games, refreshGamesTab)
		sidebar.Refresh()
	}

	refreshConsolesTab = func() {
		consoles, err := getConsoles(conn)
		if err != nil {
			log.Println("Error fetching consoles:", err)
			return
		}
		sidebar.Items[2].Content = buildConsolesTab(w, conn, consoles, refreshConsolesTab)
		sidebar.Refresh()
	}

	// Refresh accessories tab
	refreshAccessoriesTab = func() {
		accessories, err := getAccessories(conn)
		if err != nil {
			log.Println("Error fetching accessories:", err)
			return
		}
		sidebar.Items[3].Content = buildAccessoiresTab(w, conn, accessories, refreshAccessoriesTab)
		sidebar.Refresh()
	}

	// Initial load of data
	refreshGamesTab()
	refreshConsolesTab()
	refreshAccessoriesTab()

	// Run app
	w.SetContent(sidebar)
	w.Resize(fyne.NewSize(1600, 900))
	w.ShowAndRun()
}

