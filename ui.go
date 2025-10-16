package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// createActionButtons creates the Add/Edit/Delete button toolbar
func createActionButtons(editBtn, deleteBtn *widget.Button) fyne.CanvasObject {
	addBtn := widget.NewButton("Add", func() {
		fmt.Println("Add clicked")
		// TODO: Open add dialog
	})

	// Style buttons with colors
	addBtn.Importance = widget.SuccessImportance   // Green
	editBtn.Importance = widget.WarningImportance  // Yellow/Orange
	deleteBtn.Importance = widget.DangerImportance // Red

	// Start with Edit and Delete disabled
	editBtn.Disable()
	deleteBtn.Disable()

	return container.NewHBox(
		addBtn,
		editBtn,
		deleteBtn,
	)
}

// buildGamesTable creates and configures the games table
func buildGamesTable(games []Game, editBtn, deleteBtn *widget.Button) *widget.Table {
	var selectedGameID int = -1 // Track selected game

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
				label.SetText(game.ConsoleName) // Changed from Platform
			case 3:
				label.SetText(game.GenreName) // Changed from Genre
			}
		},
	)

	// Custom headers
	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		headers := []string{"ID", "Titre", "Plateforme", "Genre"}
		label.SetText(headers[id.Col])
	}

	// Handle row selection
	table.OnSelected = func(id widget.TableCellID) {
		selectedGameID = games[id.Row].GameID
		fmt.Printf("Selected game: %s (ID: %d)\n", games[id.Row].Title, selectedGameID)

		// Enable Edit and Delete buttons
		editBtn.Enable()
		deleteBtn.Enable()
	}

	// Handle deselection (clicking same row or elsewhere)
	table.OnUnselected = func(id widget.TableCellID) {
		selectedGameID = -1
		editBtn.Disable()
		deleteBtn.Disable()
	}

	// Column widths
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 400)
	table.SetColumnWidth(2, 400)
	table.SetColumnWidth(3, 100)

	table.ShowHeaderColumn = false

	return table
}

// buildConsolesTable creates and configures the consoles table
func buildConsolesTable(consoles []Console, editBtn, deleteBtn *widget.Button) *widget.Table {
	var selectedConsoleID int = -1 // Track selected console

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
				label.SetText(console.ManufacturerName) // Changed from Manufacturer
			case 3:
				// Handle nullable Generation
				if console.Generation != nil {
					label.SetText(fmt.Sprintf("%d", *console.Generation))
				} else {
					label.SetText("")
				}
			}
		},
	)

	// Custom headers
	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		headers := []string{"ID", "Nom", "Fabricant", "Gen"}
		label.SetText(headers[id.Col])
	}

	// Handle row selection
	table.OnSelected = func(id widget.TableCellID) {
		selectedConsoleID = consoles[id.Row].ConsoleID
		fmt.Printf("Selected console: %s (ID: %d)\n", consoles[id.Row].Name, selectedConsoleID)

		// Enable Edit and Delete buttons
		editBtn.Enable()
		deleteBtn.Enable()
	}

	// Handle deselection
	table.OnUnselected = func(id widget.TableCellID) {
		selectedConsoleID = -1
		editBtn.Disable()
		deleteBtn.Disable()
	}

	// Set column widths
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 300)
	table.SetColumnWidth(2, 300)
	table.SetColumnWidth(3, 50)

	table.ShowHeaderColumn = false

	return table
}

// buildJeuxTab creates the complete "Jeux" tab content
func buildJeuxTab(games []Game) fyne.CanvasObject {
	// Create Edit and Delete buttons (to be managed by table selection)
	editBtn := widget.NewButton("Edit", func() {
		fmt.Println("Edit game clicked")
		// TODO: Open edit dialog with selected game
	})

	deleteBtn := widget.NewButton("Delete", func() {
		fmt.Println("Delete game clicked")
		// TODO: Show confirmation dialog and delete
	})

	actionButtons := createActionButtons(editBtn, deleteBtn)
	table := buildGamesTable(games, editBtn, deleteBtn)

	return container.NewBorder(
		actionButtons, // top - buttons in top right
		nil,           // bottom
		nil,           // left
		nil,           // right
		table,         // center - table fills remaining space
	)
}

// buildConsolesTab creates the complete "Consoles" tab content
func buildConsolesTab(consoles []Console) fyne.CanvasObject {
	// Create Edit and Delete buttons (to be managed by table selection)
	editBtn := widget.NewButton("Edit", func() {
		fmt.Println("Edit console clicked")
		// TODO: Open edit dialog with selected console
	})

	deleteBtn := widget.NewButton("Delete", func() {
		fmt.Println("Delete console clicked")
		// TODO: Show confirmation dialog and delete
	})

	actionButtons := createActionButtons(editBtn, deleteBtn)
	table := buildConsolesTable(consoles, editBtn, deleteBtn)

	return container.NewBorder(
		actionButtons, // top - buttons
		nil,           // bottom
		nil,           // left
		nil,           // right
		table,         // center
	)
}
