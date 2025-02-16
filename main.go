package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/arcana261/lifeinuk/sliceutils"
)

const (
	colorRed    = "\033[0;31m"
	colorGreen  = "\033[0;32m"
	colorYellow = "\033[0;33m"
	colorBlue   = "\033[0;34m"
	colorNone   = "\033[0m"

	alignmentWidth    = 60
	puzzleChoiceCount = 4
)

func main() {
	highlights, err := ReadHighlights("data/highlights.txt", "scores.txt")
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(highlights.Highlights); i++ {
		highlights.Highlights[i].Content = fixAlignment(highlights.Highlights[i].Content, alignmentWidth)
	}
	copyFile("data/highlights.txt", "highlights.txt.bak")
	WriteHighlights(highlights, "data/highlights.txt")

	for {
		fmt.Printf("\n")
		fmt.Printf("  1. Print Random Card\n")
		fmt.Printf("  2. Fill Card Game\n")
		fmt.Printf("  Q. Quit\n")
		fmt.Printf("\n")

		cmd := strings.ToLower(readOne())
		switch cmd {
		case "q":
			return
		case "1":
			printRandomCard(highlights)
		case "2":
			fillCard(highlights)
		default:
		}
	}
}

func fillCard(highlights HighlightDatabase) {
	h := highlights.PickHighlight()
	hContent := []rune(h.Content)

	correctAnswers := 0
	wrongAnswers := 0
	lastI := -1
	var previousWrongs []int

	for i := 2; i < len(h.Tokens); i++ {
		currentToken := highlights.TokenMap[h.Tokens[i]]

		if currentToken.SkipPuzzle {
			continue
		}

		var skips []string
		skips = append(skips, currentToken.Content)
		for _, w := range previousWrongs {
			skips = append(skips, highlights.TokenMap[w].Content)
		}

		if strings.HasSuffix(currentToken.Content, "ies") {
			skips = append(skips, currentToken.Content[:len(currentToken.Content)-3])
		} else if strings.HasSuffix(currentToken.Content, "y") {
			skips = append(skips, fmt.Sprintf("%sies", currentToken.Content[:len(currentToken.Content)-1]))
		}

		if strings.HasSuffix(currentToken.Content, "er") {
			skips = append(skips, fmt.Sprintf("%sing", currentToken.Content[:len(currentToken.Content)-2]))
		} else if strings.HasSuffix(currentToken.Content, "ing") {
			skips = append(skips, fmt.Sprintf("%ser", currentToken.Content[:len(currentToken.Content)-3]))
		}

		if strings.HasSuffix(currentToken.Content, "s") {
			skips = append(skips, currentToken.Content[:len(currentToken.Content)-1])
		} else {
			skips = append(skips, fmt.Sprintf("%ss", currentToken.Content))
		}

		if strings.HasSuffix(currentToken.Content, "ation") {
			skips = append(skips, currentToken.Content[:len(currentToken.Content)-5])
		} else {
			skips = append(skips, fmt.Sprintf("%sation", currentToken.Content))
		}
		if strings.HasSuffix(currentToken.Content, "ration") {
			skips = append(skips, fmt.Sprintf("%ser", currentToken.Content[:len(currentToken.Content)-6]))
		} else if strings.HasSuffix(currentToken.Content, "er") {
			skips = append(skips, fmt.Sprintf("%sration", currentToken.Content[:len(currentToken.Content)-2]))
		}

		nextTokens := highlights.TokenMap[h.Tokens[i-1]].NominateNextTokens(
			highlights,
			puzzleChoiceCount-1,
			skips...,
		)
		sliceutils.Permutate(previousWrongs)
		for j := 0; j < len(previousWrongs) && len(nextTokens) < puzzleChoiceCount-1; j++ {
			nextTokens = append(nextTokens, previousWrongs[j])
		}
		if len(nextTokens) < 1 {
			continue
		}
		nextTokens = append(nextTokens, h.Tokens[i])
		sliceutils.Permutate(nextTokens)

		next := -1

		for next < 0 || next >= len(nextTokens) {
			var lineToPrint bytes.Buffer

			lineToPrint.WriteString("\n> ")
			/*
				for j := 0; j < i; j++ {
					txt := highlights.TokenMap[h.Tokens[j]].Content
					if j == lastI {
						lineToPrint.WriteString(fmt.Sprintf("%s%s%s ", colorGreen, txt, colorNone))
					} else {
						lineToPrint.WriteString(fmt.Sprintf("%s ", txt))
					}
				}*/

			if lastI > 0 {
				lineToPrint.WriteString(string(hContent[:h.TokenStarts[lastI]]))
				lineToPrint.WriteString(colorGreen)
				lineToPrint.WriteString(string(hContent[h.TokenStarts[lastI]:h.TokenStarts[lastI+1]]))
				lineToPrint.WriteString(colorNone)
				lineToPrint.WriteString(strings.TrimSuffix(string(hContent[h.TokenStarts[lastI+1]:h.TokenStarts[i]]), " "))
			} else {
				lineToPrint.WriteString(strings.TrimSuffix(string(hContent[:h.TokenStarts[i]]), " "))
			}
			lineToPrint.WriteString(fmt.Sprintf(" %s____?%s", colorYellow, colorNone))

			replacer := strings.NewReplacer(
				"\n", " ",
				"\r", " ",
				"\t", " ",
			)
			content := fixAlignment(replacer.Replace(lineToPrint.String()), alignmentWidth)

			fmt.Printf("\n%s\n\n", content)

			for j := 0; j < len(nextTokens); j++ {
				txt := highlights.TokenMap[nextTokens[j]].RealContent
				if len(txt) > 0 {
					txt = fmt.Sprintf("%s%s", strings.ToUpper(txt[:1]), txt[1:])
				}
				fmt.Printf("  %d. %s\n", (j + 1), txt)
			}
			fmt.Printf("  Q. Quit\n")
			fmt.Printf("\n")

			cmd := strings.ToLower(readOne())
			switch cmd {
			case "q":
				return
			case "1":
				next = 0
			case "2":
				next = 1
			case "3":
				next = 2
			case "4":
				next = 3
			}
		}

		selected := nextTokens[next]
		if h.Tokens[i] == selected {
			fmt.Fprintf(os.Stdout, "%sCORRECT!%s\n", colorGreen, colorNone)
			correctAnswers = correctAnswers + 1
			lastI = i
			previousWrongs = nil
		} else {
			previousWrongs = append(previousWrongs, selected)

			txt := highlights.TokenMap[selected].RealContent
			if len(txt) > 0 {
				txt = fmt.Sprintf("%s%s", strings.ToUpper(txt[:1]), txt[1:])
			}

			fmt.Fprintf(os.Stdout, "%sWRONG: %s%s\n", colorRed, txt, colorNone)
			wrongAnswers = wrongAnswers + 1
			i = i - 1
		}
	}

	totalAnswers := max(1, correctAnswers+wrongAnswers)
	h.Score.Count = h.Score.Count + 1

	score := float64(correctAnswers) / float64(totalAnswers)
	score = score * float64(h.Score.Count)
	h.Score.Sum = h.Score.Sum + score

	if lastI > 0 {
		var lineToPrint bytes.Buffer

		lineToPrint.WriteString(string(hContent[:h.TokenStarts[lastI]]))
		lineToPrint.WriteString(colorGreen)
		if lastI+1 < len(h.TokenStarts) {
			lineToPrint.WriteString(string(hContent[h.TokenStarts[lastI]:h.TokenStarts[lastI+1]]))
		} else {
			lineToPrint.WriteString(string(hContent[h.TokenStarts[lastI]:]))
		}
		lineToPrint.WriteString(colorNone)
		if lastI+1 < len(h.TokenStarts) {
			lineToPrint.WriteString(string(hContent[h.TokenStarts[lastI+1]:]))
		}

		replacer := strings.NewReplacer(
			"\n", " ",
			"\r", " ",
			"\t", " ",
		)
		content := fixAlignment(replacer.Replace(lineToPrint.String()), alignmentWidth)

		fmt.Printf("\n%s\n\n", content)
	} else {
		replacer := strings.NewReplacer(
			"\n", " ",
			"\r", " ",
			"\t", " ",
		)
		content := fixAlignment(replacer.Replace(h.Content), alignmentWidth)

		fmt.Printf("\n%s\n\n", content)
	}

	if fileExists("scores.txt") {
		copyFile("scores.txt", "scores.txt.bak")
	}
	highlights.WriteScore("scores.txt")
}

func printRandomCard(highlights HighlightDatabase) {
	h := highlights.PickHighlight()
	fmt.Printf("%s\n", h.Content)
}

func readOne() string {
	for {
		buff := make([]byte, 1)
		n, err := os.Stdin.Read(buff)
		if err != nil {
			panic(err)
		}
		if n == 1 {
			str := strings.TrimSpace(string(buff))
			if str != "" {
				return str
			}
		}
	}
}
