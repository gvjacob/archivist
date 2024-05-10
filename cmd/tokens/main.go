package main

import (
	"archivist/utils"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

func main() {
	utils.LoadDotEnv()

	tokens, err := getSpotifyAccessToken()

	if err != nil {
		log.Fatal(err)
	}

	log.Info("Tokens", "Access Token", tokens.AccessToken, "Refresh Token", tokens.RefreshToken)
}

type SpotifyAuthTokens struct {
	AccessToken  string
	RefreshToken string
}

func getSpotifyAccessToken() (SpotifyAuthTokens, error) {
	auth := url.Values{}

	auth.Set("grant_type", "authorization_code")
	auth.Set("code", os.Getenv("SPOTIFY_AUTHORIZATION_CODE"))
	auth.Set("redirect_uri", "http://localhost:8080/callback")

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(auth.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+getBase64EncodedSpotifyAuth())

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return SpotifyAuthTokens{}, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return SpotifyAuthTokens{}, errors.New(string(body))
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	accessToken := result["access_token"]
	refreshToken := result["refresh_token"]

	return SpotifyAuthTokens{
		AccessToken:  accessToken.(string),
		RefreshToken: refreshToken.(string),
	}, nil
}

func getBase64EncodedSpotifyAuth() string {
	bytes := []byte(os.Getenv("SPOTIFY_CLIENT_ID") + ":" + os.Getenv("SPOTIFY_CLIENT_SECRET"))
	return b64.URLEncoding.EncodeToString(bytes)
}
