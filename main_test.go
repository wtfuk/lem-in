package main

import (
	"reflect"
	"testing"
)

func TestRemoveComments(t *testing.T) {
	// Test cases
	testCases := []struct {
		name     string
		lines    []string
		expected []string
	}{
		{
			name:     "Valid input",
			lines:    []string{"line1", "#line2", "##end", "line3"},
			expected: []string{"line1", "##end", "line3"},
		},
		{
			name:     "Invalid input",
			lines:    []string{"line1", "#line2", "##start", "line3"},
			expected: []string{"line1", "##start", "line3"},
		},
	}

	// Run the test cases
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Call the RemoveComments function
			result := RemoveComments(tt.lines)
			// Check the result
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Unexpected result. Expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestIsRoom(t *testing.T) {
	testCases := []struct {
		input string
		want  bool
	}{
		{"room 100 200", true},
		{"room 200 100", true},
		{"room abc def", false},
		{"room 100 def", false},
		{"room abc 200", false},
		{"room 100 200 300", false},
	}

	for _, tc := range testCases {
		got := IsRoom(tc.input)
		if got != tc.want {
			t.Errorf("IsRoom(%q) = %v; want %v", tc.input, got, tc.want)
		}
	}
}

func TestDeleteStartRoom(t *testing.T) {
	testCases := []struct {
		input    []string
		expected []string
	}{
		{
			input:    []string{"##start", "room1", "room2", "room3"},
			expected: []string{"room2", "room3"},
		},
		{
			input:    []string{"room1", "##start", "room2", "room3"},
			expected: []string{"room1", "room3"},
		},
		{
			input:    []string{"room1", "room2", "room3"},
			expected: []string{"room1", "room2", "room3"},
		},
	}

	for _, tc := range testCases {
		result := DeleteStartRoom(tc.input)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Unexpected result for input %v: got %v, expected %v", tc.input, result, tc.expected)
		}
	}
}
