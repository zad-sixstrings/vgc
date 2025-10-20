package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jackc/pgx/v5"
)

// ========== UTILITY FUNCTIONS ==========

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

// createSearchBar creates a search entry that filters data as user types
// Returns a container with the search bar that has a fixed minimum width
func createSearchBar(placeholder string, onSearch func(searchText string)) *fyne.Container {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder(placeholder)

	// Filter as user types
	searchEntry.OnChanged = func(text string) {
		onSearch(text)
	}

	// Wrap in a container with padding to give it breathing room
	// Using a fixed-size container ensures consistent width
	return container.NewGridWithColumns(1, searchEntry)
}

// ========== ACTION BUTTONS ==========

// createActionButtons creates the Add/Details/Edit/Delete button toolbar
func createActionButtons(w fyne.Window, conn *pgx.Conn, entityType string, detailsBtn, editBtn, deleteBtn *widget.Button, refreshFunc func()) fyne.CanvasObject {
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
	detailsBtn.Importance = widget.HighImportance  // Blue
	editBtn.Importance = widget.WarningImportance  // Yellow/Orange
	deleteBtn.Importance = widget.DangerImportance // Red

	// Start with Details, Edit and Delete disabled
	detailsBtn.Disable()
	editBtn.Disable()
	deleteBtn.Disable()

	return container.NewHBox(
		addBtn,
		detailsBtn,
		editBtn,
		deleteBtn,
	)
}

// ========== TABLE BUILDERS ==========

// buildGamesTableWithSelection creates the games table and tracks selection
func buildGamesTableWithSelection(w fyne.Window, conn *pgx.Conn, games []Game, detailsBtn, editBtn, deleteBtn *widget.Button, selectedGameID *int, refreshFunc func()) *widget.Table {
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

	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		headers := []string{"ID", "Titre", "Plateforme", "Genre", "État"}
		label.SetText(headers[id.Col])
	}

	table.OnSelected = func(id widget.TableCellID) {
		*selectedGameID = games[id.Row].GameID
		fmt.Printf("Selected game: %s (ID: %d)\n", games[id.Row].Title, *selectedGameID)
		detailsBtn.Enable()
		editBtn.Enable()
		deleteBtn.Enable()
	}

	table.OnUnselected = func(id widget.TableCellID) {
		*selectedGameID = -1
		detailsBtn.Disable()
		editBtn.Disable()
		deleteBtn.Disable()
	}

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 400)
	table.SetColumnWidth(2, 400)
	table.SetColumnWidth(3, 200)
	table.SetColumnWidth(4, 50)
	table.ShowHeaderColumn = false

	return table
}

// buildConsolesTableWithSelection creates the consoles table and tracks selection
func buildConsolesTableWithSelection(w fyne.Window, conn *pgx.Conn, consoles []Console, detailsBtn, editBtn, deleteBtn *widget.Button, selectedConsoleID *int, refreshFunc func()) *widget.Table {
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

	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		headers := []string{"ID", "Nom", "Fabricant", "Gen", "État"}
		label.SetText(headers[id.Col])
	}

	table.OnSelected = func(id widget.TableCellID) {
		*selectedConsoleID = consoles[id.Row].ConsoleID
		fmt.Printf("Selected console: %s (ID: %d)\n", consoles[id.Row].Name, *selectedConsoleID)
		detailsBtn.Enable()
		editBtn.Enable()
		deleteBtn.Enable()
	}

	table.OnUnselected = func(id widget.TableCellID) {
		*selectedConsoleID = -1
		detailsBtn.Disable()
		editBtn.Disable()
		deleteBtn.Disable()
	}

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 300)
	table.SetColumnWidth(2, 300)
	table.SetColumnWidth(3, 50)
	table.SetColumnWidth(4, 100)
	table.ShowHeaderColumn = false

	return table
}

// buildAccessoriesTableWithSelection creates the accessories table and tracks selection
func buildAccessoriesTableWithSelection(w fyne.Window, conn *pgx.Conn, accessories []Accessory, detailsBtn, editBtn, deleteBtn *widget.Button, selectedAccessoryID *int, refreshFunc func()) *widget.Table {
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

	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		headers := []string{"ID", "Nom", "Couleur", "Type", "Fabricant", "État"}
		label.SetText(headers[id.Col])
	}

	table.OnSelected = func(id widget.TableCellID) {
		*selectedAccessoryID = accessories[id.Row].AccessoryID
		fmt.Printf("Selected accessory: %s (ID: %d)\n", accessories[id.Row].Name, *selectedAccessoryID)
		detailsBtn.Enable()
		editBtn.Enable()
		deleteBtn.Enable()
	}

	table.OnUnselected = func(id widget.TableCellID) {
		*selectedAccessoryID = -1
		detailsBtn.Disable()
		editBtn.Disable()
		deleteBtn.Disable()
	}

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 300)
	table.SetColumnWidth(2, 150)
	table.SetColumnWidth(3, 150)
	table.SetColumnWidth(4, 200)
	table.SetColumnWidth(5, 100)
	table.ShowHeaderColumn = false

	return table
}

