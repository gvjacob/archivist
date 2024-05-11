package main

import (
	"archivist/llm"
	"archivist/spotify"
	"archivist/utils"
	"fmt"

	"github.com/charmbracelet/log"
)

func main() {
	utils.LoadDotEnv()

	playlists, err := spotify.UserPlaylists()

	if err != nil {
		log.Fatal(err)
	}

	userSavedTracks, err := spotify.UserSavedTracks()

	if err != nil {
		log.Fatal(err)
	}

	for _, track := range userSavedTracks {
		playlists, err := llm.GetPlaylistsToSaveTrackInto(playlists, track.Track)

		if err != nil {
			log.Fatal(err)
		}

		for _, playlist := range playlists {
			spotify.AddTracksToPlaylist(playlist, track.Track)
			log.Info(fmt.Sprintf("Added track %s to playlist %s", track.Track.Name, playlist.Name))
		}
	}
}
