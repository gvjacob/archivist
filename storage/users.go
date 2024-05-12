package storage

type UsersTable struct {
	*Database
}

type User struct {
	ID           string
	AccessToken  string
	RefreshToken string
}

func NewUsers(db *Database) *UsersTable {
	return &UsersTable{Database: db}
}

func (u *UsersTable) GetUser() (*User, error) {
	user := &User{}
	err := u.db.QueryRow("SELECT * FROM users LIMIT 1").Scan(&user.ID, &user.AccessToken, &user.RefreshToken)

	if err != nil {
		return nil, err
	}

	return user, nil
}
