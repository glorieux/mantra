package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHadAlphaNumeric(t *testing.T) {
	type test struct {
		input  string
		result bool
	}
	tests := []test{
		test{"ablkjdalgkjfs10912380912", true},
		test{")#(@*!)!#@)(", false},
		test{"dlsakj*lkjdsalkj", false},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.result, HasAlphaNumeric(tc.input), tc.input)
	}
}
