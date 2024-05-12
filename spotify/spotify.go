package spotify

import (
	"archivist/storage"
	"archivist/utils"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
)

type SpotifyClient struct {
	User       *storage.User
	UsersTable *storage.UsersTable
}

func NewSpotifyClient(user *storage.User, userTable *storage.UsersTable) *SpotifyClient {
	return &SpotifyClient{User: user, UsersTable: userTable}
}

func (s *SpotifyClient) GetBearerAuthorizationHeader() string {
	return "Bearer " + s.User.AccessToken
}

func (s *SpotifyClient) Get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", s.GetBearerAuthorizationHeader())

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode == 401 {
		log.Warn("Refreshing token")
		s.RefreshToken()
		return s.Get(endpoint)
	}

	return resp, err
}

func (s *SpotifyClient) Post(endpoint string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", s.GetBearerAuthorizationHeader())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode == 401 {
		s.RefreshToken()
		return s.Post(endpoint, body)
	}

	return resp, err
}

func (s *SpotifyClient) RefreshToken() error {
	body := []byte("grant_type=refresh_token&refresh_token=" + s.User.RefreshToken)
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewReader(body))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", utils.GetBasicAuthorizationHeader())

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Failed to refresh token")
	}

	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	json.NewDecoder(resp.Body).Decode(&tokenResponse)
	s.User.AccessToken = tokenResponse.AccessToken

	if err := s.UsersTable.UpdateUser(s.User); err != nil {
		return err
	}

	return nil
}
