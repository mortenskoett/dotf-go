package utils

// Returns true if sl contains str.
func Contains(sl []string, str string) bool {
	for _, e := range sl {
		if str == e {
			return true
		}
	}
	return false
}