// ========== FILTER FUNCTIONS ==========

// filterGames returns games that match the search text (case-insensitive)
// Searches: Title, Console Name, Genre Name
func filterGames(games []Game, searchText string) []Game {
	if searchText == "" {
		return games
	}

	searchLower := strings.ToLower(searchText)
	var filtered []Game

	for _, game := range games {
		// Search in: Title, Console, Genre
		if strings.Contains(strings.ToLower(game.Title), searchLower) ||
			strings.Contains(strings.ToLower(game.ConsoleName), searchLower) ||
			strings.Contains(strings.ToLower(game.GenreName), searchLower) {
			filtered = append(filtered, game)
		}
	}

	return filtered
}

// filterConsoles returns consoles that match the search text (case-insensitive)
// Searches: Name, Manufacturer Name
func filterConsoles(consoles []Console, searchText string) []Console {
	if searchText == "" {
		return consoles
	}

	searchLower := strings.ToLower(searchText)
	var filtered []Console

	for _, console := range consoles {
		// Search in: Name, Manufacturer
		if strings.Contains(strings.ToLower(console.Name), searchLower) ||
			strings.Contains(strings.ToLower(console.ManufacturerName), searchLower) {
			filtered = append(filtered, console)
		}
	}

	return filtered
}

// filterAccessories returns accessories that match the search text (case-insensitive)
// Searches: Name, Type, Manufacturer, Color
func filterAccessories(accessories []Accessory, searchText string) []Accessory {
	if searchText == "" {
		return accessories
	}

	searchLower := strings.ToLower(searchText)
	var filtered []Accessory

	for _, accessory := range accessories {
		// Search in: Name, Type, Manufacturer, Color
		matchName := strings.Contains(strings.ToLower(accessory.Name), searchLower)
		matchType := strings.Contains(strings.ToLower(accessory.TypeName), searchLower)
		matchManufacturer := strings.Contains(strings.ToLower(accessory.ManufacturerName), searchLower)
		matchColor := accessory.Color != nil && strings.Contains(strings.ToLower(*accessory.Color), searchLower)

		if matchName || matchType || matchManufacturer || matchColor {
			filtered = append(filtered, accessory)
		}
	}

	return filtered
}

// ========== TAB BUILDERS ==========

// buildJeuxTab creates the complete "Jeux" tab content with search
func buildJeuxTab(w fyne.Window, conn *pgx.Conn, games []Game, refreshFunc func()) fyne.CanvasObject {
	var selectedGameID int = -1
	allGames := games

	// Create buttons
	detailsBtn := widget.NewButton("Détails", func() {
		if selectedGameID == -1 {
			return
		}
		showGameDetailDialog(w, conn, selectedGameID, func() {
			showEditGameDialog(w, conn, selectedGameID, refreshFunc)
		})
	})

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

		var gameName string
		for _, g := range allGames {
			if g.GameID == selectedGameID {
				gameName = g.Title
				break
			}
		}

		dialog.NewConfirm(
			"Supprimer le jeu",
			fmt.Sprintf("Êtes-vous sûr de vouloir supprimer '%s'? Cette action est irréversible.", gameName),
			func(confirmed bool) {
				if confirmed {
					err := deleteGame(conn, selectedGameID)
					if err != nil {
						dialog.ShowError(fmt.Errorf("échec de suppression: %w", err), w)
						return
					}
					dialog.ShowInformation("Succès", "Jeu supprimé avec succès!", w)
					refreshFunc()
				}
			},
			w,
		).Show()
	})

	actionButtons := createActionButtons(w, conn, "game", detailsBtn, editBtn, deleteBtn, refreshFunc)

	var tableContainer *fyne.Container

	rebuildTable := func(filteredGames []Game) {
		selectedGameID = -1
		detailsBtn.Disable()
		editBtn.Disable()
		deleteBtn.Disable()

		table := buildGamesTableWithSelection(w, conn, filteredGames, detailsBtn, editBtn, deleteBtn, &selectedGameID, refreshFunc)
		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
	}

	searchBar := createSearchBar("Rechercher un jeu...", func(searchText string) {
		filtered := filterGames(allGames, searchText)
		rebuildTable(filtered)
	})

	toolbar := container.NewBorder(
		nil, nil,
		actionButtons,
		nil,
		searchBar,
	)

	table := buildGamesTableWithSelection(w, conn, games, detailsBtn, editBtn, deleteBtn, &selectedGameID, refreshFunc)
	tableContainer = container.NewStack(table)

	return container.NewBorder(
		toolbar,
		nil, nil, nil,
		tableContainer,
	)
}

