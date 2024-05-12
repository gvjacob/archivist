package main

import (
	"archivist/llm"
	"archivist/spotify"
	"archivist/storage"
	"archivist/utils"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func main() {
	utils.LoadDotEnv()

	db, err := storage.NewDatabase(os.Getenv("SQLITE_FILE_PATH"))

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

	userSavedTracks, err := client.UserSavedTracks(user.LastArchived)

	if err != nil {
		log.Error("Failed to fetch user saved tracks")
		log.Fatal(err)
	}

	if len(userSavedTracks) == 0 {
		log.Warn("No new tracks to archive")
		return
	}

	var hasAddedTracks bool

	for _, track := range userSavedTracks {
		playlists, err := llm.GetPlaylistsToSaveTrackInto(playlists, track.Track)

		if err != nil {
			log.Error(err)
			continue
		}

		if len(playlists) == 0 {
			log.Info(fmt.Sprintf("No playlists selected for track '%s'", track.Track.Name))
			continue
		}

		for _, playlist := range playlists {
			client.AddTracksToPlaylist(playlist, track.Track)
			log.Info(fmt.Sprintf("Added track '%s' to playlist '%s'", track.Track.Name, playlist.Name))

			hasAddedTracks = true
		}
	}

	if hasAddedTracks {
		user.LastArchived = time.Now().UTC()
		users.UpdateUser(user)
	}
}
