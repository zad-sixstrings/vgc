package main

import "time"

// ========== Main Entity Structs ==========

// Game represents a game in the collection
type Game struct {
	GameID        int
	Title         string
	JPReleaseDate *time.Time
	USReleaseDate *time.Time
	EUReleaseDate *time.Time
	UnitsSold     *int
	Owned         bool
	BoxOwned      *bool
	Collector     *bool
	Condition     *int
	PurchaseDate  *time.Time
	PurchasePrice *float64
	Notes         *string

	// Foreign keys
	ConsoleID  *int
	GenreID    *int
	JPRatingID *int
	USRatingID *int
	EURatingID *int

	// Related data (for display in tables) - populated via JOINs
	ConsoleName string
	GenreName   string
	JPRating    string
	USRating    string
	EURating    string
	Developers  []string
	Publishers  []string
	Composers   []string
	Producers   []string
}

// Console represents a console in the collection
type Console struct {
	ConsoleID     int
	Name          string
	Generation    *int
	JPReleaseDate *time.Time
	USReleaseDate *time.Time
	EUReleaseDate *time.Time
	Discontinued  *time.Time
	PriceJPY      *int
	PriceUSD      *int
	Controllers   *int
	CPU           *string
	GPU           *string
	Memory        *string
	Audio         *string
	UnitsSold     *int
	TopGame       *string
	Predecessor   *string
	Successor     *string
	Owned         bool
	Condition     *int
	Notes         *string

	// Foreign keys
	TypeID         *int
	ManufacturerID *int

	// Related data (for display)
	TypeName         string
	ManufacturerName string
}

// Accessory represents an accessory in the collection
type Accessory struct {
	AccessoryID   int
	Name          string
	Condition     *int
	Owned         bool
	PurchaseDate  *time.Time
	PurchasePrice *float64
	Quantity      int
	Notes         *string

	// Foreign keys
	TypeID         *int
	ManufacturerID *int

	// Related data (for display)
	TypeName         string
	ManufacturerName string
	Consoles         []string // Multiple consoles via join table
}

// ========== Lookup Table Structs ==========

type Genre struct {
	GenreID int
	Name    string
}

type Developer struct {
	DeveloperID int
	Name        string
}

type Composer struct {
	ComposerID int
	Name       string
}

type Publisher struct {
	PublisherID int
	Name        string
}

type Producer struct {
	ProducerID int
	Name       string
}

type Manufacturer struct {
	ManufacturerID int
	Name           string
}

type ConsoleType struct {
	TypeID int
	Name   string
}

type AccessoryType struct {
	TypeID int
	Name   string
}

type RatingSystem struct {
	RatingID    int
	Region      string
	Code        string
	Description *string
}
