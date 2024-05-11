package main

import (
	"archivist/spotify"
	"archivist/utils"
	"fmt"
)

func main() {
	utils.LoadDotEnv()

	playlists, err := spotify.UserPlaylists()

	if err != nil {
		panic(err)
	}

	for _, playlist := range playlists {
		fmt.Println(playlist.Name)
		fmt.Println(playlist.Description)
	}
}
