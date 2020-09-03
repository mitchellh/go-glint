package dynamiccli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLongestLine(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected int
	}{
		{
			"empty",
			"",
			0,
		},

		{
			"no newline",
			"foo",
			3,
		},

		{
			"trailing newline",
			"foo\n",
			3,
		},

		{
			"middle line",
			"foo\nbarr\nbaz\n",
			4,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)
			actual := longestLine(tt.Input)
			require.Equal(tt.Expected, actual)
		})
	}
}

func TestTruncateTextHeight(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Height   int
		Expected string
	}{
		{
			"empty",
			"",
			10,
			"",
		},

		{
			"shorter than limit",
			"foo\nbar",
			5,
			"foo\nbar",
		},

		{
			"greater than limit",
			"foo\nbar\nbaz\nqux",
			3,
			"foo\nbar\nbaz",
		},

		{
			"equal to limit",
			"foo\nbar\nbaz",
			3,
			"foo\nbar\nbaz",
		},

		{
			"equal to limit with trailing newline",
			"foo\nbar\n",
			3,
			"foo\nbar\n",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)
			actual := truncateTextHeight(tt.Input, tt.Height)
			require.Equal(tt.Expected, actual)
		})
	}
}
