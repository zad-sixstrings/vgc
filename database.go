package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func dbconnect() (*pgx.Conn, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	fmt.Println("DEBUG - DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DEBUG - DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	fmt.Println("DEBUG - DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DEBUG - DB_PORT:", os.Getenv("DB_PORT"))
	fmt.Println("DEBUG - DB_NAME:", os.Getenv("DB_NAME"))

	// Make connection string from .env
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println("DEBUG - Connection string:", connStr)

	// Connect to database
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return conn, nil
}

func getGames(conn *pgx.Conn) ([]Game, error) {
	rows, err := conn.Query(context.Background(), "SELECT game_id, title, platform, genre FROM games")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var g Game
		err := rows.Scan(&g.GameID, &g.Title, &g.Platform, &g.Genre)
		if err != nil {
			return nil, err
		}
		games = append(games, g)
	}
	return games, nil
}

func getConsoles(conn *pgx.Conn) ([]Console, error) {
	rows, err := conn.Query(context.Background(), "SELECT console_id, name, manufacturer, generation FROM consoles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consoles []Console
	for rows.Next() {
		var c Console
		err := rows.Scan(&c.ConsoleID, &c.Name, &c.Manufacturer, &c.Generation)
		if err != nil {
			return nil, err
		}
		consoles = append(consoles, c)
	}
	return consoles, nil
}
