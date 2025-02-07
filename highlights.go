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

	"github.com/arcana261/lifeinuk/maputils"
	"github.com/arcana261/lifeinuk/sliceutils"
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
	items := sliceutils.Range(0, len(db.Highlights))
	minCount := sliceutils.MinFunc(db.Highlights, func(h1, h2 Highlight) int {
		return h1.Score.Count - h2.Score.Count
	}).Score.Count
	items = sliceutils.FilterFunc(items, func(idx int) bool {
		return db.Highlights[idx].Score.Count == minCount
	})
	if len(items) == 0 {
		return nil
	}
	target := rand.Float64() * db.Highlights[items[len(items)-1]].CumulativeProbability
	at := sliceutils.LowerBoundSortedFunc(items, func(idx int) int {
		return CompareFloat64(db.Highlights[idx].CumulativeProbability, target)
	})
	if at < 0 {
		return nil
	}
	return &db.Highlights[items[at]]
}

func (db HighlightDatabase) WriteScore(fname string) error {
	lines := sliceutils.MapFunc(
		sliceutils.FilterFunc(db.Highlights, func(h Highlight) bool {
			return h.Score.Count > 0
		}), func(h Highlight) string {
			return fmt.Sprintf("%s %f %d\n", h.ID, h.Score.Sum, h.Score.Count)
		},
	)
	lines = append(lines,
		sliceutils.MapFunc(
			maputils.ToEntries(db.UnmatchedScores), func(p sliceutils.Pair[string, Score]) string {
				return fmt.Sprintf("%s %f %d\n", p.Key, p.Value.Sum, p.Value.Count)
			},
		)...,
	)
	lines = sliceutils.Sort(lines)

	return os.WriteFile(fname, []byte(strings.Join(lines, "")), 0644)
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
	nexts := sliceutils.Clone(t.NextTokens)
	nexts = sliceutils.RemoveFunc(nexts, func(nt NextToken) bool {
		return db.TokenMap[nt.ID].SkipPuzzle
	})
	for len(result) < count && len(nexts) > 0 {
		target := rand.Float64() * nexts[len(nexts)-1].CumulativeProbability
		i := sliceutils.LowerBoundSortedFunc(nexts, func(nt NextToken) int {
			return CompareFloat64(nt.CumulativeProbability, target)
		})
		if i < 0 {
			fmt.Println(nexts)
		}
		result = append(result, nexts[i].ID)
		nexts = sliceutils.RemoveAt(nexts, i)
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

	entries := sliceutils.Remove(
		sliceutils.TrimSpace(
			strings.Split(string(bs), "---"),
		), "",
	)
	entryTokens := sliceutils.MapFunc(entries, tokenizeString)
	allTokens := sliceutils.ToMapFunc2(
		sliceutils.UniqueSorted(
			sliceutils.Sort(
				sliceutils.Flatten(entryTokens),
			),
		), func(index int, s string) (string, int, bool) {
			return s, index, true
		})
	entryTokensMapped := sliceutils.MapFunc(entryTokens, func(s []string) []int {
		return sliceutils.Lookup(s, allTokens)
	})
	tokenMap := maputils.MapFunc(allTokens, func(t string, id int) (int, string) {
		return id, t
	})
	result := sliceutils.MapFunc2(sliceutils.Zip(entries, entryTokens), func(index int, item sliceutils.Pair[string, []string]) (Highlight, bool) {
		entry := item.Key
		tokens := item.Value
		id := generateID(tokens)

		return Highlight{
			ID:      id,
			Content: entry,
			Index:   index,
			Tokens:  sliceutils.Lookup(tokens, allTokens),
		}, true
	})

	resultTokenMap := maputils.MapFunc(
		sliceutils.ToMultiMapFunc(
			sliceutils.Flatten(
				sliceutils.MapFunc(entryTokensMapped, func(tokens []int) []sliceutils.Pair[int, int] {
					if len(tokens) == 0 {
						return nil
					}
					return sliceutils.Zip(tokens, tokens[1:])
				}),
			), func(p sliceutils.Pair[int, int]) (int, int) {
				return p.Key, p.Value
			},
		), func(t int, nexts []int) (int, Token) {
			nextTokens := sliceutils.AccumulateFunc2(
				sliceutils.MapFunc(
					maputils.ToEntries(
						maputils.MapValuesFunc(
							sliceutils.AccumulateFunc2(nexts, make(map[int]int), func(prev map[int]int, _ int, current int) (map[int]int, bool) {
								prev[current] = prev[current] + 1
								return prev, true
							}), func(count int) float64 {
								return float64(count) / float64(len(nexts))
							},
						),
					), func(p sliceutils.Pair[int, float64]) NextToken {
						return NextToken{
							ID:                    p.Key,
							CumulativeProbability: p.Value,
						}
					},
				), make([]NextToken, 0), func(prev []NextToken, index int, current NextToken) ([]NextToken, bool) {
					if index != 0 {
						current.CumulativeProbability = current.CumulativeProbability + prev[index-1].CumulativeProbability
					}
					prev = append(prev, current)
					return prev, true
				},
			)

			return t, Token{
				ID:         t,
				Content:    tokenMap[t],
				NextTokens: nextTokens,
			}
		},
	)
	for key, value := range tokenMap {
		_, ok := resultTokenMap[key]
		if !ok {
			resultTokenMap[key] = Token{
				ID:         key,
				Content:    value,
				NextTokens: nil,
			}
		}
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
		if _, ok := highlightIDToIndex[result[i].ID]; ok {
			fmt.Println("DUPLICAT ID FOUND!:\n========\n%s\n========\n", result[i].Content)
		}
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
		"0,0", "00",
		"1,0", "10",
		"1,4", "14",
		"2,0", "20",
		"3,0", "30",
		"4,0", "40",
		"5,0", "50",
		"8,0", "80",
		"\u2018s", "s",
		"\u2019s", "s",
		"\u2013s", "s",
		"\"s", "s",
	)
	str = replacer.Replace(str)

	replacer = strings.NewReplacer(
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
		"!", " ",
		"\u2018", " ",
		"\u2019", " ",
		"\u2013", " ",
	)
	str = replacer.Replace(str)

	str = strings.ToLower(str)
	parts := strings.Split(str, " ")
	parts = sliceutils.TrimSpace(parts)
	parts = sliceutils.Remove(parts, "")

	return parts
}

func readScores(fname string) map[string]Score {
	if !fileExists(fname) {
		return nil
	}

	b, err := os.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(b), "\n")
	lines = sliceutils.TrimSpace(lines)
	parts := sliceutils.Split(lines, " ")
	parts = sliceutils.FilterFunc(parts, func(s []string) bool { return len(s) == 3 })
	pairs := sliceutils.MapFunc(parts, func(part []string) sliceutils.Pair[string, Score] {
		id := strings.TrimSpace(part[0])
		var totalSum float64
		var count int
		n, err := fmt.Sscanf(strings.TrimSpace(part[1]), "%f", &totalSum)
		if n != 1 || err != nil {
			return sliceutils.Pair[string, Score]{}
		}
		n, err = fmt.Sscanf(strings.TrimSpace(part[2]), "%d", &count)
		if n != 1 || err != nil {
			return sliceutils.Pair[string, Score]{}
		}

		return sliceutils.Pair[string, Score]{
			Key: id,
			Value: Score{
				Sum:     totalSum,
				Count:   count,
				Average: totalSum / (float64(count)*float64(count+1)/float64(2) + 1),
			},
		}
	})
	pairs = sliceutils.FilterFunc(pairs, func(p sliceutils.Pair[string, Score]) bool {
		return p.Key != ""
	})
	result := sliceutils.ToMapFunc(pairs, func(p sliceutils.Pair[string, Score]) (string, Score) {
		return p.Key, p.Value
	})
	return result
}
