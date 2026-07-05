package keystone

// Find Finds the first element that matches the given predicate
func Find[T any](slice []T, predicate func(item T) bool) (T, bool) {
	for i := range slice {
		if predicate(slice[i]) {
			return slice[i], true
		}
	}

	var result T
	return result, false
}

// FindOrElse returns the first element satisfying the predicate, or defaultValue if none match.
func FindOrElse[V any](slice []V, predicate func(V) bool, defaultValue V) V {
	for _, item := range slice {
		if predicate(item) {
			return item
		}
	}
	return defaultValue
}

// Filtered returns a new slice containing all elements from the input slice that satisfy the predicate function.
func Filtered[V any](slice []V, predicate func(V) bool) []V {
	var result []V
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}
