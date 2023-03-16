package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stripSuffixIfExist(t *testing.T) {
	assert.Equal(t, "string", stripSuffixIfExist("stringSuffix", "Suffix"))
	assert.Equal(t, "stringSuffix", stripSuffixIfExist("stringSuffix", "NoMatch"))
}

func Test_chunkSlice(t *testing.T) {
	assert.Equal(t, [][]int{{0, 1}, {2}}, chunkSlice([]int{0, 1, 2}, 2))
	assert.Equal(t, [][]int{{0, 1, 2}}, chunkSlice([]int{0, 1, 2}, 4))
}
