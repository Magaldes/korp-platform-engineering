# TRACE-REQUIREMENT — Requirement Traceability

- Project: Projeto Korp DevOps Challenge
- Operational revision: 7
- Classification: `declared/proposed/inferred` approved canonical as-built state
- Projection: `flowchart`

```mermaid
flowchart LR
    act_candidate(["Candidate Operator"])
    act_evaluator(["Korp Evaluator"])
    beh_audit_retention["Asynchronous audit retention purge"]
    cmp_audit_service["Request Audit Service"]
    cmp_metrics["Metrics Instrumentation"]
    cmp_stats_service["Request Statistics Service"]
    ctr_go_service["http-server-projeto-korp"]
    ctr_grafana["Grafana"]
    ctr_nginx["NGINX Reverse Proxy"]
    ctr_postgres["PostgreSQL 18.4"]
    ctr_postgres_migrate["PostgreSQL Migration Runner"]
    ctr_prometheus["Prometheus"]
    ctr_redis["Redis 8.8.1"]
    db_request_audit[("PostgreSQL Request Audit Store")]
    db_stats_cache[("Redis Statistics Cache")]
    if_data_http["Data Query HTTP Interface"]
    net_bridge["Docker Bridge Network"]
    tech_ansible["Ansible"]
    tech_compose["Docker Compose"]
    ac_001["Endpoint principal via NGINX"]
    ac_002["Isolamento da aplicacao"]
    ac_003["Observabilidade funcional"]
    ac_004["Provisionamento reproduzivel"]
    ac_005["Demonstracao tecnica"]
    ac_006["Persistencia de auditoria demonstravel"]
    ac_007["Cache com fallback demonstravel"]
    ac_008["Dashboard evidencia os testes de monitoramento"]
    ac_009["Provisionamento estendido reproduzivel"]
    ext_001["Adicionar camada de persistencia como extensao voluntaria"]
    ext_002["Persistir eventos de auditoria das requisicoes"]
    ext_003["Disponibilizar consulta historica e estatisticas agregadas"]
    ext_004["Adicionar cache opcional para estatisticas"]
    ext_005["Preservar independencia do endpoint original"]
    ext_006["Expandir relatorios de monitoramento no Grafana"]
    ext_007["Observar comportamento de persistencia e cache"]
    ext_008["Preservar provisionamento integral em um comando"]
    ext_009["Public audit history contract"]
    ext_010["Public request statistics contract"]
    ext_011["Audit retention policy"]
    ctr_nginx -->|"forwards HTTP requests to port 8080"| ctr_go_service
    ctr_prometheus -->|"scrapes internal GET /metrics on service port 8080"| ctr_go_service
    ctr_grafana -->|"reads time series for dashboards"| ctr_prometheus
    tech_compose -->|"declares and starts container"| ctr_nginx
    tech_compose -->|"declares and starts container"| ctr_go_service
    tech_compose -->|"declares and starts container"| ctr_prometheus
    tech_compose -->|"declares and starts container"| ctr_grafana
    ctr_nginx -->|"uses shared Docker bridge network"| net_bridge
    ctr_go_service -->|"uses shared Docker bridge network"| net_bridge
    ctr_prometheus -->|"uses shared Docker bridge network"| net_bridge
    ctr_grafana -->|"uses shared Docker bridge network"| net_bridge
    tech_compose -->|"declares the one-shot migration service"| ctr_postgres_migrate
    ctr_postgres_migrate -->|"uses the shared Docker bridge during migration"| net_bridge
    ctr_go_service -->|"uses PostgreSQL as authoritative audit/statistics data source when repository is available"| ctr_postgres
    ctr_go_service -->|"uses Redis as optional statistics cache and can serve valid cache hits without PostgreSQL"| ctr_redis
    ctr_postgres_migrate -->|"applies versioned PostgreSQL schema before application data capability is considered ready"| ctr_postgres
    ctr_postgres -->|"hosts the authoritative request audit store"| db_request_audit
    ctr_redis -->|"hosts the volatile statistics cache store"| db_stats_cache
    tech_ansible -->|"waits boundedly for pg_isready before migrations"| ctr_postgres
    tech_ansible -->|"waits boundedly for redis-cli ping before application startup"| ctr_redis
    tech_ansible -->|"runs migrations after PostgreSQL and Redis readiness checks"| ctr_postgres_migrate
    ctr_go_service -->|"O servico Go expoe as consultas voluntarias pela mesma fronteira interna."| if_data_http
    tech_compose -->|"Compose declara e executa o container PostgreSQL da extensao."| ctr_postgres
    tech_compose -->|"Compose declara e executa o container Redis da extensao."| ctr_redis
    ctr_postgres -->|"PostgreSQL permanece na rede Docker interna."| net_bridge
    ctr_redis -->|"Redis permanece na rede Docker interna."| net_bridge
    cmp_metrics -->|"Metricas observam operacoes de auditoria/persistencia."| cmp_audit_service
    cmp_metrics -->|"Metricas observam cache e consultas agregadas."| cmp_stats_service
    cmp_audit_service -->|"Audit service also observes voluntary public data-query requests, excluding /metrics."| if_data_http
    ac_001 -->|"primary verification or realization target"| ctr_go_service
    ac_001 -->|"primary verification or realization target"| ctr_nginx
    ac_002 -->|"primary verification or realization target"| ctr_go_service
    ac_002 -->|"primary verification or realization target"| ctr_nginx
    ac_002 -->|"primary verification or realization target"| net_bridge
    ac_003 -->|"primary verification or realization target"| ctr_go_service
    ac_003 -->|"primary verification or realization target"| ctr_grafana
    ac_003 -->|"primary verification or realization target"| ctr_prometheus
    ac_004 -->|"primary verification or realization target"| tech_ansible
    ac_004 -->|"primary verification or realization target"| tech_compose
    ac_005 -->|"primary verification or realization target"| act_candidate
    ac_005 -->|"primary verification or realization target"| act_evaluator
    ext_001 -->|"Proposed realization or verification target for the voluntary extension."| db_request_audit
    ext_002 -->|"Proposed realization or verification target for the voluntary extension."| cmp_audit_service
    ext_002 -->|"Proposed realization or verification target for the voluntary extension."| db_request_audit
    ext_003 -->|"Proposed realization or verification target for the voluntary extension."| if_data_http
    ext_004 -->|"Proposed realization or verification target for the voluntary extension."| cmp_stats_service
    ext_004 -->|"Proposed realization or verification target for the voluntary extension."| db_stats_cache
    ext_005 -->|"Proposed realization or verification target for the voluntary extension."| ctr_go_service
    ext_006 -->|"Proposed realization or verification target for the voluntary extension."| ctr_grafana
    ext_006 -->|"Proposed realization or verification target for the voluntary extension."| ctr_prometheus
    ext_007 -->|"Proposed realization or verification target for the voluntary extension."| cmp_metrics
    ext_011 -->|"The retention behavior materializes the seven-day audit retention requirement."| beh_audit_retention
    ext_008 -->|"Proposed realization or verification target for the voluntary extension."| tech_ansible
    ext_008 -->|"Proposed realization or verification target for the voluntary extension."| tech_compose
    ac_006 -->|"Proposed realization or verification target for the voluntary extension."| cmp_audit_service
    ac_006 -->|"Proposed realization or verification target for the voluntary extension."| db_request_audit
    ac_007 -->|"Proposed realization or verification target for the voluntary extension."| cmp_stats_service
    ac_007 -->|"Proposed realization or verification target for the voluntary extension."| db_stats_cache
    ac_007 -->|"Proposed realization or verification target for the voluntary extension."| db_request_audit
    ac_008 -->|"Proposed realization or verification target for the voluntary extension."| ctr_grafana
    ac_008 -->|"Proposed realization or verification target for the voluntary extension."| ctr_prometheus
    ac_008 -->|"Proposed realization or verification target for the voluntary extension."| cmp_metrics
    ac_009 -->|"Proposed realization or verification target for the voluntary extension."| tech_ansible
    ac_009 -->|"Proposed realization or verification target for the voluntary extension."| tech_compose
    ac_009 -->|"Proposed realization or verification target for the voluntary extension."| db_request_audit
    ac_009 -->|"Proposed realization or verification target for the voluntary extension."| db_stats_cache
    ext_009 -->|"The public data interface carries the approved audit history HTTP contract."| if_data_http
    ext_010 -->|"The public data interface carries the approved request statistics HTTP contract."| if_data_http
    ext_011 -->|"The PostgreSQL audit store is the persistence boundary governed by the retention policy."| db_request_audit
    ext_001 -->|"PostgreSQL runtime materializes the additional persistence layer"| ctr_postgres
    ext_004 -->|"Redis runtime materializes the optional cache layer"| ctr_redis
    ac_006 -->|"PostgreSQL container and durable volume support persistence demonstration"| ctr_postgres
    ac_007 -->|"Redis runtime supports cache hit and fallback verification"| ctr_redis
    ac_009 -->|"one-shot migration runner supports reproducible extended provisioning"| ctr_postgres_migrate
```

## Operational summary

Purpose: Modelar e governar a arquitetura implementada e validada do Projeto Korp, preservando rastreabilidade entre requisitos, decisoes, componentes, containers, dados, processos operacionais e o estado as-built da aplicacao.
Compiled 113 projected items from operational revision 7.
