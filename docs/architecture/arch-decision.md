# ARCH-DECISION — Architecture Decisions

- Project: Projeto Korp DevOps Challenge
- Operational revision: 6
- Classification: `declared/proposed/inferred` pre-codebase state
- Projection: `flowchart`

```mermaid
flowchart LR
    beh_http_error["HTTP error responsibility and response flow"]
    dec_001{"HTTP ingress boundary"}
    dec_002{"Monitoring data flow"}
    dec_003{"Provisioning authority"}
    dec_004{"Docker network ownership"}
    dec_005{"Grafana provisioning"}
    dec_006{"UTC timestamp representation"}
    dec_007{"Observability contract"}
    dec_008{"NGINX application routing boundary"}
    dec_009{"HTTP route and method responsibility"}
    dec_010{"Minimum reverse proxy headers"}
    dec_011{"Runtime platform and technology baseline"}
    dec_012{"Repository module structure"}
    dec_013{"Verification and automated testing strategy"}
    dec_014{"Project endpoint HTTP success contract"}
    dec_015{"Baseline HTTP error representation"}
    dec_016{"Method not allowed contract"}
    dec_017{"External metrics ingress policy"}
    dec_018{"Application error pass through"}
    dec_019{"Gateway infrastructure error semantics"}
    dec_020{"Project endpoint handler component"}
    dec_021{"UTC clock as implementation detail"}
    dec_022{"HTTP metrics middleware separation"}
    dec_023{"Bounded unmatched route label"}
    dec_024{"Gateway errors excluded from Go request counter"}
    dec_025{"Host port 80 deployment preflight"}
    dec_026{"Error handling projection"}
    cmp_http_handler["Project Endpoint Handler"]
    cmp_metrics["Metrics Instrumentation"]
    ctr_go_service["http-server-projeto-korp"]
    ctr_grafana["Grafana"]
    ctr_nginx["NGINX Reverse Proxy"]
    ctr_prometheus["Prometheus"]
    if_http["Project HTTP Interface"]
    if_metrics["Prometheus Metrics Interface"]
    net_bridge["Docker Bridge Network"]
    tech_ansible["Ansible"]
    tech_compose["Docker Compose"]
    ctr_nginx -->|"forwards HTTP requests to port 8080"| ctr_go_service
    ctr_prometheus -->|"scrapes internal GET /metrics on service port 8080"| ctr_go_service
    ctr_grafana -->|"reads time series for dashboards"| ctr_prometheus
    ctr_go_service -->|"implements GET /projeto-korp"| if_http
    ctr_go_service -->|"publishes internal GET /metrics for Prometheus scraping"| if_metrics
    cmp_http_handler -->|"implements the public project endpoint contract"| if_http
    cmp_metrics -->|"exports metric exposition surface"| if_metrics
    tech_compose -->|"declares and starts container"| ctr_nginx
    tech_compose -->|"declares and starts container"| ctr_go_service
    tech_compose -->|"declares and starts container"| ctr_prometheus
    tech_compose -->|"declares and starts container"| ctr_grafana
    ctr_nginx -->|"uses shared Docker bridge network"| net_bridge
    ctr_go_service -->|"uses shared Docker bridge network"| net_bridge
    ctr_prometheus -->|"uses shared Docker bridge network"| net_bridge
    ctr_grafana -->|"uses shared Docker bridge network"| net_bridge
    cmp_metrics -->|"observes application HTTP traffic except /metrics"| if_http
    dec_001 -->|"keeps application port internal"| ctr_go_service
    dec_001 -->|"publishes host HTTP ingress"| ctr_nginx
    dec_002 -->|"visualizes Prometheus data"| ctr_grafana
    dec_002 -->|"collects service metrics"| ctr_prometheus
    dec_003 -->|"single provisioning entrypoint"| tech_ansible
    dec_004 -->|"Compose references external network"| tech_compose
    dec_004 -->|"network created by Ansible"| net_bridge
    dec_005 -->|"dashboard and datasource provisioning"| ctr_grafana
    dec_006 -->|"requires RFC3339 UTC in horario"| if_http
    dec_008 -->|"defines generic public application ingress and internal metrics boundary"| ctr_nginx
    dec_009 -->|"owns application route and method semantics"| ctr_go_service
    dec_009 -->|"owns proxy and infrastructure failure semantics"| ctr_nginx
    dec_010 -->|"defines minimum forwarded proxy headers"| ctr_nginx
    dec_013 -->|"owns deployment acceptance execution"| tech_ansible
    dec_013 -->|"requires unit and application behavior verification"| ctr_go_service
    dec_013 -->|"requires dashboard acceptance checks"| ctr_grafana
    dec_013 -->|"participates in HTTP integration and acceptance verification"| ctr_nginx
    dec_013 -->|"requires observability acceptance checks"| ctr_prometheus
    dec_014 -->|"defines GET and native HEAD success semantics"| if_http
    dec_015 -->|"avoids custom JSON error envelope"| if_http
    dec_016 -->|"requires 405 Allow header"| if_http
    dec_017 -->|"keeps metrics interface internal"| if_metrics
    dec_017 -->|"returns 404 for external /metrics access"| ctr_nginx
    dec_018 -->|"passes application errors through without interception"| ctr_nginx
    dec_019 -->|"defines 502 and 504 gateway semantics"| ctr_nginx
    dec_020 -->|"accepts endpoint handler as declared component"| cmp_http_handler
    dec_021 -->|"keeps UTC clock access as implementation detail"| cmp_http_handler
    dec_022 -->|"defines middleware instrumentation responsibility"| cmp_metrics
    dec_023 -->|"bounds route label cardinality with unmatched"| cmp_metrics
    dec_024 -->|"excludes NGINX gateway statuses from Go counter"| cmp_metrics
    dec_024 -->|"keeps gateway failures outside Go metrics"| ctr_nginx
    dec_025 -->|"requires non-destructive host port 80 preflight"| tech_ansible
    dec_025 -->|"requires host port 80 to be available before publication"| ctr_nginx
    dec_026 -->|"enables canonical error handling behavior projection"| beh_http_error
```

## Operational summary

Purpose: Modelar e governar a arquitetura pre-codebase do desafio Korp, preservando rastreabilidade entre requisitos, criterios de aceite, decisoes, componentes e validacao.
Compiled 88 projected items from operational revision 6.
