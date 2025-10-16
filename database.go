package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func dbconnect() (*pgx.Conn, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Make connection string from .env with trimming
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		strings.TrimSpace(os.Getenv("DB_USER")),
		strings.TrimSpace(os.Getenv("DB_PASSWORD")),
		strings.TrimSpace(os.Getenv("DB_HOST")),
		strings.TrimSpace(os.Getenv("DB_PORT")),
		strings.TrimSpace(os.Getenv("DB_NAME")),
	)

	// Connect to database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return conn, nil
}

// ========== Games Functions ==========

func getGames(conn *pgx.Conn) ([]Game, error) {
	query := `
		SELECT 
			g.game_id,
			g.title,
			COALESCE(c.name, '') as console_name,
			COALESCE(ge.name, '') as genre_name,
			g.condition
		FROM games g
		LEFT JOIN consoles c ON g.console_id = c.console_id
		LEFT JOIN genres ge ON g.genre_id = ge.genre_id
		ORDER BY g.title
	`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var g Game
		err := rows.Scan(
			&g.GameID,
			&g.Title,
			&g.ConsoleName,
			&g.GenreName,
			&g.Condition,
		)
		if err != nil {
			return nil, err
		}
		games = append(games, g)
	}
	return games, nil
}

// ========== Consoles Functions ==========

func getConsoles(conn *pgx.Conn) ([]Console, error) {
	query := `
		SELECT 
			c.console_id,
			c.name,
			COALESCE(m.name, '') as manufacturer_name,
			COALESCE(c.generation, 0) as generation,
			c.condition
		FROM consoles c
		LEFT JOIN manufacturers m ON c.manufacturer_id = m.manufacturer_id
		ORDER BY c.name
	`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consoles []Console
	for rows.Next() {
		var c Console
		var gen int
		err := rows.Scan(
			&c.ConsoleID,
			&c.Name,
			&c.ManufacturerName,
			&gen,
			&c.Condition,
		)
		if err != nil {
			return nil, err
		}
		if gen != 0 {
			c.Generation = &gen
		}
		consoles = append(consoles, c)
	}
	return consoles, nil
}

// ========== Lookup Tables Functions (for dropdowns) ==========

