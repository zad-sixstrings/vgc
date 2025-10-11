package main

type Game struct {
	GameID   int
	Title    string
	Platform string
	Genre    string
}

type Console struct {
	ConsoleID    int
	Name         string
	Manufacturer string
	Generation   int
}
