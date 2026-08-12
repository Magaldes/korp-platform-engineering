# CODE-MODULE — Modules and Dependencies

- Project: Projeto Korp DevOps Challenge
- Operational revision: 7
- Classification: `declared/proposed/inferred` approved canonical as-built state
- Projection: `flowchart`

```mermaid
flowchart LR
    if_audit_repository["Audit Repository Interface"]
    if_stats_cache["Statistics Cache Interface"]
    mod_ansible_provisioning["Ansible Provisioning Module"]
    mod_compose_stack["Compose Stack Definition"]
    mod_go_audit["Audit and Statistics Application Module"]
    mod_go_config["Application Data Configuration Module"]
    mod_go_data_http["Audit and Statistics HTTP Module"]
    mod_go_entrypoint["Go Service Entry Point"]
    mod_go_http_server["HTTP Server Module"]
    mod_go_metrics["Metrics Module"]
    mod_go_postgres_repository["PostgreSQL Repository Module"]
    mod_go_redis_cache["Redis Statistics Cache Module"]
    mod_grafana_dashboard["Korp Observability Dashboard"]
    mod_grafana_dashboard_provider["Grafana Dashboard Provider"]
    mod_grafana_datasource["Grafana Datasource Provisioning"]
    mod_nginx_config["NGINX Ingress Configuration"]
    mod_postgres_migrations["PostgreSQL Migration Scripts"]
    mod_prometheus_config["Prometheus Scrape Configuration"]
    tech_go_redis_v9["go-redis v9 Client"]
    tech_pgx_v5["pgx v5 PostgreSQL Driver"]
    mod_go_entrypoint -->|"constructs the application HTTP handler"| mod_go_http_server
    mod_go_entrypoint -->|"wraps the application handler with metrics instrumentation"| mod_go_metrics
    mod_compose_stack -->|"mounts the NGINX ingress configuration"| mod_nginx_config
    mod_compose_stack -->|"mounts the Prometheus scrape configuration"| mod_prometheus_config
    mod_compose_stack -->|"mounts Grafana datasource provisioning"| mod_grafana_datasource
    mod_compose_stack -->|"mounts Grafana dashboard provider configuration"| mod_grafana_dashboard_provider
    mod_compose_stack -->|"mounts the versioned observability dashboard"| mod_grafana_dashboard
    mod_grafana_dashboard_provider -->|"loads the versioned dashboard from the dashboard directory"| mod_grafana_dashboard
    mod_grafana_datasource -->|"connects Grafana to the Prometheus service configured by the stack"| mod_prometheus_config
    mod_nginx_config -->|"routes public traffic to the Go service on internal port 8080"| mod_go_entrypoint
    mod_prometheus_config -->|"scrapes the metrics endpoint exposed by the Go metrics module"| mod_go_metrics
    mod_ansible_provisioning -->|"validates, builds and starts the approved Compose stack"| mod_compose_stack
    mod_go_entrypoint -->|"loads data/cache configuration"| mod_go_config
    mod_go_entrypoint -->|"constructs the audit/statistics service independently from repository availability"| mod_go_audit
    mod_go_entrypoint -->|"creates optional PostgreSQL repository when startup connection succeeds"| mod_go_postgres_repository
    mod_go_entrypoint -->|"creates Redis statistics cache regardless of PostgreSQL availability"| mod_go_redis_cache
    mod_go_http_server -->|"routes public data-query endpoints to their HTTP handlers"| mod_go_data_http
    mod_go_data_http -->|"delegates audit history and statistics behavior"| mod_go_audit
    mod_go_audit -->|"depends on the repository abstraction"| if_audit_repository
    mod_go_audit -->|"depends on the cache abstraction"| if_stats_cache
    mod_go_postgres_repository -->|"uses pgx v5.10.0 and pgxpool"| tech_pgx_v5
    mod_go_redis_cache -->|"uses go-redis/v9 v9.21.0"| tech_go_redis_v9
    mod_ansible_provisioning -->|"runs versioned migrations after bounded readiness checks"| mod_postgres_migrations
```

## Operational summary

Purpose: Modelar e governar a arquitetura implementada e validada do Projeto Korp, preservando rastreabilidade entre requisitos, decisoes, componentes, containers, dados, processos operacionais e o estado as-built da aplicacao.
Compiled 43 projected items from operational revision 7.
