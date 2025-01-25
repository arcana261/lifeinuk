package main

import (
	"cmp"
	"io/ioutil"
	"math/rand"
	"os"
	"slices"
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

func permutateSlice[E ~[]T, T any](s E) {
	for j := 0; j < len(s); j++ {
		k := j + rand.Intn(len(s)-j)
		if j != k {
			temp := s[j]
			s[j] = s[k]
			s[k] = temp
		}
	}
}

func findInSlice[E ~[]T, T cmp.Ordered](s E, item T) int {
	for i := 0; i < len(s); i++ {
		if s[i] == item {
			return i
		}
	}
	return -1
}

func uniqueSlice[E ~[]T, T cmp.Ordered](s E) E {
	end := len(s) - 1
	slices.Sort(s)

	start := 0
	put := 0
	for start <= end {
		until := 1
		for until <= end {
			if s[until] != s[start] {
				break
			}
			until = until + 1
		}
		s[put] = s[start]
		start = until
		put = put + 1
	}

	return s[:put]
}

func lowerBoundFunc[E ~[]T, T any, V any](s E, v V, cmp func(T, V) int) int {
	lo := 0
	hi := len(s) - 1
	for lo < hi {
		mid := (lo + hi) / 2
		c := cmp(s[lo], v)
		if c < 0 {
			lo = mid + 1
		} else if c == 0 {
			return lo
		} else {
			hi = mid
		}
	}
	return lo
}

func cloneSlice[T any](s []T) []T {
	result := make([]T, len(s))
	result = append(result, s...)
	return result
}

func removeAt[T any](s []T, index int) []T {
	for i := index; i < len(s)-1; i++ {
		s[i] = s[i+1]
	}
	s = s[:(len(s) - 1)]
	return s
}

func removeFunc[T any](s []T, fn func(T) bool) []T {
	put := 0
	for i := 0; i < len(s); i++ {
		if !fn(s[i]) {
			s[put] = s[i]
			put = put + 1
		}
	}
	s = s[:put]
	return s
}
