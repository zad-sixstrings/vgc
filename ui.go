package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jackc/pgx/v5"
)

// Condition as stars notation
func conditionToStars(condition *int) string {
	if condition == nil {
		return ""
	}

	stars := ""
	for i := 0; i < *condition; i++ {
		stars += "★"
	}
	for i := *condition; i < 5; i++ {
		stars += "☆"
	}
	return stars
}

// createActionButtons creates the Add/Edit/Delete button toolbar
func createActionButtons(w fyne.Window, conn *pgx.Conn, entityType string, editBtn, deleteBtn *widget.Button, refreshFunc func()) fyne.CanvasObject {
	addBtn := widget.NewButton("Ajouter", func() {
		if entityType == "game" {
			showAddGameDialog(w, conn, refreshFunc)
		} else if entityType == "console" {
			showAddConsoleDialog(w, conn, refreshFunc)
		} else if entityType == "accessory" {
			showAddAccessoryDialog(w, conn, refreshFunc)
		}
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

// buildGamesTableWithSelection creates the games table and tracks selection
func buildGamesTableWithSelection(games []Game, editBtn, deleteBtn *widget.Button, selectedGameID *int) *widget.Table {
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(games), 5
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
				label.SetText(game.ConsoleName)
			case 3:
				label.SetText(game.GenreName)
			case 4:
				label.SetText(conditionToStars(game.Condition))
			}
		},
	)

	// Custom headers
	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		headers := []string{"ID", "Titre", "Plateforme", "Genre", "État"}
		label.SetText(headers[id.Col])
	}

	// Handle row selection
	table.OnSelected = func(id widget.TableCellID) {
		*selectedGameID = games[id.Row].GameID
		fmt.Printf("Selected game: %s (ID: %d)\n", games[id.Row].Title, *selectedGameID)

		// Enable Edit and Delete buttons
		editBtn.Enable()
		deleteBtn.Enable()
	}

	// Handle deselection
	table.OnUnselected = func(id widget.TableCellID) {
		*selectedGameID = -1
		editBtn.Disable()
		deleteBtn.Disable()
	}

	// Column widths
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 400)
	table.SetColumnWidth(2, 400)
	table.SetColumnWidth(3, 200)
	table.SetColumnWidth(4, 50)

	table.ShowHeaderColumn = false

	return table
}

// buildConsolesTableWithSelection - creates the consoles table and tracks selection
func buildConsolesTableWithSelection(consoles []Console, editBtn, deleteBtn *widget.Button, selectedConsoleID *int) *widget.Table {
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(consoles), 5
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
				label.SetText(console.ManufacturerName)
			case 3:
				if console.Generation != nil {
					label.SetText(fmt.Sprintf("%d", *console.Generation))
				} else {
					label.SetText("")
				}
			case 4:
				label.SetText(conditionToStars(console.Condition))
			}
		},
	)

	// Custom headers
	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		headers := []string{"ID", "Nom", "Fabricant", "Gen", "État"}
		label.SetText(headers[id.Col])
	}

	// Handle row selection
	table.OnSelected = func(id widget.TableCellID) {
		*selectedConsoleID = consoles[id.Row].ConsoleID
		fmt.Printf("Selected console: %s (ID: %d)\n", consoles[id.Row].Name, *selectedConsoleID)
		editBtn.Enable()
		deleteBtn.Enable()
	}

	// Handle deselection
	table.OnUnselected = func(id widget.TableCellID) {
		*selectedConsoleID = -1
		editBtn.Disable()
		deleteBtn.Disable()
	}

	// Set column widths
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 300)
	table.SetColumnWidth(2, 300)
	table.SetColumnWidth(3, 50)
	table.SetColumnWidth(4, 100)

	table.ShowHeaderColumn = false

	return table
}

// buildAccessoriesTableWithSelection creates and configures the accessories table
func buildAccessoriesTableWithSelection(accessories []Accessory, editBtn, deleteBtn *widget.Button, selectedAccessoryID *int) *widget.Table {
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(accessories), 6
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			accessory := accessories[id.Row]

			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%d", accessory.AccessoryID))
			case 1:
				label.SetText(accessory.Name)
			case 2:
				if accessory.Color != nil {
					label.SetText(*accessory.Color)
				} else {
					label.SetText("")
				}
			case 3:
				label.SetText(accessory.TypeName)
			case 4:
				label.SetText(accessory.ManufacturerName)
			case 5:
				label.SetText(conditionToStars(accessory.Condition))
			}
		},
	)

	// Custom headers
	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		headers := []string{"ID", "Nom", "Couleur", "Type", "Fabricant", "État"}
		label.SetText(headers[id.Col])
	}

	// Handle row selection
	table.OnSelected = func(id widget.TableCellID) {
		*selectedAccessoryID = accessories[id.Row].AccessoryID
		fmt.Printf("Selected accessory: %s (ID: %d)\n", accessories[id.Row].Name, *selectedAccessoryID)
		editBtn.Enable()
		deleteBtn.Enable()
	}

	// Handle deselection
	table.OnUnselected = func(id widget.TableCellID) {
		*selectedAccessoryID = -1
		editBtn.Disable()
		deleteBtn.Disable()
	}

	// Column widths
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 300)
	table.SetColumnWidth(2, 150)
	table.SetColumnWidth(3, 150)
	table.SetColumnWidth(4, 200)
	table.SetColumnWidth(5, 100)

	table.ShowHeaderColumn = false

	return table
}

