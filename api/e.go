package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"k8s.io/client-go/util/jsonpath"
	"log"
	"strings"
)

// ExtractMultipleJSONPaths extracts key-value pairs from JSON data using multiple JSONPath expressions.
// It returns a flat map where each JSONPath expression yields a key-value pair.
// If a JSONPath does not match any data, a warning is logged.
func ExtractMultipleJSONPaths(data []byte, paths []string) (map[string]interface{}, error) {
	// Unmarshal JSON data into a generic map
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Result map to hold key-value pairs
	results := make(map[string]interface{})

	for _, path := range paths {
		// Generate a key name based on the JSONPath, e.g., "commit.author" for "{.commit.author}"
		key := strings.Trim(path, "{}.")
		jp := jsonpath.New("extractor")
		if err := jp.Parse(path); err != nil {
			return nil, fmt.Errorf("failed to parse JSONPath '%s': %v", path, err)
		}

		// Buffer to hold output
		var resultBuffer bytes.Buffer
		if err := jp.Execute(&resultBuffer, jsonData); err != nil {
			// Log a warning if JSONPath yields no results
			log.Printf("Warning: No results found for JSONPath '%s'\n", path)
			continue
		}

		// Try to decode the result as JSON, fallback to string if it fails
		var extractedValue interface{}
		rawOutput := resultBuffer.String()

		// Check if JSONPath result is valid JSON array or object
		if json.Unmarshal([]byte(rawOutput), &extractedValue) != nil {
			// If it fails, treat raw output as a single string result
			extractedValue = rawOutput
		}

		results[key] = extractedValue
	}

	return results, nil
}
