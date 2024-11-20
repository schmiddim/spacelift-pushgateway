package api

import (
	"strings"
)

type SpaceLiftPayload struct {
	Branch      string   `json:"branch"`
	Commit      Commit   `json:"commit"`
	Labels      []string `json:"labels"`
	Name        string   `json:"name"`
	Namespace   string   `json:"namespace"`
	ProjectRoot string   `json:"projectRoot"`
	Repository  string   `json:"repository"`
	StackID     string   `json:"stackId"`
	State       string   `json:"state"`
}

type Commit struct {
	Author    string `json:"author"`
	Branch    string `json:"branch"`
	CreatedAt int64  `json:"createdAt"`
	Hash      string `json:"hash"`
	IssueID   string `json:"issueId"`
	Message   string `json:"message"`
	URL       string `json:"url"`
}

func (s *SpaceLiftPayload) GetLabels() map[string]string {
	labelMap := make(map[string]string)
	for _, label := range s.Labels {
		parts := strings.SplitN(label, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			labelMap[key] = value
		}
	}
	return labelMap
}
