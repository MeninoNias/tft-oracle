CREATE TABLE matches (
    match_id            TEXT PRIMARY KEY,
    data_version        TEXT NOT NULL DEFAULT '',
    game_datetime       BIGINT NOT NULL DEFAULT 0,
    game_length         REAL NOT NULL DEFAULT 0,
    game_version        TEXT NOT NULL DEFAULT '',
    queue_id            INT NOT NULL DEFAULT 0,
    tft_set_number      INT NOT NULL DEFAULT 0,
    tft_game_type       TEXT NOT NULL DEFAULT '',
    participant_puuids  TEXT[] NOT NULL DEFAULT '{}',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_matches_game_datetime ON matches(game_datetime DESC);
