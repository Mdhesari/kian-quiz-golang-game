package validation

func Name(name string) bool {
	return len(name) > 3
}