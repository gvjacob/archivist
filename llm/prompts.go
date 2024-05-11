package llm

import (
	"archivist/spotify"
	"bytes"
	"math"
	"strings"
	"text/template"
)

func ChoosePlaylistsPrompt(playlists []spotify.Playlist, track spotify.Track) (string, error) {
	template, err := template.New("choose_playlists.tmpl").Funcs(
		map[string]interface{}{
			"msToMin": func(ms int) int {
				return int(math.Round(float64(ms) / 60000))
			},

			"join": strings.Join,

			"removeArchivistDescriptionTag": func(description string) string {
				return strings.Replace(description, "Archivist: ", "", 1)
			},
		},
	).ParseFiles("llm/templates/choose_playlists.tmpl")

	if err != nil {
		return "", err
	}

	var prompt bytes.Buffer

	err = template.Execute(&prompt, struct {
		Playlists []spotify.Playlist
		Track     spotify.Track
	}{
		Playlists: playlists,
		Track:     track,
	})

	if err != nil {
		return "", err
	}

	return string(prompt.Bytes()), err
}