// buildConsolesTab creates the complete "Consoles" tab content with search
func buildConsolesTab(w fyne.Window, conn *pgx.Conn, consoles []Console, refreshFunc func()) fyne.CanvasObject {
	var selectedConsoleID int = -1
	allConsoles := consoles

	detailsBtn := widget.NewButton("Détails", func() {
		if selectedConsoleID == -1 {
			return
		}
		showConsoleDetailDialog(w, conn, selectedConsoleID, func() {
			showEditConsoleDialog(w, conn, selectedConsoleID, refreshFunc)
		})
	})

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

		var consoleName string
		for _, c := range allConsoles {
			if c.ConsoleID == selectedConsoleID {
				consoleName = c.Name
				break
			}
		}

		dialog.NewConfirm(
			"Supprimer la console",
			fmt.Sprintf("Êtes-vous sûr de vouloir supprimer '%s'? Cette action est irréversible.", consoleName),
			func(confirmed bool) {
				if confirmed {
					err := deleteConsole(conn, selectedConsoleID)
					if err != nil {
						dialog.ShowError(fmt.Errorf("échec de suppression: %w", err), w)
						return
					}
					dialog.ShowInformation("Succès", "Console supprimée avec succès!", w)
					refreshFunc()
				}
			},
			w,
		).Show()
	})

	actionButtons := createActionButtons(w, conn, "console", detailsBtn, editBtn, deleteBtn, refreshFunc)

	var tableContainer *fyne.Container

	rebuildTable := func(filteredConsoles []Console) {
		selectedConsoleID = -1
		detailsBtn.Disable()
		editBtn.Disable()
		deleteBtn.Disable()

		table := buildConsolesTableWithSelection(w, conn, filteredConsoles, detailsBtn, editBtn, deleteBtn, &selectedConsoleID, refreshFunc)
		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
	}

	searchBar := createSearchBar("Rechercher une console...", func(searchText string) {
		filtered := filterConsoles(allConsoles, searchText)
		rebuildTable(filtered)
	})

	toolbar := container.NewBorder(
		nil, nil,
		actionButtons,
		nil,
		searchBar,
	)

	table := buildConsolesTableWithSelection(w, conn, consoles, detailsBtn, editBtn, deleteBtn, &selectedConsoleID, refreshFunc)
	tableContainer = container.NewStack(table)

	return container.NewBorder(
		toolbar,
		nil, nil, nil,
		tableContainer,
	)
}

// buildAccessoiresTab creates the complete "Accessoires" tab content with search
func buildAccessoiresTab(w fyne.Window, conn *pgx.Conn, accessories []Accessory, refreshFunc func()) fyne.CanvasObject {
	var selectedAccessoryID int = -1
	allAccessories := accessories

	detailsBtn := widget.NewButton("Détails", func() {
		if selectedAccessoryID == -1 {
			return
		}
		showAccessoryDetailDialog(w, conn, selectedAccessoryID, func() {
			showEditAccessoryDialog(w, conn, selectedAccessoryID, refreshFunc)
		})
	})

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

		var accessoryName string
		for _, a := range allAccessories {
			if a.AccessoryID == selectedAccessoryID {
				accessoryName = a.Name
				break
			}
		}

		dialog.NewConfirm(
			"Supprimer l'accessoire",
			fmt.Sprintf("Êtes-vous sûr de vouloir supprimer '%s'? Cette action est irréversible.", accessoryName),
			func(confirmed bool) {
				if confirmed {
					err := deleteAccessory(conn, selectedAccessoryID)
					if err != nil {
						dialog.ShowError(fmt.Errorf("échec de suppression: %w", err), w)
						return
					}
					dialog.ShowInformation("Succès", "Accessoire supprimé avec succès!", w)
					refreshFunc()
				}
			},
			w,
		).Show()
	})

	actionButtons := createActionButtons(w, conn, "accessory", detailsBtn, editBtn, deleteBtn, refreshFunc)

	var tableContainer *fyne.Container

	rebuildTable := func(filteredAccessories []Accessory) {
		selectedAccessoryID = -1
		detailsBtn.Disable()
		editBtn.Disable()
		deleteBtn.Disable()

		table := buildAccessoriesTableWithSelection(w, conn, filteredAccessories, detailsBtn, editBtn, deleteBtn, &selectedAccessoryID, refreshFunc)
		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
	}

	searchBar := createSearchBar("Rechercher un accessoire...", func(searchText string) {
		filtered := filterAccessories(allAccessories, searchText)
		rebuildTable(filtered)
	})

	toolbar := container.NewBorder(
		nil, nil,
		actionButtons,
		nil,
		searchBar,
	)

	table := buildAccessoriesTableWithSelection(w, conn, accessories, detailsBtn, editBtn, deleteBtn, &selectedAccessoryID, refreshFunc)
	tableContainer = container.NewStack(table)

	return container.NewBorder(
		toolbar,
		nil, nil, nil,
		tableContainer,
	)
}
