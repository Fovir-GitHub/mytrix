// Package utils provides utility functions.
package utils

// Filter returns a new slice containing only the elements for which keep returns true.
// It is a generic function that works with any slice type.
func Filter[T any](s []T, keep func(*T) bool) []T {
	var result []T
	for _, v := range s {
		if keep(&v) {
			result = append(result, v)
		}
	}
	return result
}
