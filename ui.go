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

// ========== CUSTOM TAPPABLE LABEL FOR CONTEXT MENU ==========

// tappableLabel is a label that can detect right-clicks for context menus
type tappableLabel struct {
	widget.Label
	onRightClick func(pos fyne.Position)
}

// newTappableLabel creates a label that supports right-click
func newTappableLabel(onRightClick func(pos fyne.Position)) *tappableLabel {
	label := &tappableLabel{
		onRightClick: onRightClick,
	}
	label.ExtendBaseWidget(label)
	return label
}

// TappedSecondary handles right-click events
func (t *tappableLabel) TappedSecondary(ev *fyne.PointEvent) {
	if t.onRightClick != nil {
		t.onRightClick(ev.AbsolutePosition)
	}
}

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

// ========== TABLE BUILDERS ==========

// buildGamesTableWithSelection creates the games table and tracks selection
func buildGamesTableWithSelection(w fyne.Window, conn *pgx.Conn, games []Game, editBtn, deleteBtn *widget.Button, selectedGameID *int, refreshFunc func()) *widget.Table {
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(games), 5
		},
		func() fyne.CanvasObject {
			// Create tappable label with right-click handler
			return newTappableLabel(nil)
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*tappableLabel)
			game := games[id.Row]

			// Update label text
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

			// Set up right-click handler for this cell
			label.onRightClick = func(pos fyne.Position) {
				// Create context menu
				detailsItem := fyne.NewMenuItem("Détails", func() {
					showGameDetailDialog(w, conn, game.GameID, func() {
						showEditGameDialog(w, conn, game.GameID, refreshFunc)
					})
				})

				editItem := fyne.NewMenuItem("Éditer", func() {
					showEditGameDialog(w, conn, game.GameID, refreshFunc)
				})

				deleteItem := fyne.NewMenuItem("Supprimer", func() {
					dialog.NewConfirm(
						"Supprimer le jeu",
						fmt.Sprintf("Êtes-vous sûr de vouloir supprimer '%s'? Cette action est irréversible.", game.Title),
						func(confirmed bool) {
							if confirmed {
								err := deleteGame(conn, game.GameID)
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

				// Create and show popup menu
				menu := fyne.NewMenu("", detailsItem, editItem, deleteItem)
				popup := widget.NewPopUpMenu(menu, w.Canvas())
				popup.ShowAtPosition(pos)
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

// buildConsolesTableWithSelection creates the consoles table and tracks selection
func buildConsolesTableWithSelection(w fyne.Window, conn *pgx.Conn, consoles []Console, editBtn, deleteBtn *widget.Button, selectedConsoleID *int, refreshFunc func()) *widget.Table {
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(consoles), 5
		},
		func() fyne.CanvasObject {
			return newTappableLabel(nil)
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*tappableLabel)
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

			label.onRightClick = func(pos fyne.Position) {
				detailsItem := fyne.NewMenuItem("Détails", func() {
					showConsoleDetailDialog(w, conn, console.ConsoleID, func() {
						showEditConsoleDialog(w, conn, console.ConsoleID, refreshFunc)
					})
				})

				editItem := fyne.NewMenuItem("Éditer", func() {
					showEditConsoleDialog(w, conn, console.ConsoleID, refreshFunc)
				})

				deleteItem := fyne.NewMenuItem("Supprimer", func() {
					dialog.NewConfirm(
						"Supprimer la console",
						fmt.Sprintf("Êtes-vous sûr de vouloir supprimer '%s'? Cette action est irréversible.", console.Name),
						func(confirmed bool) {
							if confirmed {
								err := deleteConsole(conn, console.ConsoleID)
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

				menu := fyne.NewMenu("", detailsItem, editItem, deleteItem)
				popup := widget.NewPopUpMenu(menu, w.Canvas())
				popup.ShowAtPosition(pos)
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
		editBtn.Enable()
		deleteBtn.Enable()
	}

	table.OnUnselected = func(id widget.TableCellID) {
		*selectedConsoleID = -1
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
func buildAccessoriesTableWithSelection(w fyne.Window, conn *pgx.Conn, accessories []Accessory, editBtn, deleteBtn *widget.Button, selectedAccessoryID *int, refreshFunc func()) *widget.Table {
	table := widget.NewTableWithHeaders(
		func() (int, int) {
			return len(accessories), 6
		},
		func() fyne.CanvasObject {
			return newTappableLabel(nil)
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*tappableLabel)
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

			label.onRightClick = func(pos fyne.Position) {
				detailsItem := fyne.NewMenuItem("Détails", func() {
					showAccessoryDetailDialog(w, conn, accessory.AccessoryID, func() {
						showEditAccessoryDialog(w, conn, accessory.AccessoryID, refreshFunc)
					})
				})

				editItem := fyne.NewMenuItem("Éditer", func() {
					showEditAccessoryDialog(w, conn, accessory.AccessoryID, refreshFunc)
				})

				deleteItem := fyne.NewMenuItem("Supprimer", func() {
					dialog.NewConfirm(
						"Supprimer l'accessoire",
						fmt.Sprintf("Êtes-vous sûr de vouloir supprimer '%s'? Cette action est irréversible.", accessory.Name),
						func(confirmed bool) {
							if confirmed {
								err := deleteAccessory(conn, accessory.AccessoryID)
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

				menu := fyne.NewMenu("", detailsItem, editItem, deleteItem)
				popup := widget.NewPopUpMenu(menu, w.Canvas())
				popup.ShowAtPosition(pos)
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
		editBtn.Enable()
		deleteBtn.Enable()
	}

	table.OnUnselected = func(id widget.TableCellID) {
		*selectedAccessoryID = -1
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
	allGames := games // Keep original unfiltered data

	// Create buttons
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
		for _, g := range allGames {
			if g.GameID == selectedGameID {
				gameName = g.Title
				break
			}
		}

		// Show confirmation dialog
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

	actionButtons := createActionButtons(w, conn, "game", editBtn, deleteBtn, refreshFunc)

	// Create container for table (will be updated by search)
	var tableContainer *fyne.Container

	// Function to rebuild table with filtered data
	rebuildTable := func(filteredGames []Game) {
		selectedGameID = -1 // Reset selection when filtering
		editBtn.Disable()
		deleteBtn.Disable()

		table := buildGamesTableWithSelection(w, conn, filteredGames, editBtn, deleteBtn, &selectedGameID, refreshFunc)
		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
	}

	// Create search bar
	searchBar := createSearchBar("Rechercher un jeu...", func(searchText string) {
		filtered := filterGames(allGames, searchText)
		rebuildTable(filtered)
	})

	// Create toolbar with action buttons and search
	toolbar := container.NewBorder(
		nil, nil,
		actionButtons,
		nil,
		searchBar, // Search bar in center gets more space
	)

	// Initial table
	table := buildGamesTableWithSelection(w, conn, games, editBtn, deleteBtn, &selectedGameID, refreshFunc)
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
	allConsoles := consoles // Keep original unfiltered data

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
		for _, c := range allConsoles {
			if c.ConsoleID == selectedConsoleID {
				consoleName = c.Name
				break
			}
		}

		// Show confirmation dialog
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

	actionButtons := createActionButtons(w, conn, "console", editBtn, deleteBtn, refreshFunc)

	// Create container for table (will be updated by search)
	var tableContainer *fyne.Container

	// Function to rebuild table with filtered data
	rebuildTable := func(filteredConsoles []Console) {
		selectedConsoleID = -1 // Reset selection when filtering
		editBtn.Disable()
		deleteBtn.Disable()

		table := buildConsolesTableWithSelection(w, conn, filteredConsoles, editBtn, deleteBtn, &selectedConsoleID, refreshFunc)
		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
	}

	// Create search bar
	searchBar := createSearchBar("Rechercher une console...", func(searchText string) {
		filtered := filterConsoles(allConsoles, searchText)
		rebuildTable(filtered)
	})

	// Create toolbar with action buttons and search
	toolbar := container.NewBorder(
		nil, nil,
		actionButtons,
		nil,
		searchBar, // Search bar in center gets more space
	)

	// Initial table
	table := buildConsolesTableWithSelection(w, conn, consoles, editBtn, deleteBtn, &selectedConsoleID, refreshFunc)
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
	allAccessories := accessories // Keep original unfiltered data

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
		for _, a := range allAccessories {
			if a.AccessoryID == selectedAccessoryID {
				accessoryName = a.Name
				break
			}
		}

		// Show confirmation dialog
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

	actionButtons := createActionButtons(w, conn, "accessory", editBtn, deleteBtn, refreshFunc)

	// Create container for table (will be updated by search)
	var tableContainer *fyne.Container

	// Function to rebuild table with filtered data
	rebuildTable := func(filteredAccessories []Accessory) {
		selectedAccessoryID = -1 // Reset selection when filtering
		editBtn.Disable()
		deleteBtn.Disable()

		table := buildAccessoriesTableWithSelection(w, conn, filteredAccessories, editBtn, deleteBtn, &selectedAccessoryID, refreshFunc)
		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
	}

	// Create search bar
	searchBar := createSearchBar("Rechercher un accessoire...", func(searchText string) {
		filtered := filterAccessories(allAccessories, searchText)
		rebuildTable(filtered)
	})

	// Create toolbar with action buttons and search
	toolbar := container.NewBorder(
		nil, nil,
		actionButtons,
		nil,
		searchBar, // Search bar in center gets more space
	)

	// Initial table
	table := buildAccessoriesTableWithSelection(w, conn, accessories, editBtn, deleteBtn, &selectedAccessoryID, refreshFunc)
	tableContainer = container.NewStack(table)

	return container.NewBorder(
		toolbar,
		nil, nil, nil,
		tableContainer,
	)
}
