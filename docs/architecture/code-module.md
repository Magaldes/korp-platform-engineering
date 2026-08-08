# CODE-MODULE — Modules and Dependencies

- Project: Projeto Korp DevOps Challenge
- Operational revision: 6
- Classification: `declared/proposed/inferred` pre-codebase state
- Projection: `flowchart`

```mermaid
flowchart LR
    cmp_http_handler["Project Endpoint Handler"]
    cmp_metrics["Metrics Instrumentation"]
    if_http["Project HTTP Interface"]
    if_metrics["Prometheus Metrics Interface"]
    mod_ansible_provisioning["Ansible Provisioning Module"]
    mod_compose_stack["Compose Stack Definition"]
    mod_docker_build["Application Docker Build"]
    mod_go_entrypoint["Go Service Entry Point"]
    mod_go_http_server["HTTP Server Module"]
    mod_go_metrics["Metrics Module"]
    mod_grafana_dashboard["Korp Observability Dashboard"]
    mod_grafana_dashboard_provider["Grafana Dashboard Provider"]
    mod_grafana_datasource["Grafana Datasource Provisioning"]
    mod_nginx_config["NGINX Ingress Configuration"]
    mod_prometheus_config["Prometheus Scrape Configuration"]
    mod_test_http_integration["HTTP Proxy Integration Test"]
    mod_test_observability_integration["Observability Integration Test"]
    tech_prometheus_go_client["Prometheus Go Client"]
    cmp_http_handler -->|"implements the public project endpoint contract"| if_http
    cmp_metrics -->|"exports metric exposition surface"| if_metrics
    cmp_metrics -->|"observes application HTTP traffic except /metrics"| if_http
    mod_go_entrypoint -->|"constructs the application HTTP handler"| mod_go_http_server
    mod_go_entrypoint -->|"wraps the application handler with metrics instrumentation"| mod_go_metrics
    mod_go_http_server -->|"materializes the Project Endpoint Handler component"| cmp_http_handler
    mod_go_metrics -->|"materializes the Metrics Instrumentation component"| cmp_metrics
    mod_go_metrics -->|"uses the Prometheus Go client library"| tech_prometheus_go_client
    mod_compose_stack -->|"builds the application service image from the Dockerfile"| mod_docker_build
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
    mod_test_http_integration -->|"verifies the public HTTP ingress behavior"| mod_nginx_config
    mod_test_observability_integration -->|"verifies Prometheus observability behavior"| mod_prometheus_config
    mod_test_observability_integration -->|"verifies Grafana datasource integration"| mod_grafana_datasource
    mod_test_observability_integration -->|"verifies the provisioned dashboard contract"| mod_grafana_dashboard
```

## Operational summary

Purpose: Modelar e governar a arquitetura pre-codebase do desafio Korp, preservando rastreabilidade entre requisitos, criterios de aceite, decisoes, componentes e validacao.
Compiled 41 projected items from operational revision 6.
