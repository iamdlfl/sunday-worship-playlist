package main

type PlanList struct {
	Links    Links          `json:"links"`
	Included []interface{}  `json:"included"`
	Meta     Meta           `json:"meta"`
	Data     []PlanListData `json:"data"`
}

type PlanListData struct {
	Type          string             `json:"type"`
	Id            string             `json:"id"`
	Links         Links              `json:"links"`
	Attributes    PlanListAttributes `json:"attributes"`
	Relationships Relationships      `json:"relationships"`
}

type PlanListAttributes struct {
	CanViewOrder         bool        `json:"can_view_order"`
	CreatedAt            string      `json:"created_at"`
	Dates                string      `json:"dates"`
	FilesExpireAt        string      `json:"files_expire_at"`
	ItemCount            int         `json:"items_count"`
	LastTimeAt           string      `json:"last_time_at"`
	MultiDay             bool        `json:"multi_day"`
	NeededPositionsCount int         `json:"needed_positions_count"`
	OtherTimeCount       int         `json:"other_time_count"`
	Permissions          string      `json:"permissions"`
	PlanNotesCount       int         `json:"plan_notes_count"`
	PlanPeopleCount      int         `json:"plan_people_count"`
	PlanningCenterUrl    string      `json:"planning_center_url"`
	PreferOrderView      bool        `json:"prefers_order_view"`
	Public               bool        `json:"public"`
	Rehearsable          bool        `json:"rehearsable"`
	RehearsalTimeCount   int         `json:"rehearsal_time_count"`
	RemindersDisabled    bool        `json:"reminders_disabled"`
	SeriesTitle          interface{} `json:"series_title"`
	ServiceTimeCount     int         `json:"service_time_count"`
	TotalLength          int         `json:"total_length"`
	ShortDates           string      `json:"short_dates"`
	SortDate             string      `json:"sort_date"`
	UpdatedAt            string      `json:"updated_at"`
	Title                interface{} `json:"title"`
}
