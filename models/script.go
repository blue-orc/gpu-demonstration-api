package models

type Script struct {
	ID           int
	Name         string
	Processor    string
	Description  string
	LocationPath string
	LocationType string
}
