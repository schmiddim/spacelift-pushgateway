package api

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"net/http"
	"time"
)

type PushGateway struct {
	pushGatewayURL string
	//fieldsToExtract []string
	targetMetric     string
	targetMetricHelp string
	jobName          string
}

func NewPushGateway(pushGatewayURL string, targetMetric string, targetMetricHelp string, jobName string) *PushGateway {
	return &PushGateway{
		pushGatewayURL:   pushGatewayURL,
		targetMetric:     targetMetric,
		jobName:          jobName,
		targetMetricHelp: targetMetricHelp,
		//fieldsToExtract: fieldsToExtract,
	}
}
func (p *PushGateway) CheckPushGatewayStatus() error {
	client := &http.Client{Timeout: 5 * time.Second}

	// Anfrage an den /metrics-Endpoint senden
	resp, err := client.Get(fmt.Sprintf("%s/metrics", p.pushGatewayURL))
	if err != nil {
		return fmt.Errorf("failed to connect to Pushgateway: %v", err)
	}
	defer resp.Body.Close()

	// Überprüfen, ob der HTTP-Status OK (200) ist
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Pushgateway responded with unexpected status: %s", resp.Status)
	}
	return nil
}

func (p *PushGateway) ValidateLabels(labels map[string]interface{}) (bool, []error) {
	var errors []error
	var hasErrors = false
	for i, l := range keys(labels) {
		counter := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("dummy_metric_total_%d", i),
				Help: "Dummy metric for label validation",
			},
			[]string{l},
		)
		err := prometheus.Register(counter)
		prometheus.Unregister(counter)
		if err != nil {
			errors = append(errors, fmt.Errorf("invalid label: %s", l))
		}

	}
	if len(errors) > 0 {
		hasErrors = true
	}
	return hasErrors, errors

}
func keys(labels map[string]interface{}) []string {
	var keys []string
	for k := range labels {
		keys = append(keys, k)
	}
	return keys
}

func (p *PushGateway) PushMetrics(labelPairs map[string]interface{}) error {
	output := make(map[string]string)

	//@todo flatten check
	//@todo created_at is not parsed
	//@todo check booleans
	for key, value := range labelPairs {
		strValue, _ := value.(string)

		output[key] = strValue
	}

	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        p.targetMetric,
		Help:        p.targetMetricHelp,
		ConstLabels: output,
	})

	metric.Set(1)
	err := push.New(p.pushGatewayURL, p.jobName).
		Collector(metric).
		Push()
	if err != nil {
		return fmt.Errorf("failed to push to Pushgateway: %v", err)
	}

	return nil
}

// Deprecated: Use NewFunction with improved parameters.
func (p *PushGateway) PushToGateway(commitInfo SpaceLiftPayload) error {
	labelMap := commitInfo.GetLabels()
	labelMap["name"] = commitInfo.Name
	labelMap["branch"] = commitInfo.Branch
	labelMap["state"] = commitInfo.State

	buildStatus := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "spacelift_build_status",
		Help:        "Build status of the application",
		ConstLabels: labelMap,
	})
	buildStatus.Set(1)

	createdAt := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "spacelift_build_created_at",
		Help:        "Build creation timestamp",
		ConstLabels: labelMap,
	})
	createdAt.Set(float64(commitInfo.Commit.CreatedAt / 1e9)) // Unix-Zeit in Sekunden
	//createdAt.Set(float64(time.Now().UnixMilli() / 1e9))

	err := push.New(p.pushGatewayURL, "spacelift_build").
		Collector(buildStatus).
		Collector(createdAt).
		Push()
	if err != nil {
		return fmt.Errorf("failed to push to Pushgateway: %v", err)
	}

	return nil
}
