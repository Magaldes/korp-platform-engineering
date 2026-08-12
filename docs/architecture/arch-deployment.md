# ARCH-DEPLOYMENT — Deployment

- Project: Projeto Korp DevOps Challenge
- Operational revision: 7
- Classification: `declared/proposed/inferred` approved canonical as-built state
- Projection: `flowchart`

```mermaid
flowchart LR
    act_candidate(["Candidate Operator"])
    ctr_go_service["http-server-projeto-korp"]
    ctr_grafana["Grafana"]
    ctr_nginx["NGINX Reverse Proxy"]
    ctr_postgres["PostgreSQL 18.4"]
    ctr_postgres_migrate["PostgreSQL Migration Runner"]
    ctr_prometheus["Prometheus"]
    ctr_redis["Redis 8.8.1"]
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
    node_linux_host -->|"runs as one-shot Docker container"| ctr_postgres_migrate
    tech_compose -->|"declares the one-shot migration service"| ctr_postgres_migrate
    ctr_postgres_migrate -->|"uses the shared Docker bridge during migration"| net_bridge
    ctr_go_service -->|"uses PostgreSQL as authoritative audit/statistics data source when repository is available"| ctr_postgres
    ctr_go_service -->|"uses Redis as optional statistics cache and can serve valid cache hits without PostgreSQL"| ctr_redis
    ctr_postgres_migrate -->|"applies versioned PostgreSQL schema before application data capability is considered ready"| ctr_postgres
    tech_ansible -->|"waits boundedly for pg_isready before migrations"| ctr_postgres
    tech_ansible -->|"waits boundedly for redis-cli ping before application startup"| ctr_redis
    tech_ansible -->|"runs migrations after PostgreSQL and Redis readiness checks"| ctr_postgres_migrate
    tech_compose -->|"Compose declara e executa o container PostgreSQL da extensao."| ctr_postgres
    tech_compose -->|"Compose declara e executa o container Redis da extensao."| ctr_redis
    ctr_postgres -->|"PostgreSQL permanece na rede Docker interna."| net_bridge
    ctr_redis -->|"Redis permanece na rede Docker interna."| net_bridge
    node_linux_host -->|"O host Linux executa PostgreSQL como container."| ctr_postgres
    node_linux_host -->|"O host Linux executa Redis como container."| ctr_redis
```

## Operational summary

Purpose: Modelar e governar a arquitetura implementada e validada do Projeto Korp, preservando rastreabilidade entre requisitos, decisoes, componentes, containers, dados, processos operacionais e o estado as-built da aplicacao.
Compiled 47 projected items from operational revision 7.
