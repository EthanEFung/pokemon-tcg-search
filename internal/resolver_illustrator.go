package internal

import (
	"context"
	"fmt"
	"sort"
)

type SimilarIllustratorResolver struct {
	NGrams    [][]Token
	Q         *Queries
	Stopwords map[string]struct{}
}

func NewSimilarIllustratorResolver(queries *Queries, nGrams [][]Token, stopwords map[string]struct{}) *SimilarIllustratorResolver {
	return &SimilarIllustratorResolver{
		NGrams:    nGrams,
		Q:         queries,
		Stopwords: stopwords,
	}
}

func (r SimilarIllustratorResolver) Resolve() ([]EntityCandidate, error) {
	ctx := context.Background()
	res := []EntityCandidate{}
	candidates := make(map[string]EntityCandidate, 0)
	rows, err := r.Q.ListCardIllustrators(ctx)
	if err != nil {
		return res, fmt.Errorf("queries.ListCardIllustrators err: %w", err)
	}

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
		for _, row := range rows {
			if !row.Valid {
				continue
			}
			ec := EntityCandidate{
				Type:      "illustrators",
				Candidate: row.String,
				Query:     phrase,
			}
			existing, ok := candidates[row.String]
			if !ok || LessCandidate(existing, ec) {
				candidates[row.String] = ec
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
