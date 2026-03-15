CREATE TABLE IF NOT EXISTS traits (
    api_name    TEXT NOT NULL,
    set_number  INT NOT NULL REFERENCES sets(number) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    icon_url    TEXT NOT NULL DEFAULT '',
    effects     JSONB NOT NULL DEFAULT '[]',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (api_name, set_number)
);

CREATE INDEX idx_traits_set_number ON traits(set_number);
