CREATE TABLE raw_tier_data (
    id               SERIAL PRIMARY KEY,
    source           VARCHAR(50) NOT NULL,
    patch            VARCHAR(20) NOT NULL,
    composition_name VARCHAR(255) NOT NULL,
    tier             VARCHAR(5),
    win_rate         DECIMAL(5,2),
    play_rate        DECIMAL(5,2),
    avg_placement    DECIMAL(3,2),
    champion_ids     TEXT[] NOT NULL DEFAULT '{}',
    core_items       JSONB NOT NULL DEFAULT '{}',
    scraped_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_raw_tier_source ON raw_tier_data(source);
CREATE INDEX idx_raw_tier_patch ON raw_tier_data(patch);
CREATE INDEX idx_raw_tier_source_patch ON raw_tier_data(source, patch);
