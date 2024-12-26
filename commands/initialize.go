package commands

import (
	"archivist/storage"
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
	_ "modernc.org/sqlite"
)

type client struct {
	database *storage.Database
}

type spotifyAuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Initialize() {
	db, _ := storage.NewDatabase(os.Getenv("SQLITE_FILE_PATH"))
	c := &client{database: db}

	if err := c.createTables(); err != nil {
		log.Fatal(err)
	}

	log.Info("Successfully created tables")

	if err := c.seedUserData(); err != nil {
		log.Fatal(err)
	}

	log.Info("Successfully seeded user data")
}

func (c client) createTables() error {
	log.Info("Creating tables...")

	const createTablesSQL = `
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

	if _, err := c.database.Exec(createTablesSQL); err != nil {
		return err
	}

	return nil
}

func (c client) seedUserData() error {
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

	insertUserSQL := `
    DELETE FROM users;
    INSERT INTO users (access_token, refresh_token) VALUES (?, ?);
  `

	_, err = c.database.Exec(insertUserSQL, authTokens.AccessToken, authTokens.RefreshToken)

	if err != nil {
		return err
	}

	return nil
}

func retrieveAuthTokensFromDotEnv() (*spotifyAuthTokens, error) {
	accessToken := os.Getenv("SPOTIFY_ACCESS_TOKEN")
	refreshToken := os.Getenv("SPOTIFY_REFRESH_TOKEN")

	if accessToken == "" || refreshToken == "" {
		return nil, errors.New("No Spotify access or refresh token found")
	}

	return &spotifyAuthTokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func authorize() (*spotifyAuthTokens, error) {
	params := url.Values{}
	params.Set("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	params.Set("response_type", "code")
	params.Set("redirect_uri", os.Getenv("SPOTIFY_REDIRECT_URI"))
	params.Set("scope", "playlist-read-private playlist-modify-private playlist-modify-public user-library-read")

	url := "https://accounts.spotify.com/authorize?" + params.Encode()
	exec.Command("open", url).Run()

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    os.Getenv("LOCAL_ADDRESS"),
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

func getSpotifyAccessToken(code string) (*spotifyAuthTokens, error) {
	auth := url.Values{}

	auth.Set("grant_type", "authorization_code")
	auth.Set("code", code)
	auth.Set("redirect_uri", os.Getenv("SPOTIFY_REDIRECT_URI"))

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

	var tokens spotifyAuthTokens
	json.Unmarshal([]byte(body), &tokens)

	return &tokens, nil
}
