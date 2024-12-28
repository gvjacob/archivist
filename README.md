# Archivist

**Archivist** creates smart AI-powered playlists for your Spotify. Housekeeping is almost an impossible task when you've got hundreds of playlists and little to no patience for organizing your favorite tracks. With Archivist, housekeeping is a one-click chore.

**⚠️ This is an experimental AI program currently suited to my use case. Expect bugs and missing features!**

## How does it work?

Archivist lets you create "smart playlists" with the help of GPT-4. These playlists can describe, with natural language, the type of songs that should be included. When you save songs to your **Liked Songs** library, Archivist will automatically save them into the correct playlists based on their descriptions, which are editable via Spotify. Since the Spotify API has no concept of webhooks, Archivist should be run as a cron job.

## Creating Archivist playlists

Playlists are marked for Archivist if their description starts with the `Archivist:` prefix. Here's an example description for a playlist for Arcane music:

```
Archivist: Favorite tracks from the Netflix show Arcane League of Legends
```

Generally, the more descriptive the descriptions are, the more accurate Archivist will be when filing tracks. **If a track matches multiple playlists, it will be filed to all of them.**

Archivist keeps a record of tracks that it has queried so it knows when there are new tracks saved into **Liked Songs**. These are the only tracks that are queried against for the next run.

## How does Archivist decide?

Currently, Archivist can only view tracks' information from the [**Get Several Tracks** endpoint](https://developer.spotify.com/documentation/web-api/reference/get-several-tracks). This includes the following properties (with examples):

```
Name: Let It Be
Artists: The Beatles
Album: Let It Be (Remastered)
Album Type: album
Album Release Date: 1970
Genres: Rock, Classic Rock
Duration Minutes: 4
Explicit: false
```

Playlists are referenced by their name and description:

```
Name: Beatlesmania
Description: Best of Beatles
```

### Spotify Audio Features

Spotify analyzes useful properties like a track's `energy`, `danceability` and `tempo`. These could've added extra dimensions to how Archivist matches tracks to playlists. Unforunately, Spotify has **deprecated** (womp womp) the [Audio Features endpoint](https://developer.spotify.com/documentation/web-api/reference/get-audio-analysis).

## Development

### Requirements

- [Go](https://go.dev/) 1.22+
- [SQLite](https://www.sqlite.org/)
- Spotify account and [API access](https://developer.spotify.com/)
- OpenAI account

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
