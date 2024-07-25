package main

import (
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/iamdlfl/pco"
	"github.com/iamdlfl/spotify"
)

var client http.Client = http.Client{}

var configFileName string = ".settings.ini"
var spotifyUrl = "https://api.spotify.com/v1"
var userId = "onthe_dl"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	pcoClient := pco.NewPcoClient(configFileName, nil)

	spotifyClient, err := spotify.NewSpotifyClient(configFileName, spotifyUrl, nil)
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
	planNumber, err := pcoClient.GetPlanNumberPco(formattedSunday)
	if err != nil {
		log.Println(err)
		return
	}

	// Remove when done testing
	planNumber = "72665641"
	// planNumber = "72665642"
	songs, err := pcoClient.GetSongsPco(planNumber)
	if err != nil {
		log.Println(err)
		return
	}
	newSongs := make([]pco.SongInfo, 0, len(songs))
	for _, song := range songs {
		newSong, err := pcoClient.GetSongInfoPco(song)
		if err != nil {
			log.Println(err)
		}
		newSongs = append(newSongs, *newSong)
	}

	spotifyIds := make([]string, 0, len(newSongs))
	for _, song := range newSongs {
		// set up search and do it
		search := "track:" + song.Name
		result, err := spotifyClient.DoSpotifySearch(search, "track")
		if err != nil {
			log.Println(err)
		}

		// create variables to track which Spotify song has the most
		// artist matches to the Author(s) in Planning Center
		songCheck := make(map[string]int)
		numToBeat := -1
		trackId := ""
		// Spotify does not sort the tracks by popularity really, though supposedly they try to
		// sort by a combination of match and popularity. We will sort the results ourselves to be sure.
		// This way the most popular versions of songs are first (which are generally the examples we use)
		slices.SortFunc(result.Tracks.Items, func(a, b spotify.ItemSearch) int {
			return int(a.Popularity) - int(b.Popularity)
		})
		for _, item := range result.Tracks.Items {
			// set songcheck to 0 for this item (spotify song) ID
			songCheck[item.ID] = 0
			// iterate through the Spotify artists
			for _, artist := range item.Artists {
				// process the PCO song authors, splitting on comma and " and"
				// note space included in " and", which is neccessary so that
				// "Chandler Moore" (for instance) doesn't get split between "Ch" and "ler"
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
						strings.Contains(a, "Shane Barnard") { // Shane Barnard is often one of the Authors in PCO
						songCheck[item.ID] += 1
					}
				}
			}
			// We will only check if something is GREATER than the greatest number
			// of author matches (not greater or equal to). This way we can preserve
			// our order of preference from above sorting if there is a tie.
			if songCheck[item.ID] > numToBeat {
				// we don't really want instrumental versions of the songs (I think)
				if strings.Contains(item.Name, "Instrumental") {
					continue
				}
				log.Println("REPLACING ITEM")
				log.Println(item.Name)
				log.Println(item.Artists)
				trackId = item.ID
				numToBeat = songCheck[item.ID]
			}
		}
		if trackId != "" {
			spotifyIds = append(spotifyIds, trackId)
		}
	}

	log.Println(spotifyIds)
	playlistName := "Sunday Worship - " + "2024-07-07" //+ formattedSunday
	existingPlaylists, err := spotifyClient.GetExistingPlaylists()
	playListId := ""

	for _, pl := range existingPlaylists.Items {
		if playlistName == pl.Name {
			playListId = pl.ID
			if pl.Tracks.Total >= 4 {
				log.Println("Playlist already created")
				emailer.SendMessage("Playlist has already been created, and has 4 or more songs. Exit successfully.")
				return
			}
		}
	}

	if playListId == "" {
		pl, err := spotifyClient.CreateSpotifyPlaylist(playlistName, userId)
		if err != nil {
			log.Println(err)
			emailer.SendMessage(err.Error() + "\n\nThere was an error creating the playlist")
		}
		playListId = pl.ID
	}

	tracksString := make([]string, 0)
	for _, track := range spotifyIds {
		tracksString = append(tracksString, "spotify:track:"+track)
	}

	err = spotifyClient.AddSongsToPlaylist(playListId, tracksString)
	if err != nil {
		emailer.SendMessage("Could not add songs to playlist!")
		return
	}

	emailer.SendMessage("Successfully set up playlist!")
}
