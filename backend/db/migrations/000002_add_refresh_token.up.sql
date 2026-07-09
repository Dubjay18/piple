CREATE TABLE refresh_tokens (
    id          TEXT        PRIMARY KEY,       -- random UUID
    user_id     TEXT        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT        NOT NULL UNIQUE,   -- hash of the token, not the token itself
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked     BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);