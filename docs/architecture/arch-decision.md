# ARCH-DECISION — Architecture Decisions

- Project: Projeto Korp DevOps Challenge
- Operational revision: 7
- Classification: `declared/proposed/inferred` approved canonical as-built state
- Projection: `flowchart`

```mermaid
flowchart LR
    beh_audit_retention["Asynchronous audit retention purge"]
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
    dec_027{"Voluntary data extension boundary"}
    dec_028{"Request audit persistence use case"}
    dec_029{"Statistics cache optional optimization"}
    dec_030{"Data service deployment boundary"}
    dec_031{"Application persistence and cache separation"}
    dec_032{"Application-boundary dependency observability"}
    dec_033{"Grafana monitoring validation expansion"}
    dec_034{"Dedicated migration step"}
    dec_035{"Audit scope and data minimization"}
    dec_036{"Public history and statistics contracts"}
    dec_037{"Audit retention policy"}
    dec_038{"Development credential handling"}
    dec_039{"PostgreSQL runtime baseline"}
    dec_040{"PostgreSQL application driver"}
    dec_041{"PostgreSQL readiness and endpoint independence"}
    dec_042{"Request audit physical identity and timing"}
    dec_043{"Redis runtime and Go client baseline"}
    dec_044{"Statistics cache TTL"}
    dec_045{"Statistics cache key and payload semantics"}
    dec_046{"Application-owned audit retention purge"}
    cmp_audit_service["Request Audit Service"]
    cmp_data_handler["Data Query HTTP Handler"]
    cmp_http_handler["Project Endpoint Handler"]
    cmp_metrics["Metrics Instrumentation"]
    cmp_postgres_repository["PostgreSQL Audit Repository Adapter"]
    cmp_redis_cache["Redis Statistics Cache Adapter"]
    cmp_stats_service["Request Statistics Service"]
    ctr_go_service["http-server-projeto-korp"]
    ctr_grafana["Grafana"]
    ctr_nginx["NGINX Reverse Proxy"]
    ctr_postgres["PostgreSQL 18.4"]
    ctr_prometheus["Prometheus"]
    ctr_redis["Redis 8.8.1"]
    db_request_audit[("PostgreSQL Request Audit Store")]
    db_stats_cache[("Redis Statistics Cache")]
    ent_request_audit["Request Audit Event"]
    if_audit_repository["Audit Repository Interface"]
    if_data_http["Data Query HTTP Interface"]
    if_http["Project HTTP Interface"]
    if_metrics["Prometheus Metrics Interface"]
    if_stats_cache["Statistics Cache Interface"]
    net_bridge["Docker Bridge Network"]
    tech_ansible["Ansible"]
    tech_compose["Docker Compose"]
    tech_go_redis_v9["go-redis v9 Client"]
    tech_pgx_v5["pgx v5 PostgreSQL Driver"]
    tech_postgresql["PostgreSQL"]
    tech_redis["Redis"]
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
    ctr_go_service -->|"uses PostgreSQL as authoritative audit/statistics data source when repository is available"| ctr_postgres
    ctr_go_service -->|"uses Redis as optional statistics cache and can serve valid cache hits without PostgreSQL"| ctr_redis
    ctr_postgres -->|"hosts the authoritative request audit store"| db_request_audit
    ctr_redis -->|"hosts the volatile statistics cache store"| db_stats_cache
    tech_ansible -->|"waits boundedly for pg_isready before migrations"| ctr_postgres
    tech_ansible -->|"waits boundedly for redis-cli ping before application startup"| ctr_redis
    cmp_postgres_repository -->|"adapter reaches the PostgreSQL runtime on the internal Docker network"| ctr_postgres
    cmp_redis_cache -->|"adapter reaches the Redis runtime on the internal Docker network"| ctr_redis
    ctr_postgres -->|"runs PostgreSQL 18.4"| tech_postgresql
    ctr_redis -->|"runs Redis 8.8.1"| tech_redis
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
    dec_027 -->|"governs"| ctr_go_service
    dec_027 -->|"governs"| if_http
    dec_028 -->|"governs"| cmp_audit_service
    dec_028 -->|"governs"| db_request_audit
    dec_029 -->|"governs"| db_stats_cache
    dec_029 -->|"governs"| cmp_stats_service
    dec_030 -->|"governs"| db_stats_cache
    dec_030 -->|"governs internal-only PostgreSQL runtime and durable persistence boundary"| ctr_postgres
    dec_030 -->|"governs internal-only volatile Redis runtime"| ctr_redis
    dec_030 -->|"governs"| db_request_audit
    dec_030 -->|"governs"| net_bridge
    dec_031 -->|"governs"| cmp_data_handler
    dec_031 -->|"governs"| cmp_postgres_repository
    dec_031 -->|"governs"| cmp_redis_cache
    dec_032 -->|"governs"| cmp_metrics
    dec_032 -->|"governs"| ctr_prometheus
    dec_033 -->|"governs"| ctr_grafana
    dec_034 -->|"governs"| tech_ansible
    dec_034 -->|"governs"| tech_compose
    dec_035 -->|"governs"| cmp_audit_service
    dec_035 -->|"governs"| ent_request_audit
    dec_036 -->|"governs"| if_data_http
    dec_037 -->|"governs"| db_request_audit
    dec_038 -->|"governs"| tech_postgresql
    dec_039 -->|"fixes postgres:18.4-trixie runtime and persistent volume"| ctr_postgres
    dec_039 -->|"Defines internal port and durable volume boundary."| db_request_audit
    dec_039 -->|"Fixes PostgreSQL runtime/image/volume baseline."| tech_postgresql
    dec_040 -->|"Requires pgx-backed adapter implementation."| cmp_postgres_repository
    dec_040 -->|"Fixes pgx v5.10.0 and pgxpool."| tech_pgx_v5
    dec_041 -->|"Requires bounded PostgreSQL readiness and migration acceptance."| tech_ansible
    dec_041 -->|"defines bounded readiness and startup-degraded semantics"| ctr_postgres
    dec_041 -->|"allows valid Redis cache hits to serve statistics without repository"| ctr_redis
    dec_041 -->|"Defines database readiness semantics."| db_request_audit
    dec_041 -->|"Preserves Go startup and project endpoint independence from PostgreSQL."| ctr_go_service
    dec_042 -->|"Duration is captured before best-effort persistence."| cmp_audit_service
    dec_042 -->|"Defines physical identity and timing semantics."| ent_request_audit
    dec_043 -->|"Requires the Redis adapter to use the approved client behind the internal interface."| cmp_redis_cache
    dec_043 -->|"Defines volatile no-host-port cache deployment."| db_stats_cache
    dec_043 -->|"Fixes go-redis/v9 v9.21.0."| tech_go_redis_v9
    dec_043 -->|"fixes redis:8.8.1-trixie runtime with no durable volume"| ctr_redis
    dec_043 -->|"Fixes Redis 8.8.1 runtime and internal deployment boundary."| tech_redis
    dec_044 -->|"Requires Redis writes to use the configured TTL."| cmp_redis_cache
    dec_044 -->|"Defines configurable cache TTL with a 60 second default."| if_stats_cache
    dec_045 -->|"Defines versioned key and complete aggregate payload semantics."| db_stats_cache
    dec_045 -->|"Defines canonical cache-key construction and cache-hit response semantics."| cmp_stats_service
    dec_046 -->|"Assigns asynchronous retention ownership to the Request Audit Service."| cmp_audit_service
    dec_046 -->|"Defines the accepted asynchronous retention behavior."| beh_audit_retention
    dec_046 -->|"Extends the repository contract with purge-before-cutoff semantics."| if_audit_repository
    ctr_go_service -->|"O servico Go expoe as consultas voluntarias pela mesma fronteira interna."| if_data_http
    cmp_data_handler -->|"O handler de dados implementa o contrato HTTP voluntario."| if_data_http
    cmp_data_handler -->|"O handler delega estatisticas ao servico de aplicacao."| cmp_stats_service
    cmp_data_handler -->|"O handler consulta historico por meio do servico de auditoria."| cmp_audit_service
    cmp_audit_service -->|"Auditoria depende do contrato interno de repository, nao do adapter PostgreSQL concreto."| if_audit_repository
    cmp_stats_service -->|"Estatisticas consultam a fonte persistente por contrato interno independente do driver."| if_audit_repository
    cmp_stats_service -->|"Estatisticas consultam e atualizam cache por contrato interno independente do Redis."| if_stats_cache
    cmp_postgres_repository -->|"Adapter PostgreSQL acessa o store persistente."| db_request_audit
    cmp_redis_cache -->|"Adapter Redis acessa o cache."| db_stats_cache
    db_request_audit -->|"O banco persiste eventos de auditoria."| ent_request_audit
    db_request_audit -->|"O store utiliza PostgreSQL como tecnologia candidata."| tech_postgresql
    db_stats_cache -->|"O cache utiliza Redis como tecnologia candidata."| tech_redis
    tech_compose -->|"Compose declara e executa o container PostgreSQL da extensao."| ctr_postgres
    tech_compose -->|"Compose declara e executa o container Redis da extensao."| ctr_redis
    ctr_postgres -->|"PostgreSQL permanece na rede Docker interna."| net_bridge
    ctr_redis -->|"Redis permanece na rede Docker interna."| net_bridge
    cmp_metrics -->|"Metricas observam operacoes de auditoria/persistencia."| cmp_audit_service
    cmp_metrics -->|"Metricas observam cache e consultas agregadas."| cmp_stats_service
    cmp_postgres_repository -->|"PostgreSQL adapter implements the application repository contract."| if_audit_repository
    cmp_redis_cache -->|"Redis adapter implements the application cache contract."| if_stats_cache
    cmp_audit_service -->|"Audit service observes original application HTTP requests after status is known, excluding /metrics."| if_http
    cmp_audit_service -->|"Audit service also observes voluntary public data-query requests, excluding /metrics."| if_data_http
    cmp_postgres_repository -->|"PostgreSQL adapter uses pgx v5.10.0 and pgxpool."| tech_pgx_v5
    tech_pgx_v5 -->|"pgx connects the application adapter to PostgreSQL."| tech_postgresql
    cmp_redis_cache -->|"Redis Statistics Cache Adapter uses the approved Go Redis client."| tech_go_redis_v9
```

## Operational summary

Purpose: Modelar e governar a arquitetura implementada e validada do Projeto Korp, preservando rastreabilidade entre requisitos, decisoes, componentes, containers, dados, processos operacionais e o estado as-built da aplicacao.
Compiled 209 projected items from operational revision 7.
