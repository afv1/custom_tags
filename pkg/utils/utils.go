package utils

// AnyOf checks if subject exists in slice of any type
func AnyOf[T comparable](what T, where ...T) bool {
	for _, item := range where {
		if item == what {
			return true
		}
	}

	return false
}
