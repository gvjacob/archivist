package main

import (
	"archivist/llm"
	"archivist/spotify"
	"archivist/utils"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

func main() {
	utils.LoadDotEnv()

	client := spotify.NewSpotifyClient(os.Getenv("SPOTIFY_ACCESS_TOKEN"))
	playlists, err := client.UserPlaylists()

	if err != nil || len(playlists) == 0 {
		log.Error("Failed to fetch playlists")
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
