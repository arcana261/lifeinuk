package main

import (
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"
)

const (
	colorRed   = "\033[0;31m"
	colorGreen = "\033[0;32m"
	colorNone  = "\033[0m"
)

func main() {
	// switch stdin into 'raw' mode
	/*
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)*/

	highlights, err := ReadHighlights("highlights.txt")
	if err != nil {
		panic(err)
	}

	//fmt.Println(highlights)
	//fmt.Println(highlights.TokenMap[45].NominateNextTokens(4))

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
	idx := rand.Intn(len(highlights.Highlights))
	h := highlights.Highlights[idx]

	shouldSkip := func(t int) bool {
		switch highlights.TokenMap[t].Content {
		case "are":
			return true
		case "the":
			return true
		case "is":
			return true
		case "that":
			return true
		case "which":
			return true
		case "not":
			return true
		case "in":
			return true
		case "for":
			return true
		case "they":
			return true
		case "was":
			return true
		case "be":
			return true
		case "and":
			return true
		case "there":
			return true
		case "a":
			return true
		case "it":
			return true
		case "what":
			return true
		default:
		}

		return false
	}

	for i := 1; i < len(h.Tokens); i++ {
		if shouldSkip(h.Tokens[i]) {
			continue
		}

		var nextTokens []int
		nextTokens = append(nextTokens, h.Tokens[i])
		nextTokens = append(nextTokens, highlights.TokenMap[h.Tokens[i-1]].NominateNextTokens(4)...)
		slices.Sort(nextTokens)

		var nextTokensUnique []int
		for j := 0; j < len(nextTokens); j++ {
			if j == 0 || nextTokens[j] != nextTokens[j-1] {
				nextTokensUnique = append(nextTokensUnique, nextTokens[j])
			}
		}
		nextTokens = nextTokensUnique

		for j := 0; j < len(nextTokens); j++ {
			k := j + rand.Intn(len(nextTokens)-j)
			if j != k {
				temp := nextTokens[j]
				nextTokens[j] = nextTokens[k]
				nextTokens[k] = temp
			}
		}
		correctIndex := -1
		for j := 0; j < len(nextTokens); j++ {
			if nextTokens[j] == h.Tokens[i] {
				correctIndex = j
				break
			}
		}
		if len(nextTokens) < 2 {
			continue
		}

		if correctIndex == len(nextTokens)-1 {
			j := rand.Intn(len(nextTokens) - 1)
			temp := nextTokens[j]
			nextTokens[j] = nextTokens[correctIndex]
			nextTokens[correctIndex] = temp
		}
		nextTokens = nextTokens[:(len(nextTokens) - 1)]

		if len(nextTokens) < 2 {
			continue
		}
		next := -1

		for next < 0 || next >= len(nextTokens) {
			fmt.Printf("\n> ")
			for j := 0; j < i; j++ {
				fmt.Printf("%s ", highlights.TokenMap[h.Tokens[j]].Content)
			}
			fmt.Printf("____?\n\n")

			for j := 0; j < len(nextTokens); j++ {
				fmt.Printf("  %d. %s\n", (j + 1), highlights.TokenMap[nextTokens[j]].Content)
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
		} else {
			fmt.Fprintf(os.Stdout, "%sWRONG: %s%s\n", colorRed, highlights.TokenMap[h.Tokens[i]].Content, colorNone)
		}
	}

	fmt.Printf("%s\n", highlights.Highlights[idx].Content)
}

func printRandomCard(highlights HighlightDatabase) {
	idx := rand.Intn(len(highlights.Highlights))
	fmt.Printf("%s\n", highlights.Highlights[idx].Content)
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
