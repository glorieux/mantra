// Package slice contains commonly used slices methods
package slice

// IntDeepEqual compares two integer slices
func IntDeepEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	if (a == nil) != (b == nil) {
		return false
	}

	b = b[:len(a)]
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

// Remove removes one element from an array of strings
func Remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

// Contains returns wether or not a string item is contained in a slice of strings
func Contains(item string, list []string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
