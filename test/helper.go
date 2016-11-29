package test

func sum(items []int) int {
	result := 0
	for _, item := range items {
		result += item
	}
	return result
}
