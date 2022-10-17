package repocounter

import "github.com/lindell/multi-gitter/internal/scm"

// log10 is an integer version of log10 (number of digits)
func log10(num int) int {
	ret := 0
	for num != 0 {
		ret++
		num /= 10
	}
	return ret
}

func shortenRepoName(repo scm.Repository, maxLength int) string {
	// TODO: Treat org name and repository differently

	replaceStr := "..."

	name := repo.FullName()
	if len(name) <= maxLength {
		return name
	}

	return name[:maxLength-len(replaceStr)] + replaceStr
}

type iterator[T any] interface {
	Next() (T, bool)
}

func iterateTimes[T any](times int, iter iterator[T], fn func(T)) {
	i := 0
	val, ok := iter.Next()
	for ok && i < times {
		fn(val)
		val, ok = iter.Next()
		i++
	}
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func maxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}
