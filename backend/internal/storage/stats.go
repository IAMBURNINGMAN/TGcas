package storage

import "database/sql"

// UserStats — статистика игрока для личного кабинета.
type UserStats struct {
	GamesPlayed  int64
	Wins         int64
	TotalWagered int64
	TotalWon     int64
	BestWin      int64
	MemberSince  string
}

// GetStats — получить статистику пользователя.
func GetStats(db *sql.DB, userID int64) (*UserStats, error) {
	s := &UserStats{}
	err := db.QueryRow(`
		SELECT
			COALESCE(gs.games_played, 0),
			COALESCE(gs.wins, 0),
			COALESCE(gs.total_wagered, 0),
			COALESCE(gs.total_won, 0),
			COALESCE(gs.best_win, 0),
			TO_CHAR(u.created_at, 'DD.MM.YYYY')
		FROM users u
		LEFT JOIN game_stats gs ON gs.user_id = u.id
		WHERE u.id = $1`, userID,
	).Scan(&s.GamesPlayed, &s.Wins, &s.TotalWagered, &s.TotalWon, &s.BestWin, &s.MemberSince)
	return s, err
}

// UpsertStats — обновить статистику внутри существующей транзакции.
// won > 0 = победа, won == 0 = проигрыш.
func UpsertStats(tx *sql.Tx, userID int64, wagered int64, won int64) error {
	var winCount int64
	if won > 0 {
		winCount = 1
	}
	_, err := tx.Exec(`
		INSERT INTO game_stats (user_id, games_played, wins, total_wagered, total_won, best_win)
		VALUES ($1, 1, $2, $3, $4, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			games_played  = game_stats.games_played + 1,
			wins          = game_stats.wins + EXCLUDED.wins,
			total_wagered = game_stats.total_wagered + EXCLUDED.total_wagered,
			total_won     = game_stats.total_won + EXCLUDED.total_won,
			best_win      = GREATEST(game_stats.best_win, EXCLUDED.best_win)`,
		userID, winCount, wagered, won,
	)
	return err
}
