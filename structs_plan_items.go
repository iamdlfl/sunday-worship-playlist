package main

type PlanItems struct {
	Links    Links           `json:"links"`
	Included []interface{}   `json:"included"`
	Meta     Meta            `json:"meta"`
	Data     []PlanItemsData `json:"data"`
}

type PlanItemsData struct {
	Type          string                 `json:"type"`
	Id            string                 `json:"id"`
	Relationships Relationships          `json:"relationships"`
	Links         Links                  `json:"links"`
	Attributes    PlanItemDataAttributes `json:"attributes"`
}

type PlanItemDataAttributes struct {
	CreatedAt                      string      `json:"created_at"`
	CustomArrangementSequence      interface{} `json:"custom_arrangement_sequence"`
	CustomArrangementSequenceFull  interface{} `json:"custom_arrangement_sequence_full"`
	CustomArrangementSequenceShort interface{} `json:"custom_arrangement_sequence_short"`
	Description                    string      `json:"description"`
	HtmlDetails                    interface{} `json:"html_details"`
	ItemType                       string      `json:"item_type"`
	KeyName                        string      `json:"key_name"`
	Length                         int         `json:"length"`
	Sequence                       int         `json:"sequence"`
	ServicePosition                string      `json:"service_position"`
	Title                          string      `json:"title"`
	UpdatedAt                      string      `json:"updated_at"`
}
