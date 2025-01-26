package sliceutils

import (
	"cmp"
	"math/rand"
	"slices"
)

func ForEach2[E ~[]T, T any](s E, fn func(int, T) bool) {
	c := len(s)
	for i := 0; i < c; i++ {
		ok := fn(i, s[i])
		if !ok {
			return
		}
	}
}

func ForEach[E ~[]T, T any](s E, fn func(T)) {
	ForEach2(s, func(_ int, t T) bool {
		fn(t)
		return true
	})
}

func AccumulateFunc2[E ~[]T, T any, V any](s E, seed V, fn func(V, int, T) (V, bool)) V {
	var ok bool

	ForEach2(s, func(i int, t T) bool {
		seed, ok = fn(seed, i, t)
		return ok
	})

	return seed
}

func UniqueSortedFunc2[E ~[]T, T any](s E, fn func(int, T, int, T) (bool, bool)) []T {
	return AccumulateFunc2(s, Pair[Pair[int, T], []T]{}, func(p Pair[Pair[int, T], []T], i int, t T) (Pair[Pair[int, T], []T], bool) {
		if i == 0 {
			return Pair[Pair[int, T], []T]{Key: Pair[int, T]{0, t}, Value: []T{t}}, true
		}

		cmp, ok := fn(p.Key.Key, p.Key.Value, i, t)
		if cmp {
			return p, ok
		}

		p.Value = append(p.Value, t)
		p.Key = Pair[int, T]{Key: i, Value: t}
		return p, ok
	}).Value
}

func UniqueSortedFunc[E ~[]T, T any](s E, fn func(T, T) bool) []T {
	return UniqueSortedFunc2(s, func(_ int, t1 T, _ int, t2 T) (bool, bool) {
		return fn(t1, t2), true
	})
}

func UniqueSorted[E ~[]T, T cmp.Ordered](s E) []T {
	return UniqueSortedFunc(s, func(t1, t2 T) bool {
		return !(t1 < t2) && !(t2 < t1)
	})
}

func MapFunc2[E ~[]T, T any, V any](s E, fn func(int, T) (V, bool)) []V {
	return AccumulateFunc2(s, make([]V, 0), func(v []V, i int, t T) ([]V, bool) {
		mapped, ok := fn(i, t)
		return append(v, mapped), ok
	})
}

func MapFunc[E ~[]T, T any, V any](s E, fn func(T) V) []V {
	return MapFunc2(s, func(_ int, t T) (V, bool) {
		return fn(t), true
	})
}

func Clone[E ~[]T, T any](s E) []T {
	return MapFunc(s, func(t T) T {
		return t
	})
}

func FilterFunc2[E ~[]T, T any](s E, fn func(int, T) (bool, bool)) []T {
	return AccumulateFunc2(s, make([]T, 0), func(t1 []T, i int, t2 T) ([]T, bool) {
		include, ok := fn(i, t2)
		if include {
			t1 = append(t1, t2)
		}
		return t1, ok
	})
}

func FilterFunc[E ~[]T, T any](s E, fn func(T) bool) []T {
	return FilterFunc2(s, func(_ int, t T) (bool, bool) {
		return fn(t), true
	})
}

func IndexOfFunc2[E ~[]T, T any](s E, fn func(int, T) (bool, bool)) int {
	return AccumulateFunc2(s, -1, func(last int, index int, t T) (int, bool) {
		found, ok := fn(index, t)
		if found {
			return index, false
		}
		return last, ok
	})
}

func IndexOfFunc[E ~[]T, T any](s E, fn func(T) bool) int {
	return IndexOfFunc2(s, func(_ int, t T) (bool, bool) {
		return fn(t), true
	})
}

func IndexOf[E ~[]T, T cmp.Ordered](s E, value T) int {
	return IndexOfFunc(s, func(t T) bool {
		return !(value < t) && !(t < value)
	})
}

func RemoveFunc2[E ~[]T, T any](s E, fn func(int, T) (bool, bool)) []T {
	return FilterFunc2(s, func(i int, t T) (bool, bool) {
		cmp, ok := fn(i, t)
		return !cmp, ok
	})
}

func RemoveFunc[E ~[]T, T any](s E, fn func(T) bool) []T {
	return RemoveFunc2(s, func(_ int, t T) (bool, bool) {
		return fn(t), true
	})
}

func Remove[E ~[]T, T cmp.Ordered](s E, value T) []T {
	return RemoveFunc(s, func(t T) bool {
		return !(value < t) && !(t < value)
	})
}

func RemoveAt[E ~[]T, T any](s E, index int) []T {
	return RemoveFunc2(s, func(i int, t T) (bool, bool) {
		return index == i, true
	})
}

type Pair[K any, V any] struct {
	Key   K
	Value V
}

type Tuple3[A any, B any, C any] struct {
	A A
	B B
	C C
}

func MaxFunc2[E ~[]T, T any](s E, fn func(int, T, int, T) (int, bool)) (int, T) {
	return MinFunc2(s, func(i int, x T, j int, y T) (int, bool) {
		cmp, ok := fn(i, x, j, y)
		return -cmp, ok
	})
}

func MaxFunc[E ~[]T, T any](s E, fn func(T, T) int) T {
	_, result := MaxFunc2(s, func(_ int, t1 T, _ int, t2 T) (int, bool) {
		return fn(t1, t2), true
	})
	return result
}

func Max[E ~[]T, T cmp.Ordered](s E) T {
	return MaxFunc(s, func(t1, t2 T) int {
		if t1 < t2 {
			return -1
		}
		if t2 < t1 {
			return 1
		}
		return 0
	})
}

