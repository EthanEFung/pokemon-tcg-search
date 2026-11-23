-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS boosters (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,

    logo TEXT,
    artwork_front TEXT,
    artwork_back TEXT
);
CREATE TABLE IF NOT EXISTS set_boosters (
    set_id TEXT NOT NULL REFERENCES sets(id) ON DELETE CASCADE,
    booster_id TEXT NOT NULL REFERENCES boosters(id) ON DELETE CASCADE,

    PRIMARY KEY (set_id, booster_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS boosters;
DROP TABLE IF EXISTS set_boosters;
-- +goose StatementEnd
