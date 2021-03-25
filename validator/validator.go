package validator

func EmptyString(s string) bool {
	if s == "" {
		return true
	}
	return false
}
