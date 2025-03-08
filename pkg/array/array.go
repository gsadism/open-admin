package array

func Merge[T any](arr ...[]T) []T {
	list := make([]T, 0)
	for _, v := range arr {
		list = append(list, v...)
	}
	return list
}
