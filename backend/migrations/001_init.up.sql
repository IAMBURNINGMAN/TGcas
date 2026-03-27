CREATE TABLE users (
    id         BIGINT PRIMARY KEY,
    username   TEXT,
    balance    BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE transactions (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id),
    amount     BIGINT NOT NULL,
    reason     TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE promo_codes (
    code       TEXT PRIMARY KEY,
    amount     BIGINT NOT NULL,
    used_by    BIGINT REFERENCES users(id),
    used_at    TIMESTAMPTZ
);
