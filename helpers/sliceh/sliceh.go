package sliceh

import "math/rand"

func Chunk[T any](items []T, size int) (chunks [][]T) {
	for size < len(items) {
		items, chunks = items[size:], append(chunks, items[0:size:size])
	}
	return append(chunks, items)
}

func Filter[T any](items []T, callback func(T) bool) []T {
	matchItems := []T{}
	for _, item := range items {
		if callback(item) {
			matchItems = append(matchItems, item)
		}
	}
	return matchItems
}

func Random[T any](slice []T) T {
	if len(slice) > 1 {
		return slice[rand.Intn(len(slice)-1)]
	} else {
		return slice[0]
	}
}
