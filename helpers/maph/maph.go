package maph

// Map mapping a map 
func Map[KEY, KEY2 comparable, VAL, VAL2 any](
	items map[KEY]VAL,
	callback func(KEY, VAL) (KEY2, VAL2),
) map[KEY2]VAL2 {
	var resultItems map[KEY2]VAL2
	for key, val := range items {
		key2, val2 := callback(key, val)
		resultItems[key2] = val2
	}
	return resultItems
}
