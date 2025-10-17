package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/jackc/pgx/v5"
)

// gameFormData holds all the form fields and their data
type gameFormData struct {
	// Form widgets
	titleEntry         *widget.Entry
	consoleSelect      *widget.Select
	genreSelect        *widget.Select
	jpReleaseDateEntry *widget.Entry
	usReleaseDateEntry *widget.Entry
	euReleaseDateEntry *widget.Entry
	jpRatingSelect     *widget.Select
	usRatingSelect     *widget.Select
	euRatingSelect     *widget.Select
	unitsSoldEntry     *widget.Entry
	ownedCheck         *widget.Check
	boxOwnedCheck      *widget.Check
	collectorCheck     *widget.Check
	conditionSlider    *widget.Slider
	conditionLabel     *widget.Label
	purchaseDateEntry  *widget.Entry
	purchasePriceEntry *widget.Entry
	notesEntry         *widget.Entry

	// Many-to-many data
	selectedDevelopers   []string
	selectedDeveloperIDs []int
	developersList       *widget.Label
	selectedComposers    []string
	selectedComposerIDs  []int
	composersList        *widget.Label
	selectedPublishers   []string
	selectedPublisherIDs []int
	publishersList       *widget.Label
	selectedProducers    []string
	selectedProducerIDs  []int
	producersList        *widget.Label

	// Lookup maps
	consoleMap map[string]int
	genreMap   map[string]int
	ratingMap  map[string]int

	// The complete form container
	form *fyne.Container
}

