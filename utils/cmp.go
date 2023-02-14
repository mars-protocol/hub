package utils

// Contains check whether an array of strings contain an element.
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
