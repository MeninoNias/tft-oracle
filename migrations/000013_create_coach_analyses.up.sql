-- Coach analysis cache tables.
-- Stores AI coaching results to avoid re-analyzing the same data.

CREATE TABLE coach_match_analyses (
    match_id    TEXT NOT NULL,
    puuid       TEXT NOT NULL,
    response    JSONB NOT NULL,
    model       VARCHAR(50) NOT NULL DEFAULT 'gpt-4o-mini',
    tokens_used INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (match_id, puuid)
);

CREATE TABLE coach_history_analyses (
    puuid       TEXT NOT NULL,
    match_count INT NOT NULL,
    match_ids   TEXT[] NOT NULL DEFAULT '{}',
    response    JSONB NOT NULL,
    model       VARCHAR(50) NOT NULL DEFAULT 'gpt-4o-mini',
    tokens_used INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (puuid, match_count)
);

CREATE INDEX idx_coach_match_puuid ON coach_match_analyses(puuid);
CREATE INDEX idx_coach_history_created ON coach_history_analyses(created_at);
