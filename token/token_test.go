package token_test

import (
	"testing"

	"github.com/dywoq/miniasm/token"
)

func TestIsIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"__foo", true},
		{"foo_123", true},
		{"___", true},
		{"##", false},
		{"1234", false},
	}

	// to prevent a lot of x and save readability, we append this test case manually.
	str := []byte{}
	for range 256 {
		str = append(str, 'x')
	}
	tests = append(tests, struct {
		input string
		want  bool
	}{string(str), false})

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := token.IsIdentifier(test.input)
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
