-- name: GetCard :one
SELECT * FROM cards where id = $1 limit 1;

-- name: ListSerieIDs :many
SELECT id from series;

-- name: UpsertSeries :exec
INSERT INTO series (id, name, logo)
VALUES ($1, $2, $3)
ON CONFLICT (id)
DO UPDATE SET
	name = $2,
	logo = $3;

-- name: ListSetIDs :many
SELECT id from sets;

-- name: InsertSet :exec
INSERT INTO sets (id, name, logo, symbol, card_count_total, card_count_official, card_count_reverse, card_count_holo, card_count_first_edition, serie_id, tcg_online, release_date, legal_standard, legal_expanded)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);

-- name: ListSetCardIDs :many
SELECT id
FROM cards
WHERE set_id = $1;

-- name: InsertSetCard :exec
INSERT INTO
	cards (id, name, image_url, category, illustrator, rarity, updated_at, variant_normal, variant_reverse, variant_holo, variant_first_edition, attributes, set_id)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);

-- name: ListAllUncompleteCardSetIDs :many
SELECT
	s.id,
	s.name AS set_name,
	s.card_count_total,
	COUNT(c.*) AS created_total
FROM
	sets AS s
	LEFT JOIN cards AS c ON c.set_id = s.id
GROUP BY
	s.id
HAVING
	COUNT(c.*) != s.card_count_total
ORDER BY
	s.release_date DESC;

-- name: CreatePokemon :copyfrom
INSERT INTO pokemon (id, name) VALUES($1, $2);

-- name: QuerySimilarIllustrators :many
SELECT DISTINCT
	cards.illustrator,
	similarity(cards.illustrator, $1) as sim
FROM
	cards
WHERE
	cards.illustrator % $1
ORDER BY
	sim DESC
;

-- name: QuerySimilarPokemon :many
SELECT
	pokemon.name,
	similarity (pokemon.name, $1) AS sim
FROM
	pokemon
WHERE
	pokemon.name % $1
ORDER BY
	sim DESC;

-- name: QuerySimilarSets :many
SELECT
	sets.name,
	similarity (sets.name, $1) AS sim
FROM
	sets
WHERE
	sets.name % $1
ORDER BY
	sim DESC;

-- name: ListCardIllustrators :many
SELECT DISTINCT cards.illustrator FROM cards WHERE cards.illustrator IS NOT NULL;

-- name: ListPokemonNames :many
SELECT pokemon.id, pokemon.name from pokemon;

-- name: ListSetNames :many
SELECT sets.id, sets.name from sets;
