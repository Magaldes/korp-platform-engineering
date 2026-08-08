# ARCH-DEPLOYMENT — Deployment

- Project: Projeto Korp DevOps Challenge
- Operational revision: 6
- Classification: `declared/proposed/inferred` pre-codebase state
- Projection: `flowchart`

```mermaid
flowchart LR
    act_candidate(["Candidate Operator"])
    ctr_go_service["http-server-projeto-korp"]
    ctr_grafana["Grafana"]
    ctr_nginx["NGINX Reverse Proxy"]
    ctr_prometheus["Prometheus"]
    env_demo["Linux Demonstration Environment"]
    net_bridge["Docker Bridge Network"]
    node_linux_host["Linux Host"]
    tech_ansible["Ansible"]
    tech_compose["Docker Compose"]
    tech_docker["Docker Engine"]
    ctr_nginx -->|"forwards HTTP requests to port 8080"| ctr_go_service
    ctr_prometheus -->|"scrapes internal GET /metrics on service port 8080"| ctr_go_service
    ctr_grafana -->|"reads time series for dashboards"| ctr_prometheus
    env_demo -->|"hosts deployment node"| node_linux_host
    node_linux_host -->|"runs as Docker container"| ctr_nginx
    node_linux_host -->|"runs as Docker container"| ctr_go_service
    node_linux_host -->|"runs as Docker container"| ctr_prometheus
    node_linux_host -->|"runs as Docker container"| ctr_grafana
    tech_ansible -->|"installs and configures the environment"| node_linux_host
    node_linux_host -->|"Docker execution engine"| tech_docker
    tech_compose -->|"declares and starts container"| ctr_nginx
    tech_compose -->|"declares and starts container"| ctr_go_service
    tech_compose -->|"declares and starts container"| ctr_prometheus
    tech_compose -->|"declares and starts container"| ctr_grafana
    ctr_nginx -->|"uses shared Docker bridge network"| net_bridge
    ctr_go_service -->|"uses shared Docker bridge network"| net_bridge
    ctr_prometheus -->|"uses shared Docker bridge network"| net_bridge
    ctr_grafana -->|"uses shared Docker bridge network"| net_bridge
```

## Operational summary

Purpose: Modelar e governar a arquitetura pre-codebase do desafio Korp, preservando rastreabilidade entre requisitos, criterios de aceite, decisoes, componentes e validacao.
Compiled 29 projected items from operational revision 6.
