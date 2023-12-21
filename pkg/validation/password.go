package validation

func Password(pass string) bool {
    return len(pass) > 5
}