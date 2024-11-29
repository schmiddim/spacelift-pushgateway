# Spacelift Metrics to Prometheus PushGateway
This project allows you to send Spacelift job metrics to a Prometheus PushGateway. The metrics are collected from Spacelift job payloads and pushed to Prometheus, enabling you to monitor your infrastructure and application states effectively.

[![codecov](https://codecov.io/gh/schmiddim/spacelift-pushgateway/graph/badge.svg?token=lxCOCj9JPi)](https://codecov.io/gh/schmiddim/spacelift-pushgateway)
[![Docker Pulls](https://img.shields.io/docker/pulls/schmiddim/spacelift-pushgateway.svg)](https://hub.docker.com/r/schmiddim/spacelift-pushgateway/)

## Todo
- Prometheus Metrics for Pushgateway Communication & Webserver


# Setup
## Step 1: Install Prometheus & PushGateway using Helm

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm upgrade --install  prommi prometheus-community/prometheus -n spacelift --create-namespace --set  server.image.tag=v3.0.0
```

## Port Forwarding
To access Prometheus and PushGateway locally, set up port forwarding to your local machine:


- Prometheus
```
 kubectl port-forward -n spacelift services/prommi-prometheus-server 9090:80
```
- Push Gateway
```
kubectl port-forward -n spacelift  services/prommi-prometheus-pushgateway 9091:9091
```
This allows you to access the services on your local machine:

- Prometheus: http://localhost:9090
- PushGateway: http://localhost:9091


## Send Data with Curl 

Sending Metrics to PushGateway
You can send metrics to the PushGateway using a POST request with curl. The payload should be a JSON object, and the API key is required for authentication.
```
curl -d @example-payload.json -H "Authorization: Bearer your-api-key" http://localhost:8080/push
```