func MinFunc2[E ~[]T, T any](s E, fn func(int, T, int, T) (int, bool)) (int, T) {
	res := AccumulateFunc2(s, Pair[int, T]{Key: -1}, func(t1 Pair[int, T], i int, t2 T) (Pair[int, T], bool) {
		if i == 0 {
			return Pair[int, T]{Key: 0, Value: t2}, true
		}

		c, ok := fn(t1.Key, t1.Value, i, t2)
		if c <= 0 {
			return t1, ok
		}
		return Pair[int, T]{Key: i, Value: t2}, ok
	})
	return res.Key, res.Value
}

func MinFunc[E ~[]T, T any](s E, fn func(T, T) int) T {
	_, result := MinFunc2(s, func(_ int, t1 T, _ int, t2 T) (int, bool) {
		return fn(t1, t2), true
	})
	return result
}

func Min[E ~[]T, T cmp.Ordered](s E) T {
	return MinFunc(s, func(t1, t2 T) int {
		if t1 < t2 {
			return -1
		}
		if t2 < t1 {
			return 1
		}
		return 0
	})
}

func LowerBoundSortedFunc2[E ~[]T, T any](s E, fn func(int, T) (int, bool)) int {
	lo := 0
	hi := len(s) - 1
	if hi < 0 {
		return -1
	}

	for lo < hi {
		mid := (lo + hi) / 2
		c, ok := fn(lo, s[mid])
		if !ok {
			return -1
		}
		if c < 0 {
			lo = mid + 1
		} else if c == 0 {
			return lo
		} else {
			hi = mid
		}
	}

	c, ok := fn(lo, s[lo])
	if !ok {
		return -1
	}
	if c < 0 {
		return -1
	}
	return lo
}

func LowerBoundSortedFunc[E ~[]T, T any](s E, fn func(T) int) int {
	return LowerBoundSortedFunc2(s, func(_ int, t T) (int, bool) {
		return fn(t), true
	})
}

func ToMapFunc2[E ~[]T, T any, K comparable, V any](s E, fn func(int, T) (K, V, bool)) map[K]V {
	return AccumulateFunc2(s, make(map[K]V), func(m map[K]V, i int, t T) (map[K]V, bool) {
		k, v, ok := fn(i, t)
		m[k] = v
		return m, ok
	})
}

func ToMapFunc[E ~[]T, T any, K comparable, V any](s E, fn func(T) (K, V)) map[K]V {
	return ToMapFunc2(s, func(_ int, t T) (K, V, bool) {
		k, v := fn(t)
		return k, v, true
	})
}

func ToMultiMapFunc2[E ~[]T, T any, K comparable, V any](s E, fn func(int, T) (K, V, bool)) map[K][]V {
	return AccumulateFunc2(s, make(map[K][]V), func(m map[K][]V, i int, t T) (map[K][]V, bool) {
		k, v, ok := fn(i, t)
		m[k] = append(m[k], v)
		return m, ok
	})
}

func ToMultiMapFunc[E ~[]T, T any, K comparable, V any](s E, fn func(T) (K, V)) map[K][]V {
	return ToMultiMapFunc2(s, func(_ int, t T) (K, V, bool) {
		k, v := fn(t)
		return k, v, true
	})
}

func Permutate[E ~[]T, T any](s E) {
	for j := 0; j < len(s); j++ {
		k := j + rand.Intn(len(s)-j)
		if j != k {
			temp := s[j]
			s[j] = s[k]
			s[k] = temp
		}
	}
}

func Flatten[E ~[]T, T ~[]V, V any](s E) []V {
	return AccumulateFunc2(s, make([]V, 0), func(last []V, _ int, current T) ([]V, bool) {
		last = append(last, current...)
		return last, true
	})
}

func Zip[E ~[]T, T any, X ~[]Y, Y any](q E, w X) []Pair[T, Y] {
	return AccumulateFunc2(q, make([]Pair[T, Y], 0), func(p []Pair[T, Y], i int, t T) ([]Pair[T, Y], bool) {
		if i >= len(w) {
			return p, false
		}
		p = append(p, Pair[T, Y]{Key: t, Value: w[i]})
		return p, true
	})
}

func Lookup2[E ~[]T, T comparable, V any](s E, m map[T]V) []Pair[V, bool] {
	return MapFunc(s, func(t T) Pair[V, bool] {
		item, ok := m[t]
		return Pair[V, bool]{Key: item, Value: ok}
	})
}

func Lookup[E ~[]T, T comparable, V any](s E, m map[T]V) []V {
	items := Lookup2(s, m)
	items = FilterFunc(items, func(p Pair[V, bool]) bool { return p.Value })
	return MapFunc(items, func(p Pair[V, bool]) V { return p.Key })
}

func Sort[S ~[]E, E cmp.Ordered](x S) []E {
	result := Clone(x)
	slices.Sort(result)
	return result
}

func SortFunc[S ~[]E, E any](x S, cmp func(a, b E) int) []E {
	result := Clone(x)
	slices.SortFunc(result, cmp)
	return result
}

func AccumulateFunc[E ~[]T, T any](s E, fn func(T, T) T) T {
	var defaultValue T
	return AccumulateFunc2(s, defaultValue, func(prev T, index int, current T) (T, bool) {
		if index == 0 {
			return current, true
		}
		return fn(prev, current), true
	})
}

func SumInt[E ~[]int](s E) int {
	return AccumulateFunc(s, func(prev, current int) int {
		return prev + current
	})
}

func Range(s, e int) []int {
	var result []int
	for s < e {
		result = append(result, s)
		s = s + 1
	}
	return result
}
