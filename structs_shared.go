package main

type Links struct {
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

type Meta struct {
	TotalCount int        `json:"total_count"`
	Count      int        `json:"count"`
	Next       MetaNext   `json:"next"`
	CanOrderBy []string   `json:"can_order_by"`
	CanQueryBy []string   `json:"can_query_by"`
	CanInclude []string   `json:"can_include"`
	CanFilter  []string   `json:"can_filter"`
	Parent     MetaParent `json:"parent"`
}

type MetaNext struct {
	Offset int `json:"offset"`
}

type MetaParent struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type RelationshipData struct {
	Data RelationshipDataAttributes `json:"data"`
}

type RelationshipDataAttributes struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type Relationships struct {
	Plan                    RelationshipData `json:"plan"`
	Song                    RelationshipData `json:"song"`
	Arrangement             RelationshipData `json:"arrangement"`
	Key                     RelationshipData `json:"key"`
	SelectedLayout          RelationshipData `json:"selected_layout"`
	SelectedBackground      RelationshipData `json:"selected_background"`
	UpdatedBy               RelationshipData `json:"updated_by"`
	CreatedBy               RelationshipData `json:"created_by"`
	ServiceType             RelationshipData `json:"service_type"`
	PreviousPlan            RelationshipData `json:"previous_plan"`
	NextPlan                RelationshipData `json:"next_plan"`
	Series                  RelationshipData `json:"series"`
	LinkedPublishingEpisode RelationshipData `json:"linked_publishing_episode"`
	AttachmentTypes         struct {
		Data []interface{} `json:"data"`
	} `json:"attachment_types"`
}
