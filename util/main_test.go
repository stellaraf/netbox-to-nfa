package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsProbablySHA256(t *testing.T) {
	type Case struct {
		string
		bool
	}
	cases := []Case{
		{"f41239e0d77cd0d5dc1d90e97c26c52991b76bd6074aca764ea8064bd81c84a3", true},
		{"this is not a sha256 hash", false},
		{"", false},
		{"f41239e0d77cd0d5dc1d90e97c26c52991b76bd6074aca764ea8064", false},
		{"f41239e0d77cd0d5dc1d90e97c26c52991b76bd6074aca764ea8064bd81c84a3123", false},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			t.Parallel()
			result := IsProbablySHA256(c.string)
			assert.Equal(t, c.bool, result)
		})
	}
}
