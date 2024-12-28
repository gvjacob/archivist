# Archivist

**Archivist** creates smart AI-powered playlists for your Spotify. Housekeeping is almost an impossible task when you've got hundreds of playlists and little to no patience for organizing your favorite tracks. With Archivist, housekeeping is a one-click chore.

## How does it work?

Archivist lets you create "smart playlists". These playlists can describe, with natural language, the type of songs that should be included. When you save songs to your **Liked Songs** library, Archivist will automatically save them into the correct playlists based on their descriptions, which are editable via Spotify.

## Development

### Requirements

- [Go](https://go.dev/) 1.22+
- Spotify account and [API access](https://developer.spotify.com/)

### Getting Started

1. Copy `.env.sample` into `.env`.

   ```
   cp .env.sample .env
   ```

1. Fill in the environment variables in `.env`. You will need to register a Spotify application through the [API portal](https://developer.spotify.com/) and use the Spotify secrets.

1. Install dependencies

   ```
   go get
   ```

1. Initialize the program, which will create a local SQLite database, and seed it with your Spotify user data. Seeding user data may require a browser Spotify login, unless the `SPOTIFY_ACCESS_TOKEN` and `SPOTIFY_REFRESH_TOKEN` are provided in the `.env` file.

   ```
   go run main.go init
   ```

1. Once database has been created and seeded, running the program will find tracks saved into Spotify's **Liked Songs** library _since account creation_, check whether they match any Archivist playlists, and save them into the playlists if so.

   ```
   go run main.go
   ```
