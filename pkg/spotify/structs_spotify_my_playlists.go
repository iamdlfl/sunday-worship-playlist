// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    myPlayLists, err := UnmarshalMyPlayLists(bytes)
//    bytes, err = myPlayLists.Marshal()

package spotify

import "encoding/json"

func UnmarshalMyPlayLists(data []byte) (MyPlayLists, error) {
	var r MyPlayLists
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *MyPlayLists) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MyPlayLists struct {
	Href     string      `json:"href"`
	Limit    int64       `json:"limit"`
	Next     string      `json:"next"`
	Offset   int64       `json:"offset"`
	Previous interface{} `json:"previous"`
	Total    int64       `json:"total"`
	Items    []ItemMy    `json:"items"`
}

type ItemMy struct {
	Collaborative bool           `json:"collaborative"`
	Description   string         `json:"description"`
	ExternalUrls  ExternalUrlsMy `json:"external_urls"`
	Href          string         `json:"href"`
	ID            string         `json:"id"`
	Images        []ImageMy      `json:"images"`
	Name          string         `json:"name"`
	Owner         OwnerMy        `json:"owner"`
	PrimaryColor  interface{}    `json:"primary_color"`
	Public        bool           `json:"public"`
	SnapshotID    string         `json:"snapshot_id"`
	Tracks        TracksMy       `json:"tracks"`
	Type          ItemTypeMy     `json:"type"`
	URI           string         `json:"uri"`
}

type ExternalUrlsMy struct {
	Spotify string `json:"spotify"`
}

type ImageMy struct {
	Height *int64 `json:"height"`
	URL    string `json:"url"`
	Width  *int64 `json:"width"`
}

type OwnerMy struct {
	DisplayName  DisplayNameMy  `json:"display_name"`
	ExternalUrls ExternalUrlsMy `json:"external_urls"`
	Href         string         `json:"href"`
	ID           IDMy           `json:"id"`
	Type         OwnerTypeMy    `json:"type"`
	URI          URIMy          `json:"uri"`
}

type TracksMy struct {
	Href  string `json:"href"`
	Total int64  `json:"total"`
}

type DisplayNameMy string

const (
	DisplayNameOntheDL DisplayNameMy = "onthe_dl"
	HillsidePlaylists  DisplayNameMy = "Hillside Playlists"
)

type IDMy string

const (
	IDOntheDL                 IDMy = "onthe_dl"
	Nsosz32W7004R00Cx9Rilkscq IDMy = "nsosz32w7004r00cx9rilkscq"
)

type OwnerTypeMy string

const (
	User OwnerTypeMy = "user"
)

type URIMy string

const (
	SpotifyUserNsosz32W7004R00Cx9Rilkscq URIMy = "spotify:user:nsosz32w7004r00cx9rilkscq"
	SpotifyUserOntheDL                   URIMy = "spotify:user:onthe_dl"
)

type ItemTypeMy string

const (
	Playlist ItemTypeMy = "playlist"
)
