// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    spotifyStructs, err := UnmarshalSpotifyStructs(bytes)
//    bytes, err = spotifyStructs.Marshal()

package spotify

import "encoding/json"

func UnmarshalSpotifyStructs(data []byte) (SpotifyStructs, error) {
	var r SpotifyStructs
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SpotifyStructs) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SpotifyStructs struct {
	Tracks TracksSearch `json:"tracks"`
}

type TracksSearch struct {
	Href     string       `json:"href"`
	Items    []ItemSearch `json:"items"`
	Limit    int64        `json:"limit"`
	Next     string       `json:"next"`
	Offset   int64        `json:"offset"`
	Previous interface{}  `json:"previous"`
	Total    int64        `json:"total"`
}

type ItemSearch struct {
	Album            AlbumSearch           `json:"album"`
	Artists          []ArtistElementSearch `json:"artists"`
	AvailableMarkets []string              `json:"available_markets"`
	DiscNumber       int64                 `json:"disc_number"`
	DurationMS       int64                 `json:"duration_ms"`
	Explicit         bool                  `json:"explicit"`
	ExternalIDS      ExternalIDSSearch     `json:"external_ids"`
	ExternalUrls     ExternalUrlsSearch    `json:"external_urls"`
	Href             string                `json:"href"`
	ID               string                `json:"id"`
	IsLocal          bool                  `json:"is_local"`
	Name             string                `json:"name"`
	Popularity       int64                 `json:"popularity"`
	PreviewURL       *string               `json:"preview_url"`
	TrackNumber      int64                 `json:"track_number"`
	Type             string                `json:"type"`
	URI              string                `json:"uri"`
}

type AlbumSearch struct {
	AlbumType            string                `json:"album_type"`
	Artists              []ArtistElementSearch `json:"artists"`
	AvailableMarkets     []string              `json:"available_markets"`
	ExternalUrls         ExternalUrlsSearch    `json:"external_urls"`
	Href                 string                `json:"href"`
	ID                   string                `json:"id"`
	Images               []ImageSearch         `json:"images"`
	Name                 string                `json:"name"`
	ReleaseDate          string                `json:"release_date"`
	ReleaseDatePrecision string                `json:"release_date_precision"`
	TotalTracks          int64                 `json:"total_tracks"`
	Type                 string                `json:"type"`
	URI                  string                `json:"uri"`
}

type ArtistElementSearch struct {
	ExternalUrls ExternalUrlsSearch `json:"external_urls"`
	Href         string             `json:"href"`
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Type         Type               `json:"type"`
	URI          string             `json:"uri"`
}

type ExternalUrlsSearch struct {
	Spotify string `json:"spotify"`
}

type ImageSearch struct {
	Height int64  `json:"height"`
	URL    string `json:"url"`
	Width  int64  `json:"width"`
}

type ExternalIDSSearch struct {
	Isrc string `json:"isrc"`
}

type Type string

const (
	ArtistSearch Type = "artist"
)
