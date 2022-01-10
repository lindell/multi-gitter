package repocounter

// log10 is an integer version of log10 (number of digits)
func log10(num int) int {
	ret := 0
	for num != 0 {
		ret += 1
		num /= 10
	}
	return ret
}
