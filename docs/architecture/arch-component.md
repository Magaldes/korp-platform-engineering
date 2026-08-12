# ARCH-COMPONENT — Components

- Project: Projeto Korp DevOps Challenge
- Operational revision: 7
- Classification: `declared/proposed/inferred` approved canonical as-built state
- Projection: `flowchart`

```mermaid
flowchart LR
    cmp_audit_service["Request Audit Service"]
    cmp_data_handler["Data Query HTTP Handler"]
    cmp_http_handler["Project Endpoint Handler"]
    cmp_metrics["Metrics Instrumentation"]
    cmp_postgres_repository["PostgreSQL Audit Repository Adapter"]
    cmp_redis_cache["Redis Statistics Cache Adapter"]
    cmp_stats_service["Request Statistics Service"]
    if_audit_repository["Audit Repository Interface"]
    if_data_http["Data Query HTTP Interface"]
    if_http["Project HTTP Interface"]
    if_metrics["Prometheus Metrics Interface"]
    if_stats_cache["Statistics Cache Interface"]
    cmp_http_handler -->|"implements the public project endpoint contract"| if_http
    cmp_metrics -->|"exports metric exposition surface"| if_metrics
    cmp_metrics -->|"observes application HTTP traffic except /metrics"| if_http
    cmp_data_handler -->|"O handler de dados implementa o contrato HTTP voluntario."| if_data_http
    cmp_data_handler -->|"O handler delega estatisticas ao servico de aplicacao."| cmp_stats_service
    cmp_data_handler -->|"O handler consulta historico por meio do servico de auditoria."| cmp_audit_service
    cmp_audit_service -->|"Auditoria depende do contrato interno de repository, nao do adapter PostgreSQL concreto."| if_audit_repository
    cmp_stats_service -->|"Estatisticas consultam a fonte persistente por contrato interno independente do driver."| if_audit_repository
    cmp_stats_service -->|"Estatisticas consultam e atualizam cache por contrato interno independente do Redis."| if_stats_cache
    cmp_metrics -->|"Metricas observam operacoes de auditoria/persistencia."| cmp_audit_service
    cmp_metrics -->|"Metricas observam cache e consultas agregadas."| cmp_stats_service
    cmp_postgres_repository -->|"PostgreSQL adapter implements the application repository contract."| if_audit_repository
    cmp_redis_cache -->|"Redis adapter implements the application cache contract."| if_stats_cache
    cmp_audit_service -->|"Audit service observes original application HTTP requests after status is known, excluding /metrics."| if_http
    cmp_audit_service -->|"Audit service also observes voluntary public data-query requests, excluding /metrics."| if_data_http
```

## Operational summary

Purpose: Modelar e governar a arquitetura implementada e validada do Projeto Korp, preservando rastreabilidade entre requisitos, decisoes, componentes, containers, dados, processos operacionais e o estado as-built da aplicacao.
Compiled 27 projected items from operational revision 7.
