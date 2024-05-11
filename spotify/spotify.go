package spotify

import (
	"bytes"
	"net/http"
)

type SpotifyClient struct {
	AccessToken string
}

func NewSpotifyClient(accessToken string) *SpotifyClient {
	return &SpotifyClient{AccessToken: accessToken}
}

func (s *SpotifyClient) GetAuthorizationHeader() string {
	return "Bearer " + s.AccessToken
}

func (s *SpotifyClient) Get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", s.GetAuthorizationHeader())

	client := &http.Client{}
	return client.Do(req)
}

func (s *SpotifyClient) Post(endpoint string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", s.GetAuthorizationHeader())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}