func getGenres(conn *pgx.Conn) ([]Genre, error) {
	rows, err := conn.Query(context.Background(), "SELECT genre_id, name FROM genres ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []Genre
	for rows.Next() {
		var g Genre
		err := rows.Scan(&g.GenreID, &g.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, nil
}

func getDevelopers(conn *pgx.Conn) ([]Developer, error) {
	rows, err := conn.Query(context.Background(), "SELECT developer_id, name FROM developers ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var developers []Developer
	for rows.Next() {
		var d Developer
		err := rows.Scan(&d.DeveloperID, &d.Name)
		if err != nil {
			return nil, err
		}
		developers = append(developers, d)
	}
	return developers, nil
}

func getComposers(conn *pgx.Conn) ([]Composer, error) {
	rows, err := conn.Query(context.Background(), "SELECT composer_id, name FROM composers ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var composers []Composer
	for rows.Next() {
		var c Composer
		err := rows.Scan(&c.ComposerID, &c.Name)
		if err != nil {
			return nil, err
		}
		composers = append(composers, c)
	}
	return composers, nil
}

func getPublishers(conn *pgx.Conn) ([]Publisher, error) {
	rows, err := conn.Query(context.Background(), "SELECT publisher_id, name FROM publishers ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var publishers []Publisher
	for rows.Next() {
		var p Publisher
		err := rows.Scan(&p.PublisherID, &p.Name)
		if err != nil {
			return nil, err
		}
		publishers = append(publishers, p)
	}
	return publishers, nil
}

func getProducers(conn *pgx.Conn) ([]Producer, error) {
	rows, err := conn.Query(context.Background(), "SELECT producer_id, name FROM producers ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var producers []Producer
	for rows.Next() {
		var p Producer
		err := rows.Scan(&p.ProducerID, &p.Name)
		if err != nil {
			return nil, err
		}
		producers = append(producers, p)
	}
	return producers, nil
}

func getManufacturers(conn *pgx.Conn) ([]Manufacturer, error) {
	rows, err := conn.Query(context.Background(), "SELECT manufacturer_id, name FROM manufacturers ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var manufacturers []Manufacturer
	for rows.Next() {
		var m Manufacturer
		err := rows.Scan(&m.ManufacturerID, &m.Name)
		if err != nil {
			return nil, err
		}
		manufacturers = append(manufacturers, m)
	}
	return manufacturers, nil
}

func getConsoleTypes(conn *pgx.Conn) ([]ConsoleType, error) {
	rows, err := conn.Query(context.Background(), "SELECT type_id, name FROM console_types ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []ConsoleType
	for rows.Next() {
		var t ConsoleType
		err := rows.Scan(&t.TypeID, &t.Name)
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, nil
}

func getRatingSystems(conn *pgx.Conn) ([]RatingSystem, error) {
	rows, err := conn.Query(context.Background(), "SELECT rating_id, region, code, description FROM rating_systems ORDER BY region, code")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ratings []RatingSystem
	for rows.Next() {
		var r RatingSystem
		err := rows.Scan(&r.RatingID, &r.Region, &r.Code, &r.Description)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, r)
	}
	return ratings, nil
}

// getGameByID fetches a complete game record with all relationships
func getGameByID(conn *pgx.Conn, gameID int) (*Game, error) {
	// Fetch main game data
	query := `
		SELECT 
			g.game_id, g.title, g.console_id, g.genre_id,
			g.jp_release_date, g.us_release_date, g.eu_release_date,
			g.jp_rating_id, g.us_rating_id, g.eu_rating_id,
			g.units_sold, g.owned, g.box_owned, g.collector, g.condition,
			g.purchase_date, g.purchase_price, g.notes,
			COALESCE(c.name, '') as console_name,
			COALESCE(ge.name, '') as genre_name
		FROM games g
		LEFT JOIN consoles c ON g.console_id = c.console_id
		LEFT JOIN genres ge ON g.genre_id = ge.genre_id
		WHERE g.game_id = $1
	`

	var game Game
	err := conn.QueryRow(context.Background(), query, gameID).Scan(
		&game.GameID, &game.Title, &game.ConsoleID, &game.GenreID,
		&game.JPReleaseDate, &game.USReleaseDate, &game.EUReleaseDate,
		&game.JPRatingID, &game.USRatingID, &game.EURatingID,
		&game.UnitsSold, &game.Owned, &game.BoxOwned, &game.Collector, &game.Condition,
		&game.PurchaseDate, &game.PurchasePrice, &game.Notes,
		&game.ConsoleName, &game.GenreName,
	)
	if err != nil {
		return nil, err
	}

	// Fetch developers
	devRows, _ := conn.Query(context.Background(), `
		SELECT d.name, d.developer_id
		FROM game_developers gd
		JOIN developers d ON gd.developer_id = d.developer_id
		WHERE gd.game_id = $1
	`, gameID)
	defer devRows.Close()

	for devRows.Next() {
		var name string
		var id int
		devRows.Scan(&name, &id)
		game.Developers = append(game.Developers, name)
	}

	// Fetch composers
	compRows, _ := conn.Query(context.Background(), `
		SELECT c.name
		FROM game_composers gc
		JOIN composers c ON gc.composer_id = c.composer_id
		WHERE gc.game_id = $1
	`, gameID)
	defer compRows.Close()

	for compRows.Next() {
		var name string
		compRows.Scan(&name)
		game.Composers = append(game.Composers, name)
	}

	// Fetch publishers
	pubRows, _ := conn.Query(context.Background(), `
		SELECT p.name
		FROM game_publishers gp
		JOIN publishers p ON gp.publisher_id = p.publisher_id
		WHERE gp.game_id = $1
	`, gameID)
	defer pubRows.Close()

	for pubRows.Next() {
		var name string
		pubRows.Scan(&name)
		game.Publishers = append(game.Publishers, name)
	}

	// Fetch producers
	prodRows, _ := conn.Query(context.Background(), `
		SELECT p.name
		FROM game_producers gpr
		JOIN producers p ON gpr.producer_id = p.producer_id
		WHERE gpr.game_id = $1
	`, gameID)
	defer prodRows.Close()

	for prodRows.Next() {
		var name string
		prodRows.Scan(&name)
		game.Producers = append(game.Producers, name)
	}

	return &game, nil
}

// deleteGame deletes a game and all its relationships
func deleteGame(conn *pgx.Conn, gameID int) error {
	// Foreign key constraints will cascade delete join table entries
	_, err := conn.Exec(context.Background(), "DELETE FROM games WHERE game_id = $1", gameID)
	return err
}

// ========== Accessories Functions ==========

func getAccessories(conn *pgx.Conn) ([]Accessory, error) {
	query := `
		SELECT 
			a.accessory_id,
			a.name,
			COALESCE(a.color, '') as color,
			COALESCE(m.name, '') as manufacturer_name,
			COALESCE(at.name, '') as type_name,
			COALESCE(a.quantity, 1) as quantity,
			a.condition
		FROM accessories a
		LEFT JOIN manufacturers m ON a.manufacturer_id = m.manufacturer_id
		LEFT JOIN accessory_types at ON a.type_id = at.type_id
		ORDER BY a.name
	`

	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accessories []Accessory
	for rows.Next() {
		var a Accessory
		var color string
		err := rows.Scan(
			&a.AccessoryID,
			&a.Name,
			&color,
			&a.ManufacturerName,
			&a.TypeName,
			&a.Quantity,
			&a.Condition,
		)
		if err != nil {
			return nil, err
		}
		if color != "" {
			a.Color = &color
		}
		accessories = append(accessories, a)
	}
	return accessories, nil
}

func getAccessoryByID(conn *pgx.Conn, accessoryID int) (*Accessory, error) {
	query := `
		SELECT 
			a.accessory_id, a.name, a.color, a.type_id, a.manufacturer_id,
			a.condition, a.owned, a.purchase_date, a.purchase_price,
			a.quantity, a.notes,
			COALESCE(m.name, '') as manufacturer_name,
			COALESCE(at.name, '') as type_name
		FROM accessories a
		LEFT JOIN manufacturers m ON a.manufacturer_id = m.manufacturer_id
		LEFT JOIN accessory_types at ON a.type_id = at.type_id
		WHERE a.accessory_id = $1
	`

	var accessory Accessory
	err := conn.QueryRow(context.Background(), query, accessoryID).Scan(
		&accessory.AccessoryID, &accessory.Name, &accessory.Color,
		&accessory.TypeID, &accessory.ManufacturerID,
		&accessory.Condition, &accessory.Owned, &accessory.PurchaseDate,
		&accessory.PurchasePrice, &accessory.Quantity, &accessory.Notes,
		&accessory.ManufacturerName, &accessory.TypeName,
	)
	if err != nil {
		return nil, err
	}

	// Fetch associated consoles
	consoleRows, _ := conn.Query(context.Background(), `
		SELECT c.name
		FROM accessory_consoles ac
		JOIN consoles c ON ac.console_id = c.console_id
		WHERE ac.accessory_id = $1
	`, accessoryID)
	defer consoleRows.Close()

	for consoleRows.Next() {
		var consoleName string
		consoleRows.Scan(&consoleName)
		accessory.Consoles = append(accessory.Consoles, consoleName)
	}

	return &accessory, nil
}

func deleteAccessory(conn *pgx.Conn, accessoryID int) error {
	conn.Exec(context.Background(), "DELETE FROM accessory_consoles WHERE accessory_id = $1", accessoryID)
	_, err := conn.Exec(context.Background(), "DELETE FROM accessories WHERE accessory_id = $1", accessoryID)
	return err
}

func getAccessoryTypes(conn *pgx.Conn) ([]AccessoryType, error) {
	rows, err := conn.Query(context.Background(), "SELECT type_id, name FROM accessory_types ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []AccessoryType
	for rows.Next() {
		var t AccessoryType
		err := rows.Scan(&t.TypeID, &t.Name)
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, nil
}

func getConsoleByID(conn *pgx.Conn, consoleID int) (*Console, error) {
	query := `
		SELECT 
			c.console_id, c.name, c.generation, c.type_id, c.manufacturer_id,
			c.jp_release_date, c.us_release_date, c.eu_release_date, c.discontinued,
			c.price_jpy, c.price_usd, c.controllers, c.cpu, c.gpu, c.memory, c.audio,
			c.units_sold, c.top_game, c.predecessor, c.successor,
			c.owned, c.condition, c.notes,
			COALESCE(m.name, '') as manufacturer_name,
			COALESCE(ct.name, '') as type_name
		FROM consoles c
		LEFT JOIN manufacturers m ON c.manufacturer_id = m.manufacturer_id
		LEFT JOIN console_types ct ON c.type_id = ct.type_id
		WHERE c.console_id = $1
	`

	var console Console
	err := conn.QueryRow(context.Background(), query, consoleID).Scan(
		&console.ConsoleID, &console.Name, &console.Generation, &console.TypeID, &console.ManufacturerID,
		&console.JPReleaseDate, &console.USReleaseDate, &console.EUReleaseDate, &console.Discontinued,
		&console.PriceJPY, &console.PriceUSD, &console.Controllers, &console.CPU, &console.GPU, &console.Memory, &console.Audio,
		&console.UnitsSold, &console.TopGame, &console.Predecessor, &console.Successor,
		&console.Owned, &console.Condition, &console.Notes,
		&console.ManufacturerName, &console.TypeName,
	)
	if err != nil {
		return nil, err
	}

	return &console, nil
}

func deleteConsole(conn *pgx.Conn, consoleID int) error {
	_, err := conn.Exec(context.Background(), "DELETE FROM consoles WHERE console_id = $1", consoleID)
	return err
}
