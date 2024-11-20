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
}

func NewPushGateway(pushGatewayURL string) *PushGateway {
	return &PushGateway{pushGatewayURL: pushGatewayURL}
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
