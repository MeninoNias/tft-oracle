CREATE TABLE players (
    puuid           TEXT PRIMARY KEY,
    game_name       TEXT NOT NULL DEFAULT '',
    tag_line        TEXT NOT NULL DEFAULT '',
    region          TEXT NOT NULL DEFAULT '',
    platform        TEXT NOT NULL DEFAULT '',
    summoner_id     TEXT NOT NULL DEFAULT '',
    profile_icon_id INT NOT NULL DEFAULT 0,
    summoner_level  BIGINT NOT NULL DEFAULT 0,
    tier            TEXT NOT NULL DEFAULT '',
    rank            TEXT NOT NULL DEFAULT '',
    league_points   INT NOT NULL DEFAULT 0,
    wins            INT NOT NULL DEFAULT 0,
    losses          INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_players_riot_id ON players(game_name, tag_line);
