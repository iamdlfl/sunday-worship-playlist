package pco

type SongPCO struct {
	Included []interface{} `json:"included"`
	Meta     MetaPCO       `json:"meta"`
	Data     SongDataPCO   `json:"data"`
}

type SongDataPCO struct {
	Type       string                `json:"type"`
	Id         string                `json:"id"`
	Links      LinksPCO              `json:"links"`
	Attributes SongDataAttributesPCO `json:"attributes"`
}

type SongDataAttributesPCO struct {
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

type ArrangementPCO struct {
	Included []interface{}      `json:"included"`
	Meta     MetaPCO            `json:"meta"`
	Data     ArrangementDataPCO `json:"data"`
}

type ArrangementDataPCO struct {
	Type          string                       `json:"type"`
	Id            string                       `json:"id"`
	Links         LinksPCO                     `json:"links"`
	Relationships Relationships                `json:"relationships"`
	Attributes    ArrangementDataAttributesPCO `json:"attributes"`
}

type ArrangementDataAttributesPCO struct {
	UpdatedAt            string            `json:"updated_at"`
	SequenceShort        []string          `json:"sequence_short"`
	SequenceFull         []FullSequencePCO `json:"sequence_full"`
	Sequence             []string          `json:"sequence"`
	ArchivedAt           string            `json:"archived_at"`
	Bpm                  int               `json:"bpm"`
	ChordChart           string            `json:"chord_chart"`
	ChordChartChordColor int               `json:"chord_chart_chord_color"`
	ChordChartColumns    int               `json:"chord_chart_columns"`
	ChordChartFont       string            `json:"chord_chart_font"`
	ChordChartFontSize   int               `json:"chord_chart_font_size"`
	ChordChartKey        string            `json:"chord_chart_key"`
	CreatedAt            string            `json:"created_at"`
	HasChordChart        bool              `json:"has_chord_chart"`
	HasChords            bool              `json:"has_chords"`
	Length               int               `json:"length"`
	Lyrics               string            `json:"lyrics"`
	LyricsEnabled        bool              `json:"lyrics_enabled"`
	Meter                string            `json:"meter"`
	ArrangementName      string            `json:"name"`
	Notes                string            `json:"notes"`
	NumberChartEnabled   bool              `json:"number_chart_enabled"`
	NumeralChartEnabled  bool              `json:"numeral_chart_enabled"`
	PrintMargin          string            `json:"print_margin"`
	PrintOrientation     string            `json:"print_orientation"`
	PrintPageSize        string            `json:"print_page_size"`
}

type FullSequencePCO struct {
	Label  string      `json:"label"`
	Number interface{} `json:"number"`
	Time   string      `json:"t"`
	Sid    int         `json:"sid"`
}
