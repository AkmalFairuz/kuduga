package pointer

func Make[T any](v T) *T {
	return &v
}

func DoubleMake[T any](v T) **T {
	return Make(&v)
}

func Get[T any](v *T) T {
	return *v
}
