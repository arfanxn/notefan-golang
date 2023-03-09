package sliceh

import "math/rand"

// Chunk will separate an array by the given number of size
func Chunk[T any](items []T, size int) (chunks [][]T) {
	for size < len(items) {
		items, chunks = items[size:], append(chunks, items[0:size:size])
	}
	return append(chunks, items)
}

// Filter returns only items that satisfy the given predicate (return the true predicate condition only)
func Filter[T any](items []T, callback func(T) bool) []T {
	matchItems := []T{}
	for _, item := range items {
		if callback(item) {
			matchItems = append(matchItems, item)
		}
	}
	return matchItems
}

// Map mapping slice of T
func Map[T1, T2 any](items []T1, callback func(T1) T2) []T2 {
	var resultItems []T2
	for _, item := range items {
		resultItems = append(resultItems, callback(item))
	}
	return resultItems
}

// Random return a random T from the given slice of T
func Random[T any](slice []T) T {
	if len(slice) > 1 {
		return slice[rand.Intn(len(slice)-1)]
	} else {
		return slice[0]
	}
}
