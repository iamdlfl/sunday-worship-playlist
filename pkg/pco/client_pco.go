package pco

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"gopkg.in/ini.v1"
)

type pco struct {
	ConfigFileName string
	client         *http.Client
}

func NewPcoClient(configFile string, client *http.Client) *pco {
	var c http.Client
	if client == nil {
		c = http.Client{}
	} else {
		c = *client
	}
	p := pco{ConfigFileName: configFile, client: &c}
	return &p
}

func (p pco) GetSongInfoPco(songInfo SongInfo) (*SongInfo, error) {
	newSong := SongInfo{
		Name:         songInfo.Name,
		Id:           songInfo.Id,
		ArrangmentId: songInfo.ArrangmentId,
		Key:          songInfo.Key,
	}
	log.Println("Getting info for song: " + songInfo.Name)
	req := p.getPcoRequest("https://api.planningcenteronline.com/services/v2/songs/" + songInfo.Id)
	arrReq := p.getPcoRequest("https://api.planningcenteronline.com/services/v2/songs/" + songInfo.Id + "/arrangements/" + songInfo.ArrangmentId)
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	arrResp, err := p.client.Do(arrReq)
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
	newSong.Lyrics = StripLyrics(arrangement.Data.Attributes.Lyrics)
	newSong.SongAdmin = song.Data.Attributes.Admin
	newSong.SongNotes = song.Data.Attributes.Notes
	return &newSong, nil
}

func (p pco) GetSongsPco(planNumber string) ([]SongInfo, error) {
	items, err := p.GetItemsPco(planNumber)
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
	if len(results) == 0 {
		return nil, errors.New("there were no songs in planning center")
	}
	return results, nil
}

func (p pco) GetItemsPco(planNumber string) (*PlanItemsPCO, error) {
	req := p.getPcoRequest("https://api.planningcenteronline.com/services/v2/service_types/6096/plans/" + planNumber + "/items")
	resp, err := p.client.Do(req)
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

func (p pco) GetPlanNumberPco(sundayDate string) (string, error) {
	req := p.getPcoRequest("https://api.planningcenteronline.com/services/v2/service_types/6096/plans?order=-sort_date&per_page=25")
	resp, err := p.client.Do(req)
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

func (p pco) getPcoRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		temp := http.Request{}
		req = &temp
	}
	req.SetBasicAuth(p.getPcoAuth())
	return req
}

func (p pco) getPcoAuth() (string, string) {
	cfg, err := ini.Load(p.ConfigFileName)
	if err != nil {
		log.Println(err)
		return "", ""
	}
	un := cfg.Section("pco").Key("un").String()
	pw := cfg.Section("pco").Key("pw").String()
	return un, pw
}

func StripLyrics(s string) string {
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
