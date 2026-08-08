# TRACE-REQUIREMENT — Requirement Traceability

- Project: Projeto Korp DevOps Challenge
- Operational revision: 6
- Classification: `declared/proposed/inferred` pre-codebase state
- Projection: `flowchart`

```mermaid
flowchart LR
    act_candidate(["Candidate Operator"])
    act_evaluator(["Korp Evaluator"])
    ctr_go_service["http-server-projeto-korp"]
    ctr_grafana["Grafana"]
    ctr_nginx["NGINX Reverse Proxy"]
    ctr_prometheus["Prometheus"]
    net_bridge["Docker Bridge Network"]
    node_linux_host["Linux Host"]
    sys_korp["Korp Challenge System"]
    tech_ansible["Ansible"]
    tech_compose["Docker Compose"]
    ac_001["Endpoint principal via NGINX"]
    ac_002["Isolamento da aplicacao"]
    ac_003["Observabilidade funcional"]
    ac_004["Provisionamento reproduzivel"]
    ac_005["Demonstracao tecnica"]
    act_evaluator -->|"accesses and evaluates"| sys_korp
    act_candidate -->|"provisions and demonstrates"| sys_korp
    ctr_nginx -->|"forwards HTTP requests to port 8080"| ctr_go_service
    ctr_prometheus -->|"scrapes internal GET /metrics on service port 8080"| ctr_go_service
    ctr_grafana -->|"reads time series for dashboards"| ctr_prometheus
    node_linux_host -->|"runs as Docker container"| ctr_nginx
    node_linux_host -->|"runs as Docker container"| ctr_go_service
    node_linux_host -->|"runs as Docker container"| ctr_prometheus
    node_linux_host -->|"runs as Docker container"| ctr_grafana
    tech_ansible -->|"installs and configures the environment"| node_linux_host
    tech_compose -->|"declares and starts container"| ctr_nginx
    tech_compose -->|"declares and starts container"| ctr_go_service
    tech_compose -->|"declares and starts container"| ctr_prometheus
    tech_compose -->|"declares and starts container"| ctr_grafana
    ctr_nginx -->|"uses shared Docker bridge network"| net_bridge
    ctr_go_service -->|"uses shared Docker bridge network"| net_bridge
    ctr_prometheus -->|"uses shared Docker bridge network"| net_bridge
    ctr_grafana -->|"uses shared Docker bridge network"| net_bridge
    ac_001 -->|"primary verification or realization target"| ctr_go_service
    ac_001 -->|"primary verification or realization target"| ctr_nginx
    ac_002 -->|"primary verification or realization target"| ctr_go_service
    ac_002 -->|"primary verification or realization target"| ctr_nginx
    ac_002 -->|"primary verification or realization target"| net_bridge
    ac_003 -->|"primary verification or realization target"| ctr_go_service
    ac_003 -->|"primary verification or realization target"| ctr_grafana
    ac_003 -->|"primary verification or realization target"| ctr_prometheus
    ac_004 -->|"primary verification or realization target"| node_linux_host
    ac_004 -->|"primary verification or realization target"| tech_ansible
    ac_004 -->|"primary verification or realization target"| tech_compose
    ac_005 -->|"primary verification or realization target"| act_candidate
    ac_005 -->|"primary verification or realization target"| act_evaluator
    ac_005 -->|"primary verification or realization target"| sys_korp
```

## Operational summary

Purpose: Modelar e governar a arquitetura pre-codebase do desafio Korp, preservando rastreabilidade entre requisitos, criterios de aceite, decisoes, componentes e validacao.
Compiled 48 projected items from operational revision 6.
