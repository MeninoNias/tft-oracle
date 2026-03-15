CREATE TABLE IF NOT EXISTS sets (
    number      INT PRIMARY KEY,
    name        TEXT NOT NULL,
    mutator     TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
