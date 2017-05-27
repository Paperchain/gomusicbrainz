package gomusicbrainz

type Artist struct {
	Disambiguation string `json:"disambiguation"`
	ID             string `json:"id"`
	SortName       string `json:"sort-name"`
	Name           string `json:"name"`
}

type ArtistCredit struct {
	Name       string `json:"name"`
	Artist     Artist `json:"artist"`
	JoinPhrase string `json:"joinphrase"`
}

// Recording is a struct that encapsulates MusicBraniz's Recording entity
type Recording struct {
	Title          string         `json:"title"`
	Length         int            `json:"length"`
	ID             string         `json:"id"`
	Disambiguation string         `json:"disambiguation"`
	ISRCs          []string       `json:"isrcs"`
	IsVideo        bool           `json:"video"`
	ArtistCredits  []ArtistCredit `json:"artist-credit"`
}

type ISRC struct {
	ISRCID     string      `json:"isrc"`
	Recordings []Recording `json:"recordings"`
}

type ISWC struct {
	Works      []Work `json:"works"`
	WorkOffset int    `json:"work-offset"`
	WorkCount  int    `json:"work-count"`
}

// Attribute is a struct that encapsulate metadata of an entity
type Attribute struct {
	Type   string `json:"type"`
	Value  string `json:"value"`
	TypeID string `json:"type-id"`
}

// Work is
type Work struct {
	Language       string      `json:"language"`
	TypeID         string      `json:"type-id"`
	Languages      []string    `json:"languages"`
	Type           string      `json:"type"`
	ISWCs          []string    `json:"iswcs"`
	Disambiguation string      `json:"disambiguation"`
	ID             string      `json:"id"`
	Title          string      `json:"title"`
	Attributes     []Attribute `json:"attributes"`
	Aliases        []string    `json:"aliases"`
}
