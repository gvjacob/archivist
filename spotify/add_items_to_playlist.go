package spotify

import (
	"encoding/json"
	"errors"
	"net/http"
)

type AddTracksToPlaylistRequest struct {
	URIs []string `json:"uris"`
}

type AddTracksToPlaylistResponse struct {
	SnapshotID string `json:"snapshot_id"`
}

func AddTracksToPlaylist(playlist Playlist, track Track) (string, error) {
	requestBody := AddTracksToPlaylistRequest{
		URIs: []string{track.URI},
	}

	jsonBody, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}

	resp, err := Post("https://api.spotify.com/v1/playlists/"+playlist.ID+"/tracks", jsonBody)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", errors.New("Failed to add track to playlist. Status code: " + resp.Status)
	}

	addTracksToPlaylistResponse := AddTracksToPlaylistResponse{}
	err = json.NewDecoder(resp.Body).Decode(&addTracksToPlaylistResponse)

	if err != nil {
		return "", err
	}

	return addTracksToPlaylistResponse.SnapshotID, nil
}
