package uuid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_IsValid(t *testing.T) {
	type testCase struct {
		name     string
		uuid     string
		expected bool
		validate func(t *testing.T, result bool, expected bool)
	}

	tests := []testCase{
		{
			name:     "returns boolean determining if uuid is valid",
			uuid:     "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			expected: true,
			validate: func(t *testing.T, result bool, expected bool) {
				assert.Equal(t, result, expected)
			},
		},
		{
			name:     "works if the uuid is all the same number",
			uuid:     "00000000-0000-0000-0000-000000000000",
			expected: true,
			validate: func(t *testing.T, result bool, expected bool) {
				assert.Equal(t, result, expected)
			},
		},
		{
			name:     "fails when the uuid is empty",
			uuid:     "",
			expected: false,
			validate: func(t *testing.T, result bool, expected bool) {
				assert.Equal(t, result, expected)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsValid(tc.uuid)
			tc.validate(t, result, tc.expected)
		})
	}
}
