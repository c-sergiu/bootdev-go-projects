package repl

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct{
		input string
		expected []string
	}{
		{
			input: "hello world",
			expected: []string{ "hello", "world" },
		},
			{
			input: " hello world",
			expected: []string{ "hello", "world" },
		},
		{
			input: "hello world ",
			expected: []string{ "hello", "world" },
		},
		{
			input: " hello world ",
			expected: []string{ "hello", "world" },
		},
		{
			input: "  heLLo worLd   ",
			expected: []string{ "hello", "world" },
		},
		{
			input: "     ",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) { 
			t.Errorf("expected length: %v, actual: %v", len(c.expected), len(actual)) 
			t.Fatalf("Test failed with fatal error.")
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("expected: %v, actual %v", c.expected[i], actual[i])
			}
		}
	}
}
