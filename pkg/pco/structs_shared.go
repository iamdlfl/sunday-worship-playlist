package pco

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
	SpotifyTrackId   string
}
type LinksPCO struct {
	Self              string `json:"self"`
	Next              string `json:"next"`
	Html              string `json:"html"`
	Arrangements      string `json:"arrangements"`
	AssignTags        string `json:"assign_tags"`
	Attachments       string `json:"attachments"`
	LastScheduledItem string `json:"last_scheduled_item"`
	SongSchedules     string `json:"song_schedules"`
	Tags              string `json:"tags"`
	Sections          string `json:"sections"`
	Keys              string `json:"keys"`
}

type MetaPCO struct {
	TotalCount int           `json:"total_count"`
	Count      int           `json:"count"`
	Next       MetaNextPCO   `json:"next"`
	CanOrderBy []string      `json:"can_order_by"`
	CanQueryBy []string      `json:"can_query_by"`
	CanInclude []string      `json:"can_include"`
	CanFilter  []string      `json:"can_filter"`
	Parent     MetaParentPCO `json:"parent"`
}

type MetaNextPCO struct {
	Offset int `json:"offset"`
}

type MetaParentPCO struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type RelationshipDataPCO struct {
	Data RelationshipDataAttributesPCO `json:"data"`
}

type RelationshipDataAttributesPCO struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type Relationships struct {
	Plan                    RelationshipDataPCO `json:"plan"`
	Song                    RelationshipDataPCO `json:"song"`
	Arrangement             RelationshipDataPCO `json:"arrangement"`
	Key                     RelationshipDataPCO `json:"key"`
	SelectedLayout          RelationshipDataPCO `json:"selected_layout"`
	SelectedBackground      RelationshipDataPCO `json:"selected_background"`
	UpdatedBy               RelationshipDataPCO `json:"updated_by"`
	CreatedBy               RelationshipDataPCO `json:"created_by"`
	ServiceType             RelationshipDataPCO `json:"service_type"`
	PreviousPlan            RelationshipDataPCO `json:"previous_plan"`
	NextPlan                RelationshipDataPCO `json:"next_plan"`
	Series                  RelationshipDataPCO `json:"series"`
	LinkedPublishingEpisode RelationshipDataPCO `json:"linked_publishing_episode"`
	AttachmentTypes         struct {
		Data []interface{} `json:"data"`
	} `json:"attachment_types"`
}