// buildGameForm creates the game form, optionally pre-populated with existing data
func buildGameForm(w fyne.Window, conn *pgx.Conn, existingGame *Game) *gameFormData {
	formData := &gameFormData{
		consoleMap: make(map[string]int),
		genreMap:   make(map[string]int),
		ratingMap:  make(map[string]int),
	}

	// Fetch lookup data
	consoles, _ := getConsoles(conn)
	genres, _ := getGenres(conn)
	ratings, _ := getRatingSystems(conn)
	developers, _ := getDevelopers(conn)
	composers, _ := getComposers(conn)
	publishers, _ := getPublishers(conn)
	producers, _ := getProducers(conn)

	// ========== Basic Info ==========
	formData.titleEntry = widget.NewEntry()
	formData.titleEntry.SetPlaceHolder("Titre (requis)")
	if existingGame != nil {
		formData.titleEntry.SetText(existingGame.Title)
	}

	// Console dropdown
	consoleOptions := []string{""}
	var selectedConsoleName string
	for _, c := range consoles {
		consoleOptions = append(consoleOptions, c.Name)
		formData.consoleMap[c.Name] = c.ConsoleID
		if existingGame != nil && existingGame.ConsoleID != nil && c.ConsoleID == *existingGame.ConsoleID {
			selectedConsoleName = c.Name
		}
	}
	formData.consoleSelect = widget.NewSelect(consoleOptions, nil)
	formData.consoleSelect.PlaceHolder = "Plateforme (requis)"
	if selectedConsoleName != "" {
		formData.consoleSelect.SetSelected(selectedConsoleName)
	}

	// Genre dropdown
	genreOptions := []string{""}
	var selectedGenreName string
	for _, g := range genres {
		genreOptions = append(genreOptions, g.Name)
		formData.genreMap[g.Name] = g.GenreID
		if existingGame != nil && existingGame.GenreID != nil && g.GenreID == *existingGame.GenreID {
			selectedGenreName = g.Name
		}
	}
	formData.genreSelect = widget.NewSelect(genreOptions, nil)
	formData.genreSelect.PlaceHolder = "Genre"
	if selectedGenreName != "" {
		formData.genreSelect.SetSelected(selectedGenreName)
	}

	// ========== Release Dates ==========
	formData.jpReleaseDateEntry = widget.NewEntry()
	formData.jpReleaseDateEntry.SetPlaceHolder("AAAA-MM-JJ")
	if existingGame != nil && existingGame.JPReleaseDate != nil {
		formData.jpReleaseDateEntry.SetText(existingGame.JPReleaseDate.Format("2006-01-02"))
	}

	formData.usReleaseDateEntry = widget.NewEntry()
	formData.usReleaseDateEntry.SetPlaceHolder("AAAA-MM-JJ")
	if existingGame != nil && existingGame.USReleaseDate != nil {
		formData.usReleaseDateEntry.SetText(existingGame.USReleaseDate.Format("2006-01-02"))
	}

	formData.euReleaseDateEntry = widget.NewEntry()
	formData.euReleaseDateEntry.SetPlaceHolder("AAAA-MM-JJ")
	if existingGame != nil && existingGame.EUReleaseDate != nil {
		formData.euReleaseDateEntry.SetText(existingGame.EUReleaseDate.Format("2006-01-02"))
	}

	// ========== Ratings ==========
	jpRatings := []string{""}
	usRatings := []string{""}
	euRatings := []string{""}
	var selectedJPRating, selectedUSRating, selectedEURating string

	for _, r := range ratings {
		label := fmt.Sprintf("%s - %s", r.Code, r.Region)
		formData.ratingMap[label] = r.RatingID

		switch r.Region {
		case "JP":
			jpRatings = append(jpRatings, label)
			if existingGame != nil && existingGame.JPRatingID != nil && r.RatingID == *existingGame.JPRatingID {
				selectedJPRating = label
			}
		case "US":
			usRatings = append(usRatings, label)
			if existingGame != nil && existingGame.USRatingID != nil && r.RatingID == *existingGame.USRatingID {
				selectedUSRating = label
			}
		case "EU":
			euRatings = append(euRatings, label)
			if existingGame != nil && existingGame.EURatingID != nil && r.RatingID == *existingGame.EURatingID {
				selectedEURating = label
			}
		}
	}

	formData.euRatingSelect = widget.NewSelect(euRatings, nil)
	formData.euRatingSelect.PlaceHolder = "Classification EU"
	if selectedEURating != "" {
		formData.euRatingSelect.SetSelected(selectedEURating)
	}

	formData.usRatingSelect = widget.NewSelect(usRatings, nil)
	formData.usRatingSelect.PlaceHolder = "Classification US"
	if selectedUSRating != "" {
		formData.usRatingSelect.SetSelected(selectedUSRating)
	}

	formData.jpRatingSelect = widget.NewSelect(jpRatings, nil)
	formData.jpRatingSelect.PlaceHolder = "Classification JP"
	if selectedJPRating != "" {
		formData.jpRatingSelect.SetSelected(selectedJPRating)
	}

	// ========== Units Sold ==========
	formData.unitsSoldEntry = widget.NewEntry()
	formData.unitsSoldEntry.SetPlaceHolder("Total des copies vendues")
	if existingGame != nil && existingGame.UnitsSold != nil {
		formData.unitsSoldEntry.SetText(fmt.Sprintf("%d", *existingGame.UnitsSold))
	}

	// ========== Collection Info ==========
	formData.ownedCheck = widget.NewCheck("Possédé", nil)
	if existingGame != nil {
		formData.ownedCheck.Checked = existingGame.Owned
	} else {
		formData.ownedCheck.Checked = true
	}

	formData.boxOwnedCheck = widget.NewCheck("Boîte possédée", nil)
	if existingGame != nil && existingGame.BoxOwned != nil {
		formData.boxOwnedCheck.Checked = *existingGame.BoxOwned
	}

	formData.collectorCheck = widget.NewCheck("Édition collector", nil)
	if existingGame != nil && existingGame.Collector != nil {
		formData.collectorCheck.Checked = *existingGame.Collector
	}

	// Condition
	formData.conditionSlider = widget.NewSlider(1, 5)
	formData.conditionSlider.Step = 1
	if existingGame != nil && existingGame.Condition != nil {
		formData.conditionSlider.Value = float64(*existingGame.Condition)
	}
	formData.conditionLabel = widget.NewLabel("État: -")
	if existingGame != nil && existingGame.Condition != nil {
		formData.conditionLabel.SetText(fmt.Sprintf("État: %d", *existingGame.Condition))
	}
	formData.conditionSlider.OnChanged = func(value float64) {
		formData.conditionLabel.SetText(fmt.Sprintf("État: %d", int(value)))
	}

	// ========== Purchase Info ==========
	formData.purchaseDateEntry = widget.NewEntry()
	formData.purchaseDateEntry.SetPlaceHolder("AAAA-MM-JJ")
	if existingGame != nil && existingGame.PurchaseDate != nil {
		formData.purchaseDateEntry.SetText(existingGame.PurchaseDate.Format("2006-01-02"))
	}

	formData.purchasePriceEntry = widget.NewEntry()
	formData.purchasePriceEntry.SetPlaceHolder("Prix d'achat")
	if existingGame != nil && existingGame.PurchasePrice != nil {
		formData.purchasePriceEntry.SetText(fmt.Sprintf("%.2f", *existingGame.PurchasePrice))
	}

	// ========== Notes ==========
	formData.notesEntry = widget.NewMultiLineEntry()
	formData.notesEntry.SetPlaceHolder("Notes")
	formData.notesEntry.SetMinRowsVisible(3)
	if existingGame != nil && existingGame.Notes != nil {
		formData.notesEntry.SetText(*existingGame.Notes)
	}

	// ========== Many-to-Many: Developers ==========
	formData.developersList = widget.NewLabel("-")
	developerOptions := []string{}
	developerNameToID := make(map[string]int)
	for _, d := range developers {
		developerOptions = append(developerOptions, d.Name)
		developerNameToID[d.Name] = d.DeveloperID
	}

	// Pre-populate if editing
	if existingGame != nil {
		formData.selectedDevelopers = existingGame.Developers
		for _, name := range existingGame.Developers {
			if id, ok := developerNameToID[name]; ok {
				formData.selectedDeveloperIDs = append(formData.selectedDeveloperIDs, id)
			}
		}
		if len(formData.selectedDevelopers) > 0 {
			formData.developersList.SetText(formatList(formData.selectedDevelopers))
		}
	}

	developerSelect := widget.NewSelect(developerOptions, nil)
	developerSelect.PlaceHolder = "Sélectionner développeur(s)"

	addDeveloperBtn := widget.NewButton("Ajouter", func() {
		if developerSelect.Selected != "" && !contains(formData.selectedDevelopers, developerSelect.Selected) {
			formData.selectedDevelopers = append(formData.selectedDevelopers, developerSelect.Selected)
			formData.selectedDeveloperIDs = append(formData.selectedDeveloperIDs, developerNameToID[developerSelect.Selected])
			formData.developersList.SetText(formatList(formData.selectedDevelopers))
			developerSelect.SetSelected("")
		}
	})

	clearDevelopersBtn := widget.NewButton("Effacer", func() {
		formData.selectedDevelopers = []string{}
		formData.selectedDeveloperIDs = []int{}
		formData.developersList.SetText("-")
	})

	newDeveloperBtn := widget.NewButton("+", func() {
		showAddDeveloperDialog(w, conn, func() {
			developers, _ = getDevelopers(conn)
			developerOptions = []string{}
			developerNameToID = make(map[string]int)
			for _, d := range developers {
				developerOptions = append(developerOptions, d.Name)
				developerNameToID[d.Name] = d.DeveloperID
			}
			developerSelect.Options = developerOptions
			developerSelect.Refresh()
		})
	})

	// ========== Many-to-Many: Composers ==========
	formData.composersList = widget.NewLabel("-")
	composerOptions := []string{}
	composerNameToID := make(map[string]int)
	for _, c := range composers {
		composerOptions = append(composerOptions, c.Name)
		composerNameToID[c.Name] = c.ComposerID
	}

	if existingGame != nil {
		formData.selectedComposers = existingGame.Composers
		for _, name := range existingGame.Composers {
			if id, ok := composerNameToID[name]; ok {
				formData.selectedComposerIDs = append(formData.selectedComposerIDs, id)
			}
		}
		if len(formData.selectedComposers) > 0 {
			formData.composersList.SetText(formatList(formData.selectedComposers))
		}
	}

	composerSelect := widget.NewSelect(composerOptions, nil)
	composerSelect.PlaceHolder = "Sélectionner compositeur(s)"

	addComposerBtn := widget.NewButton("Ajouter", func() {
		if composerSelect.Selected != "" && !contains(formData.selectedComposers, composerSelect.Selected) {
			formData.selectedComposers = append(formData.selectedComposers, composerSelect.Selected)
			formData.selectedComposerIDs = append(formData.selectedComposerIDs, composerNameToID[composerSelect.Selected])
			formData.composersList.SetText(formatList(formData.selectedComposers))
			composerSelect.SetSelected("")
		}
	})

	clearComposersBtn := widget.NewButton("Effacer", func() {
		formData.selectedComposers = []string{}
		formData.selectedComposerIDs = []int{}
		formData.composersList.SetText(".")
	})

	newComposerBtn := widget.NewButton("+", func() {
		showAddComposerDialog(w, conn, func() {
			composers, _ = getComposers(conn)
			composerOptions = []string{}
			composerNameToID = make(map[string]int)
			for _, c := range composers {
				composerOptions = append(composerOptions, c.Name)
				composerNameToID[c.Name] = c.ComposerID
			}
			composerSelect.Options = composerOptions
			composerSelect.Refresh()
		})
	})

	// ========== Many-to-Many: Publishers ==========
	formData.publishersList = widget.NewLabel("-")
	publisherOptions := []string{}
	publisherNameToID := make(map[string]int)
	for _, p := range publishers {
		publisherOptions = append(publisherOptions, p.Name)
		publisherNameToID[p.Name] = p.PublisherID
	}

	if existingGame != nil {
		formData.selectedPublishers = existingGame.Publishers
		for _, name := range existingGame.Publishers {
			if id, ok := publisherNameToID[name]; ok {
				formData.selectedPublisherIDs = append(formData.selectedPublisherIDs, id)
			}
		}
		if len(formData.selectedPublishers) > 0 {
			formData.publishersList.SetText(formatList(formData.selectedPublishers))
		}
	}

	publisherSelect := widget.NewSelect(publisherOptions, nil)
	publisherSelect.PlaceHolder = "Sélectionner distributeur"

	addPublisherBtn := widget.NewButton("Ajouter", func() {
		if publisherSelect.Selected != "" && !contains(formData.selectedPublishers, publisherSelect.Selected) {
			formData.selectedPublishers = append(formData.selectedPublishers, publisherSelect.Selected)
			formData.selectedPublisherIDs = append(formData.selectedPublisherIDs, publisherNameToID[publisherSelect.Selected])
			formData.publishersList.SetText(formatList(formData.selectedPublishers))
			publisherSelect.SetSelected("")
		}
	})

	clearPublishersBtn := widget.NewButton("Effacer", func() {
		formData.selectedPublishers = []string{}
		formData.selectedPublisherIDs = []int{}
		formData.publishersList.SetText("-")
	})

	newPublisherBtn := widget.NewButton("+", func() {
		showAddPublisherDialog(w, conn, func() {
			publishers, _ = getPublishers(conn)
			publisherOptions = []string{}
			publisherNameToID = make(map[string]int)
			for _, p := range publishers {
				publisherOptions = append(publisherOptions, p.Name)
				publisherNameToID[p.Name] = p.PublisherID
			}
			publisherSelect.Options = publisherOptions
			publisherSelect.Refresh()
		})
	})

	// ========== Many-to-Many: Producers ==========
	formData.producersList = widget.NewLabel("-")
	producerOptions := []string{}
	producerNameToID := make(map[string]int)
	for _, p := range producers {
		producerOptions = append(producerOptions, p.Name)
		producerNameToID[p.Name] = p.ProducerID
	}

	if existingGame != nil {
		formData.selectedProducers = existingGame.Producers
		for _, name := range existingGame.Producers {
			if id, ok := producerNameToID[name]; ok {
				formData.selectedProducerIDs = append(formData.selectedProducerIDs, id)
			}
		}
		if len(formData.selectedProducers) > 0 {
			formData.producersList.SetText(formatList(formData.selectedProducers))
		}
	}

	producerSelect := widget.NewSelect(producerOptions, nil)
	producerSelect.PlaceHolder = "Sélectionner producteur"

	addProducerBtn := widget.NewButton("Ajouter", func() {
		if producerSelect.Selected != "" && !contains(formData.selectedProducers, producerSelect.Selected) {
			formData.selectedProducers = append(formData.selectedProducers, producerSelect.Selected)
			formData.selectedProducerIDs = append(formData.selectedProducerIDs, producerNameToID[producerSelect.Selected])
			formData.producersList.SetText(formatList(formData.selectedProducers))
			producerSelect.SetSelected("")
		}
	})

	clearProducersBtn := widget.NewButton("Effacer", func() {
		formData.selectedProducers = []string{}
		formData.selectedProducerIDs = []int{}
		formData.producersList.SetText("-")
	})

	newProducerBtn := widget.NewButton("+", func() {
		showAddProducerDialog(w, conn, func() {
			producers, _ = getProducers(conn)
			producerOptions = []string{}
			producerNameToID = make(map[string]int)
			for _, p := range producers {
				producerOptions = append(producerOptions, p.Name)
				producerNameToID[p.Name] = p.ProducerID
			}
			producerSelect.Options = producerOptions
			producerSelect.Refresh()
		})
	})

	// ========== Build Form Layout ==========
	formData.form = container.NewVBox(
		widget.NewLabel("Titre *"),
		formData.titleEntry,
		widget.NewLabel("Plateforme *"),
		formData.consoleSelect,
		widget.NewLabel("Genre"),
		formData.genreSelect,

		widget.NewSeparator(),
		widget.NewLabel("Date de sortie"),
		widget.NewLabel("Europe:"),
		formData.euReleaseDateEntry,
		widget.NewLabel("USA:"),
		formData.usReleaseDateEntry,
		widget.NewLabel("Japon:"),
		formData.jpReleaseDateEntry,
		


		widget.NewSeparator(),
		widget.NewLabel("Classifications"),
		formData.euRatingSelect,
		formData.usRatingSelect,
		formData.jpRatingSelect,

		widget.NewSeparator(),
		widget.NewLabel("Total des copies vendues"),
		formData.unitsSoldEntry,

		widget.NewSeparator(),
		widget.NewLabel("Dévelopeur(s)"),
		container.NewBorder(nil, nil, nil, container.NewHBox(addDeveloperBtn, newDeveloperBtn), developerSelect),
		formData.developersList,
		clearDevelopersBtn,

		widget.NewSeparator(),
		widget.NewLabel("Compositeur(s)"),
		container.NewBorder(nil, nil, nil, container.NewHBox(addComposerBtn, newComposerBtn), composerSelect),
		formData.composersList,
		clearComposersBtn,

		widget.NewSeparator(),
		widget.NewLabel("Distributeur(s)"),
		container.NewBorder(nil, nil, nil, container.NewHBox(addPublisherBtn, newPublisherBtn), publisherSelect),
		formData.publishersList,
		clearPublishersBtn,

		widget.NewSeparator(),
		widget.NewLabel("Producteur(s)"),
		container.NewBorder(nil, nil, nil, container.NewHBox(addProducerBtn, newProducerBtn), producerSelect),
		formData.producersList,
		clearProducersBtn,

		widget.NewSeparator(),
		widget.NewLabel("Informations de collection"),
		formData.ownedCheck,
		formData.boxOwnedCheck,
		formData.collectorCheck,
		formData.conditionLabel,
		formData.conditionSlider,

		widget.NewSeparator(),
		widget.NewLabel("Informations d'achat"),
		widget.NewLabel("Date d'achat:"),
		formData.purchaseDateEntry,
		widget.NewLabel("Prix d'achat:"),
		formData.purchasePriceEntry,

		widget.NewSeparator(),
		widget.NewLabel("Notes"),
		formData.notesEntry,
	)

	return formData
}

