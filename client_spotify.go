package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/ini.v1"
)

func NewSpotifyClient(configFile string) (*spotify, error) {
	c := spotify{ConfigFileName: configFile}
	cfg, err := ini.Load(configFile)
	if err != nil {
		log.Println(err)
		return &c, err
	}
	cid := cfg.Section("spotify").Key("cid").String()
	cs := cfg.Section("spotify").Key("cs").String()
	rt := cfg.Section("spotify").Key("rt").String()

	c.cid = cid
	c.cs = cs
	c.rt = rt

	err = c.getSpotifyToken()
	if err != nil {
		return &c, err
	}
	return &c, nil
}

type spotify struct {
	cid            string
	cs             string
	rt             string
	token          string
	ConfigFileName string
}

func (s spotify) doSpotifySearch(searchTerm, searchType string) (SpotifyStructs, error) {
	escaped := url.PathEscape(searchTerm)
	req := s.getSpotifyRequest(spotifyApiUrl + "search?type=" + searchType + "&limit=25&q=" + escaped)
	resp, err := client.Do(req)
	if err != nil {
		return SpotifyStructs{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SpotifyStructs{}, err
	}
	spotifySearchResult, err := UnmarshalSpotifyStructs(body)
	if err != nil {
		return SpotifyStructs{}, err
	}
	return spotifySearchResult, nil
}

func (s spotify) getSpotifyRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		temp := http.Request{}
		req = &temp
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	return req
}

func (s spotify) getSpotifyToken() error {
	u := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", s.rt)
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(s.cid, s.cs)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Println("Bad response getting token")
		log.Println(string(body))
		return errors.New("response not 200")
	}
	var jsonBody map[string]interface{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		return err
	}
	token, ok := jsonBody["access_token"]
	if !ok {
		return errors.New("no access token, something is wrong")
	}
	tokenString, ok := token.(string)
	if !ok {
		return errors.New("access token not string")
	}
	ref, ok := jsonBody["refresh_token"]
	if !ok {
		log.Println("no refresh token")
	} else {
		ref_tok, ok := ref.(string)
		if ok {
			log.Println("refresh token equal?")
			log.Println(ref_tok == s.rt)
		}
	}
	s.token = tokenString
	return nil
}
