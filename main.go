package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

var client http.Client = http.Client{}

var configFileName string = ".settings.ini"
var spotifyApiUrl = "https://api.spotify.com/v1/"
var userId = "onthe_dl"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	pcoClient := pco{ConfigFileName: configFileName}

	spotifyClient, err := NewSpotifyClient(configFileName)
	if err != nil {
		log.Println(err)
		return
	}

	emailer, err := NewMailer(configFileName)
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
	planNumber, err := pcoClient.getPlanNumberPco(formattedSunday)
	if err != nil {
		log.Println(err)
		return
	}

	// Remove when done testing
	// planNumber = "72665641"
	planNumber = "72665642"
	songs, err := pcoClient.getSongsPco(planNumber)
	if err != nil {
		log.Println(err)
		return
	}
	newSongs := make([]SongInfo, 0, len(songs))
	for _, song := range songs {
		newSong, err := pcoClient.getSongInfoPco(song)
		if err != nil {
			log.Println(err)
		}
		newSongs = append(newSongs, *newSong)
	}

	spotifyIds := make([]string, 0, len(newSongs))
	for _, song := range newSongs {
		// set up search and do it
		search := "track:" + song.Name
		result, err := spotifyClient.doSpotifySearch(search, "track")
		if err != nil {
			log.Println(err)
		}

		// create variables to track which Spotify song has the most
		// artist matches to the Author(s) in Planning Center
		songCheck := make(map[string]int)
		numToBeat := -1
		trackId := ""
		slices.SortFunc(result.Tracks.Items, func(a, b ItemSearch) int {
			return int(a.Popularity) - int(b.Popularity)
		})
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

					// if there's a match, increase the match tracking by one
					if strings.EqualFold(artist.Name, a) {
						songCheck[item.ID] += 1
					} else if strings.Contains(artist.Name, "Getty") && strings.Contains(a, "Getty") { // Keith & Kirsten Getty are hard to parse/match
						songCheck[item.ID] += 1
					} else if (strings.Contains(artist.Name, "Shane & Shane") ||
						strings.Contains(artist.Name, "Shane and Shane")) &&
						strings.Contains(a, "Shane Barnard") {
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
	existingPlReq := spotifyClient.getSpotifyRequest(spotifyApiUrl + "me/playlists?limit=50")
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

	for i := 0; i < 10 && existingPlaylists.Next != ""; i++ {
		for _, pl := range existingPlaylists.Items {
			if playlistName == pl.Name {
				playListId = pl.ID
				if pl.Tracks.Total >= 4 {
					log.Println("Playlist already created")
					return
				}
			}
		}
		existingPlReq := spotifyClient.getSpotifyRequest(existingPlaylists.Next)
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

		plReq.Header.Set("Authorization", "Bearer "+spotifyClient.token)
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
	addReq.Header.Set("Authorization", "Bearer "+spotifyClient.token)
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
	emailer.SendMessage("Successfully set up playlist!")
}
