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
			input:    "     Hello World     ",
			expected: []string{"hello", "world"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("mismatch length len: %v, expected: %v", len(actual), len(c.expected))
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expetedWord := c.expected[i]
			if word != expetedWord {
				t.Errorf("words not equal: word: %v expected: %v", word, expetedWord)
				t.Fail()
			}
		}
	}
}
