package api

import (
	"reflect"
	"testing"
)

func TestRenameKey(t *testing.T) {
	tests := []struct {
		name          string
		subject       string
		renames       []Rename
		expected      string
		expectedError bool
	}{
		{
			name:          "Regex Fail",
			subject:       "hello Error",
			renames:       []Rename{{Key: "hello(", To: "gopher"}},
			expected:      "hello Error",
			expectedError: true,
		},
		{
			name:          "Single rename",
			subject:       "hello world",
			renames:       []Rename{{Key: "world", To: "gopher"}},
			expected:      "hello gopher",
			expectedError: false,
		},
		{
			name:          "Multiple renames",
			subject:       "hello world, welcome world",
			renames:       []Rename{{Key: "world", To: "gopher"}, {Key: "hello", To: "hi"}},
			expected:      "hi gopher, welcome gopher",
			expectedError: false,
		},
		{
			name:          "No renames",
			subject:       "hello world",
			renames:       []Rename{},
			expected:      "hello world",
			expectedError: false,
		},
		{
			name:          "No matching key",
			subject:       "goodbye world",
			renames:       []Rename{{Key: "hello", To: "hi"}},
			expected:      "goodbye world", // No "hello" to replace
			expectedError: false,
		},
		{
			name:          "Partial match in word",
			subject:       "helloworld",
			renames:       []Rename{{Key: "hello", To: "hi"}},
			expected:      "hiworld", // Only "hello" will be replaced with "hi"
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := renameKey(tt.subject, tt.renames)
			if result != tt.expected {
				t.Errorf("expected %q, but got %q", tt.expected, result)
			}
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error %v, but got error %v", tt.expectedError, err)
			}
		})
	}
}

func TestRenameKeys(t *testing.T) {
	tests := []struct {
		name     string
		inputMap map[string]interface{}
		renames  []Rename
		expected map[string]interface{}
		hasError bool
	}{
		{
			name: "simple rename",
			inputMap: map[string]interface{}{
				"oldKey1":    "value1",
				"anotherKey": "value2",
				"keyToKeep":  "value3",
			},
			renames: []Rename{
				{Key: "^oldKey1$", To: "newKey1"},
				{Key: "^anotherKey$", To: "updatedKey"},
			},
			expected: map[string]interface{}{
				"newKey1":    "value1",
				"updatedKey": "value2",
				"keyToKeep":  "value3",
			},
			hasError: false,
		},
		{
			name: "rename with regex pattern",
			inputMap: map[string]interface{}{
				"user_id":   "123",
				"user_name": "john_doe",
				"keyToKeep": "value3",
			},
			renames: []Rename{
				{Key: "^user_", To: "account_"},
			},
			expected: map[string]interface{}{
				"account_id":   "123",
				"account_name": "john_doe",
				"keyToKeep":    "value3",
			},
			hasError: false,
		},
		{
			name: "invalid regex pattern",
			inputMap: map[string]interface{}{
				"someKey": "value1",
			},
			renames: []Rename{
				{Key: "[invalid", To: "shouldNotMatter"},
			},
			expected: map[string]interface{}{
				"someKey": "value1",
			},
			hasError: true,
		},

		{
			name: "Empty regex pattern",
			inputMap: map[string]interface{}{
				"someKey": "value1",
			},
			renames: []Rename{
				{Key: " ", To: "shouldNotMatter"},
			},
			expected: map[string]interface{}{
				"someKey": "value1",
			},
			hasError: true,
		},
		{
			name: "Empty To ",
			inputMap: map[string]interface{}{
				"someKey": "value1",
			},
			renames: []Rename{
				{Key: "asfd", To: ""},
			},
			expected: map[string]interface{}{
				"someKey": "value1",
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run the RenameKeys function with test case parameters
			outputMap, err := RenameKeys(tt.inputMap, tt.renames)

			// Check if an error was expected
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
				return
			}

			// Check for unexpected error
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify if the output matches the expected result
			if !reflect.DeepEqual(outputMap, tt.expected) {
				t.Errorf("Test '%s' failed. Expected: %v, Got: %v", tt.name, tt.expected, outputMap)
			}
		})
	}
}
