package optional

type Optional[T any] struct {
	value *T
}

func (opt Optional[T]) Exists() bool {
	return opt.value != nil
}

func (opt Optional[T]) Value() T {
	return *opt.value
}

func NewOptional[T any](value T) Optional[T] {
	return Optional[T]{
		value: &value,
	}
}

func NewEmptyOptional[T any]() Optional[T] {
	return Optional[T]{}
}
