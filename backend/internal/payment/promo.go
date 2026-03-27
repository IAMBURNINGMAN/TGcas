package payment

import (
	"database/sql"
	"errors"

	"tgcasino/internal/wallet"
)

var (
	ErrPromoNotFound = errors.New("промокод не найден")
	ErrPromoUsed     = errors.New("промокод уже использован")
)

func ApplyPromo(db *sql.DB, userID int64, code string) (int64, error) {
	var amount int64
	var usedBy sql.NullInt64

	err := db.QueryRow(
		`SELECT amount, used_by FROM promo_codes WHERE code = $1`, code,
	).Scan(&amount, &usedBy)
	if err == sql.ErrNoRows {
		return 0, ErrPromoNotFound
	}
	if err != nil {
		return 0, err
	}
	if usedBy.Valid {
		return 0, ErrPromoUsed
	}

	if _, err = db.Exec(
		`UPDATE promo_codes SET used_by = $1, used_at = NOW() WHERE code = $2`,
		userID, code,
	); err != nil {
		return 0, err
	}

	return amount, wallet.Credit(db, userID, amount, "promo:"+code)
}
