package arrayext

// ToMap converts an array to a map, using the keyField as the key
func ToMap[T interface{}](array []*T, keyFunc func(int) string) map[string]*T {
	m := make(map[string]*T)
	for i, item := range array {
		// use keyFunc to get the key
		key := keyFunc(i)
		m[key] = item
	}
	return m
}

func RemoveDuplicates[T interface{}](arr []T) []T {
	seen := make(map[interface{}]bool)
	var result []T
	for _, v := range arr {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
