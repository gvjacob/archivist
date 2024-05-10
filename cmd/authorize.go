package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/lpernett/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	params := url.Values{}
	params.Set("client_id", os.Getenv("SPOTIFY_CLIENT_ID"))
	params.Set("response_type", "code")
	params.Set("redirect_uri", "http://localhost:8080/callback")
	params.Set("scope", "playlist-read-private playlist-modify-private playlist-modify-public user-library-read")

	url := "https://accounts.spotify.com/authorize?" + params.Encode()

	exec.Command("open", url).Run()

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		fmt.Fprintf(w, "Authorization Code:\n%s", code)
	})

	http.ListenAndServe(":8080", nil)
}