// showAddGameDialog shows the dialog to add a new game
func showAddGameDialog(w fyne.Window, conn *pgx.Conn, onSuccess func()) {
	formData := buildGameForm(w, conn, nil) // nil = no existing game

	saveBtn := widget.NewButton("Enregistrer", func() {
		gameID, err := saveGame(conn, formData, 0) // 0 = new game
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Save many-to-many relationships
		saveManyToManyRelationships(conn, gameID, formData)

		dialog.ShowInformation("Enregistré", "Jeu ajouté à la base de données", w)
		if onSuccess != nil {
			onSuccess()
		}
	})

	saveBtn.Importance = widget.HighImportance

	formWithSave := container.NewBorder(
		nil, saveBtn, nil, nil,
		container.NewScroll(formData.form),
	)

	d := dialog.NewCustom("Ajouter", "Annuler", formWithSave, w)
	d.Resize(fyne.NewSize(600, 700))
	d.Show()
}

// showEditGameDialog shows the dialog to edit an existing game
func showEditGameDialog(w fyne.Window, conn *pgx.Conn, gameID int, onSuccess func()) {
	// Fetch existing game
	existingGame, err := getGameByID(conn, gameID)
	if err != nil {
		dialog.ShowError(fmt.Errorf("échec de chargement du jeu: %w", err), w)
		return
	}

	formData := buildGameForm(w, conn, existingGame) // Pre-populate with existing data

	saveBtn := widget.NewButton("Enregistrer", func() {
		_, err := saveGame(conn, formData, gameID) // gameID != 0 = update
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Delete old relationships and save new ones
		conn.Exec(context.Background(), "DELETE FROM game_developers WHERE game_id = $1", gameID)
		conn.Exec(context.Background(), "DELETE FROM game_composers WHERE game_id = $1", gameID)
		conn.Exec(context.Background(), "DELETE FROM game_publishers WHERE game_id = $1", gameID)
		conn.Exec(context.Background(), "DELETE FROM game_producers WHERE game_id = $1", gameID)

		saveManyToManyRelationships(conn, gameID, formData)

		dialog.ShowInformation("Mise à jour", "Jeu mis à jour dans la base de données.", w)
		if onSuccess != nil {
			onSuccess()
		}
	})

	saveBtn.Importance = widget.HighImportance

	formWithSave := container.NewBorder(
		nil, saveBtn, nil, nil,
		container.NewScroll(formData.form),
	)

	d := dialog.NewCustom("Enregistrer", "Annuler", formWithSave, w)
	d.Resize(fyne.NewSize(600, 700))
	d.Show()
}

