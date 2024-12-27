package storage

import "time"

type ArchivedTracksTable struct {
	*Database
}

type ArchivedTrack struct {
	UserID     string
	PlaylistID string
}

func NewArchivedTracks(db *Database) *ArchivedTracksTable {
	return &ArchivedTracksTable{Database: db}
}

func (a *ArchivedTracksTable) Insert(tracks []ArchivedTrack) error {
	createdAt := time.Now().Unix()

	tx, err := a.Begin()

	if err != nil {
		return err
	}

	for _, track := range tracks {
		_, err := tx.Exec(`
      INSERT INTO archived_tracks (user_id, playlist_id, created_at) VALUES ($1, $2, $3)
    `, track.UserID, track.PlaylistID, createdAt)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
