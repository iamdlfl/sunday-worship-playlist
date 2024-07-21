// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    spotifyPlaylist, err := UnmarshalSpotifyPlaylist(bytes)
//    bytes, err = spotifyPlaylist.Marshal()

package main

import "encoding/json"

func UnmarshalSpotifyPlaylist(data []byte) (SpotifyPlaylist, error) {
	var r SpotifyPlaylist
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SpotifyPlaylist) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SpotifyPlaylist struct {
	Collaborative bool         `json:"collaborative"`
	Description   string       `json:"description"`
	ExternalUrls  ExternalUrls `json:"external_urls"`
	Followers     Followers    `json:"followers"`
	Href          string       `json:"href"`
	ID            string       `json:"id"`
	Images        []Image      `json:"images"`
	Name          string       `json:"name"`
	Owner         Owner        `json:"owner"`
	Public        bool         `json:"public"`
	SnapshotID    string       `json:"snapshot_id"`
	Tracks        Tracks       `json:"tracks"`
	Type          string       `json:"type"`
	URI           string       `json:"uri"`
}

type ExternalUrls struct {
	Spotify string `json:"spotify"`
}

type Followers struct {
	Href  string `json:"href"`
	Total int64  `json:"total"`
}

type Image struct {
	URL    string `json:"url"`
	Height int64  `json:"height"`
	Width  int64  `json:"width"`
}

type Owner struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Followers    *Followers   `json:"followers,omitempty"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
	DisplayName  *string      `json:"display_name,omitempty"`
	Name         *string      `json:"name,omitempty"`
}

type Tracks struct {
	Href     string `json:"href"`
	Limit    int64  `json:"limit"`
	Next     string `json:"next"`
	Offset   int64  `json:"offset"`
	Previous string `json:"previous"`
	Total    int64  `json:"total"`
	Items    []Item `json:"items"`
}

type Item struct {
	AddedAt string `json:"added_at"`
	AddedBy Owner  `json:"added_by"`
	IsLocal bool   `json:"is_local"`
	Track   Track  `json:"track"`
}

type Track struct {
	Album            Album        `json:"album"`
	Artists          []Artist     `json:"artists"`
	AvailableMarkets []string     `json:"available_markets"`
	DiscNumber       int64        `json:"disc_number"`
	DurationMS       int64        `json:"duration_ms"`
	Explicit         bool         `json:"explicit"`
	ExternalIDS      ExternalIDS  `json:"external_ids"`
	ExternalUrls     ExternalUrls `json:"external_urls"`
	Href             string       `json:"href"`
	ID               string       `json:"id"`
	IsPlayable       bool         `json:"is_playable"`
	LinkedFrom       LinkedFrom   `json:"linked_from"`
	Restrictions     Restrictions `json:"restrictions"`
	Name             string       `json:"name"`
	Popularity       int64        `json:"popularity"`
	PreviewURL       string       `json:"preview_url"`
	TrackNumber      int64        `json:"track_number"`
	Type             string       `json:"type"`
	URI              string       `json:"uri"`
	IsLocal          bool         `json:"is_local"`
}

type Album struct {
	AlbumType            string       `json:"album_type"`
	TotalTracks          int64        `json:"total_tracks"`
	AvailableMarkets     []string     `json:"available_markets"`
	ExternalUrls         ExternalUrls `json:"external_urls"`
	Href                 string       `json:"href"`
	ID                   string       `json:"id"`
	Images               []Image      `json:"images"`
	Name                 string       `json:"name"`
	ReleaseDate          string       `json:"release_date"`
	ReleaseDatePrecision string       `json:"release_date_precision"`
	Restrictions         Restrictions `json:"restrictions"`
	Type                 string       `json:"type"`
	URI                  string       `json:"uri"`
	Artists              []Owner      `json:"artists"`
}

type Restrictions struct {
	Reason string `json:"reason"`
}

type Artist struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Followers    Followers    `json:"followers"`
	Genres       []string     `json:"genres"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Images       []Image      `json:"images"`
	Name         string       `json:"name"`
	Popularity   int64        `json:"popularity"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

type ExternalIDS struct {
	Isrc string `json:"isrc"`
	Ean  string `json:"ean"`
	Upc  string `json:"upc"`
}

type LinkedFrom struct {
}
