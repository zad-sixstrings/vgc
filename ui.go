package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Create and configures the games table
func buildGamesTable(games []Game) *widget.Table {
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(games), 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
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

	// Custom headers
	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		headers := []string{"ID", "Titre", "Plateforme", "Genre"}
		label.SetText(headers[id.Col])
	}

	// Column widths
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 400)
	table.SetColumnWidth(2, 400)
	table.SetColumnWidth(3, 100)

	table.ShowHeaderColumn = false

	return table
}

// Create and configures the consoles table
func buildConsolesTable(consoles []Console) *widget.Table {
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(consoles), 4
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
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

	// Custom headers
	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		headers := []string{"ID", "Nom", "Frabriquant", "Gen"}
		label.SetText(headers[id.Col])
	}

	// Column widths
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 300)
	table.SetColumnWidth(2, 300)
	table.SetColumnWidth(3, 50)

	table.ShowHeaderColumn = false

	return table
}

// Create the games tab content
func buildJeuxTab(games []Game) fyne.CanvasObject {
	table := buildGamesTable(games)
	return container.NewBorder(
		widget.NewLabel("Jeux"),
		nil, nil, nil,
		table,
	)
}

// Create the consoles tab content
func buildConsolesTab(consoles []Console) fyne.CanvasObject {
	table := buildConsolesTable(consoles)
	return container.NewBorder(
		widget.NewLabel("Consoles"),
		nil, nil, nil,
		table,
	)
}
