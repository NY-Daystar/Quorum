package utils

import (
	"testing"
)

// For testing
// $ cd utils
// $ go test -v

// Test string conversion to int
func TestGetName(t *testing.T) {
	var stringRef = "sample.pdf"
	var expected = "sample"

	var result = GetName(stringRef)

	if expected != result {
		t.Errorf("Test case failed with result: %s - expected: %s", result, expected)
	}

}

func TestGetNames(t *testing.T) {
	tests := []struct {
		target   string
		expected string
	}{
		{
			target:   "samples.html",
			expected: "samples",
		},
		{
			target:   "sample",
			expected: "sample",
		},
		{
			target:   "sample.",
			expected: "sample",
		},
		{
			target:   "sample//",
			expected: "sample",
		},
		{
			target:   "sample..",
			expected: "sample.",
		},
		{
			target:   "19d4d7226dcf0d68.eml",
			expected: "19d4d7226dcf0d68",
		},
	}

	for _, testCase := range tests {
		var result = GetName(testCase.target)

		if testCase.expected != result {
			t.Errorf("Test case failed with case: %s result: %s - expected: %s", testCase.target, result, testCase.expected)
		}
	}
}
