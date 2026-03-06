package teams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNumberOfTeams(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected []int
	}{
		{"6 players", 6, []int{2, 3, 4, 5}},
		{"3 players", 3, []int{2}},
		{"too few", 2, []int{}},
		{"too many", 26, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetNumberOfTeams(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRandomizeTeams(t *testing.T) {
	users := []string{"user 1", "user 2", "user 3", "user 4", "user 5", "user 6", "user 7", "user 8"}
	result, err := RandomizeTeams(3, users)

	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Len(t, result[0], 3)
	assert.Len(t, result[1], 3)
	assert.Len(t, result[2], 2)
}