// saveGame saves or updates a game (INSERT if gameID=0, UPDATE if gameID>0)
func saveGame(conn *pgx.Conn, formData *gameFormData, gameID int) (int, error) {
	// Validate
	if formData.titleEntry.Text == "" {
		return 0, fmt.Errorf("titre requis")
	}
	if formData.consoleSelect.Selected == "" {
		return 0, fmt.Errorf("plateforme requise")
	}

	// Prepare data
	var consoleID *int
	if formData.consoleSelect.Selected != "" {
		id := formData.consoleMap[formData.consoleSelect.Selected]
		consoleID = &id
	}

	var genreID *int
	if formData.genreSelect.Selected != "" {
		id := formData.genreMap[formData.genreSelect.Selected]
		genreID = &id
	}

	// Parse dates
	var jpReleaseDate, usReleaseDate, euReleaseDate, purchaseDate *string
	if formData.jpReleaseDateEntry.Text != "" {
		jpReleaseDate = &formData.jpReleaseDateEntry.Text
	}
	if formData.usReleaseDateEntry.Text != "" {
		usReleaseDate = &formData.usReleaseDateEntry.Text
	}
	if formData.euReleaseDateEntry.Text != "" {
		euReleaseDate = &formData.euReleaseDateEntry.Text
	}
	if formData.purchaseDateEntry.Text != "" {
		purchaseDate = &formData.purchaseDateEntry.Text
	}

	// Parse ratings
	var jpRatingID, usRatingID, euRatingID *int
	if formData.jpRatingSelect.Selected != "" {
		id := formData.ratingMap[formData.jpRatingSelect.Selected]
		jpRatingID = &id
	}
	if formData.usRatingSelect.Selected != "" {
		id := formData.ratingMap[formData.usRatingSelect.Selected]
		usRatingID = &id
	}
	if formData.euRatingSelect.Selected != "" {
		id := formData.ratingMap[formData.euRatingSelect.Selected]
		euRatingID = &id
	}

	// Parse units sold
	var unitsSold *int
	if formData.unitsSoldEntry.Text != "" {
		var units int
		fmt.Sscanf(formData.unitsSoldEntry.Text, "%d", &units)
		unitsSold = &units
	}

	// Parse purchase price
	var purchasePrice *float64
	if formData.purchasePriceEntry.Text != "" {
		var price float64
		fmt.Sscanf(formData.purchasePriceEntry.Text, "%f", &price)
		purchasePrice = &price
	}

	// Parse condition
	var condition *int
	if formData.conditionSlider.Value > 0 {
		c := int(formData.conditionSlider.Value)
		condition = &c
	}

	// Parse notes
	var notes *string
	if formData.notesEntry.Text != "" {
		notes = &formData.notesEntry.Text
	}

	// INSERT or UPDATE
	if gameID == 0 {
		// INSERT new game
		query := `
			INSERT INTO games (
				title, console_id, genre_id, 
				jp_release_date, us_release_date, eu_release_date,
				jp_rating_id, us_rating_id, eu_rating_id,
				units_sold, owned, box_owned, collector, condition,
				purchase_date, purchase_price, notes
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
			RETURNING game_id
		`

		err := conn.QueryRow(context.Background(), query,
			formData.titleEntry.Text, consoleID, genreID,
			jpReleaseDate, usReleaseDate, euReleaseDate,
			jpRatingID, usRatingID, euRatingID,
			unitsSold, formData.ownedCheck.Checked, formData.boxOwnedCheck.Checked,
			formData.collectorCheck.Checked, condition, purchaseDate, purchasePrice, notes,
		).Scan(&gameID)

		if err != nil {
			return 0, fmt.Errorf("échec d'ajout du jeu: %w", err)
		}
		return gameID, nil
	} else {
		// UPDATE existing game
		query := `
			UPDATE games SET
				title = $1, console_id = $2, genre_id = $3,
				jp_release_date = $4, us_release_date = $5, eu_release_date = $6,
				jp_rating_id = $7, us_rating_id = $8, eu_rating_id = $9,
				units_sold = $10, owned = $11, box_owned = $12, collector = $13,
				condition = $14, purchase_date = $15, purchase_price = $16, notes = $17
			WHERE game_id = $18
		`

		_, err := conn.Exec(context.Background(), query,
			formData.titleEntry.Text, consoleID, genreID,
			jpReleaseDate, usReleaseDate, euReleaseDate,
			jpRatingID, usRatingID, euRatingID,
			unitsSold, formData.ownedCheck.Checked, formData.boxOwnedCheck.Checked,
			formData.collectorCheck.Checked, condition, purchaseDate, purchasePrice, notes,
			gameID,
		)

		if err != nil {
			return 0, fmt.Errorf("échec de mise à jour du jeu: %w", err)
		}
		return gameID, nil
	}
}

// saveManyToManyRelationships saves all many-to-many relationships for a game
func saveManyToManyRelationships(conn *pgx.Conn, gameID int, formData *gameFormData) {
	// Developers
	for _, devID := range formData.selectedDeveloperIDs {
		conn.Exec(context.Background(),
			"INSERT INTO game_developers (game_id, developer_id) VALUES ($1, $2)",
			gameID, devID)
	}

	// Composers
	for _, compID := range formData.selectedComposerIDs {
		conn.Exec(context.Background(),
			"INSERT INTO game_composers (game_id, composer_id) VALUES ($1, $2)",
			gameID, compID)
	}

	// Publishers
	for _, pubID := range formData.selectedPublisherIDs {
		conn.Exec(context.Background(),
			"INSERT INTO game_publishers (game_id, publisher_id) VALUES ($1, $2)",
			gameID, pubID)
	}

	// Producers
	for _, prodID := range formData.selectedProducerIDs {
		conn.Exec(context.Background(),
			"INSERT INTO game_producers (game_id, producer_id) VALUES ($1, $2)",
			gameID, prodID)
	}
}

// ========== Helper Functions ==========

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func formatList(items []string) string {
	if len(items) == 0 {
		return "Aucune sélection"
	}
	result := ""
	for i, item := range items {
		if i > 0 {
			result += ", "
		}
		result += item
	}
	return result
}

// ========== Add New Lookup Entry Dialogs ==========

func showAddDeveloperDialog(w fyne.Window, conn *pgx.Conn, onSuccess func()) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nom du développeur")

	form := container.NewVBox(
		widget.NewLabel("Ajouter nouveau développeur"),
		nameEntry,
	)

	d := dialog.NewCustomConfirm("Ajouter", "Enregistrer", "Annuler", form, func(save bool) {
		if save && nameEntry.Text != "" {
			_, err := conn.Exec(context.Background(),
				"INSERT INTO developers (name) VALUES ($1)",
				nameEntry.Text)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation("Ajouté", "Développeur ajouté à la liste des développeurs.", w)
			if onSuccess != nil {
				onSuccess()
			}
		}
	}, w)

	d.Show()
}

func showAddComposerDialog(w fyne.Window, conn *pgx.Conn, onSuccess func()) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nom du compositeur")

	form := container.NewVBox(
		widget.NewLabel("Ajouter nouveau compositeur"),
		nameEntry,
	)

	d := dialog.NewCustomConfirm("Ajouter", "Enregistrer", "Annuler", form, func(save bool) {
		if save && nameEntry.Text != "" {
			_, err := conn.Exec(context.Background(),
				"INSERT INTO composers (name) VALUES ($1)",
				nameEntry.Text)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation("Ajouté", "Compositeur ajouté à la liste des compositeurs.", w)
			if onSuccess != nil {
				onSuccess()
			}
		}
	}, w)

	d.Show()
}

func showAddPublisherDialog(w fyne.Window, conn *pgx.Conn, onSuccess func()) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nom de l'éditeur")

	form := container.NewVBox(
		widget.NewLabel("Ajouter nouvel éditeur"),
		nameEntry,
	)

	d := dialog.NewCustomConfirm("Ajouter", "Enregistrer", "Annuler", form, func(save bool) {
		if save && nameEntry.Text != "" {
			_, err := conn.Exec(context.Background(),
				"INSERT INTO publishers (name) VALUES ($1)",
				nameEntry.Text)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation("Ajouté", "Editeur ajouté à la liste des éditeurs.", w)
			if onSuccess != nil {
				onSuccess()
			}
		}
	}, w)

	d.Show()
}

func showAddProducerDialog(w fyne.Window, conn *pgx.Conn, onSuccess func()) {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nom du producteur")

	form := container.NewVBox(
		widget.NewLabel("Ajouter nouveau producteur"),
		nameEntry,
	)

	d := dialog.NewCustomConfirm("Ajouter", "Enregistrer", "Annuler", form, func(save bool) {
		if save && nameEntry.Text != "" {
			_, err := conn.Exec(context.Background(),
				"INSERT INTO producers (name) VALUES ($1)",
				nameEntry.Text)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation("Ajouté", "Producteur ajouté à la liste des producteurs.", w)
			if onSuccess != nil {
				onSuccess()
			}
		}
	}, w)

	d.Show()
}

// ========== ACCESSORIES CRUD ==========

// accessoryFormData holds all the form fields for accessories
type accessoryFormData struct {
	// Form widgets
	nameEntry          *widget.Entry
	colorEntry         *widget.Entry
	typeSelect         *widget.Select
	manufacturerSelect *widget.Select
	quantityEntry      *widget.Entry
	ownedCheck         *widget.Check
	conditionSlider    *widget.Slider
	conditionLabel     *widget.Label
	purchaseDateEntry  *widget.Entry
	purchasePriceEntry *widget.Entry
	notesEntry         *widget.Entry

	// Many-to-many consoles
	selectedConsoles   []string
	selectedConsoleIDs []int
	consolesList       *widget.Label

	// Lookup maps
	typeMap         map[string]int
	manufacturerMap map[string]int

	// The complete form container
	form *fyne.Container
}

