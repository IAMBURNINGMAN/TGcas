package wallet

import (
	"database/sql"
	"errors"

	"tgcasino/internal/storage"
)

var ErrInsufficientFunds = errors.New("недостаточно сабинок")

// Debit — списать amount со счёта. Блокирует строку на время транзакции (FOR UPDATE).
func Debit(db *sql.DB, userID int64, amount int64, reason string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance int64
	err = tx.QueryRow(
		`SELECT balance FROM users WHERE id = $1 FOR UPDATE`, userID,
	).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < amount {
		return ErrInsufficientFunds
	}

	if _, err = tx.Exec(
		`UPDATE users SET balance = balance - $1 WHERE id = $2`, amount, userID,
	); err != nil {
		return err
	}

	if _, err = tx.Exec(
		`INSERT INTO transactions (user_id, amount, reason) VALUES ($1, $2, $3)`,
		userID, -amount, reason,
	); err != nil {
		return err
	}

	return tx.Commit()
}

// Credit — начислить amount на счёт (для возврата ставок).
func Credit(db *sql.DB, userID int64, amount int64, reason string) error {
	if err := storage.UpdateBalance(db, userID, amount); err != nil {
		return err
	}
	return storage.AddTransaction(db, userID, amount, reason)
}

// RecordResult — зачислить выигрыш (если won > 0) и обновить статистику атомично.
// wagered — сумма ставки, уже списанная через Debit.
func RecordResult(db *sql.DB, userID int64, wagered int64, won int64, reason string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if won > 0 {
		if _, err = tx.Exec(
			`UPDATE users SET balance = balance + $1 WHERE id = $2`, won, userID,
		); err != nil {
			return err
		}
		if _, err = tx.Exec(
			`INSERT INTO transactions (user_id, amount, reason) VALUES ($1, $2, $3)`,
			userID, won, reason,
		); err != nil {
			return err
		}
	}

	if err = storage.UpsertStats(tx, userID, wagered, won); err != nil {
		return err
	}

	return tx.Commit()
}
