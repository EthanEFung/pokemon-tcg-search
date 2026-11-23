-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS pg_trgm;

DO $$ 
	BEGIN
		CREATE TYPE card_category AS ENUM('pokemon', 'trainer', 'energy');
	EXCEPTION WHEN duplicate_object 
		THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS cards (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    image_url TEXT,
    category card_category NOT NULL DEFAULT 'pokemon',
    illustrator TEXT,
    rarity TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

	variant_normal BOOLEAN NOT NULL,
	variant_reverse BOOLEAN NOT NULL,
	variant_holo BOOLEAN NOT NULL,
	variant_first_edition BOOLEAN NOT NULL,

	attributes JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cards;
DROP TYPE IF EXISTS card_category;
DROP EXTENSION IF EXISTS pg_trgm;
-- +goose StatementEnd
