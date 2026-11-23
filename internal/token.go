package internal

import "strings"

type Token struct {
	Start int
	End   int
	Text  string
}

func ExtractPhrase(tokens []Token) string {
	words := make([]string, len(tokens))
	for i, t := range tokens {
		words[i] = t.Text
	}
	return strings.Join(words, " ")
}
