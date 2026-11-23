-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sets (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,

    logo TEXT,
    symbol TEXT,

    card_count_total INTEGER NOT NULL,
    card_count_official INTEGER NOT NULL,
    card_count_reverse INTEGER NOT NULL,
    card_count_holo INTEGER NOT NULL,
    card_count_first_edition INTEGER NOT NULL,

    serie_id TEXT REFERENCES series(id) ON DELETE SET NULL,

    tcg_online TEXT,

    release_date DATE NOT NULL,

    legal_standard BOOLEAN NOT NULL,
    legal_expanded BOOLEAN NOT NULL
);
CREATE TABLE IF NOT EXISTS set_cards (
    set_id TEXT NOT NULL REFERENCES sets(id) ON DELETE CASCADE,
    card_id TEXT NOT NULL REFERENCES cards(id) ON DELETE CASCADE,

    PRIMARY KEY (set_id, card_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS set_cards;
DROP TABLE IF EXISTS sets;
-- +goose StatementEnd