// buildAccessoryForm creates the accessory form, optionally pre-populated
func buildAccessoryForm(w fyne.Window, conn *pgx.Conn, existingAccessory *Accessory) *accessoryFormData {
	formData := &accessoryFormData{
		typeMap:         make(map[string]int),
		manufacturerMap: make(map[string]int),
	}

	// Fetch lookup data
	types, _ := getAccessoryTypes(conn)
	manufacturers, _ := getManufacturers(conn)
	consoles, _ := getConsoles(conn)

	// ========== Basic Info ==========
	formData.nameEntry = widget.NewEntry()
	formData.nameEntry.SetPlaceHolder("Nom de l'accessoire (requis)")
	if existingAccessory != nil {
		formData.nameEntry.SetText(existingAccessory.Name)
	}

	formData.colorEntry = widget.NewEntry()
	formData.colorEntry.SetPlaceHolder("Couleur")
	if existingAccessory != nil && existingAccessory.Color != nil {
		formData.colorEntry.SetText(*existingAccessory.Color)
	}

	// Type dropdown (required)
	typeOptions := []string{""}
	var selectedTypeName string
	for _, t := range types {
		typeOptions = append(typeOptions, t.Name)
		formData.typeMap[t.Name] = t.TypeID
		if existingAccessory != nil && existingAccessory.TypeID != nil && t.TypeID == *existingAccessory.TypeID {
			selectedTypeName = t.Name
		}
	}
	formData.typeSelect = widget.NewSelect(typeOptions, nil)
	formData.typeSelect.PlaceHolder = "Type d'accessoire (requis)"
	if selectedTypeName != "" {
		formData.typeSelect.SetSelected(selectedTypeName)
	}

	// Manufacturer dropdown (optional)
	manufacturerOptions := []string{""}
	var selectedManufacturerName string
	for _, m := range manufacturers {
		manufacturerOptions = append(manufacturerOptions, m.Name)
		formData.manufacturerMap[m.Name] = m.ManufacturerID
		if existingAccessory != nil && existingAccessory.ManufacturerID != nil && m.ManufacturerID == *existingAccessory.ManufacturerID {
			selectedManufacturerName = m.Name
		}
	}
	formData.manufacturerSelect = widget.NewSelect(manufacturerOptions, nil)
	formData.manufacturerSelect.PlaceHolder = "Fabricant"
	if selectedManufacturerName != "" {
		formData.manufacturerSelect.SetSelected(selectedManufacturerName)
	}

	// ========== Consoles (Many-to-Many) ==========
	formData.consolesList = widget.NewLabel("-")
	consoleOptions := []string{}
	consoleNameToID := make(map[string]int)
	for _, c := range consoles {
		consoleOptions = append(consoleOptions, c.Name)
		consoleNameToID[c.Name] = c.ConsoleID
	}

	// Pre-populate if editing
	if existingAccessory != nil {
		formData.selectedConsoles = existingAccessory.Consoles
		for _, name := range existingAccessory.Consoles {
			if id, ok := consoleNameToID[name]; ok {
				formData.selectedConsoleIDs = append(formData.selectedConsoleIDs, id)
			}
		}
		if len(formData.selectedConsoles) > 0 {
			formData.consolesList.SetText(formatList(formData.selectedConsoles))
		}
	}

	consoleSelect := widget.NewSelect(consoleOptions, nil)
	consoleSelect.PlaceHolder = "Plateforme"

	addConsoleBtn := widget.NewButton("Ajouter", func() {
		if consoleSelect.Selected != "" && !contains(formData.selectedConsoles, consoleSelect.Selected) {
			formData.selectedConsoles = append(formData.selectedConsoles, consoleSelect.Selected)
			formData.selectedConsoleIDs = append(formData.selectedConsoleIDs, consoleNameToID[consoleSelect.Selected])
			formData.consolesList.SetText(formatList(formData.selectedConsoles))
			consoleSelect.SetSelected("")
		}
	})

	clearConsolesBtn := widget.NewButton("Effacer", func() {
		formData.selectedConsoles = []string{}
		formData.selectedConsoleIDs = []int{}
		formData.consolesList.SetText("-")
	})

	// ========== Quantity ==========
	formData.quantityEntry = widget.NewEntry()
	formData.quantityEntry.SetPlaceHolder("Quantité")
	if existingAccessory != nil {
		formData.quantityEntry.SetText(fmt.Sprintf("%d", existingAccessory.Quantity))
	} else {
		formData.quantityEntry.SetText("1") // Default to 1
	}

	// ========== Collection Info ==========
	formData.ownedCheck = widget.NewCheck("Possédé", nil)
	if existingAccessory != nil {
		formData.ownedCheck.Checked = existingAccessory.Owned
	} else {
		formData.ownedCheck.Checked = true
	}

	// Condition
	formData.conditionSlider = widget.NewSlider(1, 5)
	formData.conditionSlider.Step = 1
	if existingAccessory != nil && existingAccessory.Condition != nil {
		formData.conditionSlider.Value = float64(*existingAccessory.Condition)
	}
	formData.conditionLabel = widget.NewLabel("État: -")
	if existingAccessory != nil && existingAccessory.Condition != nil {
		formData.conditionLabel.SetText(fmt.Sprintf("État: %d", *existingAccessory.Condition))
	}
	formData.conditionSlider.OnChanged = func(value float64) {
		formData.conditionLabel.SetText(fmt.Sprintf("État: %d", int(value)))
	}

	// ========== Purchase Info ==========
	formData.purchaseDateEntry = widget.NewEntry()
	formData.purchaseDateEntry.SetPlaceHolder("AAAA-MM-JJ")
	if existingAccessory != nil && existingAccessory.PurchaseDate != nil {
		formData.purchaseDateEntry.SetText(existingAccessory.PurchaseDate.Format("2006-01-02"))
	}

	formData.purchasePriceEntry = widget.NewEntry()
	formData.purchasePriceEntry.SetPlaceHolder("Prix d'achat")
	if existingAccessory != nil && existingAccessory.PurchasePrice != nil {
		formData.purchasePriceEntry.SetText(fmt.Sprintf("%.2f", *existingAccessory.PurchasePrice))
	}

	// ========== Notes ==========
	formData.notesEntry = widget.NewMultiLineEntry()
	formData.notesEntry.SetPlaceHolder("Notes")
	formData.notesEntry.SetMinRowsVisible(3)
	if existingAccessory != nil && existingAccessory.Notes != nil {
		formData.notesEntry.SetText(*existingAccessory.Notes)
	}

	// ========== Build Form Layout ==========
	formData.form = container.NewVBox(
		widget.NewLabel("Nom *"),
		formData.nameEntry,

		widget.NewLabel("Couleur"),
		formData.colorEntry,

		widget.NewLabel("Type *"),
		formData.typeSelect,

		widget.NewLabel("Fabricant"),
		formData.manufacturerSelect,

		widget.NewSeparator(),
		widget.NewLabel("Plateforme(s)"),
		container.NewBorder(nil, nil, nil, addConsoleBtn, consoleSelect),
		formData.consolesList,
		clearConsolesBtn,

		widget.NewSeparator(),
		widget.NewLabel("Quantité"),
		formData.quantityEntry,

		widget.NewSeparator(),
		widget.NewLabel("Informations de collection"),
		formData.ownedCheck,
		formData.conditionLabel,
		formData.conditionSlider,

		widget.NewSeparator(),
		widget.NewLabel("Informations d'achat"),
		widget.NewLabel("Date d'achat:"),
		formData.purchaseDateEntry,
		widget.NewLabel("Prix d'achat:"),
		formData.purchasePriceEntry,

		widget.NewSeparator(),
		widget.NewLabel("Notes"),
		formData.notesEntry,
	)

	return formData
}

