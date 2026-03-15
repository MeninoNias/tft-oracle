CREATE TABLE IF NOT EXISTS champions (
    api_name        TEXT NOT NULL,
    set_number      INT NOT NULL REFERENCES sets(number) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    cost            INT NOT NULL DEFAULT 0,
    role            TEXT NOT NULL DEFAULT '',
    stats           JSONB NOT NULL DEFAULT '{}',
    ability         JSONB NOT NULL DEFAULT '{}',
    icon_url        TEXT NOT NULL DEFAULT '',
    square_icon_url TEXT NOT NULL DEFAULT '',
    tile_icon_url   TEXT NOT NULL DEFAULT '',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (api_name, set_number)
);

CREATE INDEX idx_champions_set_number ON champions(set_number);
CREATE INDEX idx_champions_cost ON champions(cost);
