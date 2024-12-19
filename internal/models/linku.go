package models

type WordsData map[string]WordData
type SignsData map[string]SignData

// type LinkuData map[string]WordData
type LinkuData struct {
	Words WordsData
	Signs SignsData
}

var UsageCategories = map[string]int{
	"core":     4,
	"common":   3,
	"uncommon": 2,
	"obscure":  1,
	"sandbox":  0,
}

// Helper for rendering pu data in templates.
type PuData struct {
	PartOfSpeech string
	Definition   string
}

// Helpers for rendering etymology data in templates.
type EtymologyData struct {
	Source  string
	Entries []EtymologyEntry
}

type EtymologyEntry struct {
	Word       string
	Alt        string
	Definition string
	Language   string
}