// showAddAccessoryDialog shows the dialog to add a new accessory
func showAddAccessoryDialog(w fyne.Window, conn *pgx.Conn, onSuccess func()) {
	formData := buildAccessoryForm(w, conn, nil)

	saveBtn := widget.NewButton("Enregistrer", func() {
		accessoryID, err := saveAccessory(conn, formData, 0)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Save console relationships
		for _, consoleID := range formData.selectedConsoleIDs {
			conn.Exec(context.Background(),
				"INSERT INTO accessory_consoles (accessory_id, console_id) VALUES ($1, $2)",
				accessoryID, consoleID)
		}

		dialog.ShowInformation("Ajouté", "Accessoire ajouté à la base de données.", w)
		if onSuccess != nil {
			onSuccess()
		}
	})

	saveBtn.Importance = widget.HighImportance

	formWithSave := container.NewBorder(
		nil, saveBtn, nil, nil,
		container.NewScroll(formData.form),
	)

	d := dialog.NewCustom("Ajouter", "Annuler", formWithSave, w)
	d.Resize(fyne.NewSize(600, 700))
	d.Show()
}

// showEditAccessoryDialog shows the dialog to edit an existing accessory
func showEditAccessoryDialog(w fyne.Window, conn *pgx.Conn, accessoryID int, onSuccess func()) {
	existingAccessory, err := getAccessoryByID(conn, accessoryID)
	if err != nil {
		dialog.ShowError(fmt.Errorf("échec du chargement de l'accessoire: %w", err), w)
		return
	}

	formData := buildAccessoryForm(w, conn, existingAccessory)

	saveBtn := widget.NewButton("Enregistrer", func() {
		_, err := saveAccessory(conn, formData, accessoryID)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Delete old console relationships and save new ones
		conn.Exec(context.Background(), "DELETE FROM accessory_consoles WHERE accessory_id = $1", accessoryID)
		for _, consoleID := range formData.selectedConsoleIDs {
			conn.Exec(context.Background(),
				"INSERT INTO accessory_consoles (accessory_id, console_id) VALUES ($1, $2)",
				accessoryID, consoleID)
		}

		dialog.ShowInformation("Enregistré", "Accessoire mis à jour dans la base de données.", w)
		if onSuccess != nil {
			onSuccess()
		}
	})

	saveBtn.Importance = widget.HighImportance

	formWithSave := container.NewBorder(
		nil, saveBtn, nil, nil,
		container.NewScroll(formData.form),
	)

	d := dialog.NewCustom("Éditer", "Annuler", formWithSave, w)
	d.Resize(fyne.NewSize(600, 700))
	d.Show()
}

// saveAccessory saves or updates an accessory
func saveAccessory(conn *pgx.Conn, formData *accessoryFormData, accessoryID int) (int, error) {
	// Validate
	if formData.nameEntry.Text == "" {
		return 0, fmt.Errorf("nom requis")
	}
	if formData.typeSelect.Selected == "" {
		return 0, fmt.Errorf("type requis")
	}

	// Prepare data
	var typeID *int
	if formData.typeSelect.Selected != "" {
		id := formData.typeMap[formData.typeSelect.Selected]
		typeID = &id
	}

	var manufacturerID *int
	if formData.manufacturerSelect.Selected != "" {
		id := formData.manufacturerMap[formData.manufacturerSelect.Selected]
		manufacturerID = &id
	}

	var color *string
	if formData.colorEntry.Text != "" {
		color = &formData.colorEntry.Text
	}

	// Parse quantity
	var quantity int = 1
	if formData.quantityEntry.Text != "" {
		fmt.Sscanf(formData.quantityEntry.Text, "%d", &quantity)
	}

	// Parse condition
	var condition *int
	if formData.conditionSlider.Value > 0 {
		c := int(formData.conditionSlider.Value)
		condition = &c
	}

	// Parse purchase date
	var purchaseDate *string
	if formData.purchaseDateEntry.Text != "" {
		purchaseDate = &formData.purchaseDateEntry.Text
	}

	// Parse purchase price
	var purchasePrice *float64
	if formData.purchasePriceEntry.Text != "" {
		var price float64
		fmt.Sscanf(formData.purchasePriceEntry.Text, "%f", &price)
		purchasePrice = &price
	}

	// Parse notes
	var notes *string
	if formData.notesEntry.Text != "" {
		notes = &formData.notesEntry.Text
	}

	// INSERT or UPDATE
	if accessoryID == 0 {
		// INSERT
		query := `
			INSERT INTO accessories (
				name, color, type_id, manufacturer_id, quantity,
				condition, owned, purchase_date, purchase_price, notes
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING accessory_id
		`

		err := conn.QueryRow(context.Background(), query,
			formData.nameEntry.Text, color, typeID, manufacturerID, quantity,
			condition, formData.ownedCheck.Checked, purchaseDate, purchasePrice, notes,
		).Scan(&accessoryID)

		if err != nil {
			return 0, fmt.Errorf("échec de l'ajout de l'accessoire: %w", err)
		}
		return accessoryID, nil
	} else {
		// UPDATE
		query := `
			UPDATE accessories SET
				name = $1, color = $2, type_id = $3, manufacturer_id = $4,
				quantity = $5, condition = $6, owned = $7,
				purchase_date = $8, purchase_price = $9, notes = $10
			WHERE accessory_id = $11
		`

		_, err := conn.Exec(context.Background(), query,
			formData.nameEntry.Text, color, typeID, manufacturerID, quantity,
			condition, formData.ownedCheck.Checked, purchaseDate, purchasePrice, notes,
			accessoryID,
		)

		if err != nil {
			return 0, fmt.Errorf("échec de la mise à jour de l'accessoire: %w", err)
		}
		return accessoryID, nil
	}
}

// ========== CONSOLES CRUD ==========

// consoleFormData holds all the form fields for consoles
type consoleFormData struct {
	// Form widgets
	nameEntry          *widget.Entry
	generationEntry    *widget.Entry
	jpReleaseDateEntry *widget.Entry
	usReleaseDateEntry *widget.Entry
	euReleaseDateEntry *widget.Entry
	discontinuedEntry  *widget.Entry
	priceJPYEntry      *widget.Entry
	priceUSDEntry      *widget.Entry
	controllersEntry   *widget.Entry
	cpuEntry           *widget.Entry
	gpuEntry           *widget.Entry
	memoryEntry        *widget.Entry
	audioEntry         *widget.Entry
	unitsSoldEntry     *widget.Entry
	topGameEntry       *widget.Entry
	predecessorEntry   *widget.Entry
	successorEntry     *widget.Entry
	ownedCheck         *widget.Check
	conditionSlider    *widget.Slider
	conditionLabel     *widget.Label
	notesEntry         *widget.Entry
	typeSelect         *widget.Select
	manufacturerSelect *widget.Select

	// Lookup maps
	typeMap         map[string]int
	manufacturerMap map[string]int

	// The complete form container
	form *fyne.Container
}

