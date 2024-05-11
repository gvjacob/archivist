package main

import (
	"archivist/spotify"
	"archivist/utils"
	"fmt"
)

func main() {
	utils.LoadDotEnv()

	// playlists, err := spotify.UserPlaylists()

	// if err != nil {
	// 	panic(err)
	// }

	userSavedTracks, err := spotify.UserSavedTracks()

	if err != nil {
		panic(err)
	}

	for _, track := range userSavedTracks {
		fmt.Println(track.AddedAt, track.Track.Name)
	}
}
