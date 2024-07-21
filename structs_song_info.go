package main

type SongInfo struct {
	Name             string
	Id               string
	ArrangmentId     string
	Key              string
	Lyrics           string
	CopyrightInfo    string
	Author           string
	CCLI             string
	SongNotes        string
	SongAdmin        string
	ArrangementName  string
	ArrangementNotes string
}

type Song struct {
	Included []interface{} `json:"included"`
	Meta     Meta          `json:"meta"`
	Data     SongData      `json:"data"`
}

type SongData struct {
	Type       string             `json:"type"`
	Id         string             `json:"id"`
	Links      Links              `json:"links"`
	Attributes SongDataAttributes `json:"attributes"`
}

type SongDataAttributes struct {
	Admin                  string `json:"admin"`
	Author                 string `json:"author"`
	CcliNum                int    `json:"ccli_number"`
	Copyright              string `json:"copyright"`
	CreatedAt              string `json:"created_at"`
	Hidden                 bool   `json:"hidden"`
	LastScheduledAt        string `json:"last_scheduled_at"`
	LastScheduledShortDate string `json:"last_scheduled_short_dates"`
	Notes                  string `json:"notes"`
	Themes                 string `json:"themes"`
	Title                  string `json:"title"`
	UpdatedAt              string `json:"updated_at"`
}

type Arrangement struct {
	Included []interface{}   `json:"included"`
	Meta     Meta            `json:"meta"`
	Data     ArrangementData `json:"data"`
}

type ArrangementData struct {
	Type          string                    `json:"type"`
	Id            string                    `json:"id"`
	Links         Links                     `json:"links"`
	Relationships Relationships             `json:"relationships"`
	Attributes    ArrangementDataAttributes `json:"attributes"`
}

type ArrangementDataAttributes struct {
	UpdatedAt            string         `json:"updated_at"`
	SequenceShort        []string       `json:"sequence_short"`
	SequenceFull         []FullSequence `json:"sequence_full"`
	Sequence             []string       `json:"sequence"`
	ArchivedAt           string         `json:"archived_at"`
	Bpm                  int            `json:"bpm"`
	ChordChart           string         `json:"chord_chart"`
	ChordChartChordColor int            `json:"chord_chart_chord_color"`
	ChordChartColumns    int            `json:"chord_chart_columns"`
	ChordChartFont       string         `json:"chord_chart_font"`
	ChordChartFontSize   int            `json:"chord_chart_font_size"`
	ChordChartKey        string         `json:"chord_chart_key"`
	CreatedAt            string         `json:"created_at"`
	HasChordChart        bool           `json:"has_chord_chart"`
	HasChords            bool           `json:"has_chords"`
	Length               int            `json:"length"`
	Lyrics               string         `json:"lyrics"`
	LyricsEnabled        bool           `json:"lyrics_enabled"`
	Meter                string         `json:"meter"`
	ArrangementName      string         `json:"name"`
	Notes                string         `json:"notes"`
	NumberChartEnabled   bool           `json:"number_chart_enabled"`
	NumeralChartEnabled  bool           `json:"numeral_chart_enabled"`
	PrintMargin          string         `json:"print_margin"`
	PrintOrientation     string         `json:"print_orientation"`
	PrintPageSize        string         `json:"print_page_size"`
}

type FullSequence struct {
	Label  string      `json:"label"`
	Number interface{} `json:"number"`
	Time   string      `json:"t"`
	Sid    int         `json:"sid"`
}