// buildConsoleForm creates the console form, optionally pre-populated
func buildConsoleForm(w fyne.Window, conn *pgx.Conn, existingConsole *Console) *consoleFormData {
	formData := &consoleFormData{
		typeMap:         make(map[string]int),
		manufacturerMap: make(map[string]int),
	}

	// Fetch lookup data
	types, _ := getConsoleTypes(conn)
	manufacturers, _ := getManufacturers(conn)

	// ========== Basic Info ==========
	formData.nameEntry = widget.NewEntry()
	formData.nameEntry.SetPlaceHolder("Plateforme (requis)")
	if existingConsole != nil {
		formData.nameEntry.SetText(existingConsole.Name)
	}

	// Type dropdown (required)
	typeOptions := []string{""}
	var selectedTypeName string
	for _, t := range types {
		typeOptions = append(typeOptions, t.Name)
		formData.typeMap[t.Name] = t.TypeID
		if existingConsole != nil && existingConsole.TypeID != nil && t.TypeID == *existingConsole.TypeID {
			selectedTypeName = t.Name
		}
	}
	formData.typeSelect = widget.NewSelect(typeOptions, nil)
	formData.typeSelect.PlaceHolder = "Type (requis)"
	if selectedTypeName != "" {
		formData.typeSelect.SetSelected(selectedTypeName)
	}

	// Manufacturer dropdown (required)
	manufacturerOptions := []string{""}
	var selectedManufacturerName string
	for _, m := range manufacturers {
		manufacturerOptions = append(manufacturerOptions, m.Name)
		formData.manufacturerMap[m.Name] = m.ManufacturerID
		if existingConsole != nil && existingConsole.ManufacturerID != nil && m.ManufacturerID == *existingConsole.ManufacturerID {
			selectedManufacturerName = m.Name
		}
	}
	formData.manufacturerSelect = widget.NewSelect(manufacturerOptions, nil)
	formData.manufacturerSelect.PlaceHolder = "Fabricant (requis)"
	if selectedManufacturerName != "" {
		formData.manufacturerSelect.SetSelected(selectedManufacturerName)
	}

	// Generation
	formData.generationEntry = widget.NewEntry()
	formData.generationEntry.SetPlaceHolder("Génération")
	if existingConsole != nil && existingConsole.Generation != nil {
		formData.generationEntry.SetText(fmt.Sprintf("%d", *existingConsole.Generation))
	}

	// ========== Release Dates ==========
	formData.jpReleaseDateEntry = widget.NewEntry()
	formData.jpReleaseDateEntry.SetPlaceHolder("AAAA-MM-JJ")
	if existingConsole != nil && existingConsole.JPReleaseDate != nil {
		formData.jpReleaseDateEntry.SetText(existingConsole.JPReleaseDate.Format("2006-01-02"))
	}

	formData.usReleaseDateEntry = widget.NewEntry()
	formData.usReleaseDateEntry.SetPlaceHolder("AAAA-MM-JJ")
	if existingConsole != nil && existingConsole.USReleaseDate != nil {
		formData.usReleaseDateEntry.SetText(existingConsole.USReleaseDate.Format("2006-01-02"))
	}

	formData.euReleaseDateEntry = widget.NewEntry()
	formData.euReleaseDateEntry.SetPlaceHolder("AAAA-MM-JJ")
	if existingConsole != nil && existingConsole.EUReleaseDate != nil {
		formData.euReleaseDateEntry.SetText(existingConsole.EUReleaseDate.Format("2006-01-02"))
	}

	formData.discontinuedEntry = widget.NewEntry()
	formData.discontinuedEntry.SetPlaceHolder("AAAA-MM-JJ")
	if existingConsole != nil && existingConsole.Discontinued != nil {
		formData.discontinuedEntry.SetText(existingConsole.Discontinued.Format("2006-01-02"))
	}

	// ========== Prices ==========
	formData.priceUSDEntry = widget.NewEntry()
	formData.priceUSDEntry.SetPlaceHolder("Prix de lancement ($)")
	if existingConsole != nil && existingConsole.PriceUSD != nil {
		formData.priceUSDEntry.SetText(fmt.Sprintf("%d", *existingConsole.PriceUSD))
	}

	formData.priceJPYEntry = widget.NewEntry()
	formData.priceJPYEntry.SetPlaceHolder("Prix de lancement (¥)")
	if existingConsole != nil && existingConsole.PriceJPY != nil {
		formData.priceJPYEntry.SetText(fmt.Sprintf("%d", *existingConsole.PriceJPY))
	}



	// ========== Hardware Specs ==========
	formData.controllersEntry = widget.NewEntry()
	formData.controllersEntry.SetPlaceHolder("Ports contrôleurs")
	if existingConsole != nil && existingConsole.Controllers != nil {
		formData.controllersEntry.SetText(fmt.Sprintf("%d", *existingConsole.Controllers))
	}

	formData.cpuEntry = widget.NewEntry()
	formData.cpuEntry.SetPlaceHolder("CPU")
	if existingConsole != nil && existingConsole.CPU != nil {
		formData.cpuEntry.SetText(*existingConsole.CPU)
	}

	formData.gpuEntry = widget.NewEntry()
	formData.gpuEntry.SetPlaceHolder("GPU")
	if existingConsole != nil && existingConsole.GPU != nil {
		formData.gpuEntry.SetText(*existingConsole.GPU)
	}

	formData.memoryEntry = widget.NewEntry()
	formData.memoryEntry.SetPlaceHolder("Mémoire")
	if existingConsole != nil && existingConsole.Memory != nil {
		formData.memoryEntry.SetText(*existingConsole.Memory)
	}

	formData.audioEntry = widget.NewEntry()
	formData.audioEntry.SetPlaceHolder("Processeur audio")
	if existingConsole != nil && existingConsole.Audio != nil {
		formData.audioEntry.SetText(*existingConsole.Audio)
	}

	// ========== Other Info ==========
	formData.unitsSoldEntry = widget.NewEntry()
	formData.unitsSoldEntry.SetPlaceHolder("Nombre d'unités vendues")
	if existingConsole != nil && existingConsole.UnitsSold != nil {
		formData.unitsSoldEntry.SetText(fmt.Sprintf("%d", *existingConsole.UnitsSold))
	}

	formData.topGameEntry = widget.NewEntry()
	formData.topGameEntry.SetPlaceHolder("Top vente")
	if existingConsole != nil && existingConsole.TopGame != nil {
		formData.topGameEntry.SetText(*existingConsole.TopGame)
	}

	formData.predecessorEntry = widget.NewEntry()
	formData.predecessorEntry.SetPlaceHolder("Prédécesseur")
	if existingConsole != nil && existingConsole.Predecessor != nil {
		formData.predecessorEntry.SetText(*existingConsole.Predecessor)
	}

	formData.successorEntry = widget.NewEntry()
	formData.successorEntry.SetPlaceHolder("SSuccesseur")
	if existingConsole != nil && existingConsole.Successor != nil {
		formData.successorEntry.SetText(*existingConsole.Successor)
	}

	// ========== Collection Info ==========
	formData.ownedCheck = widget.NewCheck("Possédé", nil)
	if existingConsole != nil {
		formData.ownedCheck.Checked = existingConsole.Owned
	} else {
		formData.ownedCheck.Checked = true
	}

	// Condition
	formData.conditionSlider = widget.NewSlider(1, 5)
	formData.conditionSlider.Step = 1
	if existingConsole != nil && existingConsole.Condition != nil {
		formData.conditionSlider.Value = float64(*existingConsole.Condition)
	}
	formData.conditionLabel = widget.NewLabel("État: -")
	if existingConsole != nil && existingConsole.Condition != nil {
		formData.conditionLabel.SetText(fmt.Sprintf("État: %d", *existingConsole.Condition))
	}
	formData.conditionSlider.OnChanged = func(value float64) {
		formData.conditionLabel.SetText(fmt.Sprintf("État: %d", int(value)))
	}

	// ========== Notes ==========
	formData.notesEntry = widget.NewMultiLineEntry()
	formData.notesEntry.SetPlaceHolder("Notes")
	formData.notesEntry.SetMinRowsVisible(3)
	if existingConsole != nil && existingConsole.Notes != nil {
		formData.notesEntry.SetText(*existingConsole.Notes)
	}

	// ========== Build Form Layout ==========
	formData.form = container.NewVBox(
		widget.NewLabel("Nom *"),
		formData.nameEntry,

		widget.NewLabel("Type *"),
		formData.typeSelect,

		widget.NewLabel("Fabricant *"),
		formData.manufacturerSelect,

		widget.NewLabel("Génération"),
		formData.generationEntry,

		widget.NewSeparator(),
		widget.NewLabel("Dates de sortie"),
		widget.NewLabel("Europe:"),
		formData.euReleaseDateEntry,
		widget.NewLabel("United States:"),
		formData.usReleaseDateEntry,
		widget.NewLabel("Japan:"),
		formData.jpReleaseDateEntry,
		widget.NewLabel("Fin de production:"),
		formData.discontinuedEntry,

		widget.NewSeparator(),
		widget.NewLabel("Prix de lancement"),
		widget.NewLabel("USA (USD):"),
		formData.priceUSDEntry,
		widget.NewLabel("Japon (JPY):"),
		formData.priceJPYEntry,

		widget.NewSeparator(),
		widget.NewLabel("Caractéristiques techniques"),
		widget.NewLabel("Ports contrôleurs:"),
		formData.controllersEntry,
		widget.NewLabel("CPU:"),
		formData.cpuEntry,
		widget.NewLabel("GPU:"),
		formData.gpuEntry,
		widget.NewLabel("Mémoire:"),
		formData.memoryEntry,
		widget.NewLabel("Audio:"),
		formData.audioEntry,

		widget.NewSeparator(),
		widget.NewLabel("Ventes & Histoire"),
		widget.NewLabel("UUnités vendues:"),
		formData.unitsSoldEntry,
		widget.NewLabel("Top ventes:"),
		formData.topGameEntry,
		widget.NewLabel("Prédécesseur:"),
		formData.predecessorEntry,
		widget.NewLabel("Successeur:"),
		formData.successorEntry,

		widget.NewSeparator(),
		widget.NewLabel("Informations de collection"),
		formData.ownedCheck,
		formData.conditionLabel,
		formData.conditionSlider,

		widget.NewSeparator(),
		widget.NewLabel("Notes"),
		formData.notesEntry,
	)

	return formData
}

