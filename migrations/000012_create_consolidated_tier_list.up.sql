CREATE TABLE consolidated_tier_list (
    id                  SERIAL PRIMARY KEY,
    patch               VARCHAR(20) NOT NULL,
    composition_name    VARCHAR(255) NOT NULL,
    consolidated_tier   VARCHAR(5) NOT NULL,
    consolidated_score  DECIMAL(5,2) NOT NULL,
    confidence          VARCHAR(20) NOT NULL DEFAULT 'low',
    consensus           VARCHAR(30) NOT NULL DEFAULT 'unknown',
    metatft_tier        VARCHAR(5),
    tftactics_tier      VARCHAR(5),
    mobalytics_tier     VARCHAR(5),
    avg_win_rate        DECIMAL(5,2),
    avg_play_rate       DECIMAL(5,2),
    avg_placement       DECIMAL(3,2),
    core_champions      TEXT[] NOT NULL DEFAULT '{}',
    recommended_items   JSONB NOT NULL DEFAULT '{}',
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (patch, composition_name)
);

CREATE INDEX idx_consolidated_patch ON consolidated_tier_list(patch);
CREATE INDEX idx_consolidated_tier ON consolidated_tier_list(consolidated_tier);
