package utils

func SliceContainsEntryFromSlice(slice1 []string, slice2 []string) bool {
	slice1Map := make(map[string]*string, len(slice1))
	for _, v := range slice1 {
		slice1Map[v] = nil
	}

	for _, v := range slice2 {
		if _, ok := slice1Map[v]; ok {
			return true
		}
	}

	return false
}
