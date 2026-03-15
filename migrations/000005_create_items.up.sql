CREATE TABLE IF NOT EXISTS items (
    api_name            TEXT NOT NULL,
    set_number          INT NOT NULL REFERENCES sets(number) ON DELETE CASCADE,
    name                TEXT NOT NULL,
    description         TEXT NOT NULL DEFAULT '',
    composition         TEXT[] NOT NULL DEFAULT '{}',
    effects             JSONB NOT NULL DEFAULT '{}',
    icon_url            TEXT NOT NULL DEFAULT '',
    associated_traits   TEXT[] NOT NULL DEFAULT '{}',
    incompatible_traits TEXT[] NOT NULL DEFAULT '{}',
    tags                TEXT[] NOT NULL DEFAULT '{}',
    is_unique           BOOLEAN NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (api_name, set_number)
);

CREATE INDEX idx_items_set_number ON items(set_number);
