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
			COALESCE(ge.name, '') as genre_name
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
			COALESCE(c.generation, 0) as generation
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
