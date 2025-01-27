package box

func Wrap[T any](v T) *T {
	return &v
}
