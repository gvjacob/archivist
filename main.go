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
		prompt, err := llm.ChoosePlaylistsPrompt(playlists, track.Track)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(prompt)
	}
}
