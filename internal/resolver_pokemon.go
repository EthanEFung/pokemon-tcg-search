package internal

import (
	"context"
	"fmt"
	"sort"
)

type SimilarPokemonResolver struct {
	NGrams    [][]Token
	Q         *Queries
	Stopwords map[string]struct{}
}

func NewSimilarPokemonResolver(queries *Queries, nGrams [][]Token, stopwords map[string]struct{}) *SimilarPokemonResolver {
	return &SimilarPokemonResolver{
		NGrams:    nGrams,
		Q:         queries,
		Stopwords: stopwords,
	}
}

func (r SimilarPokemonResolver) Resolve() ([]EntityCandidate, error) {
	ctx := context.Background()
	res := []EntityCandidate{}
	candidates := make(map[string]EntityCandidate)
	for _, ngram := range r.NGrams {
		tokens := []Token{}
		for _, t := range ngram {
			if _, ok := r.Stopwords[t.Text]; !ok {
				tokens = append(tokens, t)
			}
		}
		if len(tokens) == 0 {
			continue
		}
		phrase := ExtractPhrase(tokens)
		rows, err := r.Q.QuerySimilarPokemon(ctx, phrase)
		if err != nil {
			return nil, fmt.Errorf("QuerySimilarPokemon: %w", err)
		}
		for _, row := range rows {
			ec := EntityCandidate{
				Sim:       row.Sim,
				Type:      "pokemon",
				Candidate: row.Name,
				Query:     phrase,
			}
			existing, ok := candidates[row.Name]

			if !ok || LessCandidate(existing, ec) {
				candidates[row.Name] = ec
			}
		}
	}
	for _, c := range candidates {
		res = append(res, c)
	}
	sort.Slice(res, func(i, j int) bool {
		x, y := res[i], res[j]
		return !LessCandidate(x, y)
	})
	return res, nil
}
