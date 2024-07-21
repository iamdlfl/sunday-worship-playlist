package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var client http.Client = http.Client{}

func main() {
	today := time.Now()
	dayNumber := int(today.Weekday())
	daysToAdd := 7 - dayNumber
	if daysToAdd == 7 {
		daysToAdd = 0
	}

	upcomingSunday := today.Add(time.Hour * 24 * time.Duration(daysToAdd))
	formattedSunday := upcomingSunday.Format("2006-01-02")
	planNumber, err := getPlanNumber(formattedSunday)
	if err != nil {
		log.Panic(err)
	}

	// Remove when done testing
	// planNumber = "72665641"
	planNumber = "72665642"
	songs, err := getSongs(planNumber)
	if err != nil {
		log.Panic(err)
	}
	for _, song := range songs {
		getSongInfo(&song)
	}
}

func getSongInfo(songInfo *SongInfo) error {
	log.Println("Getting info for song: " + songInfo.Name)
	req := getRequest("https://api.planningcenteronline.com/services/v2/songs/" + songInfo.Id)
	arrReq := getRequest("https://api.planningcenteronline.com/services/v2/songs/" + songInfo.Id + "/arrangements/" + songInfo.ArrangmentId)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	arrResp, err := client.Do(arrReq)
	if err != nil {
		return err
	}
	defer arrResp.Body.Close()
	arrBody, err := io.ReadAll(arrResp.Body)
	if err != nil {
		return err
	}
	var song Song
	var arrangement Arrangement

	if jsonErr, jsonErr2 := json.Unmarshal(body, &song), json.Unmarshal(arrBody, &arrangement); jsonErr != nil || jsonErr2 != nil {
		log.Println(jsonErr)
		log.Println(jsonErr2)
		return errors.New("There was an error unmarshalling into song or arrangement")
	}
	songInfo.ArrangementName = arrangement.Data.Attributes.ArrangementName
	songInfo.ArrangementNotes = arrangement.Data.Attributes.Notes
	songInfo.Author = song.Data.Attributes.Author
	songInfo.CCLI = fmt.Sprintf("%d", song.Data.Attributes.CcliNum)
	songInfo.CopyrightInfo = song.Data.Attributes.Copyright
	songInfo.Lyrics = processLyrics(arrangement.Data.Attributes.Lyrics)
	songInfo.SongAdmin = song.Data.Attributes.Admin
	songInfo.SongNotes = song.Data.Attributes.Notes
	return nil
}

func getSongs(planNumber string) ([]SongInfo, error) {
	items, err := getItems(planNumber)
	if err != nil {
		return nil, err
	}

	results := make([]SongInfo, 0, len((*items).Data))
	for _, item := range (*items).Data {
		if item.Attributes.ItemType == "song" || item.Relationships.Song.Data != (RelationshipDataAttributes{}) {
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

func getItems(planNumber string) (*PlanItems, error) {
	req := getRequest("https://api.planningcenteronline.com/services/v2/service_types/6096/plans/" + planNumber + "/items")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var jsonBody PlanItems
	if jsonErr := json.Unmarshal(body, &jsonBody); jsonErr != nil {
		return nil, jsonErr
	}
	return &jsonBody, nil
}

func getPlanNumber(sundayDate string) (string, error) {
	req := getRequest("https://api.planningcenteronline.com/services/v2/service_types/6096/plans?order=-sort_date&per_page=25")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var jsonBody PlanList

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

func getRequest(url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		temp := http.Request{}
		req = &temp
	}
	req.SetBasicAuth(getAuth())
	return req
}

func getAuth() (string, string) {
	return "", ""
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
