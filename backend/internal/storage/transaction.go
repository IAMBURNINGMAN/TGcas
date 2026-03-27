package storage

import "database/sql"

func AddTransaction(db *sql.DB, userID int64, amount int64, reason string) error {
	_, err := db.Exec(
		`INSERT INTO transactions (user_id, amount, reason) VALUES ($1, $2, $3)`,
		userID, amount, reason,
	)
	return err
}
