-- +goose Up
-- +goose StatementBegin
ALTER TABLE cards ADD COLUMN set_id TEXT;

-- Move the existing association into the card row before removing the
-- many-to-many join table. A card cannot belong to more than one set in the
-- new model.
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM set_cards
        GROUP BY card_id
        HAVING COUNT(*) > 1
    ) THEN
        RAISE EXCEPTION 'cannot migrate cards associated with multiple sets';
    END IF;
END $$;

UPDATE cards AS c
SET set_id = sc.set_id
FROM set_cards AS sc
WHERE sc.card_id = c.id;

ALTER TABLE cards
    ALTER COLUMN set_id SET NOT NULL,
    ADD CONSTRAINT cards_set_id_fkey
        FOREIGN KEY (set_id) REFERENCES sets(id) ON DELETE CASCADE;

DROP TABLE set_cards;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE set_cards (
    set_id TEXT NOT NULL REFERENCES sets(id) ON DELETE CASCADE,
    card_id TEXT NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
    PRIMARY KEY (set_id, card_id)
);

INSERT INTO set_cards (set_id, card_id)
SELECT set_id, id FROM cards;

ALTER TABLE cards
    DROP CONSTRAINT cards_set_id_fkey,
    DROP COLUMN set_id;
-- +goose StatementEnd
