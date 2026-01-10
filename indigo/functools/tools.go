package functools

func MapWithError[T, V any](f func(T) (V, error), xs []T) ([]V, error) {
	result := make([]V, len(xs))
	for i, t := range xs {
		v, err := f(t)
		if err != nil {
			return result, err
		}
		result[i] = v
	}
	return result, nil
}
