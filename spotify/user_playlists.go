package spotify

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Playlist struct {
	ID          string
	Name        string
	Description string
}

type CurrentUserPlaylistsResponse struct {
	Next  string     `json:"next"`
	Total int        `json:"total"`
	Items []Playlist `json:"items"`
}

func (c *SpotifyClient) UserPlaylists() ([]Playlist, error) {
	resp, err := c.Get("https://api.spotify.com/v1/me/playlists?limit=50")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to fetch playlists. Status code: " + resp.Status)
	}

	userPlaylistsResponse := CurrentUserPlaylistsResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userPlaylistsResponse)

	if err != nil {
		return nil, err
	}

	return filterPlaylistsByArchivist(userPlaylistsResponse.Items), nil
}

func filterPlaylistsByArchivist(playlists []Playlist) []Playlist {
	archivistPlaylists := []Playlist{}

	for _, playlist := range playlists {
		if strings.HasPrefix(playlist.Description, "Archivist:") {
			archivistPlaylists = append(archivistPlaylists, playlist)
		}
	}

	return archivistPlaylists
}
