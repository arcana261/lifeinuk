package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func copyFile(src string, dst string) {
	// Read all content of src to data, may cause OOM for a large file.
	data, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err)
	}
	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		panic(err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func fixAlignment(str string, width int) string {
	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.ReplaceAll(str, "\r", " ")
	str = strings.ReplaceAll(str, "\t", " ")
	parts := strings.Split(str, " ")
	var result []string
	var current []string
	var currentLength int
	for i := 0; i <= len(parts); i++ {
		var part string

		if i < len(parts) {
			part = strings.TrimSpace(parts[i])
		}
		if i == len(parts) || currentLength+len(part) > width {
			result = append(result, strings.Join(current, " "))
			current = nil
			currentLength = 0
		}
		if part != "" {
			current = append(current, part)
			currentLength = currentLength + len(part)
		}
	}
	if len(current) > 0 {
		panic("missing string here in fixAlignment")
	}

	return strings.Join(result, "\n")
}
