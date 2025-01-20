package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"os"
	"slices"
	"strings"
)

type HighlightDatabase struct {
	Highlights []Highlight
	TokenMap   map[int]Token
}

type Token struct {
	ID         int
	Content    string
	NextTokens []NextToken
}

func (t Token) NominateNextTokens(count int) []int {
	var result []int
	mark := make([]bool, len(t.NextTokens))
	for len(result) < len(t.NextTokens) && len(result) < count {
		target := rand.Float64()
		lo := 0
		hi := len(t.NextTokens) - 1
		for lo < hi && mark[lo] {
			lo = lo + 1
		}
		for lo < hi && mark[hi] {
			hi = hi - 1
		}
		for lo < hi {
			mid := (lo + hi) / 2
			if t.NextTokens[mid].CumulativeProbability < target {
				lo = mid + 1
			} else if t.NextTokens[mid].CumulativeProbability > target {
				hi = mid
			}
		}
		if !mark[lo] {
			result = append(result, t.NextTokens[lo].ID)
			mark[lo] = true
		}
	}
	return result
}

type NextToken struct {
	ID                    int
	CumulativeProbability float64
}

type Highlight struct {
	ID      string
	Content string
	Tokens  []int
}

func ReadHighlights(fname string) (HighlightDatabase, error) {
	bs, err := os.ReadFile(fname)
	if err != nil {
		return HighlightDatabase{}, err
	}

	var result []Highlight
	allTokens := make(map[string]int)
	tokenMap := make(map[int]string)
	resultTokenMap := make(map[int]Token)

	nextToken := make(map[int]map[int]int)

	for _, line := range strings.Split(string(bs), "---") {
		line = strings.TrimSpace(line)
		tokens := tokenizeString(line)
		id := generateID(tokens)
		var lineTokens []int

		for _, token := range tokens {
			tokenID, ok := allTokens[token]
			if !ok {
				tokenID = len(allTokens)
				allTokens[token] = tokenID
				tokenMap[tokenID] = token
				nextToken[tokenID] = make(map[int]int)
			}
			lineTokens = append(lineTokens, tokenID)
		}
		for i := 1; i < len(tokens); i++ {
			prevToken := allTokens[tokens[i-1]]
			token := allTokens[tokens[i]]
			nextToken[prevToken][token] = nextToken[prevToken][token] + 1
		}

		result = append(result, Highlight{
			ID:      id,
			Content: strings.TrimSpace(line),
			Tokens:  lineTokens,
		})
	}

	for tokenID, tokenContent := range tokenMap {
		currentNextToken := nextToken[tokenID]
		total := 0
		for _, count := range currentNextToken {
			total = total + count
		}

		var nextTokens []NextToken
		for nextTokenID, nextTokenCount := range currentNextToken {
			nextTokens = append(nextTokens, NextToken{
				ID:                    nextTokenID,
				CumulativeProbability: float64(nextTokenCount) / float64(total),
			})
		}
		slices.SortFunc(nextTokens, func(left, right NextToken) int {
			if left.CumulativeProbability < right.CumulativeProbability {
				return -1
			}
			if left.CumulativeProbability > right.CumulativeProbability {
				return 1
			}
			return 0
		})
		for i := 1; i < len(nextTokens); i++ {
			nextTokens[i].CumulativeProbability = nextTokens[i].CumulativeProbability + nextTokens[i-1].CumulativeProbability
		}

		token := Token{
			ID:         tokenID,
			Content:    tokenContent,
			NextTokens: nextTokens,
		}
		resultTokenMap[tokenID] = token
	}

	return HighlightDatabase{
		TokenMap:   resultTokenMap,
		Highlights: result,
	}, nil
}

func generateID(tokens []string) string {
	var buff bytes.Buffer
	for _, token := range tokens {
		buff.WriteString(token)
		buff.WriteString(" ")
	}

	h := sha256.New()
	h.Write(buff.Bytes())
	sum := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}

func tokenizeString(str string) []string {
	replacer := strings.NewReplacer(
		"'s", "s",
		".", " ",
		"'", " ",
		"*", "",
		",", " ",
		":", " ",
		"\n", " ",
		"\r", " ",
		"\t", " ",
		";", " ",
		"\"", " ",
		"-", " ",
		"_", " ",
		"=", " ",
		"+", " ",
		"#", " ",
		"(", " ",
		")", " ",
		"[", " ",
		"]", " ",
		"{", " ",
		"}", " ",
	)

	str = replacer.Replace(str)
	parts := strings.Split(strings.ToLower(str), " ")

	var result []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}
