package util

import "strings"

// SplitRemoveEmpty splits a string by a separator and removes any empty strings from the result.
func SplitRemoveEmpty(s string, sep string) (a []string) {
	parts := strings.Split(s, sep)
	for _, p := range parts {
		if p != "" {
			a = append(a, p)
		}
	}
	return a
}

// CompareArrays determines if any elements of one array exist in another array.
func CompareArrays(one []string, two []string) bool {
	for _, a := range one {
		for _, b := range two {
			if a == b {
				return true
			}
		}
	}
	return false
}
