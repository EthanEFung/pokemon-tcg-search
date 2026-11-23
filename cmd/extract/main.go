package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ethanefung/pokemon-tcg-search/internal"
	"github.com/ethanefung/pokemon-tcg-search/internal/tcgdex"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func main() {
	// Setup dependencies
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=ethanfung dbname=pokemon-tcg-search sslmode=verify-full")
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer conn.Close(ctx)

	queries := internal.New(conn)

	api := tcgdex.NewTCGDex()
	// Dependencies are set

	serieBriefs, err := api.GetSeries()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	for _, brief := range serieBriefs {
		params := internal.UpsertSeriesParams{
			ID:   brief.ID,
			Name: brief.Name,
		}
		if brief.LogoURL != nil {
			params.Logo = pgtype.Text{
				String: *brief.LogoURL,
				Valid:  true,
			}
		}
		queries.UpsertSeries(ctx, params)
		fmt.Printf("Upserted Series '%s'\n", params.Name)
	}

	// Start with a "cache" layer. By querying to see what sets have been stored
	setIDs, err := queries.ListSetIDs(ctx)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	knownSetIDs := make(map[string]struct{})
	for _, setID := range setIDs {
		knownSetIDs[setID] = struct{}{}
	}

	// then ask the api for a list of all sets
	setBriefs, err := api.GetSets()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	// store the set ids that we need to retrieve and store
	newSetIDs := make(map[string]struct{})
	for _, brief := range setBriefs {
		if _, has := knownSetIDs[brief.ID]; !has {
			newSetIDs[brief.ID] = struct{}{}
		}
	}

	for newSetID := range newSetIDs {
		// ask the api for the full details of the set
		tcgdexSet, err := api.GetSet(newSetID)
		if err != nil {
			log.Fatalf("%v\n", err)
		}

		// store in our database the information about the new set.
		setParams := internal.InsertSetParams{
			ID:                    tcgdexSet.ID,
			Name:                  tcgdexSet.Name,
			CardCountTotal:        int32(tcgdexSet.CardCount.Total),
			CardCountOfficial:     int32(tcgdexSet.CardCount.Official),
			CardCountReverse:      int32(tcgdexSet.CardCount.Reverse),
			CardCountHolo:         int32(tcgdexSet.CardCount.Holo),
			CardCountFirstEdition: int32(tcgdexSet.CardCount.FirstEd),
			LegalStandard:         tcgdexSet.Legal.Standard,
			LegalExpanded:         tcgdexSet.Legal.Expanded,
		}
		fmt.Printf("series %v\n", tcgdexSet.Serie)
		setParams.SerieID = pgtype.Text{
			String: tcgdexSet.Serie.ID,
			Valid:  true,
		}
		if tcgdexSet.LogoURL != nil {
			setParams.Logo = pgtype.Text{
				String: *tcgdexSet.LogoURL,
				Valid:  true,
			}
		}
		if tcgdexSet.SymbolURL != nil {
			setParams.Symbol = pgtype.Text{
				String: *tcgdexSet.SymbolURL,
				Valid:  true,
			}
		}
		if tcgdexSet.TCGOnline != nil {
			setParams.TcgOnline = pgtype.Text{
				String: *tcgdexSet.TCGOnline,
				Valid:  true,
			}
		}

		releaseDate, err := time.Parse("2006-01-02", tcgdexSet.ReleaseDate)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		setParams.ReleaseDate = pgtype.Date{
			Time:  releaseDate,
			Valid: tcgdexSet.ReleaseDate != "",
		}

		err = queries.InsertSet(ctx, setParams)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		fmt.Printf("Inserted Set '%s'\n", setParams.Name)

		existingCardIDList, err := queries.ListSetCardIDs(ctx, setParams.ID)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		existingCardIDs := make(map[string]struct{}, len(existingCardIDList))
		for _, cardID := range existingCardIDList {
			existingCardIDs[cardID] = struct{}{}
		}

		// loop over the card briefs to find any non-existent cards in our db
		for _, cardBrief := range tcgdexSet.Cards {
			fmt.Printf("checking %s\n", cardBrief.Name)
			if _, ok := existingCardIDs[cardBrief.ID]; ok {
				continue
			}
			// request data from the api
			b, err := api.GetCard(cardBrief.ID)
			if err != nil {
				log.Fatalf("%v\n", err)
			}
			fmt.Printf("Requested %s\n", cardBrief.ID)
			var common tcgdex.TCGDexCardCommon
			if err := json.Unmarshal(b, &common); err != nil {
				log.Fatalf("%v\n", err)
			}
			updatedAt, err := time.Parse(time.RFC3339, common.UpdatedAt)
			if err != nil {
				log.Fatalf("%v\n", err)
			}
			cardParams := internal.InsertSetCardParams{
				ID:                  common.ID,
				Name:                common.Name,
				Category:            internal.CardCategory(strings.ToLower(common.Category)),
				VariantNormal:       common.Variants.Normal,
				VariantReverse:      common.Variants.Reverse,
				VariantHolo:         common.Variants.Holo,
				VariantFirstEdition: common.Variants.FirstEd,
				SetID:               setParams.ID,
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
					log.Fatalf("%v\n", err)
				}
				attrs, err = json.Marshal(pkmn)
			case "Energy":
				var energy tcgdex.TCGDexCardEnergy
				if err := json.Unmarshal(b, &energy); err != nil {
					log.Fatalf("%v\n", err)
				}
				attrs, err = json.Marshal(energy)
			case "Trainer":
				var trainer tcgdex.TCGDexCardTrainer
				if err := json.Unmarshal(b, &trainer); err != nil {
					log.Fatalf("%v\n", err)
				}
				attrs, err = json.Marshal(trainer)
			default:
				log.Fatalf("unhandled category '%s'", common.Category)
			}
			if err != nil {
				log.Fatalf("%v\n", err)
			}
			cardParams.Attributes = attrs
			err = queries.InsertSetCard(ctx, cardParams)
			if err != nil {
				log.Fatalf("%v\n", err)
			}
			fmt.Printf("Inserted %s\n", cardParams.Name)
		}
	}

	// Now we're getting somewhere!

	fmt.Println("OK")
}
