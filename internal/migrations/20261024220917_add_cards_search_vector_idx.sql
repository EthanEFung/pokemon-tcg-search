-- +goose Up
-- +goose StatementBegin
ALTER TABLE cards
ADD COLUMN tsv tsvector GENERATED ALWAYS AS (
	to_tsvector(
		'simple',
		COALESCE(cards.illustrator, '') || '''s ' || COALESCE(cards.name) || ' from ' || COALESCE(cards.attributes::jsonb #>> '{set,name}', '') || '. ' || COALESCE(cards.attributes::jsonb ->> 'description', '') || ' Types: ' || COALESCE(cards.attributes::jsonb ->> 'types', '')
	)
) STORED;

CREATE INDEX tsv_idx ON cards USING GIN (tsv);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX tsv_idx;

ALTER TABLE cards
DROP COLUMN tsv;
-- +goose StatementEnd
