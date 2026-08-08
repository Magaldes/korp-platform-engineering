# PROC-OPERATIONAL — Operational Process

- Project: Projeto Korp DevOps Challenge
- Operational revision: 6
- Classification: `declared/proposed/inferred` pre-codebase state
- Projection: `flowchart`

```mermaid
flowchart TB
    step_start(["Start Ansible command"])
    step_port80{"Check host TCP port 80 availability"}
    step_port80_fail["Fail deployment and report port 80 conflict without changing existing services"]
    step_docker["Install and configure Docker"]
    step_net["Create Docker bridge network"]
    step_build["Build Go service image"]
    step_compose["Create and run containers with Docker Compose"]
    step_nginx["Configure NGINX reverse proxy"]
    step_mon["Configure Prometheus and Grafana"]
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
    step_build --> step_compose
    step_compose --> step_nginx
    step_nginx --> step_mon
    step_mon --> step_verify
    step_verify --> step_output
    step_output --> step_end
```

## Operational summary

Purpose: Modelar e governar a arquitetura pre-codebase do desafio Korp, preservando rastreabilidade entre requisitos, criterios de aceite, decisoes, componentes e validacao.
Compiled 25 projected items from operational revision 6.
