package api

import (
	"encoding/json"
	"fmt"
	"github.com/yalp/jsonpath"
	"strings"
)

func TransformJsonValues(data []byte, jsonPath string, separator string) ([]byte, error) {
	// Parse the JSON into a generic map
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Extract the value at the given JSONPath
	labelsData, err := jsonpath.Read(jsonData, jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to extract data from JSON using path %s: %v", jsonPath, err)
	}

	// Check if "labels" is of type slice
	labels, found := labelsData.([]interface{})
	if !found {
		return nil, fmt.Errorf("data at %s is not a slice", jsonPath)
	}

	// Initialize a map to hold the transformed label data
	labelMap := make(map[string]string)

	// Iterate over the labels and split each into key-value pairs based on the separator
	for _, label := range labels {
		labelStr, ok := label.(string)
		if !ok {
			continue // Skip non-string labels
		}
		// Split the string at the first occurrence of the separator
		parts := strings.SplitN(labelStr, separator, 2)
		if len(parts) == 2 {
			labelMap[parts[0]] = parts[1]
		}
	}

	// Replace the "labels" field with the transformed label map
	jsonDataMap, ok := jsonData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected JSON to be a map but found %T", jsonData)
	}
	jsonDataMap["labels"] = labelMap

	// Marshal the modified jsonData back into JSON
	transformedJSON, err := json.MarshalIndent(jsonDataMap, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed JSON: %v", err)
	}

	return transformedJSON, nil
}
