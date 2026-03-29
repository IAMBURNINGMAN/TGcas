CREATE TABLE game_stats (
    user_id       BIGINT PRIMARY KEY REFERENCES users(id),
    games_played  BIGINT NOT NULL DEFAULT 0,
    wins          BIGINT NOT NULL DEFAULT 0,
    total_wagered BIGINT NOT NULL DEFAULT 0,
    total_won     BIGINT NOT NULL DEFAULT 0,
    best_win      BIGINT NOT NULL DEFAULT 0
);
