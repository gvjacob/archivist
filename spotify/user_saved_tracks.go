package spotify

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
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
	URI        string   `json:"uri"`
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

func (c *SpotifyClient) LastArchiveTime() (time.Time, error) {
	var lastArchivedTimeStamp int

	c.UsersTable.QueryRow(`
    SELECT created_at
    FROM archived_tracks
    WHERE user_id == $1
    ORDER BY created_at DESC
    LIMIT 1
  `, c.User.ID).Scan(&lastArchivedTimeStamp)

	if lastArchivedTimeStamp == 0 {
		return time.Now(), nil
	}

	return time.Unix(int64(lastArchivedTimeStamp), 0), nil
}

func (c *SpotifyClient) UserSavedTracksSinceLastArchive() ([]SavedTrack, error) {
	since, err := c.LastArchiveTime()

	if err != nil {
		return nil, err
	}

	// TODO: Reconsider limit and add pagination support
	resp, err := c.Get("https://api.spotify.com/v1/me/tracks?limit=20")

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

	savedTracks := userSavedTracksResponse.Items
	var tracksAfter []SavedTrack

	for _, track := range savedTracks {
		addedAt, err := time.Parse(time.RFC3339, track.AddedAt)

		if err == nil && since.Before(addedAt) {
			tracksAfter = append(tracksAfter, track)
		}
	}

	return tracksAfter, nil
}
