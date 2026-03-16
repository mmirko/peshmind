package peshmind

func elementInSlice(element string, slice []string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}
