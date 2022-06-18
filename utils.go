package graph

func empty[T any]() T {
	var t T
	return t
}

func isOK[A any, B bool](_ A, ok B) B {
	return ok
}
