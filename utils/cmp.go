package utils

// Contains checks whether an array of strings contains an element.
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
