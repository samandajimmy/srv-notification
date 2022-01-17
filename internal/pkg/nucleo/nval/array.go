package nval

func InArrayString(value string, list []string) bool {
	for _, listValue := range list {
		if listValue == value {
			return true
		}
	}
	return false
}
