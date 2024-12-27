package main

import (
	"archivist/commands"
	"archivist/llm"
	"archivist/spotify"
	"archivist/storage"
	"archivist/utils"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

func main() {
	utils.LoadDotEnv()

	command, found := utils.SafeGet(os.Args, 1)

	if found && command == "init" {
		commands.Initialize()
	} else {
		archive()
	}
}

func archive() {
	db, err := storage.NewDatabase(os.Getenv("SQLITE_FILE_PATH"))

	if err != nil {
		log.Fatal(err)
	}

	users := storage.NewUsers(db)
	user, err := users.GetUser()

	if err != nil {
		log.Fatal(err)
	}

	client := spotify.NewSpotifyClient(user, users)
	playlists, err := client.UserPlaylists()

	if err != nil {
		log.Fatal(err)
	}

	if len(playlists) == 0 {
		log.Warn("No playlists found. Exiting...")
		os.Exit(0)
	}

	userSavedTracks, err := client.UserSavedTracksSinceLastArchive()

	if err != nil {
		log.Fatal(err)
	}

	if len(userSavedTracks) == 0 {
		log.Warn("No new tracks to archive")
		return
	}

	var archivedTracks []storage.ArchivedTrack

	for _, track := range userSavedTracks {
		playlists, err := llm.GetPlaylistsToSaveTrackInto(playlists, track.Track)

		if err != nil {
			log.Error(err)
			continue
		}

		if len(playlists) == 0 {
			log.Warn(fmt.Sprintf("No playlists selected for track '%s'", track.Track.Name))

			archivedTracks = append(archivedTracks, storage.ArchivedTrack{
				UserID:  user.ID,
				TrackID: track.Track.ID,
			})

			continue
		}

		for _, playlist := range playlists {
			_, err := client.AddTracksToPlaylist(playlist, track.Track)

			if err != nil {
				log.Error(err)
				continue
			}

			log.Info(fmt.Sprintf("Added track '%s' to playlist '%s'", track.Track.Name, playlist.Name))

			archivedTracks = append(archivedTracks, storage.ArchivedTrack{
				UserID:     user.ID,
				TrackID:    track.Track.ID,
				PlaylistID: playlist.ID,
			})
		}
	}

	if len(archivedTracks) > 0 {
		archivedTracksTable := storage.NewArchivedTracks(db)

		err = archivedTracksTable.Insert(archivedTracks)

		if err != nil {
			log.Fatal(err)
		}

		log.Info("Archived tracks successfully")
	}
}
