package main

import (
	"archivist/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
	_ "modernc.org/sqlite"
)

const (
	localAddress = "localhost:3000"
	redirectURI  = "http://localhost:3000/callback"
)

type SpotifyAuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Initialize() {
	if err := createTables(); err != nil {
		log.Fatal(err)
	}

	if err := seedUserData(); err != nil {
		log.Fatal(err)
	}

	log.Info("Successfully seeded user data")
}

func createTables() error {
	log.Info("Creating tables...")

	const createTablesQuery = `
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      access_token TEXT NOT NULL,
      refresh_token TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS archived_tracks (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      user_id INTEGER NOT NULL,
      track_id TEXT NOT NULL,
      playlist_id TEXT NULL,
      created_at INTEGER NOT NULL,
      FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
    );
  `

	db, err := sql.Open("sqlite", os.Getenv("SQLITE_FILE_PATH"))

	if err != nil {
		return err
	}

	if _, err := db.Exec(createTablesQuery); err != nil {
		return err
	}

	return nil
}

func seedUserData() error {
	log.Info("Seeding user data...")

	authTokens, err := retrieveAuthTokensFromDotEnv()

	if err == nil {
		log.Info("Retrieved Spotify access and refresh tokens from environment")
	} else {
		log.Warn("No Spotify access or refresh token found from environment. Authorizing...")
		authTokens, err = authorize()
	}

	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite", os.Getenv("SQLITE_FILE_PATH"))

	if err != nil {
		return err
	}

	insertUserQuery := `
    DELETE FROM users;
    INSERT INTO users (access_token, refresh_token) VALUES (?, ?);
  `

	_, err = db.Exec(insertUserQuery, authTokens.AccessToken, authTokens.RefreshToken)

	if err != nil {
		return err
	}

	return nil
}

func retrieveAuthTokensFromDotEnv() (*SpotifyAuthTokens, error) {
	accessToken := os.Getenv("SPOTIFY_ACCESS_TOKEN")
	refreshToken := os.Getenv("SPOTIFY_REFRESH_TOKEN")

	if accessToken == "" || refreshToken == "" {
		return nil, errors.New("No Spotify access or refresh token found")
	}

	return &SpotifyAuthTokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func authorize() (*SpotifyAuthTokens, error) {
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
		return nil, err
	}

	if code == "" {
		return nil, errors.New("No authorization code found")
	}

	log.Info("Retrieved authorization code")
	tokens, err := getSpotifyAccessToken(code)

	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func getSpotifyAccessToken(code string) (*SpotifyAuthTokens, error) {
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
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	var tokens SpotifyAuthTokens
	json.Unmarshal([]byte(body), &tokens)

	return &tokens, nil
}
