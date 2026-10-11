package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "commAnd",
			expected: []string{"command"},
		},
		{
			input:    "  hello     world    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  hello world    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  ",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Actual length doesn't match the expected")
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("Word doesn't match the expected\na: %s \ne: %s\n", word, expectedWord)
			}
		}
	}
}
