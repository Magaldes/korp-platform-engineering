# PROC-OPERATIONAL — Operational Process

- Project: Projeto Korp DevOps Challenge
- Operational revision: 7
- Classification: `declared/proposed/inferred` approved canonical as-built state
- Projection: `flowchart`

```mermaid
flowchart TB
    step_start(["Start Ansible command"])
    step_port80{"Check host TCP port 80 availability"}
    step_port80_fail["Fail deployment and report port 80 conflict without changing existing services"]
    step_docker["Install and configure Docker"]
    step_net["Create Docker bridge network"]
    step_build["Validate Compose and build Go service image"]
    step_data["Start PostgreSQL and Redis"]
    step_pg_ready{"Wait for PostgreSQL readiness with bounded retries"}
    step_pg_fail["Fail extended provisioning after PostgreSQL readiness timeout"]
    step_redis_ready{"Wait for Redis readiness with bounded retries"}
    step_redis_fail["Fail extended provisioning after Redis readiness timeout"]
    step_migrate{"Run versioned PostgreSQL migrations"}
    step_migrate_fail["Fail extended provisioning without destructive volume reset"]
    step_app_mon["Start application Prometheus and Grafana"]
    step_nginx["Start NGINX reverse proxy"]
    step_verify["Execute HTTP validation request through host port 80"]
    step_output["Display response in console"]
    step_end(["Environment ready for demonstration"])
    step_blocked(["Human decision required before deployment"])
    step_start --> step_port80
    step_port80 -->|"Port 80 available"| step_docker
    step_port80 -->|"Port 80 occupied"| step_port80_fail
    step_port80_fail --> step_blocked
    step_docker --> step_net
    step_net --> step_build
    step_build --> step_data
    step_data --> step_pg_ready
    step_pg_ready -->|"PostgreSQL ready"| step_redis_ready
    step_pg_ready -->|"Readiness timeout"| step_pg_fail
    step_pg_fail --> step_blocked
    step_redis_ready -->|"Redis ready"| step_migrate
    step_redis_ready -->|"Readiness timeout"| step_redis_fail
    step_redis_fail --> step_blocked
    step_migrate -->|"Migration success"| step_app_mon
    step_migrate -->|"Migration failure"| step_migrate_fail
    step_migrate_fail --> step_blocked
    step_app_mon --> step_nginx
    step_nginx --> step_verify
    step_verify --> step_output
    step_output --> step_end
```

## Operational summary

Purpose: Modelar e governar a arquitetura implementada e validada do Projeto Korp, preservando rastreabilidade entre requisitos, decisoes, componentes, containers, dados, processos operacionais e o estado as-built da aplicacao.
Compiled 40 projected items from operational revision 7.