// showAddConsoleDialog shows the dialog to add a new console
func showAddConsoleDialog(w fyne.Window, conn *pgx.Conn, onSuccess func()) {
	formData := buildConsoleForm(w, conn, nil)

	saveBtn := widget.NewButton("Enregistrer", func() {
		_, err := saveConsole(conn, formData, 0)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		dialog.ShowInformation("Enregistré", "Console ajoutée à la base de données.", w)
		if onSuccess != nil {
			onSuccess()
		}
	})

	saveBtn.Importance = widget.HighImportance

	formWithSave := container.NewBorder(
		nil, saveBtn, nil, nil,
		container.NewScroll(formData.form),
	)

	d := dialog.NewCustom("Ajouter", "Annuler", formWithSave, w)
	d.Resize(fyne.NewSize(600, 700))
	d.Show()
}

// showEditConsoleDialog shows the dialog to edit an existing console
func showEditConsoleDialog(w fyne.Window, conn *pgx.Conn, consoleID int, onSuccess func()) {
	existingConsole, err := getConsoleByID(conn, consoleID)
	if err != nil {
		dialog.ShowError(fmt.Errorf("échec du chargement de la console: %w", err), w)
		return
	}

	formData := buildConsoleForm(w, conn, existingConsole)

	saveBtn := widget.NewButton("Enregistrer", func() {
		_, err := saveConsole(conn, formData, consoleID)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		dialog.ShowInformation("Enregistré", "Console mise à jour dans la base de données.", w)
		if onSuccess != nil {
			onSuccess()
		}
	})

	saveBtn.Importance = widget.HighImportance

	formWithSave := container.NewBorder(
		nil, saveBtn, nil, nil,
		container.NewScroll(formData.form),
	)

	d := dialog.NewCustom("Éditer", "Annuler", formWithSave, w)
	d.Resize(fyne.NewSize(600, 700))
	d.Show()
}

// saveConsole saves or updates a console
func saveConsole(conn *pgx.Conn, formData *consoleFormData, consoleID int) (int, error) {
	// Validate
	if formData.nameEntry.Text == "" {
		return 0, fmt.Errorf("nom requis")
	}
	if formData.typeSelect.Selected == "" {
		return 0, fmt.Errorf("type requis")
	}
	if formData.manufacturerSelect.Selected == "" {
		return 0, fmt.Errorf("plateforme requise")
	}

	// Prepare data
	var typeID *int
	if formData.typeSelect.Selected != "" {
		id := formData.typeMap[formData.typeSelect.Selected]
		typeID = &id
	}

	var manufacturerID *int
	if formData.manufacturerSelect.Selected != "" {
		id := formData.manufacturerMap[formData.manufacturerSelect.Selected]
		manufacturerID = &id
	}

	// Parse generation
	var generation *int
	if formData.generationEntry.Text != "" {
		var gen int
		fmt.Sscanf(formData.generationEntry.Text, "%d", &gen)
		generation = &gen
	}

	// Parse dates
	var jpReleaseDate, usReleaseDate, euReleaseDate, discontinued *string
	if formData.jpReleaseDateEntry.Text != "" {
		jpReleaseDate = &formData.jpReleaseDateEntry.Text
	}
	if formData.usReleaseDateEntry.Text != "" {
		usReleaseDate = &formData.usReleaseDateEntry.Text
	}
	if formData.euReleaseDateEntry.Text != "" {
		euReleaseDate = &formData.euReleaseDateEntry.Text
	}
	if formData.discontinuedEntry.Text != "" {
		discontinued = &formData.discontinuedEntry.Text
	}

	// Parse prices
	var priceJPY, priceUSD *int
	if formData.priceJPYEntry.Text != "" {
		var price int
		fmt.Sscanf(formData.priceJPYEntry.Text, "%d", &price)
		priceJPY = &price
	}
	if formData.priceUSDEntry.Text != "" {
		var price int
		fmt.Sscanf(formData.priceUSDEntry.Text, "%d", &price)
		priceUSD = &price
	}

	// Parse hardware specs
	var controllers *int
	if formData.controllersEntry.Text != "" {
		var c int
		fmt.Sscanf(formData.controllersEntry.Text, "%d", &c)
		controllers = &c
	}

	var cpu, gpu, memory, audio *string
	if formData.cpuEntry.Text != "" {
		cpu = &formData.cpuEntry.Text
	}
	if formData.gpuEntry.Text != "" {
		gpu = &formData.gpuEntry.Text
	}
	if formData.memoryEntry.Text != "" {
		memory = &formData.memoryEntry.Text
	}
	if formData.audioEntry.Text != "" {
		audio = &formData.audioEntry.Text
	}

	// Parse other info
	var unitsSold *int
	if formData.unitsSoldEntry.Text != "" {
		var units int
		fmt.Sscanf(formData.unitsSoldEntry.Text, "%d", &units)
		unitsSold = &units
	}

	var topGame, predecessor, successor *string
	if formData.topGameEntry.Text != "" {
		topGame = &formData.topGameEntry.Text
	}
	if formData.predecessorEntry.Text != "" {
		predecessor = &formData.predecessorEntry.Text
	}
	if formData.successorEntry.Text != "" {
		successor = &formData.successorEntry.Text
	}

	// Parse condition
	var condition *int
	if formData.conditionSlider.Value > 0 {
		c := int(formData.conditionSlider.Value)
		condition = &c
	}

	// Parse notes
	var notes *string
	if formData.notesEntry.Text != "" {
		notes = &formData.notesEntry.Text
	}

	// INSERT or UPDATE
	if consoleID == 0 {
		// INSERT
		query := `
			INSERT INTO consoles (
				name, type_id, manufacturer_id, generation,
				jp_release_date, us_release_date, eu_release_date, discontinued,
				price_jpy, price_usd, controllers, cpu, gpu, memory, audio,
				units_sold, top_game, predecessor, successor,
				owned, condition, notes
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
			RETURNING console_id
		`

		err := conn.QueryRow(context.Background(), query,
			formData.nameEntry.Text, typeID, manufacturerID, generation,
			jpReleaseDate, usReleaseDate, euReleaseDate, discontinued,
			priceJPY, priceUSD, controllers, cpu, gpu, memory, audio,
			unitsSold, topGame, predecessor, successor,
			formData.ownedCheck.Checked, condition, notes,
		).Scan(&consoleID)

		if err != nil {
			return 0, fmt.Errorf("échec de l'ajout de console: %w", err)
		}
		return consoleID, nil
	} else {
		// UPDATE
		query := `
			UPDATE consoles SET
				name = $1, type_id = $2, manufacturer_id = $3, generation = $4,
				jp_release_date = $5, us_release_date = $6, eu_release_date = $7, discontinued = $8,
				price_jpy = $9, price_usd = $10, controllers = $11, cpu = $12, gpu = $13, memory = $14, audio = $15,
				units_sold = $16, top_game = $17, predecessor = $18, successor = $19,
				owned = $20, condition = $21, notes = $22
			WHERE console_id = $23
		`

		_, err := conn.Exec(context.Background(), query,
			formData.nameEntry.Text, typeID, manufacturerID, generation,
			jpReleaseDate, usReleaseDate, euReleaseDate, discontinued,
			priceJPY, priceUSD, controllers, cpu, gpu, memory, audio,
			unitsSold, topGame, predecessor, successor,
			formData.ownedCheck.Checked, condition, notes,
			consoleID,
		)

		if err != nil {
			return 0, fmt.Errorf("échec de la mise à jour de la console: %w", err)
		}
		return consoleID, nil
	}
}
