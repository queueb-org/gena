package common

// Deref returns fallback if pointer is nil.
func Deref[T any](in *T, fallback T) T {
	if in == nil {
		return fallback
	}

	return *in
}
