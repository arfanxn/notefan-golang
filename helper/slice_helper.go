package helper

func SliceChunk[T any](items []T, size int) (chunks [][]T) {
	for size < len(items) {
		items, chunks = items[size:], append(chunks, items[0:size:size])
	}
	return append(chunks, items)
}

func SliceFilter[T any](items []T, callback func(T) bool) []T {
	matchItems := []T{}
	for _, item := range items {
		if callback(item) {
			matchItems = append(matchItems, item)
		}
	}
	return matchItems
}
