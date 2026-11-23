-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pokemon (
	id INTEGER PRIMARY KEY, 
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS pokemon_cards (
	pokemon_id INTEGER REFERENCES pokemon (id) ON DELETE CASCADE, 
	card_id TEXT REFERENCES cards (id) ON DELETE CASCADE,
	PRIMARY KEY (pokemon_id, card_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pokemon_cards;

DROP TABLE IF EXISTS pokemon;
-- +goose StatementEnd
