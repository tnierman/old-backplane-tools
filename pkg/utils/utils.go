package utils

func Contains(list []string, val string) bool {
	for _, elem := range list {
		if elem == val {
			return true
		}
	}
	return false
}
