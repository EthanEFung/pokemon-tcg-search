package internal

import (
	"fmt"
	"strings"
	"unicode"
)

type Resolver interface {
	Resolve() ([]EntityCandidate, error)
}

type EntityCandidate struct {
	Type      string
	Query     string
	Candidate string
}

type TokenAlignment struct {
	Matching []string
	Missing  []string
	Extra    []string
}

type TokenMap = map[string]int

func (ec EntityCandidate) String() string {
	return fmt.Sprintf(
		"Candidate: %v\n\tType: %s\n\tQuery: %s\n",
		ec.Candidate,
		ec.Type,
		ec.Query,
	)
}

func (ec *EntityCandidate) Alignment() TokenAlignment {
	result := TokenAlignment{
		Matching: []string{},
		Missing:  []string{},
		Extra:    []string{},
	}

	cTokens := tokenize(ec.Candidate)
	qTokens := tokenize(ec.Query)

	union := make(TokenMap)
	for t, cCount := range cTokens {
		if qCount, has := qTokens[t]; has {
			union[t] = qCount
			if cCount < qCount {
				union[t] = cCount
			}
		} else {
			result.Missing = append(result.Missing, t)
		}
	}

	for t := range union {
		result.Matching = append(result.Matching, t)
	}

	for t := range qTokens {
		if _, has := cTokens[t]; has {
			result.Extra = append(result.Extra, t)
		}
	}

	return result
}

/*
LessCandidate

Let's utilize first principle thinking to break down this problem
into the fundamental building blocks.

LessCandidate compares two entity candidates and reports which one
is less of a match.

What is a match?
A match is when the candidate:
  - is the phrase that is queried
  - has resemblence to the queried phrase

How do you know if the candidate is the queried phrase?
  - The phrase is broken into normalized tokens, and contain all the
    same tokens.

```
candidate: Ken Sugimori Yusuke Ohmura     -> ken, ohmura, sugimori, yusuke
query:     yusuke ohmura and Ken Sugimori -> ken, ohmura, sugimori, yusuke
is_exact:                                 -> yes,    yes,      yes,    yes
```

How do you know if a candidate has resemblence to a queried phrase?
  - when comparing the tokens of the candidate and query, the levenshtein
    distance of a token is under a certain threshold

```
candidate: Ken Sugimori Yusuke Ohmura -> ken, ohmura, sugimori, yusuke
query: sugimory                       ->              sugimory
is_exact:                             ->  no,     no,       no,     no
distance:                             ->   8,      8,        1,      8
```

What is a lesser match?
A lesser match is one that:

  - does not match at all, and the other does

  - contains less candidate matching and candidate missing tokens than the next

    ```
    A
    candidate: ken sugimori -> ken, sugimori
    query: sugimori ken     -> ken, sugimori
    matching tokens: 2

    B
    candidate: Ken Sugimori Yusuke Ohmura -> ken, ohmura, sugimori, yusuke
    query: Ken Sugimory and Yusuke Ohmura -> ken, ohmura, sugimory, yusuke
    matching tokens: 4

    ```

    A < B because A has less query tokens matching the candidate tokens than B
    A has 2 matching and B has 4

    ---

    A
    candidate: ken sugimori -> ken, sugimori
    query: sugimori ken     -> ken, sugimori
    matching tokens: 2
    missing tokens: 0

    B
    candidate: Ken Sugimori Yusuke Ohmura -> ken, ohmura, sugimori, yusuke
    query: Sugimory and Ohmura            ->      ohmura, sugimory
    matching tokens: 2
    missing tokens: 2

    ```

    A < B because while A and B both have the same number of query tokens that
    are matching their respective candidates, B's candidate tokens are greater
    than the query tokens, while A's query does not miss any of the candidate
    tokens, and therefore while being more accurate, may not encapsulate the
    queries entire meaning as well as B.

  - contains more extra tokens if candidate matching + candidate missing
    tokens are equal

    ```
    A
    candidate: Ken Sugimori Yusuke Ohmura ->        ken, ohmura, sugimori, yusuke
    query: charz Sugimory and Ohmura      -> charz,      ohmura, sugimory
    matching tokens: 2
    missing tokens: 2
    extra tokens: 1

    B
    candidate: Ken Sugimori Yusuke Ohmura -> ken, ohmura, sugimori, yusuke
    query: Sugimory and Ohmura            ->      ohmura, sugimory
    matching tokens: 2
    missing tokens: 2

    ```

    The goal is to capture as much query meaning which is roughly equivalent to
    how much of the span is captured, and how accurately the span represents a
    candidate. Therefore A < B because in considering the same number of tokens
    that match the candidate, and miss the candidate, A also includes extra
    tokens that are not relevent to the candidate.

    Definitions:
    Missing tokens: the unmatched candidate tokens
    Extra tokens: the unmatched query tokens
*/
func LessCandidate(x, y EntityCandidate) bool {
}

func normalize(word string) string {
	word = strings.ToLower(word)
	word = strings.Map(func(ch rune) rune {
		if unicode.IsPunct(ch) {
			return -1
		}
		return ch
	}, word)
	return word
}

func tokenize(phrase string) TokenMap {
	words := strings.Fields(phrase)
	for i := range words {
		words[i] = normalize(words[i])
	}
	result := make(TokenMap)
	for _, word := range words {
		if _, ok := result[word]; !ok {
			result[word] = 0
		}
		result[word] += 1
	}
	return result
}
