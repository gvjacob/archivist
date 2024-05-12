package storage

import (
	"time"
)

type UsersTable struct {
	*Database
}

type User struct {
	ID           string
	AccessToken  string
	RefreshToken string
	LastArchived time.Time
}

func NewUsers(db *Database) *UsersTable {
	return &UsersTable{Database: db}
}

func (u *UsersTable) GetUser() (*User, error) {
	user := &User{LastArchived: time.Now().UTC()}

	var lastArchived string
	err := u.db.QueryRow("SELECT * FROM users LIMIT 1").Scan(&user.ID, &user.AccessToken, &user.RefreshToken, &lastArchived)

	if err != nil {
		return nil, err
	}

	if lastArchived == "" {
		return user, nil
	}

	if timestamp, err := time.Parse(time.RFC3339, lastArchived); err == nil {
		user.LastArchived = timestamp
	}

	return user, nil
}

func (u *UsersTable) UpdateUser(user *User) error {
	_, err := u.db.Exec(
		"UPDATE users SET access_token = $1, refresh_token = $2, last_archived = $3 WHERE id = $4",
		user.AccessToken,
		user.RefreshToken,
		user.LastArchived.Format(time.RFC3339),
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
