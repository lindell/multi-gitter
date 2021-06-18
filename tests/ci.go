// +build ci

package tests

func init() {
	skipTypes = append(skipTypes, skipTypeCI)
}
