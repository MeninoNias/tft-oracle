CREATE TABLE match_participants (
    match_id                TEXT NOT NULL REFERENCES matches(match_id) ON DELETE CASCADE,
    puuid                   TEXT NOT NULL,
    placement               INT NOT NULL DEFAULT 0,
    level                   INT NOT NULL DEFAULT 0,
    gold_left               INT NOT NULL DEFAULT 0,
    last_round              INT NOT NULL DEFAULT 0,
    time_eliminated         REAL NOT NULL DEFAULT 0,
    total_damage_to_players INT NOT NULL DEFAULT 0,
    players_eliminated      INT NOT NULL DEFAULT 0,
    augments                TEXT[] NOT NULL DEFAULT '{}',
    companion               JSONB NOT NULL DEFAULT '{}',
    traits                  JSONB NOT NULL DEFAULT '[]',
    units                   JSONB NOT NULL DEFAULT '[]',
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (match_id, puuid)
);

CREATE INDEX idx_match_participants_puuid ON match_participants(puuid);
