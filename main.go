package main

import (
	"archivist/llm"
	"archivist/spotify"
	"archivist/storage"
	"archivist/utils"
	"fmt"

	"github.com/charmbracelet/log"
)

func main() {
	utils.LoadDotEnv()

	db, err := storage.NewDatabase("archivist.db")

	if err != nil {
		log.Error("Failed to connect to database")
		log.Fatal(err)
	}

	users := storage.NewUsers(db)
	user, err := users.GetUser()

	if err != nil {
		log.Error("Failed to fetch user")
		log.Fatal(err)
	}

	client := spotify.NewSpotifyClient(user, users)
	playlists, err := client.UserPlaylists()

	if err != nil || len(playlists) == 0 {
		log.Fatal(err)
	}

	userSavedTracks, err := client.UserSavedTracks()

	if err != nil || len(userSavedTracks) == 0 {
		log.Error("Failed to fetch user saved tracks")
		log.Fatal(err)
	}

	for _, track := range userSavedTracks {
		playlists, err := llm.GetPlaylistsToSaveTrackInto(playlists, track.Track)

		if err != nil {
			log.Error(err)
			continue
		}

		if len(playlists) == 0 {
			log.Info(fmt.Sprintf("No playlists selected for track %s", track.Track.Name))
			continue
		}

		for _, playlist := range playlists {
			client.AddTracksToPlaylist(playlist, track.Track)
			log.Info(fmt.Sprintf("Added track '%s' to playlist '%s'", track.Track.Name, playlist.Name))
		}
	}
}
