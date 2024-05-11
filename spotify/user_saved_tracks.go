package spotify

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Artist struct {
	Name   string   `json:"name"`
	Genres []string `json:"genres"`
}

type Album struct {
	Name        string `json:"name"`
	Type        string `json:"album_type"`
	ReleaseDate string `json:"release_date"`
}

type Track struct {
	Name       string   `json:"name"`
	Album      Album    `json:"album"`
	Artists    []Artist `json:"artists"`
	DurationMs int      `json:"duration_ms"`
	Explicit   bool     `json:"explicit"`
}

type SavedTrack struct {
	Track   Track  `json:"track"`
	AddedAt string `json:"added_at"`
}

type UserSavedTracksResponse struct {
	Next  string       `json:"next"`
	Total int          `json:"total"`
	Items []SavedTrack `json:"items"`
}

func UserSavedTracks() ([]SavedTrack, error) {
	resp, err := Get("https://api.spotify.com/v1/me/tracks?limit=1")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to fetch saved tracks. Status code: " + resp.Status)
	}

	userSavedTracksResponse := UserSavedTracksResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userSavedTracksResponse)

	if err != nil {
		return nil, err
	}

	return userSavedTracksResponse.Items, nil
}
