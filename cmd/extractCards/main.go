package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/ethanefung/pokemon-tcg-search/internal"
	"github.com/ethanefung/pokemon-tcg-search/internal/tcgdex"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=ethanfung dbname=pokemon-tcg-search sslmode=verify-full")
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer conn.Close(ctx)
	queries := internal.New(conn)

	api := tcgdex.NewTCGDex()
	sets, err := queries.ListAllUncompleteCardSetIDs(ctx)
	if err != nil {
		log.Fatalf("queries.ListAllUncompleteCardSetIDs err: %v\n", err)
	}
	for _, row := range sets {
		setID := row.ID
		existingCardIdList, err := queries.ListSetCardIDs(ctx, setID)
		if err != nil {
			log.Printf("queries.ListSetCardIDs err: %v\n", err)
			continue
		}
		existingCardIds := make(map[string]struct{}, len(existingCardIdList))
		for _, cardID := range existingCardIdList {
			existingCardIds[cardID] = struct{}{}
		}
		cardSet, err := api.GetSet(setID)
		if err != nil {
			log.Printf("api.GetSet err: %v\n", err)
		}
		for _, cardBrief := range cardSet.Cards {
			if _, ok := existingCardIds[cardBrief.ID]; ok {
				continue
			}
			log.Printf("Processing %s of %s: %s\n", cardBrief.Name, cardSet.Name, cardBrief.ID)
			b, err := api.GetCard(cardBrief.ID)
			if err != nil {
				log.Printf("api.GetCard err: %v\n", err)
				continue
			}
			var common tcgdex.TCGDexCardCommon
			if err := json.Unmarshal(b, &common); err != nil {
				log.Fatalf("json.Unmarshal err: %v\n", err)
			}
			updatedAt, err := time.Parse(time.RFC3339, common.UpdatedAt)
			if err != nil {
				log.Fatalf("time.Parse err: %v\n", err)
			}
			cardParams := internal.InsertSetCardParams{
				ID:                  common.ID,
				Name:                common.Name,
				Category:            internal.CardCategory(strings.ToLower(common.Category)),
				VariantNormal:       common.Variants.Normal,
				VariantReverse:      common.Variants.Reverse,
				VariantHolo:         common.Variants.Holo,
				VariantFirstEdition: common.Variants.FirstEd,
				SetID:               setID,
				UpdatedAt: pgtype.Timestamptz{
					Time:  updatedAt,
					Valid: true,
				},
			}
			if common.ImageURL != nil {
				cardParams.ImageUrl = pgtype.Text{
					String: *common.ImageURL,
					Valid:  true,
				}
			}
			if common.Illustrator != nil {
				cardParams.Illustrator = pgtype.Text{
					String: *common.Illustrator,
					Valid:  true,
				}
			}
			if common.Rarity != nil {
				cardParams.Rarity = pgtype.Text{
					String: *common.Rarity,
					Valid:  true,
				}
			}
			var attrs []byte
			switch common.Category {
			case "Pokemon":
				var pkmn tcgdex.TCGDexCardPokemon
				if err := json.Unmarshal(b, &pkmn); err != nil {
					log.Fatalf("json.Unmarshal err: %v\n", err)
				}
				attrs, err = json.Marshal(pkmn)
			case "Energy":
				var energy tcgdex.TCGDexCardEnergy
				if err := json.Unmarshal(b, &energy); err != nil {
					log.Fatalf("json.Unmarshal err: %v\n", err)
				}
				attrs, err = json.Marshal(energy)
			case "Trainer":
				var trainer tcgdex.TCGDexCardTrainer
				if err := json.Unmarshal(b, &trainer); err != nil {
					log.Fatalf("json.Unmarshal err: %v\n", err)
				}
				attrs, err = json.Marshal(trainer)
			default:
				log.Fatalf("unhandled category '%s'", common.Category)
			}
			if err != nil {
				log.Fatalf("attrs json.Marshal err: %v\n", err)
			}
			cardParams.Attributes = attrs
			err = queries.InsertSetCard(ctx, cardParams)
			if err != nil {
				log.Fatalf("queries.InsertSetCard err: %v\n", err)
			}
			log.Println("Inserted")
		}
	}
	log.Println("OK")
}
