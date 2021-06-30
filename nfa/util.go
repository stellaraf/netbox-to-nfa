package nfa

import (
	"fmt"
	"strings"
)

// cleanSprintf formats a string and cleans the result based on preconfigured replacement patterns.
func cleanSprintf(s string, f ...interface{}) string {
	pairs := [][]string{{"\t", ""}, {"\n", ""}, {"\r", ""}, {"\": ", "\":"}}

	for _, p := range pairs {
		o := p[0]
		n := p[1]
		s = strings.ReplaceAll(s, o, n)
	}
	return fmt.Sprintf(s, f...)
}
