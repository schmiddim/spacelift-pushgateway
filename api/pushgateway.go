package api

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	log "github.com/sirupsen/logrus"
	"io"

	"net/http"
	"time"
)

type PushGateway struct {
	pushGatewayURL   string
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("failed to close response body: %v", err)
		}
	}(resp.Body)

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
	for key, value := range labelPairs {
		switch v := value.(type) {
		case string:
			output[key] = v
		case int, int64, float64:
			output[key] = fmt.Sprintf("%v", v) // Zahlen in String umwandeln
		case bool:
			output[key] = fmt.Sprintf("%t", v) // Boolesche Werte in "true"/"false" umwandeln
		case time.Time:
			output[key] = v.Format(time.RFC3339) // Zeitstempel als RFC3339-String speichern
		default:
			log.Printf("Warning: Ignoring unsupported label type for key '%s'", key)
		}
	}

	metric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        p.targetMetric,
		Help:        p.targetMetricHelp,
		ConstLabels: output,
	})

	metric.Set(float64(time.Now().Unix()))

	err := push.New(p.pushGatewayURL, p.jobName).
		Collector(metric).
		Push()
	if err != nil {
		return fmt.Errorf("failed to push to Pushgateway: %v", err)
	}

	return nil
}
