package api

import (
	"encoding/json"
	"fmt"
	"github.com/yalp/jsonpath"
	"log"
	"strings"
)

type Extractor struct {
	fieldsToExtract []string
}

func NewExtractor(fieldsToExtract []string) *Extractor {
	return &Extractor{fieldsToExtract}
}

func flattenMap(prefix string, value interface{}, flatMap map[string]string) {
	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			flattenMap(prefix+key+".", val, flatMap)
		}
	case []interface{}:
		for i, item := range v {
			// Prüfen, ob item ein string ist
			if str, ok := item.(string); ok {
				// Überprüfen, ob der string ":" enthält
				if strings.Contains(str, ":") {
					// String in Key und Value aufteilen
					parts := strings.SplitN(str, ":", 2)
					key := parts[0]
					value := parts[1]
					// Key-Value-Paar in flatMap hinzufügen
					flatMap[fmt.Sprintf("%s%s", prefix, key)] = value
				} else {
					// Falls kein ":", füge den ganzen string als Wert hinzu
					flatMap[fmt.Sprintf("%s%d", prefix, i)] = str
				}
			} else {
				// Rekursiver Aufruf für nicht-string Elemente
				flattenMap(fmt.Sprintf("%s%d.", prefix, i), item, flatMap)
			}
		}
	default:
		flatMap[strings.TrimSuffix(prefix, ".")] = fmt.Sprintf("%v", v)
	}
}

func (e *Extractor) Extract(jsonData []byte) (map[string]string, error) {

	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("error unmarshalling data: %v", err)
	}
	extractedData := make(map[string]string)
	for _, field := range e.fieldsToExtract {
		query := fmt.Sprintf("$.%s", field)
		value, err := jsonpath.Read(data, query)
		if err != nil {
			log.Printf("Warning: Field %s unable to extract: %v", field, err)
			continue
		}
		flattenMap(field+".", value, extractedData)
	}

	return extractedData, nil
}
