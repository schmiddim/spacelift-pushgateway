package api

import (
	"fmt"
	"regexp"
	"strings"
)

type Rename struct {
	Key string
	To  string
}

// RenameKeys renames keys in a map[string]interface{} according to the regex patterns in renames
func RenameKeys(m map[string]interface{}, renames []Rename) (map[string]interface{}, error) {
	// Create a new map to store the renamed keys
	renamedMap := make(map[string]interface{})

	// Iterate over each key in the original map
	for oldKey, value := range m {
		// Rename the key using the renameKey function
		newKey, err := renameKey(oldKey, renames)
		if err != nil {
			return nil, fmt.Errorf("error renaming key '%s': %v", oldKey, err)
		}
		// Add the renamed key to the new map
		renamedMap[newKey] = value
	}

	return renamedMap, nil
}

func renameKey(subject string, renames []Rename) (string, error) {
	for _, r := range renames {
		if strings.TrimSpace(r.Key) == "" {
			return subject, fmt.Errorf("empty regex for Key")
		}
		if strings.TrimSpace(r.To) == "" {
			return subject, fmt.Errorf("empty To Value!")
		}
		// Compile the regex for the key to be replaced
		re, err := regexp.Compile(r.Key)
		if err != nil {
			return subject, fmt.Errorf("failed to compile regex for key '%s': %v", r.Key, err)
		}
		// Replace all occurrences of the regex match in the subject with the "To" value
		subject = re.ReplaceAllString(subject, r.To)
	}
	return subject, nil
}
