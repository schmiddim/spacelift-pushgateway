package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransformLabelsJSON(t *testing.T) {
	tests := []struct {
		name          string
		inputJSON     []byte
		jsonPath      string
		separator     string
		expectedJSON  []byte
		expectedError bool
	}{
		{
			name: "valid labels transformation",
			inputJSON: []byte(`{
				"branch": "master",
				"labels": [
					"class:platform",
					"environment:prod",
					"feature:enable_log_timestamps",
					"folder:terraform/services/foo"
				]
			}`),
			jsonPath:  "$.labels",
			separator: ":",
			expectedJSON: []byte(`{
  "branch": "master",
  "labels": {
    "class": "platform",
    "environment": "prod",
    "feature": "enable_log_timestamps",
    "folder": "terraform/services/foo"
  }
}`),
			expectedError: false,
		},
		{
			name: "missing labels field",
			inputJSON: []byte(`{
				"branch": "master"
			}`),
			jsonPath:  "$.labels",
			separator: ":",
			expectedJSON: []byte(`{
  "branch": "master",
  "labels": {}
}`),
			expectedError: true,
		},
		{
			name: "non-key-value entries in labels",
			inputJSON: []byte(`{
				"branch": "master",
				"labels": [
					"another:valid",
					"still:good"
				]
			}`),
			jsonPath:  "$.labels",
			separator: ":",
			expectedJSON: []byte(`{
  "branch": "master",
  "labels": {
    "another": "valid",
    "still": "good"
  }
}`),
			expectedError: false,
		},
		{
			name: "duplicate keys in labels",
			inputJSON: []byte(`{
				"branch": "master",
				"labels": [
					"tool:terraform",
					"tool:terragrunt"
				]
			}`),
			jsonPath:  "$.labels",
			separator: ":",
			expectedJSON: []byte(`{
  "branch": "master",
  "labels": {
    "tool": "terragrunt"
  }
}`),
			expectedError: false,
		},
		{
			name: "invalid jsonpath",
			inputJSON: []byte(`{
				"branch": "master",
				"labels": [
					"class:platform"
				]
			}`),
			jsonPath:      "$.nonexistentField",
			separator:     ":",
			expectedJSON:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformedJSON, err := TransformJsonValues(tt.inputJSON, tt.jsonPath, tt.separator)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, string(tt.expectedJSON), string(transformedJSON))
			}
		})
	}
}
