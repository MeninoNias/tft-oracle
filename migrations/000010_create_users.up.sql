CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    access_key_hash TEXT NOT NULL UNIQUE,
    riot_puuid      TEXT NOT NULL DEFAULT '',
    game_name       TEXT NOT NULL DEFAULT '',
    tag_line        TEXT NOT NULL DEFAULT '',
    region          TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_puuid ON users(riot_puuid);
