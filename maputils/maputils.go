package maputils

import "github.com/arcana261/lifeinuk/sliceutils"

func ToEntries[E ~map[K]V, K comparable, V any](s E) []sliceutils.Pair[K, V] {
	var result []sliceutils.Pair[K, V]
	for k, v := range s {
		result = append(result, sliceutils.Pair[K, V]{Key: k, Value: v})
	}
	return result
}

func FromEntries[E ~[]sliceutils.Pair[K, V], K comparable, V any](s E) map[K]V {
	return sliceutils.ToMapFunc(s, func(p sliceutils.Pair[K, V]) (K, V) {
		return p.Key, p.Value
	})
}

func MapFunc2[E ~map[K]V, K comparable, V any, X comparable, Y any](s E, fn func(K, V) (X, Y, bool)) map[X]Y {
	entries := ToEntries(s)
	newEntries := sliceutils.MapFunc2(entries, func(_ int, p sliceutils.Pair[K, V]) (sliceutils.Pair[X, Y], bool) {
		x, y, ok := fn(p.Key, p.Value)
		return sliceutils.Pair[X, Y]{Key: x, Value: y}, ok
	})
	return FromEntries(newEntries)
}

func MapFunc[E ~map[K]V, K comparable, V any, X comparable, Y any](s E, fn func(K, V) (X, Y)) map[X]Y {
	return MapFunc2(s, func(k K, v V) (X, Y, bool) {
		x, y := fn(k, v)
		return x, y, true
	})
}

func MapValuesFunc2[E ~map[K]V, K comparable, V any, Y any](s E, fn func(K, V) (Y, bool)) map[K]Y {
	return MapFunc2(s, func(k K, v V) (K, Y, bool) {
		y, ok := fn(k, v)
		return k, y, ok
	})
}

func MapValuesFunc[E ~map[K]V, K comparable, V any, Y any](s E, fn func(V) Y) map[K]Y {
	return MapValuesFunc2(s, func(k K, v V) (Y, bool) {
		return fn(v), true
	})
}
