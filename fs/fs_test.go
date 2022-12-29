package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimInvalidFileNameChars(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "hello:world",
			want:  "hello_world",
		},
		{
			input: "hello?world",
			want:  "hello_world",
		},
		{
			input: "hello*world",
			want:  "hello_world",
		},
		{
			input: "hello:?*world",
			want:  "hello___world",
		},
		{
			input: " hello  world ",
			want:  "hello world",
		},
	}

	for _, test := range tests {
		assert.Equal(t, TrimInvalidFileNameChars(test.input), test.want)
	}
}
