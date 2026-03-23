package fetcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFetchWindows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		totalDays int
		expected  []fetchWindow
	}{
		{
			name:      "zero days",
			totalDays: 0,
			expected:  nil,
		},
		{
			name:      "one day",
			totalDays: 1,
			expected: []fetchWindow{
				{idx: 0, n: 1},
			},
		},
		{
			name:      "eight days",
			totalDays: 8,
			expected: []fetchWindow{
				{idx: 0, n: 8},
			},
		},
		{
			name:      "nine days",
			totalDays: 9,
			expected: []fetchWindow{
				{idx: 0, n: 8},
				{idx: 8, n: 1},
			},
		},
		{
			name:      "fifteen days",
			totalDays: 15,
			expected: []fetchWindow{
				{idx: 0, n: 8},
				{idx: 8, n: 7},
			},
		},
		{
			name:      "sixteen days",
			totalDays: 16,
			expected: []fetchWindow{
				{idx: 0, n: 8},
				{idx: 8, n: 8},
			},
		},
		{
			name:      "more than sixteen days gets capped",
			totalDays: 30,
			expected: []fetchWindow{
				{idx: 0, n: 8},
				{idx: 8, n: 8},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, buildFetchWindows(tt.totalDays))
		})
	}
}
