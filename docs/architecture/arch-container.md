# ARCH-CONTAINER — Containers

- Project: Projeto Korp DevOps Challenge
- Operational revision: 7
- Classification: `declared/proposed/inferred` approved canonical as-built state
- Projection: `flowchart`

```mermaid
flowchart LR
    ctr_go_service["http-server-projeto-korp"]
    ctr_grafana["Grafana"]
    ctr_nginx["NGINX Reverse Proxy"]
    ctr_postgres["PostgreSQL 18.4"]
    ctr_postgres_migrate["PostgreSQL Migration Runner"]
    ctr_prometheus["Prometheus"]
    ctr_redis["Redis 8.8.1"]
    ctr_nginx -->|"forwards HTTP requests to port 8080"| ctr_go_service
    ctr_prometheus -->|"scrapes internal GET /metrics on service port 8080"| ctr_go_service
    ctr_grafana -->|"reads time series for dashboards"| ctr_prometheus
    ctr_go_service -->|"uses PostgreSQL as authoritative audit/statistics data source when repository is available"| ctr_postgres
    ctr_go_service -->|"uses Redis as optional statistics cache and can serve valid cache hits without PostgreSQL"| ctr_redis
    ctr_postgres_migrate -->|"applies versioned PostgreSQL schema before application data capability is considered ready"| ctr_postgres
```

## Operational summary

Purpose: Modelar e governar a arquitetura implementada e validada do Projeto Korp, preservando rastreabilidade entre requisitos, decisoes, componentes, containers, dados, processos operacionais e o estado as-built da aplicacao.
Compiled 13 projected items from operational revision 7.
