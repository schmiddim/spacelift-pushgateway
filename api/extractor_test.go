package api

import (
	"bytes"
	"encoding/json"
	"log"

	"strings"
	"testing"
)

func TestExtractMultipleJSONPaths(t *testing.T) {
	// Beispiel-JSON-Daten für Tests
	validJsonData := []byte(`{
		"repository": {
			"url": "https://github.com/example/repo",
			"author": "hansihamster"
		},
		"commit": {
			"hash": "e9ea5a543fce1b2d52207153f2d580431933b927",
			"message": "Initial commit",
			"createdAt": 1742103798000
		}
	}`)

	// Ungültige JSON-Daten für Tests
	invalidJsonData := []byte(`{
		"repository": {
			"url": "https://github.com/example/repo",
			"author": "hansihamster",
		`)

	tests := []struct {
		name          string
		data          []byte
		paths         []string
		expected      map[string]interface{}
		expectError   bool
		expectWarning bool
	}{
		{
			name: "Valid JSONPath expressions with valid JSON",
			data: validJsonData,
			paths: []string{
				"{.repository.url}",
				"{.commit.hash}",
				"{.commit.message}",
			},
			expected: map[string]interface{}{
				"repository.url": "https://github.com/example/repo",
				"commit.hash":    "e9ea5a543fce1b2d52207153f2d580431933b927",
				"commit.message": "Initial commit",
			},
			expectError:   false,
			expectWarning: false,
		},
		{
			name: "Invalid JSONPath expressions",
			data: validJsonData,
			paths: []string{
				"{.repository.nonexistent}",
				"{.invalid.path[}",
			},
			expected:      map[string]interface{}{},
			expectError:   true, // Expect error due to invalid JSONPath syntax
			expectWarning: true,
		},
		{
			name: "Valid JSONPath expressions with invalid JSON",
			data: invalidJsonData,
			paths: []string{
				"{.repository.url}",
				"{.commit.hash}",
			},
			expected:      map[string]interface{}{},
			expectError:   true, // Expect error due to invalid JSON data
			expectWarning: false,
		},
		{
			name: "Mixed valid and invalid JSONPath expressions",
			data: validJsonData,
			paths: []string{
				"{.repository.url}",
				"{.commit.nonexistent}",
				"{.commit.createdAt}",
			},
			expected: map[string]interface{}{
				"repository.url":   "https://github.com/example/repo",
				"commit.createdAt": float64(1742103798000),
			},
			expectError:   false,
			expectWarning: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect log output to capture warnings
			var logBuffer bytes.Buffer
			log.SetOutput(&logBuffer)

			// Aufruf der zu testenden Funktion
			results, err := ExtractMultipleJSONPaths(tt.data, tt.paths)

			// Überprüfen, ob ein Fehler erwartet wurde
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got: %v", err)
				}
			}

			// Überprüfen der erwarteten Ergebnisse
			for key, expectedValue := range tt.expected {
				actualValue, found := results[key]
				if !found {
					t.Errorf("Expected key '%s' not found in results", key)
					continue
				}

				// JSON-Objekte vergleichen
				expectedJSON, _ := json.Marshal(expectedValue)
				actualJSON, _ := json.Marshal(actualValue)
				if string(expectedJSON) != string(actualJSON) {
					t.Errorf("Expected value for key '%s': %s, but got: %s", key, string(expectedJSON), string(actualJSON))
				}
			}

			// Überprüfen, ob unerwartete Einträge in den Ergebnissen sind
			for key := range results {
				if _, ok := tt.expected[key]; !ok {
					t.Errorf("Unexpected key '%s' found in results", key)
				}
			}

			// Warnungen überprüfen, wenn Pfade nicht gefunden wurden
			logContent := logBuffer.String()
			if tt.expectWarning && !strings.Contains(logContent, "Warning") {
				t.Error("Expected warning for missing JSONPath, but none was logged")
			}
			if !tt.expectWarning && strings.Contains(logContent, "Warning") {
				t.Error("Did not expect warning for missing JSONPath, but one was logged")
			}
		})
	}
}
