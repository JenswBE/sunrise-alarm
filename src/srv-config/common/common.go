package common

func ToSlice[T comparable, U any](input map[T]U) []U {
	output := make([]U, 0, len(input))
	for _, value := range input {
		output = append(output, value)
	}
	return output
}
