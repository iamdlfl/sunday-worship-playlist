package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

var client http.Client = http.Client{}
var un string
var pw string

var cid string
var cs string
var rt string

var token string

var spotifyApiUrl = "https://api.spotify.com/v1/"
var userId = "onthe_dl"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg, err := ini.Load(".settings.ini")
	if err != nil {
		log.Println(err)
		return
	}
	un = cfg.Section("pco").Key("un").String()
	pw = cfg.Section("pco").Key("pw").String()

	cid = cfg.Section("spotify").Key("cid").String()
	cs = cfg.Section("spotify").Key("cs").String()
	rt = cfg.Section("spotify").Key("rt").String()

	token, err = getSpotifyToken(cid, cs, rt)
	if err != nil {
		log.Println(err)
		return
	}

	today := time.Now()
	dayNumber := int(today.Weekday())
	daysToAdd := 7 - dayNumber
	if daysToAdd == 7 {
		daysToAdd = 0
	}

	upcomingSunday := today.Add(time.Hour * 24 * time.Duration(daysToAdd))
	formattedSunday := upcomingSunday.Format("2006-01-02")
	planNumber, err := getPlanNumberPco(formattedSunday)
	if err != nil {
		log.Println(err)
		return
	}

	// Remove when done testing
	// planNumber = "72665641"
	planNumber = "72665642"
	songs, err := getSongsPco(planNumber)
	if err != nil {
		log.Println(err)
		return
	}
	newSongs := make([]SongInfo, 0, len(songs))
	for _, song := range songs {
		newSong, err := getSongInfoPco(song)
		if err != nil {
			log.Println(err)
		}
		newSongs = append(newSongs, *newSong)
	}

	spotifyIds := make([]string, 0, len(newSongs))
	for _, song := range newSongs {
		// set up search and do it
		search := "track:" + song.Name
		result, err := doSpotifySearch(search, "track")
		if err != nil {
			log.Println(err)
		}

		// create variables to track which Spotify song has the most
		// artist matches to the Author(s) in Planning Center
		songCheck := make(map[string]int)
		numToBeat := -1
		trackId := ""
		for _, item := range result.Tracks.Items {
			// set songcheck to 0 for this item (spotify song) ID
			songCheck[item.ID] = 0
			log.Println("============")
			log.Println(item.Name)
			log.Println(item.Artists)
			// iterate through the Spotify artists
			for _, artist := range item.Artists {
				// process the PCO song authors, splitting on comma and " and"
				// note space included in " and", which is neccessary so that
				// "Chandler Moore" (for instance) doesn't get split between "Ch" and "ler"
				log.Println(song.Author)
				tempAuthors := strings.Split(song.Author, ",")
				authors := make([]string, 0)
				for _, a := range tempAuthors {
					new := strings.Split(a, " and")
					authors = append(authors, new...)
				}
				for _, author := range authors {
					if author == "" {
						continue
					}
					// remove whitespace
					a := strings.TrimSpace(author)
					if strings.EqualFold(artist.Name, a) {
						// if there's a match, increase the match tracking by one
						songCheck[item.ID] += 1
					}
				}
			}
			// I think Spotify orders the results by their best match/popular matches
			// So we will only check if something is GREATER than the greatest number
			// of author matches (not greater or equal to). This way we can preserve
			// Spotify's order of preference if there is a tie.
			if songCheck[item.ID] > numToBeat {
				trackId = item.ID
			}
		}
		if trackId != "" {
			spotifyIds = append(spotifyIds, trackId)
		}
	}

	log.Println(spotifyIds)
	playlistName := "Sunday Worship - " + "2024-07-14" //+ formattedSunday
	existingPlReq := getSpotifyRequest(spotifyApiUrl + "me/playlists?limit=50")
	existingPlResp, err := client.Do(existingPlReq)
	if err != nil {
		log.Println(err)
		return
	}
	defer existingPlResp.Body.Close()
	if existingPlResp.StatusCode >= http.StatusBadRequest {
		log.Println("Can't get my playlists")
	}
	existingPlBody, err := io.ReadAll(existingPlResp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	existingPlaylists, err := UnmarshalMyPlayLists(existingPlBody)
	if err != nil {
		log.Println(err)
		return
	}

	playListId := ""

	for i := 0; i < 6 && existingPlaylists.Next != ""; i++ {
		for _, pl := range existingPlaylists.Items {
			if playlistName == pl.Name {
				playListId = pl.ID
			}
		}
		existingPlReq := getSpotifyRequest(existingPlaylists.Next)
		existingPlResp, err := client.Do(existingPlReq)
		if err != nil {
			log.Println(err)
			return
		}
		defer existingPlResp.Body.Close()
		if existingPlResp.StatusCode >= http.StatusBadRequest {
			log.Println("Can't get my playlists")
		}
		existingPlBody, err := io.ReadAll(existingPlResp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		existingPlaylists, err = UnmarshalMyPlayLists(existingPlBody)
		if err != nil {
			log.Println(err)
			return
		}
	}

	if playListId == "" {

		plData := make(map[string]interface{})
		plData["name"] = playlistName
		plData["public"] = true
		jplBody, _ := json.Marshal(plData)
		plReq, err := http.NewRequest(http.MethodPost, spotifyApiUrl+"users/"+userId+"/playlists", bytes.NewReader(jplBody))
		if err != nil {
			log.Println(err)
			return
		}

		plReq.Header.Set("Authorization", "Bearer "+token)
		plReq.Header.Set("Content-Type", "application/json")
		plResp, err := client.Do(plReq)
		if err != nil {
			log.Println(err)
			return
		}
		defer plResp.Body.Close()
		plBody, err := io.ReadAll(plResp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		if plResp.StatusCode >= http.StatusBadRequest {
			log.Panic("bad response")
		}

		pl, err := UnmarshalSpotifyPlaylist(plBody)
		if err != nil {
			log.Println(err)
			return
		}
		playListId = pl.ID
	}

	tracksString := make([]string, 0)
	for _, track := range spotifyIds {
		tracksString = append(tracksString, "spotify:track:"+track)
	}
	addData := make(map[string]interface{})
	addData["playlist_id"] = playListId
	addData["uris"] = tracksString
	jaBody, _ := json.Marshal(addData)
	addReq, err := http.NewRequest(http.MethodPost, spotifyApiUrl+"playlists/"+playListId+"/tracks", bytes.NewReader(jaBody))
	if err != nil {
		log.Println(err)
		return
	}
	addReq.Header.Set("Authorization", "Bearer "+token)
	addReq.Header.Set("Content-Type", "application/json")

	addResp, err := client.Do(addReq)
	if err != nil {
		log.Println(err)
		return
	}

	defer addResp.Body.Close()
	addBody, err := io.ReadAll(addResp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(addBody))
	if addResp.StatusCode >= http.StatusBadRequest {
		log.Panic("bad response")
	}

}

func getSongInfoPco(songInfo SongInfo) (*SongInfo, error) {
	newSong := SongInfo{
		Name:         songInfo.Name,
		Id:           songInfo.Id,
		ArrangmentId: songInfo.ArrangmentId,
		Key:          songInfo.Key,
	}
	log.Println("Getting info for song: " + songInfo.Name)
	req := getPcoRequest("https://api.planningcenteronline.com/services/v2/songs/" + songInfo.Id)
	arrReq := getPcoRequest("https://api.planningcenteronline.com/services/v2/songs/" + songInfo.Id + "/arrangements/" + songInfo.ArrangmentId)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	arrResp, err := client.Do(arrReq)
	if err != nil {
		return nil, err
	}
	defer arrResp.Body.Close()
	arrBody, err := io.ReadAll(arrResp.Body)
	if err != nil {
		return nil, err
	}
	var song SongPCO
	var arrangement ArrangementPCO

	if jsonErr, jsonErr2 := json.Unmarshal(body, &song), json.Unmarshal(arrBody, &arrangement); jsonErr != nil || jsonErr2 != nil {
		log.Println(jsonErr)
		log.Println(jsonErr2)
		return nil, errors.New("There was an error unmarshalling into song or arrangement")
	}
	newSong.ArrangementName = arrangement.Data.Attributes.ArrangementName
	newSong.ArrangementNotes = arrangement.Data.Attributes.Notes
	newSong.Author = song.Data.Attributes.Author
	newSong.CCLI = fmt.Sprintf("%d", song.Data.Attributes.CcliNum)
	newSong.CopyrightInfo = song.Data.Attributes.Copyright
	newSong.Lyrics = processLyrics(arrangement.Data.Attributes.Lyrics)
	newSong.SongAdmin = song.Data.Attributes.Admin
	newSong.SongNotes = song.Data.Attributes.Notes
	return &newSong, nil
}

func getSongsPco(planNumber string) ([]SongInfo, error) {
	items, err := getItemsPco(planNumber)
	if err != nil {
		return nil, err
	}

	results := make([]SongInfo, 0, len((*items).Data))
	for _, item := range (*items).Data {
		if item.Attributes.ItemType == "song" || item.Relationships.Song.Data != (RelationshipDataAttributesPCO{}) {
			songInfo := SongInfo{
				Name:         item.Attributes.Title,
				Id:           item.Relationships.Song.Data.Id,
				ArrangmentId: item.Relationships.Arrangement.Data.Id,
				Key:          item.Attributes.KeyName,
			}
			results = append(results, songInfo)
		}
	}
	return results, nil
}

func getItemsPco(planNumber string) (*PlanItemsPCO, error) {
	req := getPcoRequest("https://api.planningcenteronline.com/services/v2/service_types/6096/plans/" + planNumber + "/items")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonBody PlanItemsPCO
	if jsonErr := json.Unmarshal(body, &jsonBody); jsonErr != nil {
		return nil, jsonErr
	}
	return &jsonBody, nil
}

func getPlanNumberPco(sundayDate string) (string, error) {
	req := getPcoRequest("https://api.planningcenteronline.com/services/v2/service_types/6096/plans?order=-sort_date&per_page=25")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var jsonBody PlanListPCO

	if jsonErr := json.Unmarshal(body, &jsonBody); jsonErr != nil {
		return "", jsonErr
	}

	upcomingPlanNumber := ""
	for _, plan := range jsonBody.Data {
		if strings.Contains(plan.Attributes.SortDate, sundayDate) {
			upcomingPlanNumber = plan.Id
		}
	}
	return upcomingPlanNumber, nil
}

func getPcoRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		temp := http.Request{}
		req = &temp
	}
	req.SetBasicAuth(getPcoAuth())
	return req
}

