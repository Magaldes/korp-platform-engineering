# ARCH-CONTAINER — Containers

- Project: Projeto Korp DevOps Challenge
- Operational revision: 6
- Classification: `declared/proposed/inferred` pre-codebase state
- Projection: `flowchart`

```mermaid
flowchart LR
    ctr_go_service["http-server-projeto-korp"]
    ctr_grafana["Grafana"]
    ctr_nginx["NGINX Reverse Proxy"]
    ctr_prometheus["Prometheus"]
    ctr_nginx -->|"forwards HTTP requests to port 8080"| ctr_go_service
    ctr_prometheus -->|"scrapes internal GET /metrics on service port 8080"| ctr_go_service
    ctr_grafana -->|"reads time series for dashboards"| ctr_prometheus
```

## Operational summary

Purpose: Modelar e governar a arquitetura pre-codebase do desafio Korp, preservando rastreabilidade entre requisitos, criterios de aceite, decisoes, componentes e validacao.
Compiled 7 projected items from operational revision 6.
