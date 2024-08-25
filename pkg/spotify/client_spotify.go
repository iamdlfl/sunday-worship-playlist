package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/ini.v1"
)

func NewSpotifyClient(configFile, apiUrl string, client *http.Client) (*spotify, error) {
	if client == nil {
		c := http.Client{}
		client = &c
	}
	c := spotify{ConfigFileName: configFile, ApiUrl: apiUrl, client: client}
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
	ApiUrl         string
	client         *http.Client
}

func (s spotify) DoSpotifySearch(searchTerm, searchType string) (SpotifyStructs, error) {
	escaped := url.PathEscape(searchTerm)
	req := s.getSpotifyRequest(http.MethodGet, s.ApiUrl+"/search?type="+searchType+"&limit=10&q="+escaped, nil)
	resp, err := s.client.Do(req)
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

func (s spotify) GetExistingPlaylists() (MyPlayLists, error) {
	allItems := make([]ItemMy, 0)
	existingPlReq := s.getSpotifyRequest(http.MethodGet, s.ApiUrl+"/me/playlists?limit=50", nil)
	existingPlResp, err := s.client.Do(existingPlReq)
	if err != nil {
		return MyPlayLists{}, err
	}
	defer existingPlResp.Body.Close()
	if existingPlResp.StatusCode >= http.StatusBadRequest {
		return MyPlayLists{}, errors.New("response code >= 400")
	}
	existingPlBody, err := io.ReadAll(existingPlResp.Body)
	if err != nil {
		return MyPlayLists{}, err
	}
	existingPlaylists, err := UnmarshalMyPlayLists(existingPlBody)
	if err != nil {
		return MyPlayLists{}, err
	}
	allItems = append(allItems, existingPlaylists.Items...)
	for i := 0; i < 10 && existingPlaylists.Next != ""; i++ {
		existingPlReq := s.getSpotifyRequest(http.MethodGet, existingPlaylists.Next, nil)
		existingPlResp, err := s.client.Do(existingPlReq)
		if err != nil {
			return MyPlayLists{}, err
		}
		defer existingPlResp.Body.Close()
		if existingPlResp.StatusCode >= http.StatusBadRequest {
			return MyPlayLists{}, errors.New("response code >= 400")
		}
		existingPlBody, err := io.ReadAll(existingPlResp.Body)
		if err != nil {
			return MyPlayLists{}, err
		}
		existingPlaylists, err = UnmarshalMyPlayLists(existingPlBody)
		if err != nil {
			return MyPlayLists{}, err
		}
		allItems = append(allItems, existingPlaylists.Items...)
	}
	allPlaylists := MyPlayLists{
		Total: int64(len(allItems)),
		Items: allItems,
	}
	return allPlaylists, nil
}

func (s spotify) CreateSpotifyPlaylist(plName, userId string) (SpotifyPlaylist, error) {
	plData := make(map[string]interface{})
	plData["name"] = plName
	plData["public"] = true
	jplBody, _ := json.Marshal(plData)
	plReq := s.getSpotifyRequest(http.MethodPost, s.ApiUrl+"/users/"+userId+"/playlists", bytes.NewReader(jplBody))
	plReq.Header.Set("Content-Type", "application/json")
	plResp, err := s.client.Do(plReq)
	if err != nil {
		return SpotifyPlaylist{}, err
	}
	defer plResp.Body.Close()
	plBody, err := io.ReadAll(plResp.Body)
	if err != nil {
		return SpotifyPlaylist{}, err
	}
	if plResp.StatusCode >= http.StatusBadRequest {
		return SpotifyPlaylist{}, errors.New("status code >= 400")
	}

	pl, err := UnmarshalSpotifyPlaylist(plBody)
	if err != nil {
		return SpotifyPlaylist{}, err
	}
	return pl, nil
}

func (s spotify) AddSongsToPlaylist(pid string, tracks []string) error {
	addData := make(map[string]interface{})
	addData["playlist_id"] = pid
	addData["uris"] = tracks
	jaBody, _ := json.Marshal(addData)
	addReq := s.getSpotifyRequest(http.MethodPost, s.ApiUrl+"/playlists/"+pid+"/tracks", bytes.NewReader(jaBody))
	addReq.Header.Set("Content-Type", "application/json")

	addResp, err := s.client.Do(addReq)
	if err != nil {
		return err
	}

	defer addResp.Body.Close()
	addBody, err := io.ReadAll(addResp.Body)
	if err != nil {
		return err
	}
	log.Println(string(addBody))
	if addResp.StatusCode >= http.StatusBadRequest {
		return errors.New("status code >= 400")
	}
	return nil
}

func (s spotify) GetPlaylistById(pid string) (SpotifyPlaylist, error) {
	p := SpotifyPlaylist{}
	req := s.getSpotifyRequest(http.MethodGet, s.ApiUrl+"/playlists/"+pid, nil)
	res, err := s.client.Do(req)
	if err != nil {
		return p, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		log.Println("Status code not ok: ", res.StatusCode)
		return p, errors.New("status code >= 400")
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return p, err
	}

	jsonErr := json.Unmarshal(body, &p)

	return p, jsonErr
}

func (s spotify) RemoveAllSongsFromPlaylist(pid string) error {
	type TrackToDelete struct {
		Uri string `json:"uri"`
	}
	type DeleteBody struct {
		Tracks []TrackToDelete `json:"tracks"`
	}
	playlist, err := s.GetPlaylistById(pid)
	if err != nil {
		return err
	}
	delReqBody := DeleteBody{
		Tracks: []TrackToDelete{},
	}
	if len(playlist.Tracks.Items) == 0 {
		log.Println("No tracks found to delete")
		// early return - there are no tracks on this playlist to remove
		return nil
	}
	for _, trackItem := range playlist.Tracks.Items {
		delTrack := TrackToDelete{Uri: trackItem.Track.URI}
		delReqBody.Tracks = append(delReqBody.Tracks, delTrack)
	}
	log.Printf("Deleting %d tracks from playlist %s", len(delReqBody.Tracks), pid)

	jaBody, err := json.Marshal(delReqBody)
	if err != nil {
		log.Println(err)
	}
	delReq := s.getSpotifyRequest(http.MethodDelete, s.ApiUrl+"/playlists/"+pid+"/tracks", bytes.NewReader(jaBody))
	delReq.Header.Set("Content-Type", "application/json")

	delResp, err := s.client.Do(delReq)
	if err != nil {
		return err
	}

	defer delResp.Body.Close()
	delRespBody, err := io.ReadAll(delResp.Body)
	if err != nil {
		return err
	}
	log.Println(string(delRespBody))
	if delResp.StatusCode >= http.StatusBadRequest {
		return errors.New("status code >= 400")
	}
	return nil
}

// getSpotifyRequest just returns a request with the authorization already filled out
func (s spotify) getSpotifyRequest(method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		temp := http.Request{}
		req = &temp
	}
	req.Header.Set("Authorization", "Bearer "+s.token)
	return req
}

func (s *spotify) getSpotifyToken() error {
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
	resp, err := s.client.Do(req)
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
