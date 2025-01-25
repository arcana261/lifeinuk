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
	Highlights      []Highlight
	TokenMap        map[int]Token
	UnmatchedScores map[string]Score
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
	for id, score := range db.UnmatchedScores {
		if score.Count < 1 {
			continue
		}

		buff.WriteString(fmt.Sprintf("%s %f %d\n", id, score.Sum, score.Count))
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
	nexts := cloneSlice(t.NextTokens)
	nexts = removeFunc(nexts, func(nt NextToken) bool {
		return db.TokenMap[nt.ID].SkipPuzzle
	})
	for len(result) < count && len(nexts) > 0 {
		i := lowerBoundFunc(nexts, rand.Float64(), func(nt NextToken, f float64) int {
			if nt.CumulativeProbability+1e-5 < f {
				return -1
			}
			if nt.CumulativeProbability-1e-5 > f {
				return 1
			}
			return 0
		})
		result = append(result, nexts[i].ID)
		nexts = removeAt(nexts, i)
	}
	return result
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
	Index                 int
}

type Score struct {
	Sum     float64
	Count   int
	Average float64
}

func WriteHighlights(db HighlightDatabase, fname string) {
	var order []int
	for i := 0; i < len(db.Highlights); i++ {
		order = append(order, i)
	}
	slices.SortFunc(order, func(a, b int) int {
		av := db.Highlights[a].Index
		bv := db.Highlights[b].Index
		cv := av - bv
		if cv < 0 {
			return -1
		}
		if cv > 0 {
			return 1
		}
		return 0
	})

	var buff bytes.Buffer
	for i := 0; i < len(db.Highlights); i++ {
		if i > 0 {
			buff.WriteString("\n\n---\n\n")
		}
		buff.WriteString(db.Highlights[order[i]].Content)
	}
	err := os.WriteFile(fname, buff.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
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
			Index:   len(result),
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

	unmatchedScores := make(map[string]Score)
	for id, score := range readScores(scores) {
		index, ok := highlightIDToIndex[id]
		if !ok {
			unmatchedScores[id] = score
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
		TokenMap:        resultTokenMap,
		Highlights:      result,
		UnmatchedScores: unmatchedScores,
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
