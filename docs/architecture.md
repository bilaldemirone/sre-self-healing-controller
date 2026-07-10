# Architecture

```mermaid
flowchart TD
    K8S[Kubernetes API] --> C[Pod Controller]
    P[Prometheus] --> C
    CM[ConfigMap Thresholds] --> C

    C --> A[Analyzer / Decision Engine]

    A -->|CREATE_EVENT| E[Kubernetes Event]
    A -->|DELETE_POD| Q[Remediation Queue]
    A -->|NONE| H[Healthy]

    Q --> W[Remediation Worker]
    W --> K8S

    C --> STDOUT[Controller Logs]
    STDOUT --> FB[Fluent Bit]
    FB --> ES[Elasticsearch]
    ES --> KB[Kibana]