// buildJeuxTab creates the complete "Jeux" tab content
func buildJeuxTab(w fyne.Window, conn *pgx.Conn, games []Game, refreshFunc func()) fyne.CanvasObject {
	var selectedGameID int = -1 // Track which game is selected

	editBtn := widget.NewButton("Éditer", func() {
		if selectedGameID == -1 {
			return
		}
		showEditGameDialog(w, conn, selectedGameID, refreshFunc)
	})

	deleteBtn := widget.NewButton("Supprimer", func() {
		if selectedGameID == -1 {
			return
		}

		// Find the game name for confirmation dialog
		var gameName string
		for _, g := range games {
			if g.GameID == selectedGameID {
				gameName = g.Title
				break
			}
		}

		// Show confirmation dialog
		dialog.NewConfirm(
			"Delete Game",
			fmt.Sprintf("Are you sure you want to delete '%s'? This cannot be undone.", gameName),
			func(confirmed bool) {
				if confirmed {
					err := deleteGame(conn, selectedGameID)
					if err != nil {
						dialog.ShowError(fmt.Errorf("failed to delete game: %w", err), w)
						return
					}
					dialog.ShowInformation("Success", "Game deleted successfully!", w)
					refreshFunc()
				}
			},
			w,
		).Show()
	})

	actionButtons := createActionButtons(w, conn, "game", editBtn, deleteBtn, refreshFunc)

	// Build table with selection tracking
	table := buildGamesTableWithSelection(games, editBtn, deleteBtn, &selectedGameID)

	return container.NewBorder(
		actionButtons,
		nil, nil, nil,
		table,
	)
}

// buildConsolesTab creates the complete "Consoles" tab content
func buildConsolesTab(w fyne.Window, conn *pgx.Conn, consoles []Console, refreshFunc func()) fyne.CanvasObject {
	var selectedConsoleID int = -1

	editBtn := widget.NewButton("Éditer", func() {
		if selectedConsoleID == -1 {
			return
		}
		showEditConsoleDialog(w, conn, selectedConsoleID, refreshFunc)
	})

	deleteBtn := widget.NewButton("Supprimer", func() {
		if selectedConsoleID == -1 {
			return
		}

		// Find the console name
		var consoleName string
		for _, c := range consoles {
			if c.ConsoleID == selectedConsoleID {
				consoleName = c.Name
				break
			}
		}

		// Show confirmation dialog
		dialog.NewConfirm(
			"Delete Console",
			fmt.Sprintf("Are you sure you want to delete '%s'? This cannot be undone.", consoleName),
			func(confirmed bool) {
				if confirmed {
					err := deleteConsole(conn, selectedConsoleID)
					if err != nil {
						dialog.ShowError(fmt.Errorf("failed to delete console: %w", err), w)
						return
					}
					dialog.ShowInformation("Success", "Console deleted successfully!", w)
					refreshFunc()
				}
			},
			w,
		).Show()
	})

	actionButtons := createActionButtons(w, conn, "console", editBtn, deleteBtn, refreshFunc)
	table := buildConsolesTableWithSelection(consoles, editBtn, deleteBtn, &selectedConsoleID)

	return container.NewBorder(
		actionButtons,
		nil, nil, nil,
		table,
	)
}

// buildAccessoiresTab creates the complete "Accessoires" tab content
func buildAccessoiresTab(w fyne.Window, conn *pgx.Conn, accessories []Accessory, refreshFunc func()) fyne.CanvasObject {
	var selectedAccessoryID int = -1

	editBtn := widget.NewButton("Éditer", func() {
		if selectedAccessoryID == -1 {
			return
		}
		showEditAccessoryDialog(w, conn, selectedAccessoryID, refreshFunc)
	})

	deleteBtn := widget.NewButton("Supprimer", func() {
		if selectedAccessoryID == -1 {
			return
		}

		// Find the accessory name
		var accessoryName string
		for _, a := range accessories {
			if a.AccessoryID == selectedAccessoryID {
				accessoryName = a.Name
				break
			}
		}

		// Show confirmation dialog
		dialog.NewConfirm(
			"Delete Accessory",
			fmt.Sprintf("Are you sure you want to delete '%s'? This cannot be undone.", accessoryName),
			func(confirmed bool) {
				if confirmed {
					err := deleteAccessory(conn, selectedAccessoryID)
					if err != nil {
						dialog.ShowError(fmt.Errorf("failed to delete accessory: %w", err), w)
						return
					}
					dialog.ShowInformation("Success", "Accessory deleted successfully!", w)
					refreshFunc()
				}
			},
			w,
		).Show()
	})

	actionButtons := createActionButtons(w, conn, "accessory", editBtn, deleteBtn, refreshFunc)
	table := buildAccessoriesTableWithSelection(accessories, editBtn, deleteBtn, &selectedAccessoryID)

	return container.NewBorder(
		actionButtons,
		nil, nil, nil,
		table,
	)
}
