package llm

import (
	"archivist/spotify"
	"context"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

func Complete(prompt string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_SECRET_KEY"))

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func GetPlaylistsToSaveTrackInto(playlists []spotify.Playlist, track spotify.Track) ([]spotify.Playlist, error) {
	prompt, err := ChoosePlaylistsPrompt(playlists, track)

	if err != nil {
		return nil, err
	}

	completion, err := Complete(prompt)

	if err != nil {
		return nil, err
	}

	playlistIDs := strings.Split(completion, "\n")

	var playlistsToSaveTrackInto []spotify.Playlist

	for _, p := range playlists {
		for _, id := range playlistIDs {
			if p.ID == id {
				playlistsToSaveTrackInto = append(playlistsToSaveTrackInto, p)
			}
		}
	}

	return playlistsToSaveTrackInto, nil
}
