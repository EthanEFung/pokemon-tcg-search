package cli

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/ethanefung/pokemon-tcg-search/internal"
	"github.com/jackc/pgx/v5"
)

//go:embed stopwords.txt
var stopwordsRaw string

type Token struct {
	Text  string
	Start int
	End   int
}

type Candidate struct {
	Start      int
	End        int
	TokenCount int
	Confidence float64
	EntityID   any
	Kind       string
	Value      string
}

type Row[T float32] struct {
	Sim T
}

type Querier[T float32] interface {
	Query() (Row[T], error)
}

func Search(query string) error {
	fmt.Println("Searching...")
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "user=ethanfung dbname=pokemon-tcg-search sslmode=verify-full")
	if err != nil {
		return fmt.Errorf("db connection err: %w", err)
	}
	defer conn.Close(ctx)

	queries := internal.New(conn)

	// break up the query string into tokens
	words := strings.Fields(query)

	start, end := -1, -1
	tokens := []internal.Token{}

	for _, word := range words {
		l := len(word)
		start, end = end+1, end+1+l
		token := internal.Token{
			Text:  word,
			Start: start,
			End:   end,
		}

		tokens = append(tokens, token)
	}

	ngrams := nGrams(tokens)
	stopwords := loadStopwords()

	// pkmnResolver := internal.NewSimilarPokemonResolver(queries, ngrams, stopwords)
	illustratorResolver := internal.NewSimilarIllustratorResolver(queries, ngrams, stopwords)
	// setResolver := internal.NewSimilarSetResolver(queries, ngrams, stopwords)

	// candidates, err := pkmnResolver.Resolve()
	// if err != nil {
	// 	return fmt.Errorf("SimilarPokemonResolver: %w", err)
	// }
	// fmt.Printf("%d pokemon candidate(s):\n", len(candidates))
	//
	// for _, c := range candidates {
	// 	fmt.Println(c)
	// }

	candidates, err := illustratorResolver.Resolve()
	if err != nil {
		return fmt.Errorf("SimilarIllustratorResolver: %w", err)
	}
	fmt.Printf("%d illustrator candidate(s):\n", len(candidates))

	for _, c := range candidates {
		fmt.Println(c)
	}

	// candidates, err = setResolver.Resolve()
	// if err != nil {
	// 	return fmt.Errorf("SimilarSetResolver: %w", err)
	// }
	// fmt.Printf("%d set candidate(s):\n", len(candidates))
	// for _, c := range candidates {
	// 	fmt.Println(c)
	// }

	// scan the remaining tokens to get any leftover text

	return nil
}

func loadStopwords() map[string]struct{} {
	words := strings.Fields(stopwordsRaw)
	smap := make(map[string]struct{}, len(words))
	for _, word := range words {
		smap[word] = struct{}{}
	}
	return smap
}

func nGrams(tokens []internal.Token) [][]internal.Token {
	result := [][]internal.Token{}
	for num := 1; num <= len(tokens); num++ {
		for i := 0; i+num <= len(tokens); i++ {
			nGram := tokens[i : i+num]
			result = append(result, nGram)
		}
	}
	return result
}