func getPcoAuth() (string, string) {
	return un, pw
}

func doSpotifySearch(searchTerm, searchType string) (SpotifyStructs, error) {
	escaped := url.PathEscape(searchTerm)
	req := getSpotifyRequest(spotifyApiUrl + "search?type=" + searchType + "&limit=25&q=" + escaped)
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

func getSpotifyRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		temp := http.Request{}
		req = &temp
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

func getSpotifyToken(cid, cs, rt string) (string, error) {
	u := "https://accounts.spotify.com/api/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", rt)
	req, err := http.NewRequest(http.MethodPost, u, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(cid, cs)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Println("Bad response getting token")
		log.Println(string(body))
		return "", errors.New("response not 200")
	}
	var jsonBody map[string]interface{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		return "", err
	}
	token, ok := jsonBody["access_token"]
	if !ok {
		return "", errors.New("no access token, something is wrong")
	}
	tokenString, ok := token.(string)
	if !ok {
		return "", errors.New("access token not string")
	}
	ref, ok := jsonBody["refresh_token"]
	if !ok {
		log.Println("no refresh token")
	} else {
		ref_tok, ok := ref.(string)
		if ok {
			log.Println("refresh token equal?")
			log.Println(ref_tok == rt)
		}
	}
	return tokenString, nil
}

func processLyrics(s string) string {
	remove := []string{"\n",
		"Verse:",
		"Chorus:",
		"Verse 1:",
		"Verse 2:",
		"Verse 3:",
		"Verse 4:",
		"Verse 5:",
		"Chorus 1:",
		"Chorus 2:",
		"Chorus 3:",
		"Bridge:",
		"Tag:",
		"Bridge 1:",
		"Bridge 2:",
		"Bridge 3:",
		"Bridge 4:",
		"Verse",
		"Chorus",
		"Verse 1",
		"Verse 2",
		"Verse 3",
		"Verse 4",
		"Verse 5",
		"Chorus 1",
		"Chorus 2",
		"Chorus 3",
		"Bridge",
		"Tag",
		"Bridge 1",
		"Bridge 2",
		"Bridge 3",
		"Bridge 4",
	}
	for _, r := range remove {
		s = strings.ReplaceAll(s, r, "")
	}
	return s
}
