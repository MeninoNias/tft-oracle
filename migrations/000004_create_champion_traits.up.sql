CREATE TABLE IF NOT EXISTS champion_traits (
    champion_api_name TEXT NOT NULL,
    trait_api_name    TEXT NOT NULL,
    set_number        INT NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (champion_api_name, trait_api_name, set_number),
    FOREIGN KEY (champion_api_name, set_number) REFERENCES champions(api_name, set_number) ON DELETE CASCADE,
    FOREIGN KEY (trait_api_name, set_number) REFERENCES traits(api_name, set_number) ON DELETE CASCADE
);

CREATE INDEX idx_champion_traits_champion ON champion_traits(champion_api_name, set_number);
CREATE INDEX idx_champion_traits_trait ON champion_traits(trait_api_name, set_number);
