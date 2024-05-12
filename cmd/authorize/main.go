package main

import (
	"archivist/utils"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
)

const (
	localAddress = "localhost:8080"
	redirectURI  = "http://localhost:8080/callback"
)

// Retrieve authorization code from Spotify
// and exchange it for access tokens.
func main() {
	utils.LoadDotEnv()

	params := url.Values{}
	params.Set("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	params.Set("response_type", "code")
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", "playlist-read-private playlist-modify-private playlist-modify-public user-library-read")

	url := "https://accounts.spotify.com/authorize?" + params.Encode()
	exec.Command("open", url).Run()

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    localAddress,
		Handler: mux,
	}

	var code string

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code = r.URL.Query().Get("code")
		go server.Shutdown(context.Background())
	})

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}

	if code == "" {
		log.Fatal("No authorization code found")
	}

	log.Info("Retrieved authorization code")
	tokens, err := getSpotifyAccessToken(code)

	if err != nil {
		log.Fatal(err)
	}

	log.Info("Tokens", "Access Token", tokens.AccessToken, "Refresh Token", tokens.RefreshToken)
}

type SpotifyAuthTokens struct {
	AccessToken  string
	RefreshToken string
}

func getSpotifyAccessToken(code string) (SpotifyAuthTokens, error) {
	auth := url.Values{}

	auth.Set("grant_type", "authorization_code")
	auth.Set("code", code)
	auth.Set("redirect_uri", redirectURI)

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(auth.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", utils.GetBasicAuthorizationHeader())

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
