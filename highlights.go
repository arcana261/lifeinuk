package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"os"
	"slices"
	"strings"
)

const (
	UntrustedNext = 0.9
)

type HighlightDatabase struct {
	Highlights []Highlight
	TokenMap   map[int]Token
}

func (db HighlightDatabase) PickHighlight() *Highlight {
	if len(db.Highlights) == 0 {
		return nil
	}

	minCount := -1
	for _, h := range db.Highlights {
		if minCount < 0 || h.Score.Count < minCount {
			minCount = h.Score.Count
		}
	}

	lo := 0
	hi := len(db.Highlights) - 1
	for lo < hi && db.Highlights[lo].Score.Count != minCount {
		lo = lo + 1
	}
	for lo < hi && db.Highlights[hi].Score.Count != minCount {
		hi = hi - 1
	}
	if lo == hi {
		lo = 0
		hi = len(db.Highlights) - 1
	}

	startLo := lo
	startHi := hi

	var result *Highlight

	for result == nil {
		target := rand.Float64()
		lo = startLo
		hi = startHi

		for lo < hi {
			mid := (lo + hi) / 2
			if db.Highlights[mid].CumulativeProbability < target {
				lo = mid + 1
			} else if db.Highlights[mid].CumulativeProbability > target {
				hi = mid
			}
		}

		if db.Highlights[lo].Score.Count == minCount {
			result = &db.Highlights[lo]
		}
	}

	return result
}

func (db HighlightDatabase) WriteScore(fname string) error {
	var buff bytes.Buffer

	for _, h := range db.Highlights {
		if h.Score.Count < 1 {
			continue
		}

		buff.WriteString(fmt.Sprintf("%s %f %d\n", h.ID, h.Score.Sum, h.Score.Count))
	}

	return os.WriteFile(fname, buff.Bytes(), 0644)
}

type Token struct {
	ID                int
	Content           string
	NextTokens        []NextToken
	AppearInNextToken int
	SkipPuzzle        bool
}

func (t Token) NominateNextTokens(db HighlightDatabase, count int) []int {
	var result []int
	mark := make([]bool, len(t.NextTokens))
	skipped := 0
	for (skipped+len(result)) < len(t.NextTokens) && (skipped+len(result)) < count {
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
			mark[lo] = true

			if !db.TokenMap[t.NextTokens[lo].ID].ShouldSkipPuzzle() {
				result = append(result, t.NextTokens[lo].ID)
			} else {
				skipped = skipped + 1
			}
		}
	}
	return result
}

func (t Token) ShouldSkipPuzzle() bool {
	return t.SkipPuzzle
}

type NextToken struct {
	ID                    int
	CumulativeProbability float64
}

type Highlight struct {
	ID                    string
	Content               string
	Tokens                []int
	Score                 Score
	CumulativeProbability float64
}

type Score struct {
	Sum     float64
	Count   int
	Average float64
}

func ReadHighlights(fname string, scores string) (HighlightDatabase, error) {
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

	var tokenIDs []int
	for tokenID, token := range resultTokenMap {
		for _, next := range token.NextTokens {
			temp := resultTokenMap[next.ID]
			temp.AppearInNextToken += 1
			resultTokenMap[next.ID] = temp
		}
		tokenIDs = append(tokenIDs, tokenID)
	}
	slices.SortFunc(tokenIDs, func(x, y int) int {
		xv := resultTokenMap[x].AppearInNextToken
		yv := resultTokenMap[y].AppearInNextToken
		if xv < yv {
			return -1
		}
		if xv > yv {
			return 1
		}
		return 0
	})
	if len(tokenIDs) > 0 {
		threshold := max(0, int(math.Floor(float64(len(tokenIDs))*UntrustedNext)))
		for j := len(tokenIDs) - 1; j >= threshold; j-- {
			temp := resultTokenMap[tokenIDs[j]]
			temp.SkipPuzzle = true
			resultTokenMap[tokenIDs[j]] = temp
		}
	}

	highlightIDToIndex := make(map[string]int)
	for i := 0; i < len(result); i++ {
		highlightIDToIndex[result[i].ID] = i
	}

	for id, score := range readScores(scores) {
		index, ok := highlightIDToIndex[id]
		if !ok {
			continue
		}
		result[index].Score = score
	}

	var sumScore float64
	for i := 0; i < len(result); i++ {
		f := 1.0 - result[i].Score.Average
		result[i].CumulativeProbability = f
		sumScore = sumScore + f
	}
	for i := 0; i < len(result); i++ {
		result[i].CumulativeProbability = result[i].CumulativeProbability / sumScore
	}
	slices.SortFunc(result, func(x, y Highlight) int {
		if x.CumulativeProbability < y.CumulativeProbability {
			return -1
		}
		if x.CumulativeProbability > y.CumulativeProbability {
			return -1
		}
		return 0
	})
	for i := 1; i < len(result); i++ {
		result[i].CumulativeProbability = result[i].CumulativeProbability + result[i-1].CumulativeProbability
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
		"\u2018", " ",
		"\u2019", " ",
		"\u2013", " ",
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

func readScores(fname string) map[string]Score {
	result := make(map[string]Score)

	b, err := os.ReadFile(fname)
	if err != nil {
		return result
	}

	for _, line := range strings.Split(string(b), "\n") {
		parts := strings.Split(strings.TrimSpace(line), " ")
		if len(parts) != 3 {
			continue
		}
		id := strings.TrimSpace(parts[0])
		var totalSum float64
		var count int
		n, err := fmt.Sscanf(strings.TrimSpace(parts[1]), "%f", &totalSum)
		if n != 1 || err != nil {
			continue
		}
		n, err = fmt.Sscanf(strings.TrimSpace(parts[2]), "%d", &count)
		if n != 1 || err != nil {
			continue
		}

		result[id] = Score{
			Sum:     totalSum,
			Count:   count,
			Average: totalSum / (float64(count)*float64(count+1)/float64(2) + 1),
		}
	}

	return result
}
