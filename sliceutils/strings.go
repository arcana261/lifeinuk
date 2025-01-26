package sliceutils

import "strings"

func TrimSpace[E ~[]string](s E) []string {
	return MapFunc(s, func(t string) string {
		return strings.TrimSpace(t)
	})
}

func Split[E ~[]string](s E, sep string) [][]string {
	return MapFunc(s, func(t string) []string {
		return strings.Split(t, sep)
	})
}
