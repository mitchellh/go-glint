package glint

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

		{
			"unicode character",
			"\u2584",
			1,
		},

		{
			"unicode character with newlines",
			"\u2584\n",
			1,
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
			"height zero",
			"hello\nworld\n",
			0,
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

func TestClampTextWidth(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Width    int
		Expected string
	}{
		{
			"empty",
			"",
			10,
			"",
		},

		{
			"width zero",
			"hello\nworld\n",
			0,
			"",
		},

		{
			"width fits",
			"hello world\ni fit!",
			100,
			"hello world\ni fit!",
		},

		{
			"clamped one line",
			"hello world",
			5,
			"hello",
		},

		{
			"clamped one line ends in newline",
			"hello world\n",
			5,
			"hello\n",
		},

		{
			"fits ends in newline",
			"hello\n",
			5,
			"hello\n",
		},

		{
			"clamped both lines",
			"hello world\ni fit!",
			5,
			"hello\ni fit",
		},

		{
			"unicode multi-byte character",
			"\u2584",
			1,
			"\u2584",
		},

		{
			"unicode exceeds width",
			"\u2584\u2582",
			1,
			"\u2584",
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			require := require.New(t)
			actual := clampTextWidth(tt.Input, tt.Width)
			require.Equal(tt.Expected, actual)
		})
	}
}
