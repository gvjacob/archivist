package storage

type UsersTable struct {
	*Database
}

type User struct {
	ID           string
	AccessToken  string
	RefreshToken string
	CreatedAt    int
}

func NewUsers(db *Database) *UsersTable {
	return &UsersTable{Database: db}
}

func (u *UsersTable) GetUser() (*User, error) {
	user := &User{}
	err := u.QueryRow("SELECT * FROM users LIMIT 1").Scan(&user.ID, &user.AccessToken, &user.RefreshToken, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UsersTable) UpdateUser(user *User) error {
	_, err := u.Exec(
		"UPDATE users SET access_token = $1, refresh_token = $2 WHERE id = $3",
		user.AccessToken,
		user.RefreshToken,
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
