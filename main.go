package main

import (
	"context"
	"fmt"
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

	// Build games table
	gamesTable := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(games), 4 // +1 for header row
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)

			// Data rows
			game := games[id.Row]
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%d", game.GameID))
			case 1:
				label.SetText(game.Title)
			case 2:
				label.SetText(game.Platform)
			case 3:
				label.SetText(game.Genre)
			}
		},
	)
	// Custom headers for games table
	gamesTable.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		switch id.Col {
		case 0:
			label.SetText("ID")
		case 1:
			label.SetText("Titre")
		case 2:
			label.SetText("Plateforme")
		case 3:
			label.SetText("Genre")
		}
	}

	// Build games table
	consolesTable := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(consoles), 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)

			// Data rows
			console := consoles[id.Row]
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%d", console.ConsoleID))
			case 1:
				label.SetText(console.Name)
			case 2:
				label.SetText(console.Manufacturer)
			case 3:
				label.SetText(fmt.Sprintf("%d", console.Generation))
			}
		},
	)
	// Custom headers for games table
	consolesTable.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		switch id.Col {
		case 0:
			label.SetText("ID")
		case 1:
			label.SetText("Nom")
		case 2:
			label.SetText("Frabriquant")
		case 3:
			label.SetText("Gen")
		}
	}
	// Set columns width
	gamesTable.SetColumnWidth(0, 50)  // ID column - narrow
	gamesTable.SetColumnWidth(1, 400) // Title - wide
	gamesTable.SetColumnWidth(2, 400) // Platform - medium
	gamesTable.SetColumnWidth(3, 100) // Genre - medium
	consolesTable.SetColumnWidth(0, 50)
	consolesTable.SetColumnWidth(1, 300)
	consolesTable.SetColumnWidth(2, 300)
	consolesTable.SetColumnWidth(3, 50)

	// Hide header column
	gamesTable.ShowHeaderColumn = false
	consolesTable.ShowHeaderColumn = false

	// Contents
	accueilContent := widget.NewLabel("Dashboard")
	jeuxContent := container.NewBorder(
		widget.NewLabel("Jeux"),
		nil,
		nil,
		nil,
		gamesTable,
	)
	consolesContent := container.NewBorder(
		widget.NewLabel("Consoles"),
		nil,
		nil,
		nil,
		consolesTable,
	)
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
