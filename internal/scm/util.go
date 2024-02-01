package scm

// Diff two slices and get the added and removed items compared to s1
func Diff[T comparable](s1, s2 []T) (added, removed []T) {
	s1Lookup := map[T]struct{}{}
	for _, v := range s1 {
		s1Lookup[v] = struct{}{}
	}
	s2Lookup := map[T]struct{}{}
	for _, v := range s2 {
		s2Lookup[v] = struct{}{}
	}

	for _, v := range s2 {
		if _, ok := s1Lookup[v]; !ok {
			added = append(added, v)
		}
	}
	for _, v := range s1 {
		if _, ok := s2Lookup[v]; !ok {
			removed = append(removed, v)
		}
	}

	return added, removed
}

// Map runs a function for each value in a slice and returns a slice of all function returns
func Map[T any, K any](vals []T, mapping func(T) K) []K {
	newVals := make([]K, len(vals))
	for i, v := range vals {
		newVals[i] = mapping(v)
	}
	return newVals
}
