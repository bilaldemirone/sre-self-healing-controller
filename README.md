# Kubernetes SRE Self-Healing Controller
![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)
![Kubernetes](https://img.shields.io/badge/Kubernetes-1.36-326CE5?logo=kubernetes)
![Prometheus](https://img.shields.io/badge/Prometheus-Monitoring-E6522C?logo=prometheus)
![Elasticsearch](https://img.shields.io/badge/Elasticsearch-8.x-005571?logo=elasticsearch)
![Kibana](https://img.shields.io/badge/Kibana-8.x-005571?logo=kibana)
![Fluent Bit](https://img.shields.io/badge/Fluent_Bit-Logging-49BDA5)

A production-inspired Kubernetes controller built with **Go** and **controller-runtime** that continuously monitors Pod health, evaluates runtime conditions using Prometheus metrics, generates Kubernetes Events, and performs asynchronous self-healing actions through a remediation worker.

To provide end-to-end observability, controller logs are collected by Fluent Bit, indexed in Elasticsearch, and visualized in Kibana.

---

## Features

- Kubernetes controller built with **controller-runtime**
- Real-time Pod monitoring
- Prometheus CPU & restart metrics integration
- Configurable thresholds via ConfigMap
- Automated remediation using an asynchronous worker queue
- Kubernetes Event generation
- Fluent Bit log collection
- Elasticsearch log indexing
- Kibana log visualization
- Docker & Kind based local development environment

---

## Architecture

```
                   +----------------+
                   | Kubernetes API |
                   +--------+-------+
                            |
                            |
                     Watches Pods
                            |
                            ▼
                  +------------------+
                  |  Pod Controller  |
                  +------------------+
                            |
                            ▼
                 +----------------------+
                 | Analyzer Component   |
                 +----------------------+
                            |
        +-------------------+-------------------+
        |                                       |
        ▼                                       ▼
 Kubernetes Event                    Remediation Queue
                                             |
                                             ▼
                                    Remediation Worker
                                             |
                                             ▼
                                     Kubernetes Client
                                             |
                                             ▼
                                      Delete Failed Pod


                Prometheus
                      │
                      ▼
        CPU Usage & Restart Metrics


             Controller Logs
                    │
                    ▼
               Fluent Bit
                    │
                    ▼
             Elasticsearch
                    │
                    ▼
                 Kibana
```

---

## Components

| Component | Description |
|------------|-------------|
| controller-runtime | Kubernetes controller |
| Prometheus | Metrics provider |
| Analyzer | Health decision engine |
| Queue | Asynchronous remediation queue |
| Worker | Executes remediation actions |
| Fluent Bit | Log collector |
| Elasticsearch | Log storage |
| Kibana | Log visualization |

---

## Health Rules

The analyzer currently evaluates:

- Pod Failed
- High CPU Usage
- High Restart Count

Example thresholds:

| Metric | Threshold |
|----------|-----------|
| CPU Usage | > 0.05 Core |
| Restart Count | Configurable |
| Pod Phase | Failed |

Thresholds are loaded dynamically from a Kubernetes ConfigMap.

---

## Automated Remediation

Supported remediation actions:

- Generate Kubernetes Events
- Queue remediation requests
- Delete unhealthy Pods

---

## Logging Pipeline

```
Controller
    │
stdout
    │
Fluent Bit
    │
Elasticsearch
    │
Kibana
```

---

## Project Structure

```
apps/
controller-runtime/
deploy/
docs/
scripts/
```

| Directory | Description |
|-----------|-------------|
| apps | Sample workloads |
| controller-runtime | Kubernetes controller implementation |
| deploy | Kubernetes manifests |
| docs | Documentation, diagrams and screenshots |
| scripts | Build & deployment automation |

---

## Technologies

| Category | Technologies |
|-----------|--------------|
| Language | Go |
| Kubernetes | controller-runtime, client-go |
| Monitoring | Prometheus |
| Logging | Fluent Bit, Elasticsearch, Kibana |
| Container | Docker |
| Local Cluster | Kind |
| Package Manager | Helm |

---

## Quick Start

### Build

```bash
./scripts/build.sh
```

### Deploy

```bash
./scripts/deploy.sh
```

### Generate CPU Load

```bash
kubectl exec -it deploy/dummy-app -- sh

yes > /dev/null
```

Stop the load using:

```bash
Ctrl + C
```

### Trigger Pod Remediation

```bash
kubectl delete pod -l app=dummy-app
```

### Observe

- Prometheus metrics
- Kubernetes Events
- Controller logs
- Elasticsearch logs
- Kibana dashboards

---

## Screenshots

### Controller Logs

![Controller Logs](docs/screenshots/controller-logs.png)

### Kubernetes Events

![Events](docs/screenshots/events.png)

### Prometheus Metrics

![Prometheus](docs/screenshots/prometheus.png)

### Kibana

![Kibana](docs/screenshots/kibana.png)

### Cluster Overview

![Cluster](docs/screenshots/cluster.png)

---

## Future Improvements

- Slack Integration
- Microsoft Teams Integration
- Deployment Restart
- StatefulSet Support
- Machine Learning Based Anomaly Detection
- OpenTelemetry Integration
- Grafana Dashboards
- Alertmanager Integration

---

## Known Limitations

- Only Pod remediation is currently implemented.
- Thresholds are statically configured through ConfigMap.
- Deployment restart remediation is planned but not yet implemented.
- Alert notifications are not yet integrated with external platforms.

---

## Author

**Bilal Demir**

DevOps / Site Reliability Engineer

GitHub: https://github.com/bilaldemirone