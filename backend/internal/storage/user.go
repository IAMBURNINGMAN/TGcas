package storage

import "database/sql"

type User struct {
	ID       int64
	Username string
	Balance  int64
}

func GetOrCreate(db *sql.DB, id int64, username string) (*User, error) {
	u := &User{}
	err := db.QueryRow(
		`INSERT INTO users (id, username) VALUES ($1, $2)
		 ON CONFLICT (id) DO UPDATE SET username = EXCLUDED.username
		 RETURNING id, username, balance`,
		id, username,
	).Scan(&u.ID, &u.Username, &u.Balance)
	return u, err
}

func GetBalance(db *sql.DB, userID int64) (int64, error) {
	var balance int64
	err := db.QueryRow(`SELECT balance FROM users WHERE id = $1`, userID).Scan(&balance)
	return balance, err
}

func UpdateBalance(db *sql.DB, userID int64, delta int64) error {
	_, err := db.Exec(
		`UPDATE users SET balance = balance + $1 WHERE id = $2`,
		delta, userID,
	)
	return err
}
