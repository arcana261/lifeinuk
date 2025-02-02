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

	alignmentWidth    = 50
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

	correctAnswers := 0
	wrongAnswers := 0
	lastI := -1

	for i := 2; i < len(h.Tokens); i++ {
		if highlights.TokenMap[h.Tokens[i]].SkipPuzzle {
			//fmt.Printf("skipping '%s' due to SkipPuzzle\n", highlights.TokenMap[h.Tokens[i]].Content)
			continue
		}

		nextTokens := highlights.TokenMap[h.Tokens[i-1]].NominateNextTokens(highlights, puzzleChoiceCount)
		if len(nextTokens) < 2 {
			continue
		}
		sliceutils.Permutate(nextTokens)
		if sliceutils.IndexOf(nextTokens, h.Tokens[i]) < 0 {
			nextTokens = nextTokens[1:]
			nextTokens = append(nextTokens, h.Tokens[i])
			sliceutils.Permutate(nextTokens)
		}
		if len(nextTokens) < 2 {
			//var allNextTokens []string
			//for _, nt := range highlights.TokenMap[h.Tokens[i-1]].NextTokens {
			//allNextTokens = append(allNextTokens, highlights.TokenMap[nt.ID].Content)
			//}
			//fmt.Printf("skipping '%s' due to len(nextTokens)<2 [case 1][%s]\n", highlights.TokenMap[h.Tokens[i]].Content, strings.Join(allNextTokens, ","))
			continue
		}

		next := -1

		for next < 0 || next >= len(nextTokens) {
			var lineToPrint bytes.Buffer

			lineToPrint.WriteString("\n> ")
			for j := 0; j < i; j++ {
				txt := highlights.TokenMap[h.Tokens[j]].Content
				if j == lastI {
					lineToPrint.WriteString(fmt.Sprintf("%s%s%s ", colorGreen, txt, colorNone))
				} else {
					lineToPrint.WriteString(fmt.Sprintf("%s ", txt))
				}
			}
			lineToPrint.WriteString("____?\n\n")

			fmt.Printf("\n%s\n\n", fixAlignment(lineToPrint.String(), alignmentWidth))

			for j := 0; j < len(nextTokens); j++ {
				txt := highlights.TokenMap[nextTokens[j]].Content
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
		} else {
			fmt.Fprintf(os.Stdout, "%sWRONG: %s%s\n", colorRed, highlights.TokenMap[selected].Content, colorNone)
			wrongAnswers = wrongAnswers + 1
			i = i - 1
		}
	}

	totalAnswers := max(1, correctAnswers+wrongAnswers)
	h.Score.Count = h.Score.Count + 1

	score := float64(correctAnswers) / float64(totalAnswers)
	score = score * float64(h.Score.Count)
	h.Score.Sum = h.Score.Sum + score

	fmt.Printf("%s\n", h.Content)
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
