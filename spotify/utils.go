package spotify

import (
	"net/http"
	"os"
)

func Get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("SPOTIFY_ACCESS_TOKEN"))

	client := &http.Client{}
	return client.Do(req)
}
